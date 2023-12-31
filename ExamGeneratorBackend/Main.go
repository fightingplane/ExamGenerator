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

func main() {

	initViper()
	initGin()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	logging.GetLogger().Info().Msg("Router initialized")

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run("127.0.0.1" + ":" + viper.GetString("Server.Port"))
}
