/**
 *  MindLab
 *
 *  Create by songli on 2020/10/23
 *  Copyright © 2021 imind.tech All rights reserved.
 */

package template

import (
    "os"
    "strings"
    "text/template"
)

type Data struct {
    Domain  string
    Repo    string
    Name    string
    Package string
    Service string
    Svc     string
    Suffix  string
    Date    string
    Year    string
    Tracing bool
    MQ      bool
}

func CreateFile(data *Data, tpl, path, name string) error {
    tpl = strings.ReplaceAll(tpl, "${backtick}", "`")
    t, err := template.New(name).Parse(tpl)
    if err != nil {
        return err
    }

    t.Option()

    err = os.MkdirAll(path, os.ModePerm)
    if err != nil {
        return err
    }

    fileName := path + name

    f, err := os.Create(fileName)
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

func WriteFile(tpl, path, name string) error {
    err := os.MkdirAll(path, os.ModePerm)
    if err != nil {
        return err
    }

    fileName := path + name

    err = os.WriteFile(fileName, []byte(tpl), os.ModePerm)
    if err != nil {
        return err
    }

    return nil
}
