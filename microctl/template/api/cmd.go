package api

import (
	"os"
	"text/template"

	tp "github.com/imind-lab/micro/microctl/template"
)

// 生成cmd
func CreateCmd(data *tp.Data) error {
	// 生成cmd.go
	var tpl = `package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run the gRPC {{.Svc}} Server",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
`

	t, err := template.New("cmd.go").Parse(tpl)
	if err != nil {
		return err
	}

	dir := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/cmd/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName := dir + "cmd.go"

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	// 生成server.go
	tpl = `package cmd

import (
        "fmt"
        "log"

        "github.com/spf13/cobra"
        "github.com/spf13/viper"

        "{{.Domain}}/{{.Project}}/{{.Service}}-api/server"
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

	t, err = template.New("server.go").Parse(tpl)
	if err != nil {
		return err
	}

	dir = "./" + data.Domain + "/" + data.Project + "/" + data.Service + "-api/cmd/"

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	fileName = dir + "server.go"

	f, err = os.Create(fileName)
	if err != nil {
		return err
	}
	err = t.Execute(f, data)
	if err != nil {
		return err
	}
	f.Close()

	return nil
}
