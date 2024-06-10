package service

import (
	"github.com/xissg/online-judge/internal/config"
	"github.com/xissg/online-judge/internal/repository/docker"
	"github.com/xissg/online-judge/internal/repository/elastic"
	"github.com/xissg/online-judge/internal/repository/mysql"
	"github.com/xissg/online-judge/internal/repository/rabbitmq"
	"github.com/xissg/online-judge/internal/repository/redis"
	"strconv"
	"testing"
)

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

func TestAISuggest(t *testing.T) {
	ctx := initJudgeContext("/app")
	ctx.Question.Title = "输入描述:输入包括多组数据。每组数据第一行是两个整数N、M（N<=100，M<=10000），N表示成都的大街上有几个路口，标号为1的路口是商店所在地，标号为N的路口是赛场所在地，M则表示在成都有几条路。N=M=0表示输入结束。接下来M行，每行包括3个整数A，B，C（1<=A,B<=N,1<=C<=1000）,表示在路口A与路口B之间有一条路，我们的工作人员需要C分钟的时间走过这条路。输入保证至少存在1条商店到赛场的路线。\n输出描述:\n对于每组输入，输出一行，表示工作人员从商店走到赛场的最短时间"
	ctx.Question.Code = "#include<iostream>\n#include<cstdio>\n#include<algorithm>\n#include<cstring>\n#define INF 1000000\nusing namespace std;\nint city,road;\nint mapp[205][205];\nint spand[205];\nint select[205];\nvoid dij(int city)\n{\n    int minn,k,i,j;\n    memset(select,0,sizeof(select));\n    for(i=1;i<=city;i++)\n    spand[i]=mapp[1][i];\n    spand[1]=0;select[1]=1;\n    for(i=2;i<=city;i++)\n    {\n        minn=INF;k=-1;\n        for(j=1;j<=city;j++)\n        {\n            if(!select[j]&&spand[j]<minn)\n            {\n                k=j;minn=spand[j];\n            }\n        }\n        if(k==-1)break;\n        select[k]=1;\n        for(j=1;j<=city;j++)\n        {\n            if(spand[j]>spand[k]+mapp[k][j]&&!select[j])\n            spand[j]=spand[k]+mapp[k][j];\n        }\n    }\n    printf(\"%d\\n\",spand[city]);\n}\nint main()\n{\n    int i;int x,y,z;\n    while(scanf(\"%d%d\",&city,&road)!=EOF)\n    {\n        if(city==0||road==0)\n        break;\n        memset(mapp,INF,sizeof(mapp));\n        memset(spand,INF,sizeof(spand));\n        for(i=1;i<=road;i++)\n        {\n            scanf(\"%d%d%d\",&x,&y,&z);\n            if(mapp[x][y]>z||mapp[y][x]>z)\n            {\n                mapp[x][y]=z;mapp[y][x]=z;\n            }\n        }\n        dij(city);\n    }\n    return 0;\n}\n"
	aiSuggestion(ctx)
}

//func TestChooseImage(t *testing.T) {
//	docker := docker2.NewDockerClient()
//	service := NewJudgeService(docker)
//	ctx := initJudgeContext("/app")
//	ctx.Question.Language = constant.JAVA_LANGUAGE
//	err := service.chooseImage(ctx)
//	if err != nil {
//		panic(err)
//	}
//}

func TestJudge(t *testing.T) {

	appConfig := config.LoadConfig()

	dockerClient := docker.NewDockerClient()
	mysqlClient := mysql.NewMysqlClient(appConfig.Database)
	redisClient := redis.NewRedisClient(appConfig.Redis)
	rabbitMqClient := rabbitmq.NewRabbitMQClient(appConfig.RabbitMQ)
	esClient := elastic.NewElasticSearchClient(appConfig.Elasticsearch)

	questionSvc := NewQuestionService(mysqlClient, esClient, redisClient)
	submitSvc := NewSubmitService(mysqlClient, esClient, redisClient)
	judge := NewJudgeService(dockerClient, questionSvc, submitSvc)

	rabbitMqClient.ExchangeDeclare(appConfig.RabbitMQ)
	rabbitMqClient.QueueDeclareAndBind(appConfig.RabbitMQ)
	msgs, _ := rabbitMqClient.Consume(appConfig.RabbitMQ)
	ch := make(chan struct{})
	go func() {
		for msg := range msgs {
			id, _ := strconv.Atoi(string(msg.Body))
			judge.Start(id)
		}
	}()
	<-ch
}
