/**
 *  MindLab
 *
 *  Create by songli on 2023/02/03
 *  Copyright © 2023 imind.tech All rights reserved.
 */

package srv

import (
	"github.com/imind-lab/micro/v2/microctl/template"
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
	Short: "Run the gRPC {{.Service}} Server",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
`

	path := "./" + data.Name + "/cmd/"
	name := "cmd.go"

	return template.CreateFile(data, tpl, path, name)
}
