package service

//func TestJudgeService(t *testing.T) {
//	docker := docker2.NewDockerClient()
//	judge := NewJudgeService(docker)
//	ctx := initJudgeContext("/", "golang")
//	ctx.Question.Language = constant.GO_LANGUAGE
//	ctx.Question.Code = "package main\n\nimport (\n\t\"fmt\"\n\t\"os\"\n)\n\nfunc main() {\n\targs := os.Args[1:]\n\tfor _,arg :=range args {\n\t\tfmt.Println(arg)\n\t}\n}"
//	ctx.Question.JudgeCase = []string{"hello A", "hello B"}
//	err := judge.generateFiles(ctx)
//	if err != nil {
//		panic(err)
//	}
//	err = judge.startSandbox(ctx)
//	if err != nil {
//		panic(err)
//	}
//
//	err = judge.getResult(ctx)
//	if err != nil {
//		panic(err)
//	}
//
//	err = judge.removeContainer(ctx.Config.ContainerId)
//	if err != nil {
//		panic(err)
//	}
//}
