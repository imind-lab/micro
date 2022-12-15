/**
 *  MindLab
 *
 *  Create by songli on 2022/02/27
 *  Copyright © 2022 imind.tech All rights reserved.
 */

package srv

import (
	"github.com/imind-lab/micro/microctl/template"
)

// 生成client/service.go
func CreateCmdServer(data *template.Data) error {
	var tpl = `package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"{{.Domain}}/{{.Project}}/{{.Service}}/server"
)

var cfgFile string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gRPC Merchant server",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recover error : %v\n", err)
			}
		}()
		err := server.Serve()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./conf/conf.yaml", "Start server with provided configuration file")
	rootCmd.AddCommand(serverCmd)
	cobra.OnInitialize(initConf)
}

func initConf() {
	viper.SetConfigFile(cfgFile)
	//初始化全部的配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}
`

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/cmd/"
	name := "server.go"

	return template.CreateFile(data, tpl, path, name)
}
