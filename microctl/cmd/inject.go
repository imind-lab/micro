package cmd

import (
	"fmt"
	"github.com/imind-lab/micro/microctl/inject"
	"github.com/spf13/cobra"
)

var (
	path string
)

var injectCmd = &cobra.Command{
	Use:   "inject",
	Short: "Use chopectl inject process microservice",
	Run: func(cmd *cobra.Command, args []string) {
		content, err := inject.Process(path)
		if err != nil {
			fmt.Println("[ERROR]处理出错", err)
			return
		}
		err = inject.Save(path, content)
		if err != nil {
			fmt.Println("[ERROR]保存出错", err)
			return
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&path, "path", "f", "", "inject file path")
	rootCmd.AddCommand(injectCmd)
}
