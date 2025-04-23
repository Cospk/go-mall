package dao

import (
	"fmt"
	"github.com/Cospk/go-mall/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _DBMaster *gorm.DB
var _DBSlave *gorm.DB

func DB() *gorm.DB {
	return _DBSlave
}

func DBMaster() *gorm.DB {
	return _DBMaster
}

func InitGorm() {
	_DBMaster = initDB(config.Database.Master)
	_DBSlave = initDB(config.Database.Slave)
}

// initDB 初始化数据库
func initDB(option config.DbConnectOption) *gorm.DB {
	fmt.Println(option)
	// 创建数据库连接（参数：驱动和日志记录器）
	db, err := gorm.Open(
		getDialector(option.Type, option.DSN),
		&gorm.Config{
			Logger: NewGormLogger(),
		},
	)
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	// 设置连接池参数
	sqlDB.SetMaxOpenConns(option.MaxOpenConn)
	sqlDB.SetMaxIdleConns(option.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(option.MaxLifeTime)
	if err = sqlDB.Ping(); err != nil {
		// 测试连接，若失败了直接panic
		panic(err)
	}
	return db
}

func getDialector(t, dsn string) gorm.Dialector {
	switch t {
	case "mysql":
		return mysql.Open(dsn)
	default:
		return mysql.Open(dsn)
	}
}
