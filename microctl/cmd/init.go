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
			err := srv.CreateApplicationProto(data)
			if err == nil {
				fmt.Println("生成ApplicationProto成功")
			} else {
				fmt.Println("生成ApplicationProto出错", err)
			}

			err = srv.CreateApplicationService(data)
			if err == nil {
				fmt.Println("生成ApplicationService成功")
			} else {
				fmt.Println("生成ApplicationService出错", err)
			}

			err = srv.CreateApplicationSubscriber(data)
			if err == nil {
				fmt.Println("生成ApplicationSubscriber成功")
			} else {
				fmt.Println("生成ApplicationSubscriber出错", err)
			}

			err = srv.CreateBuildDockerfile(data)
			if err == nil {
				fmt.Println("生成BuildDockerfile成功")
			} else {
				fmt.Println("生成BuildDockerfile出错", err)
			}

			err = srv.CreateBuildMakefile(data)
			if err == nil {
				fmt.Println("生成BuildMakefile成功")
			} else {
				fmt.Println("生成BuildMakefile出错", err)
			}

			err = srv.CreateClient(data)
			if err == nil {
				fmt.Println("生成Client成功")
			} else {
				fmt.Println("生成Client出错", err)
			}

			err = srv.CreateClientService(data)
			if err == nil {
				fmt.Println("生成ClientService成功")
			} else {
				fmt.Println("生成ClientService出错", err)
			}

			err = srv.CreateCmd(data)
			if err == nil {
				fmt.Println("生成Cmd成功")
			} else {
				fmt.Println("生成Cmd出错", err)
			}

			err = srv.CreateCmdServer(data)
			if err == nil {
				fmt.Println("生成CmdServer成功")
			} else {
				fmt.Println("生成CmdServer出错", err)
			}

			err = srv.CreateConf(data)
			if err == nil {
				fmt.Println("生成Conf成功")
			} else {
				fmt.Println("生成Conf出错", err)
			}

			err = srv.CreateDeploy(data)
			if err == nil {
				fmt.Println("生成Deploy成功")
			} else {
				fmt.Println("生成Deploy出错", err)
			}

			err = srv.CreateDomainConvert(data)
			if err == nil {
				fmt.Println("生成DomainConvert成功")
			} else {
				fmt.Println("生成DomainConvert出错", err)
			}

			err = srv.CreateDomainService(data)
			if err == nil {
				fmt.Println("生成DomainService成功")
			} else {
				fmt.Println("生成DomainService出错", err)
			}

			err = srv.CreateMain(data)
			if err == nil {
				fmt.Println("生成Main成功")
			} else {
				fmt.Println("生成Main出错", err)
			}

			err = srv.CreateMod(data)
			if err == nil {
				fmt.Println("生成Mod成功")
			} else {
				fmt.Println("生成Mod出错", err)
			}

			err = srv.CreatePkg(data)
			if err == nil {
				fmt.Println("生成Pkg成功")
			} else {
				fmt.Println("生成Pkg出错", err)
			}

			err = srv.CreatePkgConstantCache(data)
			if err == nil {
				fmt.Println("生成PkgConstantCache成功")
			} else {
				fmt.Println("生成PkgConstantCache出错", err)
			}

			err = srv.CreatePkgConstantOption(data)
			if err == nil {
				fmt.Println("生成PkgConstantOption成功")
			} else {
				fmt.Println("生成PkgConstantOption出错", err)
			}

			err = srv.CreatePkgUtilCache(data)
			if err == nil {
				fmt.Println("生成PkgUtilCache成功")
			} else {
				fmt.Println("生成PkgUtilCache出错", err)
			}

			err = srv.CreateRepositoryInterface(data)
			if err == nil {
				fmt.Println("生成RepositoryInterface成功")
			} else {
				fmt.Println("生成RepositoryInterface出错", err)
			}

			err = srv.CreateRepositoryModel(data)
			if err == nil {
				fmt.Println("生成RepositoryModel成功")
			} else {
				fmt.Println("生成RepositoryModel出错", err)
			}

			err = srv.CreateRepositoryOptions(data)
			if err == nil {
				fmt.Println("生成RepositoryOptions成功")
			} else {
				fmt.Println("生成RepositoryOptions出错", err)
			}

			err = srv.CreateRepositoryPersistence(data)
			if err == nil {
				fmt.Println("生成RepositoryPersistence成功")
			} else {
				fmt.Println("生成RepositoryPersistence出错", err)
			}

			err = srv.CreateServer(data)
			if err == nil {
				fmt.Println("生成Server成功")
			} else {
				fmt.Println("生成Server出错", err)
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
