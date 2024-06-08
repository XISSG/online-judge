package entity

import (
	"github.com/xissg/online-judge/internal/model/common"
)

type Context struct {
	Question *QuestionContext
	Config   *ConfigContext
	Result   []*ResultContext
}
type QuestionContext struct {
	SubmitId    int
	QuestionId  int
	Language    string
	Code        string
	Answer      []string
	Status      string
	JudgeResult []string
	JudgeCase   []string
	JudgeConfig []common.Config
}
type ConfigContext struct {
	SourceFileName  string
	SourceFileDir   string
	ShellFileName   string
	CompileDir      string
	CompileFileName string
	Image           string
	ContainerId     []string
}

type ResultContext struct {
	ExitCode    int
	ExecTime    int64
	MemoryUsage uint64
	Output      string
}
