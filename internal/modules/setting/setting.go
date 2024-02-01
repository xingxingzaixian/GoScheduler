package setting

import (
	"GoScheduler/internal/modules/global"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type DB struct {
	Engine   string
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Setting struct {
	Name   string `mapstructure:"name"`
	Queue  int    `mapstructure:"queue"`
	DBInfo DB     `mapstructure:"db"`
}

func InitConfig(configFile string) {
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicf("读取配置文件[%s]失败", configFile)
	}

	if err := v.Unmarshal(global.Setting); err != nil {
		zap.S().Panicf("配置文件【%s】格式异常", configFile)
	}

	zap.S().Infof("配置文件【%s】读取成功", configFile)
	zap.S().Infof("配置信息：%v", global.Setting)
}

func WriteConfig() {

}
