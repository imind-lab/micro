package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	tpl "github.com/imind-lab/micro/microctl/template"
	"github.com/imind-lab/micro/microctl/template/api"
	"github.com/imind-lab/micro/microctl/template/srv"
	"github.com/imind-lab/micro/util"
)

var (
	domain  string
	project string
	service string
	layer   string
	kafka   bool
	tracing bool
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
			Svc:     util.GetPascalCase(service),
			Date:    date,
			Year:    year,
			Package: strings.ReplaceAll(service, "-", "_"),
			Name:    util.GetCamelCase(service),
			MQ:      kafka,
		}

		if layer == "api" {
			data.Suffix = "-api"
			data.Package = data.Package + "_api"

			err := api.CreateApplicationProto(data)
			if err == nil {
				fmt.Println("生成Api-ApplicationProto成功")
			} else {
				fmt.Println("生成Api-ApplicationProto出错", err)
			}

			err = api.CreateApplicationService(data)
			if err == nil {
				fmt.Println("生成Api-ApplicationService成功")
			} else {
				fmt.Println("生成Api-ApplicationService出错", err)
			}

			err = api.CreateCmd(data)
			if err == nil {
				fmt.Println("生成Api-Cmd成功")
			} else {
				fmt.Println("生成Api-Cmd出错", err)
			}

			err = api.CreateCmdServer(data)
			if err == nil {
				fmt.Println("生成Api-CmdServer成功")
			} else {
				fmt.Println("生成Api-CmdServer出错", err)
			}

			err = api.CreateConf(data)
			if err == nil {
				fmt.Println("生成Api-Conf成功")
			} else {
				fmt.Println("生成Api-Conf出错", err)
			}

			err = api.CreateConfCrt(data)
			if err == nil {
				fmt.Println("生成Api-ConfCrt成功")
			} else {
				fmt.Println("生成Api-ConfCrt出错", err)
			}

			err = api.CreateConfKey(data)
			if err == nil {
				fmt.Println("生成Api-ConfKey成功")
			} else {
				fmt.Println("生成Api-ConfKey出错", err)
			}

			err = api.CreateDeploy(data)
			if err == nil {
				fmt.Println("生成Api-Deploy成功")
			} else {
				fmt.Println("生成Api-Deploy出错", err)
			}

			err = api.CreateDockerfile(data)
			if err == nil {
				fmt.Println("生成Api-BuildDockerfile成功")
			} else {
				fmt.Println("生成Api-BuildDockerfile出错", err)
			}

			err = api.CreateDomainConvert(data)
			if err == nil {
				fmt.Println("生成Api-DomainConvert成功")
			} else {
				fmt.Println("生成Api-DomainConvert出错", err)
			}

			err = api.CreateDomainDomain(data)
			if err == nil {
				fmt.Println("生成Api-DomainDomain成功")
			} else {
				fmt.Println("生成Api-DomainDomain出错", err)
			}

			err = api.CreateDomainService(data)
			if err == nil {
				fmt.Println("生成Api-DomainService成功")
			} else {
				fmt.Println("生成Api-DomainService出错", err)
			}

			err = api.CreateMain(data)
			if err == nil {
				fmt.Println("生成Api-Main成功")
			} else {
				fmt.Println("生成Api-Main出错", err)
			}

			err = api.CreateMakefile(data)
			if err == nil {
				fmt.Println("生成Api-BuildMakefile成功")
			} else {
				fmt.Println("生成Api-BuildMakefile出错", err)
			}

			err = api.CreateMod(data)
			if err == nil {
				fmt.Println("生成Api-Mod成功")
			} else {
				fmt.Println("生成Api-Mod出错", err)
			}

			err = api.CreatePkgConstantCache(data)
			if err == nil {
				fmt.Println("生成Api-PkgConstantCache成功")
			} else {
				fmt.Println("生成Api-PkgConstantCache出错", err)
			}

			err = api.CreatePkgConstantOption(data)
			if err == nil {
				fmt.Println("生成Api-PkgConstantOption成功")
			} else {
				fmt.Println("生成Api-PkgConstantOption出错", err)
			}

			err = api.CreatePkgGoogleProtos(data)
			if err == nil {
				fmt.Println("生成Api-PkgGoogleProtos成功")
			} else {
				fmt.Println("生成Api-PkgGoogleProtos出错", err)
			}

			err = api.CreatePkgUtilCache(data)
			if err == nil {
				fmt.Println("生成Api-PkgUtilCache成功")
			} else {
				fmt.Println("生成Api-PkgUtilCache出错", err)
			}

			err = api.CreateRepositoryInterface(data)
			if err == nil {
				fmt.Println("生成Api-RepositoryInterface成功")
			} else {
				fmt.Println("生成Api-RepositoryInterface出错", err)
			}

			err = api.CreateRepositoryPersistence(data)
			if err == nil {
				fmt.Println("生成Api-RepositoryPersistence成功")
			} else {
				fmt.Println("生成Api-RepositoryPersistence出错", err)
			}

			err = api.CreateRepositoryPersistenceService(data)
			if err == nil {
				fmt.Println("生成Api-RepositoryPersistenceService成功")
			} else {
				fmt.Println("生成Api-RepositoryPersistenceService出错", err)
			}

			err = api.CreateServer(data)
			if err == nil {
				fmt.Println("生成Api-Server成功")
			} else {
				fmt.Println("生成Api-Server出错", err)
			}

			err = api.CreateServerWire(data)
			if err == nil {
				fmt.Println("生成Api-ServerWire成功")
			} else {
				fmt.Println("生成Api-ServerWire出错", err)
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

			if kafka {
				err = srv.CreateApplicationSubscriber(data)
				if err == nil {
					fmt.Println("生成ApplicationSubscriber成功")
				} else {
					fmt.Println("生成ApplicationSubscriber出错", err)
				}
			}

			err = srv.CreateClient(data)
			if err == nil {
				fmt.Println("生成Client成功")
			} else {
				fmt.Println("生成Client出错", err)
			}

			err = srv.CreateCmd(data)
			if err == nil {
				fmt.Println("生成Cmd成功")
			} else {
				fmt.Println("生成Cmd出错", err)
			}

			err = srv.CreateCmdCron(data)
			if err == nil {
				fmt.Println("生成CmdCron成功")
			} else {
				fmt.Println("生成CmdCron出错", err)
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

			err = srv.CreateConfCrt(data)
			if err == nil {
				fmt.Println("生成ConfCrt成功")
			} else {
				fmt.Println("生成ConfCrt出错", err)
			}

			err = srv.CreateConfKey(data)
			if err == nil {
				fmt.Println("生成ConfKey成功")
			} else {
				fmt.Println("生成ConfKey出错", err)
			}

			err = srv.CreateDeploy(data)
			if err == nil {
				fmt.Println("生成Deploy成功")
			} else {
				fmt.Println("生成Deploy出错", err)
			}

			err = srv.CreateDockerfile(data)
			if err == nil {
				fmt.Println("生成BuildDockerfile成功")
			} else {
				fmt.Println("生成BuildDockerfile出错", err)
			}

			err = srv.CreateDomainConvert(data)
			if err == nil {
				fmt.Println("生成DomainConvert成功")
			} else {
				fmt.Println("生成DomainConvert出错", err)
			}

			err = srv.CreateDomainDomain(data)
			if err == nil {
				fmt.Println("生成DomainDomain成功")
			} else {
				fmt.Println("生成DomainDomain出错", err)
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

			err = srv.CreateMakefile(data)
			if err == nil {
				fmt.Println("生成BuildMakefile成功")
			} else {
				fmt.Println("生成BuildMakefile出错", err)
			}

			err = srv.CreateMod(data)
			if err == nil {
				fmt.Println("生成Mod成功")
			} else {
				fmt.Println("生成Mod出错", err)
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

			err = srv.CreatePkgGoogleProtos(data)
			if err == nil {
				fmt.Println("生成Pkg成功")
			} else {
				fmt.Println("生成Pkg出错", err)
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

			err = srv.CreateRepositoryPersistence(data)
			if err == nil {
				fmt.Println("生成RepositoryPersistence成功")
			} else {
				fmt.Println("生成RepositoryPersistence出错", err)
			}

			err = srv.CreateRepositoryPersistenceService(data)
			if err == nil {
				fmt.Println("生成RepositoryPersistenceService成功")
			} else {
				fmt.Println("生成RepositoryPersistenceService出错", err)
			}

			err = srv.CreateServer(data)
			if err == nil {
				fmt.Println("生成Server成功")
			} else {
				fmt.Println("生成Server出错", err)
			}

			err = srv.CreateServerWire(data)
			if err == nil {
				fmt.Println("生成ServerWire成功")
			} else {
				fmt.Println("生成ServerWire出错", err)
			}
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "git.chope.it", "company domain")
	rootCmd.PersistentFlags().StringVarP(&project, "project", "p", "backend", "project name")
	rootCmd.PersistentFlags().StringVarP(&service, "service", "s", "greeter", "service name")
	rootCmd.PersistentFlags().StringVarP(&layer, "layer", "l", "srv", "service type: src for backend service, api for api gateway")
	rootCmd.PersistentFlags().BoolVarP(&kafka, "kafka", "", true, "whether to generate kafka-related code")
	rootCmd.PersistentFlags().BoolVarP(&tracing, "tracing", "", true, "whether to generate tracing-related code")
	rootCmd.AddCommand(serverCmd)
}
