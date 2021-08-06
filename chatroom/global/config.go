package global

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	SensitiveWords []string
)

func initConfig() {
	viper.SetConfigName("chatroom")
	viper.AddConfigPath(RootDir + "/config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	SensitiveWords = viper.GetStringSlice("sensstive")

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		viper.ReadInConfig()
		SensitiveWords = viper.GetStringSlice("sensstive")
	})
}

