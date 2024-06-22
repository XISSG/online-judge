package entity

import (
	"github.com/xissg/online-judge/internal/model/common"
)

type JudgeContext struct {
	Question *QuestionContext
	Config   *ConfigContext
	Result   *ResultContext
}
type QuestionContext struct {
	SubmitId    int
	QuestionId  int
	Language    string
	Title       string
	Content     string
	Code        string
	Answer      []string
	Status      string
	JudgeResult string
	JudgeCase   []string
	JudgeConfig *common.Config
}

type ConfigContext struct {
	SourceFilePath   string //源码文件路径
	ScriptFilePath   string //脚本文件路径
	CompileDir       string //docker中编译目录
	CompileFilePath  string //编译文件路径
	ExecFilePath     string //可执行文件路径
	DockerScriptPath string //docker中脚本文件的执行路径
	ImageList        map[string]string
	Image            string
	ContainerId      string
}

type ResultContext struct {
	ExitCode    int
	ExecTime    int64
	MemoryUsage uint64
	Output      []string
}
