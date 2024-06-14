package mysql

import (
	"github.com/xissg/online-judge/internal/config"
	"testing"
)

func TestMysqlClient(t *testing.T) {
	appConfig := config.LoadConfig()
	mysql := NewMysqlClient(appConfig.Database)

	//user := &entity.User{
	//	ID: 1,
	//	//UserName:     "test",
	//	//UserPassword: "test",
	//	//AvatarURL:    "test",
	//	//CreateTime:   time.Now().Format(time.RFC3339),
	//	//UpdateTime:   time.Now().Format(time.RFC3339),
	//	IsDelete: 0,
	//	UserRole: constant.ADMIN,
	//}
	//err := mysql.CreateUser(user)
	//err := mysql.DeleteUser(1)
	//if err != nil {
	//	t.Error()
	//}
	//mysql.UpdateUser(user)
	mysql.BanUser(1)
}

func TestQuestion(t *testing.T) {
	appConfig := config.LoadConfig()
	mysql := NewMysqlClient(appConfig.Database)
	//question := &entity.Question{
	//	ID: 1,
	//Title:       "test",
	//Content:     "test",
	//Tag:         `{"Java","Python","JavaScript"}`,
	//Answer:    "test1",
	//SubmitNum: 11,
	//AcceptNum: 2,
	//JudgeCase:   "test",
	//JudgeConfig: "test",
	//UserId: 2,
	//CreateTime:  time.Now().Format(time.RFC3339),
	//UpdateTime:  time.Now().Format(time.RFC3339),
	//IsDelete:    0,
	//}
	//err := mysql.CreateQuestion(question)
	//if err != nil {
	//	t.Error()
	//}
	//questions := mysql.GetQuestionById(1)
	//fmt.Println(questions)

	//err := mysql.UpdateQuestion(question)
	//if err != nil {
	//	panic(err)
	//}
	mysql.DeleteQuestion(1)
}

func TestSubmit(t *testing.T) {
	appConfig := config.LoadConfig()
	mysql := NewMysqlClient(appConfig.Database)
	//submit := &entity.Submit{
	//	ID: 1,
	//Language:    "Java",
	//Code:        "test",
	//JudgeResult: "test2",
	//QuestionId:  1,
	//UserId:      1,
	//Status:      "ACCEPTED",
	//CreateTime:  time.Now().Format(time.RFC3339Nano),
	//UpdateTime:  time.Now().Format(time.RFC3339Nano),
	//IsDelete:    0,
	//}
	//err := mysql.CreateSubmit(submit)
	//if err != nil {
	//	t.Error()
	//}
	//submit := mysql.GetSubmitList(1, 10)
	//fmt.Printf("%v\n", submit)
	//mysql.UpdateSubmit(submit)
	mysql.DeleteSubmit(1)
}
