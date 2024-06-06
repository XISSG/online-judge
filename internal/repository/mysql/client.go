package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xissg/online-judge/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlClient struct {
	*gorm.DB
}

func NewMysqlClient(cfg config.DatabaseConfig) *MysqlClient {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(cfg.MaxIdleConns)       //设置最大连接数
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConns)       //设置最大的空闲连接数
	sqlDb.SetConnMaxLifetime(cfg.ConnMaxLifetime) //设置最大连接时间

	return &MysqlClient{
		db,
	}
}
