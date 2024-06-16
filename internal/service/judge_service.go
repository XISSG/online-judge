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
	submitId, _ := strconv.Atoi(id)

	//判题校验，已判题或正在判题的直接返回
	s.logger.Infof("start judge service")
	submit, _ := s.submit.GetSubmitById(submitId)
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
	ctx := initJudgeContext("/app")
	err = s.receiveSubmit(submitId, ctx)
	err = s.receiveQuestion(ctx.Question.QuestionId, ctx)
	s.logger.Infof("start init docker image")
	err = s.chooseImage(ctx)
	if err != nil {
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
		return
	}

	//进行判题
	s.logger.Infof("start judge")
	ok := s.judge(ctx)

	s.logger.Infof("start store judge result")
	ch := make(chan error)
	go func() {
		//更新结果
		err = s.storeJudgeResults(ok, ctx)
		ch <- err
		close(ch)
	}()
	//文件清理
	s.logger.Infof("start remove generated files")
	err = s.removeFiles(ctx)
	err = s.removeContainer(ctx.Config.ContainerId)
	err = <-ch
	if err != nil {
		return
	}
}

func initJudgeContext(compileDir string) *entity.JudgeContext {
	return &entity.JudgeContext{
		Question: &entity.QuestionContext{
			Answer:      []string{},
			JudgeConfig: &common.Config{},
			JudgeCase:   []string{},
		},
		Config: &entity.ConfigContext{
			CompileDir: compileDir,
		},
		Result: &entity.ResultContext{
			Output: []string{},
		},
	}

}

// 获取提交信息
func (s *judgeService) receiveSubmit(submitId int, ctx *entity.JudgeContext) error {
	submitEntity, err := s.submit.GetSubmitById(submitId)
	if err != nil {
		return err
	}
	submitResponse := utils.ConvertSubmitResponse(submitEntity)

	ctx.Question.SubmitId = submitId
	ctx.Question.QuestionId = submitEntity.QuestionId
	ctx.Question.Language = strings.ToLower(submitEntity.Language)
	ctx.Question.Code = submitEntity.Code
	ctx.Question.Status = submitEntity.Status
	ctx.Question.JudgeResult = submitResponse.JudgeResult

	return nil
}

// 获取题目信息
func (s *judgeService) receiveQuestion(questionId int, ctx *entity.JudgeContext) error {
	questionEntity, err := s.question.GetQuestionById(questionId)
	if err != nil {
		return err
	}
	questionResponse := utils.ConvertQuestionResponse(questionEntity)

	ctx.Question.Title = questionEntity.Title
	ctx.Question.Content = questionEntity.Content
	ctx.Question.Answer = questionResponse.Answer
	ctx.Question.JudgeConfig = questionResponse.JudgeConfig
	ctx.Question.JudgeCase = questionResponse.JudgeCase
	return nil
}

// 根据编程语言选择对应的镜像,本地不存在该镜像则立即拉取镜像
func (s *judgeService) chooseImage(ctx *entity.JudgeContext) error {
	appConfig := config.LoadConfig()
	image := appConfig.Image[ctx.Question.Language]
	if image == "" {
		return errors.New("not supported language")
	}

	var flag bool
	imageList := s.docker.ImageList()
	for _, img := range imageList {
		if img == image {
			flag = true
			break
		}
		flag = false
	}

	if flag == false {
		s.docker.ImagePull(image)
	}

	ctx.Config.Image = image
	return nil
}

// 生成源码文件和脚本文件
func (s *judgeService) generateFiles(ctx *entity.JudgeContext) error {
	absoluteDir := getDirAbsolutePath()
	//生成文件名
	sourceFileName, shellFileName, execFileName := generateFileName(ctx.Question)
	//生成文件路径
	sourceFilePath := filepath.Join(absoluteDir, sourceFileName)
	shellFilePath := filepath.Join(absoluteDir, shellFileName)

	//保存文件路径和文件名
	ctx.Config.SourceFileDir = absoluteDir
	ctx.Config.SourceFileName = sourceFileName
	ctx.Config.ShellFileName = shellFileName
	ctx.Config.ExecFileName = execFileName

	//创建文件
	fd1, err := os.Create(sourceFilePath)
	fd2, err := os.Create(shellFilePath)
	if err != nil {
		return err
	}

	defer fd1.Close()
	defer fd2.Close()

	//写入源码
	_, err = fd1.WriteString(ctx.Question.Code)
	if err != nil {
		os.Remove(sourceFilePath)
		return err
	}

	//写入编译脚本命令
	compileFilePath := filepath.Join(ctx.Config.CompileDir, sourceFileName)
	execFilePath := filepath.Join(ctx.Config.CompileDir, execFileName)

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
		os.Remove(shellFilePath)
		return err
	}
	return nil
}

// 开启容器，将源码和脚本文件复制进docker中执行脚本文件，进行编译执行
func (s *judgeService) startSandbox(ctx *entity.JudgeContext) error {
	//设置执行脚本的命令
	dockerShellPath := filepath.Join(ctx.Config.CompileDir, ctx.Config.ShellFileName)
	cmds := []string{
		"/bin/bash", dockerShellPath,
	}

	//复制文件的docker路径
	dstDir := ctx.Config.CompileDir
	timeOut := time.Second * 10
	containerId := s.docker.ContainerCreate(ctx.Config.Image, "", dstDir, cmds, timeOut)
	//源码的文件路径
	sourceFilePath := filepath.Join(ctx.Config.SourceFileDir, ctx.Config.SourceFileName)
	err := s.docker.CopyToContainer(containerId, dstDir, sourceFilePath)
	if err != nil {
		return err
	}

	//shell文件的文件路径
	shellFilePath := filepath.Join(ctx.Config.SourceFileDir, ctx.Config.ShellFileName)
	err = s.docker.CopyToContainer(containerId, dstDir, shellFilePath)
	if err != nil {
		return err
	}

	//开启docker
	err = s.docker.ContainerStart(containerId)
	if err != nil {
		return err
	}

	//保存docker ID
	ctx.Config.ContainerId = containerId
	return nil
}

// 获取判题结果,退出码，执行结果，执行时间，内存占用
func (s *judgeService) getResult(ctx *entity.JudgeContext) (err error) {
	//等待docker运行结束
	chanResponse, chanErr := s.docker.ContainerWait(ctx.Config.ContainerId)
	select {
	case <-chanResponse:
	case <-chanErr:
		return err
	}
	ctx.Result.ExitCode, ctx.Result.ExecTime = s.docker.ContainerInspect(ctx.Config.ContainerId)
	ctx.Result.MemoryUsage, err = s.docker.ContainerStats(ctx.Config.ContainerId)
	output, err := s.docker.ContainerLogs(ctx.Config.ContainerId)
	if err != nil {
		return
	}
	//解析输出结果，通过分割符将各个输出单独切割出来
	result := resultLogProcedure(output, len(ctx.Question.JudgeCase))
	ctx.Result.Output = result
	return
}

// 判题服务,通过退出码，执行结果，执行时间占用内存来判断结果
func (s *judgeService) judge(ctx *entity.JudgeContext) bool {
	ok := s.normalJudge(ctx)
	if ok {
		s.aiSuggestion(ctx)
	}
	return ok
}

// 正常判题逻辑：比对答案
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
	s.ai.SendMessage(msg)
	resp, err := s.ai.ReceiveMessage()
	if err != nil {
		return err
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
	sourceFilePath := filepath.Join(ctx.Config.SourceFileDir, ctx.Config.SourceFileName)
	shellFilePath := filepath.Join(ctx.Config.SourceFileDir, ctx.Config.ShellFileName)

	err := os.Remove(sourceFilePath)
	if err != nil {
		return err
	}
	err = os.Remove(shellFilePath)
	if err != nil {
		return err
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
		return err
	}
	updateEntity := &request.UpdateSubmit{

		ID:          ctx.Question.SubmitId,
		JudgeResult: ctx.Question.JudgeResult,
		Status:      ctx.Question.Status,
	}

	return s.submit.UpdateSubmit(updateEntity)
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
func generateFileName(ctx *entity.QuestionContext) (string, string, string) {
	randomString := utils.GenerateRandomString(16)
	var compileFileName string
	switch ctx.Language {
	case "go", "c", "cpp":
		compileFileName = randomString
	case "java":
		compileFileName = randomString + ".class"
	case "javascript":
		compileFileName = randomString + ".js"
	case "python":
		compileFileName = randomString + ".py"
	}
	return fmt.Sprintf("%v.%v", randomString, ctx.Language), randomString + ".sh", compileFileName
}

// TODO:待优化
// 根据不同编程语言选择不同的编译方式和执行方式
func chooseExecCommand(lan string, filePath string, execPath string) (compileCmd string, execCmd string) {
	switch lan {
	case constant.GO_LANGUAGE:
		compileCmd = fmt.Sprintf("go build  -o %v %v && chmod +x %v", execPath, filePath, execPath)
		execCmd = fmt.Sprintf("%v", execPath)

	case constant.JAVA_LANGUAGE:
		compileCmd = fmt.Sprintf("javac %v && chmod +x %v", filePath, execPath)
		execCmd = fmt.Sprintf("java -jar %v", execPath)
	case constant.PYTHON_LANGUAGE:
		compileCmd = ""
		execCmd = fmt.Sprintf("python %v", execPath)
	case constant.JAVASCRIPT_LANGUAGE:
		compileCmd = ""
		execCmd = fmt.Sprintf("node %v", execPath)
	case constant.CPP_LANGUAGE:
		compileCmd = fmt.Sprintf("g++ %v -o %v && chmod +x %v", filePath, execPath, execPath)
		execCmd = fmt.Sprintf("%v\n", execPath)
	case constant.C_LANGUAGE:
		compileCmd = fmt.Sprintf("gcc %v -o %v && chmod +x %v", filePath, execPath, execPath)
		execCmd = fmt.Sprintf("%v", execPath)
	default:
		compileCmd = ""
		execCmd = ""
	}
	return
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
