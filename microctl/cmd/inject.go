package cmd

import (
    "fmt"

    "github.com/spf13/cobra"

    "github.com/imind-lab/micro/v2/microctl/inject"
)

var (
    path string
)

var injectCmd = &cobra.Command{
    Use:   "inject",
    Short: "Use microctl inject process microservice",
    Run: func(cmd *cobra.Command, args []string) {
        inj := inject.NewInject(path)
        err := inj.Process()
        if err != nil {
            fmt.Println("[ERROR]处理出错", err)
            return
        }
        err = inj.Save()
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
