package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/common"
	"github.com/xissg/online-judge/internal/model/entity"
	"github.com/xissg/online-judge/internal/model/request"
	"github.com/xissg/online-judge/internal/repository/docker"
	"github.com/xissg/online-judge/internal/utils"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type JudgeService interface {
	Run(id string)
}

type judgeService struct {
	docker   *docker.DockerClient
	question QuestionService
	submit   SubmitService
	ai       AIService
	logger   *zap.SugaredLogger
}

func NewJudgeService(docker *docker.DockerClient, question QuestionService, submit SubmitService, ai AIService, logger *zap.SugaredLogger) JudgeService {
	return &judgeService{
		docker:   docker,
		question: question,
		submit:   submit,
		ai:       ai,
		logger:   logger,
	}
}

func (s *judgeService) Run(id string) {
	var err error
	submitId, err := strconv.Atoi(id)
	if err != nil {
		s.logger.Errorf("Invalid submit ID: %s", id)
		return
	}

	//判题校验，已判题或正在判题的直接返回
	s.logger.Infof("start judge service")
	submit, err := s.submit.GetSubmitById(submitId)
	if err != nil || submit == nil {
		s.logger.Errorf("Failed to get submit by ID: %d, error: %v", submitId, err)
		return
	}
	if submit.Status != constant.WATING_STATUS {
		s.logger.Infof("already judged, status: %v", submit.Status)
		return
	}

	//更新判题状态
	s.logger.Infof("update submit judge status")
	err = s.submit.UpdateSubmit(&request.UpdateSubmit{
		ID:     submitId,
		Status: constant.JUDGING_STATUS,
	})

	//开始初始化必要信息
	s.logger.Infof("start init judge context")
	ctx, err := s.initJudgeContext("/app")
	if err != nil {
		s.logger.Errorf("%v", err)
		return
	}
	err = s.getSubmitInfo(submitId, ctx)
	err = s.getQuestionInfo(ctx.Question.QuestionId, ctx)
	s.logger.Infof("start init docker image")
	err = s.chooseImage(ctx)
	if err != nil {
		s.logger.Errorf("%v", err)
		return
	}

	//生成文件，沙箱执行，获取结果
	s.logger.Infof("start generate files")
	err = s.generateFiles(ctx)
	s.logger.Infof("start start sanbox")
	err = s.startSandbox(ctx)
	s.logger.Infof("start get result")
	err = s.getResult(ctx)
	if err != nil {
		s.logger.Errorf("%v", err)
		return
	}

	//进行判题
	s.logger.Infof("start judge")
	ok := s.judge(ctx)

	//更新结果
	s.logger.Infof("start store judge result")
	err = s.storeJudgeResults(ok, ctx)
	if err != nil {
		s.logger.Errorf("Failed to store judge results: %v", err)
		return
	}

	//文件清理
	s.logger.Infof("start remove generated files")
	err = s.removeFiles(ctx)
	err = s.removeContainer(ctx.Config.ContainerId)
	if err != nil {
		s.logger.Errorf("%v", err)
	}
}

// 指定docker中编译路径
func (s *judgeService) initJudgeContext(compileDir string) (*entity.JudgeContext, error) {
	if compileDir == "" {
		return nil, fmt.Errorf("service layer:judge,init judge context, %w", constant.ErrInvalidFilePath)
	}
	imageList := config.LoadConfig().Images
	return &entity.JudgeContext{
		Question: &entity.QuestionContext{
			Answer:      []string{},
			JudgeConfig: &common.Config{},
			JudgeCase:   []string{},
		},
		Config: &entity.ConfigContext{
			CompileDir: compileDir,
			ImageList:  imageList,
		},
		Result: &entity.ResultContext{
			Output: []string{},
		},
	}, nil

}

// 获取提交信息
func (s *judgeService) getSubmitInfo(submitId int, ctx *entity.JudgeContext) error {
	submitEntity, err := s.submit.GetSubmitById(submitId)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}
	submitResponse := utils.ConvertSubmitResponse(submitEntity)

	if submitResponse == nil {
		return errors.New("service layer:judge, submit response convert error")
	}
	ctx.Question.SubmitId = submitId
	ctx.Question.QuestionId = submitEntity.QuestionId
	ctx.Question.Language = strings.ToLower(submitEntity.Language)
	ctx.Question.Code = submitEntity.Code
	ctx.Question.Status = submitEntity.Status
	ctx.Question.JudgeResult = submitResponse.JudgeResult

	return nil
}

// 获取题目信息
func (s *judgeService) getQuestionInfo(questionId int, ctx *entity.JudgeContext) error {
	questionEntity, err := s.question.GetQuestionById(questionId)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}
	questionResponse := utils.ConvertQuestionResponse(questionEntity)
	if questionResponse == nil {
		return errors.New("service layer:judge, question response convert error")
	}
	ctx.Question.Title = questionEntity.Title
	ctx.Question.Content = questionEntity.Content
	ctx.Question.Answer = questionResponse.Answer
	ctx.Question.JudgeConfig = questionResponse.JudgeConfig
	ctx.Question.JudgeCase = questionResponse.JudgeCase
	return nil
}

// 根据编程语言选择对应的镜像,本地不存在该镜像则立即拉取镜像
func (s *judgeService) chooseImage(ctx *entity.JudgeContext) error {
	//判断镜像是否存在
	image, exists := ctx.Config.ImageList[ctx.Question.Language]
	if !exists {
		return fmt.Errorf("service layer: judge, unsupported language")
	}

	var flag bool
	imageList, err := s.docker.ImageList()
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}
	for _, img := range imageList {
		if img == image {
			flag = true
			break
		}
		flag = false
	}

	if flag == false {
		err = s.docker.ImagePull(image)
		if err != nil {
			return fmt.Errorf("service layer:judge -> %w", err)
		}
	}

	ctx.Config.Image = image
	return nil
}

// 生成源码文件和脚本文件
func (s *judgeService) generateFiles(ctx *entity.JudgeContext) error {
	//获取源码和脚本文件本机保存的绝对路径
	sourceFileDir := getDirAbsolutePath()
	//根据编程语言生成文件名
	sourceFileName, scriptFileName, execFileName := generateFileName(ctx.Question)

	//生成源码文件路径、脚本文件路径、docker中编译路径、执行路径和脚本文件执行路径
	sourceFilePath := filepath.Join(sourceFileDir, sourceFileName)
	scriptFilePath := filepath.Join(sourceFileDir, scriptFileName)
	compileFilePath := filepath.Join(ctx.Config.CompileDir, sourceFileName)
	execFilePath := filepath.Join(ctx.Config.CompileDir, execFileName)
	dockerScriptPath := filepath.Join(ctx.Config.CompileDir, scriptFileName)

	//保存文件路径和文件名
	ctx.Config.SourceFilePath = sourceFilePath
	ctx.Config.ScriptFilePath = scriptFilePath
	ctx.Config.CompileFilePath = compileFilePath
	ctx.Config.ExecFilePath = execFilePath
	ctx.Config.DockerScriptPath = dockerScriptPath

	//生成源码文件和脚本文件
	err := createFiles(ctx)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}
	return nil
}

// 开启容器，将源码和脚本文件复制进docker中执行脚本文件，进行编译执行
func (s *judgeService) startSandbox(ctx *entity.JudgeContext) error {
	//设置执行脚本的命令
	dockerScriptPath := ctx.Config.DockerScriptPath
	cmds := []string{
		"/bin/bash", dockerScriptPath,
	}

	//创建docker，指定工作目录和超时时间
	dstDir := ctx.Config.CompileDir
	timeOut := time.Second * 10
	containerId, err := s.docker.ContainerCreate(ctx.Config.Image, "", dstDir, cmds, timeOut)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}

	//复制源码文件
	sourceFilePath := ctx.Config.SourceFilePath
	err = s.docker.CopyToContainer(containerId, dstDir, sourceFilePath)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}

	//复制shell文件
	shellFilePath := ctx.Config.ScriptFilePath
	err = s.docker.CopyToContainer(containerId, dstDir, shellFilePath)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}

	//开启docker
	err = s.docker.ContainerStart(containerId)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}

	//保存docker ID
	ctx.Config.ContainerId = containerId
	return nil
}

// 获取判题结果,退出码，执行结果，执行时间，内存占用
func (s *judgeService) getResult(ctx *entity.JudgeContext) error {
	//判断docker是否在运行
	running, err := s.docker.IsContainerRunning(ctx.Config.ContainerId)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}

	//获取内存信息
	if running {
		ctx.Result.MemoryUsage, err = s.docker.ContainerStats(ctx.Config.ContainerId)
		if err != nil {
			return fmt.Errorf("service layer:judge -> %w", err)
		}
	}

	//等待docker运行结束
	chanResponse, chanErr := s.docker.ContainerWait(ctx.Config.ContainerId)
	select {
	case <-chanResponse:
	case <-chanErr:
		return fmt.Errorf("service layer:judge -> %w", err)
	}

	//获取退出码和执行时间
	ctx.Result.ExitCode, ctx.Result.ExecTime = s.docker.ContainerInspect(ctx.Config.ContainerId)
	if ctx.Result.ExitCode != 0 {
		return fmt.Errorf("service layer:judge, user code compile error %w", constant.ErrCompile)
	}

	//获取输出结果
	output, err := s.docker.ContainerLogs(ctx.Config.ContainerId)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}
	//解析输出结果，通过分割符将各个输出单独切割出来
	result := resultLogProcedure(output, len(ctx.Question.JudgeCase))
	ctx.Result.Output = result
	return err
}

// 正常判题通过后，会交由ai进行代码优化建议
func (s *judgeService) judge(ctx *entity.JudgeContext) bool {
	ok := s.normalJudge(ctx)
	if ok {
		_ = s.aiSuggestion(ctx)
	}
	return ok
}

// 判题服务,通过退出码，执行结果，执行时间占用内存来判断结果
func (s *judgeService) normalJudge(ctx *entity.JudgeContext) bool {
	var flag bool
	ctx.Question.Status = constant.FAILED_STATUS
	if ctx.Result.ExitCode != 0 {
		ctx.Question.JudgeResult = constant.COMPILE_ERR_RESULT
		return false
	}
	if int64(ctx.Question.JudgeConfig.TimeLimit) < ctx.Result.ExecTime {
		ctx.Question.JudgeResult = constant.TIME_LIMIT_EXCEED_RESULT
		return false
	}
	if ctx.Question.JudgeConfig.MemoryLimit < ctx.Result.MemoryUsage {
		ctx.Question.JudgeResult = constant.MEMORY_LIMIT_EXCEED_RESULT
		return false
	}
	for i := 0; i < len(ctx.Question.Answer); i++ {
		if ctx.Question.Answer[i] != ctx.Result.Output[i] {
			ctx.Question.JudgeResult = constant.WRONG_ANSWER_RESULT
			flag = false
			continue
		}
		flag = true
	}

	if flag == true {
		ctx.Question.Status = constant.SUCCESS_STATUS
		return true
	}
	return false
}

// 给通过的用户提供代码优化建议
func (s *judgeService) aiSuggestion(ctx *entity.JudgeContext) error {
	roleStr := "现在你是一位算法工程师，你将根据后面的题目的描述信息和实现代码，从时间复杂度和空间复杂度对代码做出评价，并给出代码的优化思路建议。对于接下来我给出的题目描述信息和题目实现代码，请严格按照如下格式给出回复(待优化代码和优化建议需要根据实际代码进行替换)：\n代码时间复杂度：xxx，空间复杂度：xxx。\n待优化代码：第x行到第x行：a=b;\n c = a;\n优化建议：c=b 直接将b赋值给c减少内存拷贝\n"
	judgeCase, _ := json.Marshal(ctx.Question.JudgeCase)
	msg := roleStr + "\n" + "题目标题：" + ctx.Question.Title + "\n" + ctx.Question.Content + string(judgeCase) + "\n" + "实现代码：" + ctx.Question.Code + "\n"
	err := s.ai.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}
	resp, err := s.ai.ReceiveMessage()
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}
	ctx.Question.JudgeResult = resp
	return nil
}

// 删除容器
func (s *judgeService) removeContainer(containerId string) error {
	return s.docker.ContainerRemove(containerId)
}

// 删除源码文件和脚本文件
func (s *judgeService) removeFiles(ctx *entity.JudgeContext) error {
	sourceFilePath := ctx.Config.SourceFilePath
	shellFilePath := ctx.Config.ScriptFilePath

	err := os.Remove(sourceFilePath)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}
	err = os.Remove(shellFilePath)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}
	return nil
}

// 存储判题结果
func (s *judgeService) storeJudgeResults(ok bool, ctx *entity.JudgeContext) (err error) {
	if ok {
		err = s.question.UpdateQuestionAcceptNum(ctx.Question.QuestionId)
	} else {
		err = s.question.UpdateQuestionSubmitNum(ctx.Question.QuestionId)
	}

	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}
	updateEntity := &request.UpdateSubmit{

		ID:          ctx.Question.SubmitId,
		JudgeResult: ctx.Question.JudgeResult,
		Status:      ctx.Question.Status,
	}

	err = s.submit.UpdateSubmit(updateEntity)
	if err != nil {
		return fmt.Errorf("service layer:judge -> %w", err)
	}
	return nil
}

// 获取绝对路径
func getDirAbsolutePath() string {
	relateDir := "public/submit_code"
	//测试路径
	//relateDir := "../../public/submit_code"
	execDir, _ := os.Getwd()
	return filepath.Join(execDir, relateDir)
}

// 生成源码，脚本，执行文件文件名
func generateFileName(ctx *entity.QuestionContext) (sourceFileName string, scriptFileName string, compileFileName string) {
	randomString := utils.GenerateRandomString(16)
	extensions := map[string]string{
		"java":       ".class",
		"javascript": ".js",
		"python":     ".py",
	}

	compileFileName = randomString
	if ext, exists := extensions[ctx.Language]; exists {
		compileFileName += ext
	}

	sourceFileName = fmt.Sprintf("%v.%v", randomString, ctx.Language)
	scriptFileName = randomString + ".sh"

	return
}

// 根据代码和语言生成对应的源码文件和脚本文件
func createFiles(ctx *entity.JudgeContext) error {
	//获取源码文件和脚本文件路径
	sourceFilePath := ctx.Config.SourceFilePath
	scriptFilePath := ctx.Config.ScriptFilePath

	//创建文件
	fd1, err := os.Create(sourceFilePath)
	fd2, err := os.Create(scriptFilePath)
	if err != nil {
		return fmt.Errorf("service layer: judge, creating file error %w", constant.ErrCreateFile)
	}

	defer fd1.Close()
	defer fd2.Close()

	//写入源码
	_, err = fd1.WriteString(ctx.Question.Code)
	if err != nil {
		os.Remove(sourceFilePath)
		return fmt.Errorf("service layer:judge -> %w", err)
	}

	//获取编译文件路径和执行文件路径
	compileFilePath := ctx.Config.CompileFilePath
	execFilePath := ctx.Config.ExecFilePath

	//根据编程语言选择对应的的编译命令和执行方式
	compileCmd, execCmd := chooseExecCommand(ctx.Question.Language, compileFilePath, execFilePath)
	if compileCmd == "" && execCmd == "" {
		return errors.New("not supported language")
	}

	//写入编译脚本命令
	_, err = fd2.WriteString(fmt.Sprintf("#!bin/bash\n"))
	_, err = fd2.WriteString(compileCmd + "\n")

	//写入代码执行参数,使用分割段来区分多个输入执行的多个输出
	for i := range ctx.Question.JudgeCase {
		//多个判题用例分隔起始符
		_, err = fd2.WriteString(fmt.Sprintf("echo '===%d START==='\n", i))

		//使用管道和echo命令将输入用例输入，使用 -e 选项来解释其中的转义字符
		_, err = fd2.WriteString(fmt.Sprintf("echo -e \"%s\"| %v\n", ctx.Question.JudgeCase[i], execCmd))

		//多个判题用例结束符
		_, err = fd2.WriteString(fmt.Sprintf("echo '===%d END==='\n", i))
		if err != nil {
			return err
		}
	}
	if err != nil {
		os.Remove(scriptFilePath)
		return fmt.Errorf("service layer:judge -> %w", err)
	}
	return nil
}

// 对docker的执行结果进行处理，去除不可见字符
func resultLogProcedure(out string, length int) []string {
	var res strings.Builder
	for _, runeValue := range out {
		if unicode.IsGraphic(runeValue) || runeValue == '\n' || runeValue == '\t' || runeValue == ' ' {
			res.WriteRune(runeValue)
		}
	}
	output := res.String()

	//解析输出结果，通过分割符将各个输出单独切割出来
	var result []string
	for i := 0; i < length; i++ {
		startMarker := fmt.Sprintf("===%d START===\n", i)
		endMarker := fmt.Sprintf("===%d END===\n", i)
		startIndex := bytes.Index([]byte(output), []byte(startMarker))
		endIndex := bytes.Index([]byte(output), []byte(endMarker))
		if startIndex != -1 && endIndex != -1 && endIndex > startIndex {
			commandOutput := output[startIndex+len(startMarker) : endIndex]
			result = append(result, strings.TrimSuffix(commandOutput, "\n"))
		}
	}

	return result
}

type Command struct {
	Compile string
	Execute string
}

// 根据不同编程语言选择不同的编译方式和执行方式
func chooseExecCommand(lan string, filePath string, execPath string) (compileCmd string, execCmd string) {
	//使用map进行存储查询，方便扩展
	commands := map[string]Command{
		constant.GO_LANGUAGE: {
			Compile: fmt.Sprintf("go build -o %v %v && chmod +x %v", execPath, filePath, execPath),
			Execute: fmt.Sprintf("%v", execPath),
		},
		constant.JAVA_LANGUAGE: {
			Compile: fmt.Sprintf("javac %v && chmod +x %v", filePath, execPath),
			Execute: fmt.Sprintf("java -jar %v", execPath),
		},
		constant.PYTHON_LANGUAGE: {
			Compile: "",
			Execute: fmt.Sprintf("python %v", execPath),
		},
		constant.JAVASCRIPT_LANGUAGE: {
			Compile: "",
			Execute: fmt.Sprintf("node %v", execPath),
		},
		constant.CPP_LANGUAGE: {
			Compile: fmt.Sprintf("g++ %v -o %v && chmod +x %v", filePath, execPath, execPath),
			Execute: fmt.Sprintf("%v", execPath),
		},
		constant.C_LANGUAGE: {
			Compile: fmt.Sprintf("gcc %v -o %v && chmod +x %v", filePath, execPath, execPath),
			Execute: fmt.Sprintf("%v", execPath),
		},
	}

	if cmd, exists := commands[lan]; exists {
		compileCmd = cmd.Compile
		execCmd = cmd.Execute
	} else {
		compileCmd = ""
		execCmd = ""
	}

	return
}
