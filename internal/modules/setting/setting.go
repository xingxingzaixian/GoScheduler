package setting

import (
	"GoScheduler/internal/modules/global"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig(configFile string) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath(global.AppDir + "/conf")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AutomaticEnv()
	}

	if err := viper.ReadInConfig(); err != nil {
		zap.S().Panicf("读取配置文件[%s]失败", configFile)
	}

	zap.S().Infof("配置文件【%s】读取成功", configFile)
}

func WriteConfig() {

}
