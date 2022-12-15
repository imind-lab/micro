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
func CreateCmd(data *template.Data) error {
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

	path := "./" + data.Domain + "/" + data.Project + "/" + data.Service + "/cmd/"
	name := "cmd.go"

	return template.CreateFile(data, tpl, path, name)
}
