package main

import (
	logging "examgen/Logger"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func initViper() {
	viper.SetConfigName("ExamGeneratorConfigurations")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	fmt.Println("init viper successfully")
}

func initGin() {

	if mode := viper.GetBool("gin.isRelease"); mode {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("using Gin as release mode")
	} else {
		gin.SetMode(gin.DebugMode)
		fmt.Println("using gin as debug mode")
	}
}

var (
	ExamGenLogger *logging.ExamGenLogger
)

func initLogger() {

	config := &logging.LoggerConfig{}
	viper.UnmarshalKey("LoggerConfigurations", config)

	ExamGenLogger = logging.ConfigLogger(*config)
	ExamGenLogger.Info().Msg("Logger initialized")
}

func main() {

	initViper()
	initLogger()
	initGin()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	ExamGenLogger.Info().Msg("Router initialized")
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
