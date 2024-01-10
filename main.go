package main

import (
	"fmt"
	"github.com/danapsimer/dnsmasq-gui/dnsmasq"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
)

func init() {
	viper.SetDefault("LeaseFile", "/var/lib/misc/dnsmasq.leases")
	viper.SetConfigName("dnsmasq-gui")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("DNSMASQ_GUI")
	viper.AutomaticEnv()
}

func main() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}

	r := gin.Default()
	r.GET("/lease", func(context *gin.Context) {
		f, err := os.Open(viper.GetString("LeaseFile"))
		if err != nil {
			context.JSON(500, gin.H{
				"error": err.Error(),
			})
		} else {
			defer f.Close()
			leases, err := dnsmasq.ReadLeases(f)
			if err != nil {
				context.JSON(500, gin.H{
					"error": err.Error(),
				})
			} else {
				context.JSON(200, leases)
			}
		}
	})
	r.Run()
}
