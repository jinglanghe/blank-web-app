package dao

import (
	"fmt"
	"os"

	"gitlab.apulis.com.cn/hjl/blank-web-app/configs"
	"gitlab.apulis.com.cn/hjl/blank-web-app/internal/model"
	"gitlab.apulis.com.cn/hjl/blank-web-app/logging"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var database *gorm.DB

func Init() {
	var err error
	switch configs.Config.Datasource {
	case "sqlite":
		database, _ = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	case "mysql":
		break
	case "postgres":
		dbConf := configs.Config.Postgres
		dbPassword, exist := os.LookupEnv("POSTGRES_PASSWORD")
		if !exist {
			dbPassword = dbConf.Password
		}
		dsn := "host=" + dbConf.Host + " user=" + dbConf.Username + " password=" + dbPassword +
			" port=" + fmt.Sprintf("%d", dbConf.Port) + " sslmode=disable TimeZone=Asia/Shanghai"
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			logging.Fatal().Err(err).Msg("connect postgresql failed")
			return
		}

		dataName := ""
		result := database.Raw("SELECT u.datname  FROM pg_catalog.pg_database u where u.datname = ?", dbConf.DbName).Scan(&dataName)
		if result.Error != nil {
			logging.Fatal().Err(result.Error).Msg("")
			return
		}
		if len(dataName) == 0 {
			result = database.Exec("CREATE DATABASE " + dbConf.DbName)
			if result.Error != nil {
				logging.Fatal().Err(err).Msg("create database error")
				return
			}
		}
		dsn += " dbname=" + dbConf.DbName
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			logging.Fatal().Err(err).Msg("connect postgresql failed")
			return
		}

		// 设置连接池
		sqlDb, err := database.DB()
		if err != nil {
			logging.Fatal().Err(err).Msg("get sql db failed")
			return
		}
		sqlDb.SetMaxOpenConns(dbConf.MaxOpenConns)
		sqlDb.SetMaxIdleConns(dbConf.MaxIdleConns)
		logging.Info().Msg("PostgreSQL connected success")
	}

	database.AutoMigrate(
		&model.Setting{},
		&model.SysVersion{},
	)
	database = database.Debug()
}

func GetDB() *gorm.DB {
	return database
}
