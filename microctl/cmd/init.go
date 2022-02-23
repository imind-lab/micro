package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	tpl "github.com/imind-lab/micro/microctl/template"
	"github.com/imind-lab/micro/microctl/template/api"
	"github.com/imind-lab/micro/microctl/template/srv"
)

var (
	domain  string
	project string
	service string
	layer   string
)

var serverCmd = &cobra.Command{
	Use:   "init",
	Short: "Use microctl create new microservice",
	Run: func(cmd *cobra.Command, args []string) {
		date := time.Now().Format("2006/01/02")
		year := time.Now().Format("2006")

		data := &tpl.Data{
			Domain:  domain,
			Project: project,
			Service: service,
			Svc:     strings.Title(service),
			Date:    date,
			Year:    year,
		}

		if layer == "api" {
			err := api.CreateBuild(data)
			if err == nil {
				fmt.Println("生成API-Build成功")
			} else {
				fmt.Println("生成API-Build出错", err)
			}

			err = api.CreateConf(data)
			if err == nil {
				fmt.Println("生成API-Conf成功")
			} else {
				fmt.Println("生成API-Conf出错", err)
			}

			err = api.CreateProto(data)
			if err == nil {
				fmt.Println("生成API-Proto成功")
			} else {
				fmt.Println("生成API-Proto出错", err)
			}

			err = api.CreateService(data)
			if err == nil {
				fmt.Println("生成API-Service成功")
			} else {
				fmt.Println("生成API-Service出错", err)
			}

			err = api.CreateMain(data)
			if err == nil {
				fmt.Println("生成API-Main成功")
			} else {
				fmt.Println("生成API-Main出错", err)
			}

			err = api.CreatePkg(data)
			if err == nil {
				fmt.Println("生成API-Pkg成功")
			} else {
				fmt.Println("生成API-Pkg出错", err)
			}

			err = api.CreateCmd(data)
			if err == nil {
				fmt.Println("生成API-Cmd成功")
			} else {
				fmt.Println("生成API-Cmd出错", err)
			}

			err = api.CreateServer(data)
			if err == nil {
				fmt.Println("生成API-Server成功")
			} else {
				fmt.Println("生成API-Server出错", err)
			}
		} else {
			err := srv.CreateModel(data)
			if err == nil {
				fmt.Println("生成Model成功")
			} else {
				fmt.Println("生成Model出错", err)
			}

			err = srv.CreateRepository(data)
			if err == nil {
				fmt.Println("生成Repository成功")
			} else {
				fmt.Println("生成Repository出错", err)
			}

			err = srv.CreateProto(data)
			if err == nil {
				fmt.Println("生成Proto成功")
			} else {
				fmt.Println("生成Proto出错", err)
			}

			err = srv.CreateBuild(data)
			if err == nil {
				fmt.Println("生成Build成功")
			} else {
				fmt.Println("生成Build出错", err)
			}

			err = srv.CreateConf(data)
			if err == nil {
				fmt.Println("生成Conf成功")
			} else {
				fmt.Println("生成Conf出错", err)
			}

			err = srv.CreateDomain(data)
			if err == nil {
				fmt.Println("生成Domain成功")
			} else {
				fmt.Println("生成Domain出错", err)
			}

			err = srv.CreateService(data)
			if err == nil {
				fmt.Println("生成Service成功")
			} else {
				fmt.Println("生成Service出错", err)
			}

			err = srv.CreateSubscriber(data)
			if err == nil {
				fmt.Println("生成Subscriber成功")
			} else {
				fmt.Println("生成Subscriber出错", err)
			}

			err = srv.CreateClient(data)
			if err == nil {
				fmt.Println("生成Client成功")
			} else {
				fmt.Println("生成Client出错", err)
			}

			err = srv.CreateMain(data)
			if err == nil {
				fmt.Println("生成Main成功")
			} else {
				fmt.Println("生成Main出错", err)
			}

			err = srv.CreatePkg(data)
			if err == nil {
				fmt.Println("生成Pkg成功")
			} else {
				fmt.Println("生成Pkg出错", err)
			}

			err = srv.CreateCmd(data)
			if err == nil {
				fmt.Println("生成Cmd成功")
			} else {
				fmt.Println("生成Cmd出错", err)
			}

			err = srv.CreateServer(data)
			if err == nil {
				fmt.Println("生成Server成功")
			} else {
				fmt.Println("生成Server出错", err)
			}

			err = srv.CreateDeploy(data)
			if err == nil {
				fmt.Println("生成Deploy成功")
			} else {
				fmt.Println("生成Deploy出错", err)
			}
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "github.com", "company domain")
	rootCmd.PersistentFlags().StringVarP(&project, "project", "p", "imind-lab", "project name")
	rootCmd.PersistentFlags().StringVarP(&service, "service", "s", "greeter", "service name")
	rootCmd.PersistentFlags().StringVarP(&layer, "layer", "l", "srv", "generate service layere")
	rootCmd.AddCommand(serverCmd)
}
