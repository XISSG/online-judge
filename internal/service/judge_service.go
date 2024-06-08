package service

import (
	"fmt"
	"github.com/xissg/online-judge/internal/constant"
	"github.com/xissg/online-judge/internal/model/common"
	"github.com/xissg/online-judge/internal/model/entity"
	"github.com/xissg/online-judge/internal/repository/docker"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type JudgeService interface {
	Start(submitId int) error
}

type judgeService struct {
	docker *docker.DockerClient
	mysql  *mysql.MysqlClient
}

func NewJudgeService(docker *docker.DockerClient, mysql *mysql.MysqlClient) JudgeService {
	return &judgeService{
		docker: docker,
		mysql:  mysql,
	}
}

func (s *judgeService) Start(submitId int) error {
	//TODO:从rabbitmq中接收提交id
	ctx := initJudgeContext("/app", "my-golang-image")
	s.receiveSubmit(submitId, ctx)
	s.receiveQuestion(submitId, ctx)

	if err := s.generateFiles(ctx); err != nil {
		return err
	}
	//TODO:支持多个测试用例输入
	for i := 0; i < len(ctx.Question.JudgeCase); i++ {
		if err := s.startSandbox(ctx); err != nil {
			return err
		}
		result, err := s.getResult(ctx.Config.ContainerId[i])
		if err != nil {
			return err
		}
		ctx.Result = append(ctx.Result, result)
	}
	ok := s.judge(ctx)
	s.storeJudgeResults(ok, ctx)
	return nil
}

func initJudgeContext(compileDir string, image string) *entity.Context {
	return &entity.Context{
		Question: &entity.QuestionContext{
			SubmitId:    0,
			QuestionId:  0,
			Language:    "",
			Code:        "",
			Answer:      []string{},
			JudgeResult: []string{},
			JudgeConfig: []common.Config{},
			JudgeCase:   []string{},
		},
		Config: &entity.ConfigContext{
			SourceFileName:  "",
			ShellFileName:   "",
			SourceFileDir:   "",
			CompileFileName: "",
			CompileDir:      compileDir,
			Image:           image,
			ContainerId:     []string{},
		},
		Result: []*entity.ResultContext{},
	}

}

func (s *judgeService) receiveSubmit(submitId int, ctx *entity.Context) {
	submitEntity := s.mysql.GetSubmitById(submitId)
	submitResponse := utils.ConvertSubmitResponse(submitEntity)

	ctx.Question.SubmitId = submitId
	ctx.Question.QuestionId = submitEntity.QuestionId
	ctx.Question.Language = strings.ToLower(submitEntity.Language)
	ctx.Question.Code = submitEntity.Code
	ctx.Question.Status = submitEntity.Status
	ctx.Question.JudgeResult = submitResponse.JudgeResult
}

func (s *judgeService) receiveQuestion(questionId int, ctx *entity.Context) {
	questionEntity := s.mysql.GetQuestionById(questionId)
	questionResponse := utils.ConvertQuestionResponse(questionEntity)

	ctx.Question.Answer = questionResponse.Answer
	ctx.Question.JudgeConfig = questionResponse.JudgeConfig
	ctx.Question.JudgeCase = questionResponse.JudgeCase
}

func (s *judgeService) generateFiles(ctx *entity.Context) error {
	absoluteDir := getDirAbsolutePath()
	sourceFileName, shellFileName, compileFileName := generateFileName(ctx.Question)
	sourceFilePath := filepath.Join(absoluteDir, sourceFileName)
	shellFilePath := filepath.Join(absoluteDir, shellFileName)

	ctx.Config.SourceFileDir = absoluteDir
	ctx.Config.SourceFileName = sourceFileName
	ctx.Config.ShellFileName = shellFileName
	ctx.Config.CompileFileName = compileFileName

	fd1, err := os.Create(sourceFilePath)
	fd2, err := os.Create(shellFilePath)
	if err != nil {
		return err
	}

	defer fd1.Close()
	defer fd2.Close()

	_, err = fd1.WriteString(ctx.Question.Code)
	if err != nil {
		os.Remove(sourceFilePath)
		return err
	}

	compileFilePath := filepath.Join(ctx.Config.CompileDir, sourceFileName)
	execFilePath := filepath.Join(ctx.Config.CompileDir, compileFileName)
	//TODO:优化代码支持外部输入示例输入
	_, err = fd2.WriteString(fmt.Sprintf("#bin/bash\t\ngo build %v && %v", compileFilePath, execFilePath))
	if err != nil {
		os.Remove(shellFilePath)
		return err
	}
	return nil
}

func (s *judgeService) startSandbox(ctx *entity.Context) error {
	//TODO:策略模式
	dockerShellPath := filepath.Join(ctx.Config.CompileDir, ctx.Config.ShellFileName)
	cmds := []string{
		"/bin/bash", dockerShellPath,
	}
	containerId := s.docker.ContainerCreate(ctx.Config.Image, "", cmds, time.Second*10)
	dstDir := ctx.Config.CompileDir
	sourceFilePath := filepath.Join(ctx.Config.SourceFileDir, ctx.Config.SourceFileName)
	err := s.docker.CopyToContainer(containerId, dstDir, sourceFilePath)
	if err != nil {
		return err
	}

	shellFilePath := filepath.Join(ctx.Config.SourceFileDir, ctx.Config.ShellFileName)
	err = s.docker.CopyToContainer(containerId, dstDir, shellFilePath)
	if err != nil {
		return err
	}
	err = s.docker.ContainerStart(containerId)
	if err != nil {
		s.docker.ContainerRemove(containerId)
		panic(err)
	}

	ctx.Config.ContainerId = append(ctx.Config.ContainerId, containerId)
	return nil
}

func (s *judgeService) getResult(containerId string) (result *entity.ResultContext, err error) {
	result.ExitCode, result.ExecTime = s.docker.ContainerInspect(containerId)
	result.MemoryUsage, err = s.docker.ContainerStats(containerId)
	result.Output, err = s.docker.ContainerLogs(containerId)
	return
}

func (s *judgeService) judge(ctx *entity.Context) bool {
	//TODO:使用ai进行判题
	var flag bool
	for i := 0; i < len(ctx.Question.Answer); i++ {
		if ctx.Result[i].ExitCode != 0 {
			ctx.Question.JudgeResult[i] = constant.COMPILE_ERR_RESULT
			flag = false
			continue
		}
		if ctx.Question.Answer[i] != ctx.Result[i].Output {
			ctx.Question.JudgeResult[i] = constant.WRONG_ANSWER_RESULT
			flag = false
			continue
		}
		if int64(ctx.Question.JudgeConfig[i].TimeLimit) < ctx.Result[i].ExecTime {
			ctx.Question.JudgeResult[i] = constant.TIME_LIMIT_EXCEED_RESULT
			flag = false
			continue
		}
		if ctx.Question.JudgeConfig[i].MemoryLimit < ctx.Result[i].MemoryUsage {
			ctx.Question.JudgeResult[i] = constant.MEMORY_LIMIT_EXCEED_RESULT
			flag = false
			continue
		}
		flag = true
	}

	if flag == false {
		ctx.Question.Status = constant.FAILED_STATUS
		return false
	}

	ctx.Question.Status = constant.SUCCESS_STATUS
	return true
}

func (s *judgeService) storeJudgeResults(ok bool, ctx *entity.Context) {

}

func getDirAbsolutePath() string {
	//测试文件中执行是测试文件的路径，main中执行是项目的路径
	//relateDir := "public/submit_code"
	relateDir := "../../public/submit_code"
	execDir, _ := os.Getwd()
	return filepath.Join(execDir, relateDir)
}

func generateFileName(ctx *entity.QuestionContext) (string, string, string) {
	randomString := utils.GenerateRandomString(6)
	var compileFileName string
	switch ctx.Language {
	case "go", "c", "cpp":
		compileFileName = randomString

	case "java":
		compileFileName = randomString + ".class"
	case "python":
		compileFileName = randomString + ".py"
	}

	return fmt.Sprintf("%v.%v", randomString, ctx.Language), randomString + ".sh", compileFileName
}
