package cmd

import (
    "fmt"
    "github.com/spf13/viper"
    "strings"

    "github.com/imind-lab/micro/v2/microctl/orm"
    "github.com/spf13/cobra"
)

var (
    cfgFile string
    outPath string
    db      string
    table   string
)

var ormCmd = &cobra.Command{
    Use:   "orm",
    Short: "create orm models",
    Run: func(cmd *cobra.Command, args []string) {
        //initConf()

        tables := strings.Split(table, ",")
        if len(tables) == 0 {
            fmt.Println("请输入表名")
            return
        }
        if !strings.HasSuffix(outPath, "/") {
            outPath += "/"
        }
        err := orm.Run(db, tables, outPath)
        if err != nil {
            fmt.Println("Model生成发生错误", err)
        }
    },
}

func init() {
    rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./conf/conf.yaml", "Start server with provided configuration file")
    rootCmd.PersistentFlags().StringVarP(&outPath, "out", "o", "./repository/model/", "model output dir")
    rootCmd.PersistentFlags().StringVarP(&db, "db", "b", "imind", "db name in conf file")
    rootCmd.PersistentFlags().StringVarP(&table, "table", "t", "", "tables name")
    rootCmd.AddCommand(ormCmd)
    //cobra.OnInitialize(initConf)
}

func initConf() {
    viper.SetConfigFile(cfgFile)
    //初始化全部的配置
    err := viper.ReadInConfig()
    if err != nil {
        panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }
}
