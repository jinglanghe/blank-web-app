package configs

import (
	"encoding/json"
	"github.com/spf13/viper"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/internal/metadata"
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/logging"
)

var (
	Config   AppConfig
)

type LogConfig struct {
	LogLevel string `mapstructure:"log_level"`
}

type JWTConfig struct {
	SignAlgorithm string `mapstructure:"algo"`
	SecretKey     string `mapstructure:"secret_key"`
	PublicKeyFile string `mapstructure:"public_key"`
}

// AppConfig 服务配置
type AppConfig struct {
	Port             string    `mapstructure:"port"`
	Datasource       string    `mapstructure:"datasource"`
	LogConfig        LogConfig `mapstructure:"log"`
	JWTConfig        JWTConfig `mapstructure:"jwt"`
	NodeInfoInterval int64     `mapstructure:"node_info_interval"`

	Mysql      DbStruct    `mapstructure:"mysql"`
	Postgres   DbStruct    `mapstructure:"postgres"`
	Sqlite     DbStruct    `mapstructure:"sqlite"`
	Rabbitmq   Rabbitmq    `mapstructure:"rabbitmq"`
	Redis      RedisConfig `mapstructure:"redis"`

}

type Rabbitmq struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Topic    string `mapstructure:"topic"`
}

type DbStruct struct {
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	DbName       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Auth     string `mapstructure:"auth"`
	Database int    `mapstructure:"db"`
}

// Init 配置初始化
func Init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.AllowEmptyEnv(true)

	if err := viper.ReadInConfig(); err != nil {
		logging.Fatal().Err(err).Msg("config: viper reading config file failed")
	}

	err := viper.Unmarshal(&Config)
	if err != nil {
		logging.Fatal().Err(err).Msg("config: viper decode config failed")
	}

	queryConditionMetaData, err := metadata.Asset("query_condition.json")
	if err != nil {
		logging.Fatal().Err(err).Msg("load query-condition metadata error")
	}
	if err := json.Unmarshal(queryConditionMetaData, &QueryConditionMetaData); err != nil {
		logging.Fatal().Err(err).Msg("Fatal error unmarshal queryConditionMetaData")
	}
}