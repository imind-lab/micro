/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package api

import (
    "github.com/imind-lab/micro/v2/microctl/template"
)

// 生成client/service.go
func CreateCmdServer(data *template.Data) error {
    var tpl = `package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"gitlab.imind.tech/{{.Repo}}/{{.Service}}-api/server"
)

var cfgFile string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the gRPC {{.Svc}}-api server",
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

    path := "./" + data.Domain + "/" + data.Repo + "/" + data.Service + "-api/cmd/"
    name := "server.go"

    return template.CreateFile(data, tpl, path, name)
}
