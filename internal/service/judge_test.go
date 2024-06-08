package service

import (
	"github.com/xissg/online-judge/internal/constant"
	docker2 "github.com/xissg/online-judge/internal/repository/docker"
	"testing"
)

func TestJudgeService(t *testing.T) {
	docker := docker2.NewDockerClient()
	judge := NewJudgeService(docker)
	ctx := initJudgeContext("/app", "my-golang-image")
	ctx.Language = constant.GO_LANGUAGE
	ctx.Code = "package main\n\nimport \"fmt\"\n\nfunc main(){\n\tfmt.Println(\"hello world\")\n}\n"
	judge.generateFiles(ctx)
	judge.StartSandbox(ctx)
}
