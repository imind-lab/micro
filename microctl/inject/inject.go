package inject

import (
    "bufio"
    "bytes"
    "io"
    "os"
    "strings"
)

type Inject struct {
    path string

    header  bytes.Buffer
    body    bytes.Buffer
    content bytes.Buffer

    tag       string
    response  string
    exist     bool
    headerEnd bool
}

func NewInject(path string) *Inject {
    return &Inject{path: path}
}

func (i *Inject) Process() error {
    file, err := os.Open(i.path)
    if err != nil {
        return err
    }

    reader := bufio.NewReader(file)
    for {
        line, _, err := reader.ReadLine()
        if err != nil {
            file.Close()
            if err == io.EOF {
                break
            }
            panic(err)
        }
        i.process(line)
    }

    i.content.Write(i.header.Bytes())
    if i.exist {
        i.content.WriteString("\t\"strings\"\n")
        i.content.WriteString("\t\"github.com/imind-lab/micro/v2/status\"\n")
    }
    i.content.Write(i.body.Bytes())

    return nil
}

func (i *Inject) Save() error {
    file, err := os.OpenFile(i.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    _, err = writer.Write(i.content.Bytes())
    if err != nil {
        return err
    }
    writer.Flush()

    return nil
}

func (i *Inject) processTag(line []byte) {
    temp := string(line)

    items := strings.Split(i.tag, " ")
    if len(items) >= 3 {
        new := "omitempty\" " + items[2]
        tmp := strings.Replace(temp, "omitempty\"", new, 1)
        i.body.WriteString(tmp)
    } else {
        i.body.Write(line)
    }
    i.body.WriteString("\n")
    i.tag = ""
}

func (i *Inject) processResponse(line []byte) {
    i.body.Write(line)
    i.body.WriteString("\n")

    temp := string(line)

    if temp == "}" {
        i.body.WriteString("\n")
        i.body.WriteString(i.response)

        items := strings.Split(i.response, " ")
        if len(items) >= 3 {
            i.body.WriteString("\nfunc (x *" + items[2] + ") SetCode(code status.Code, msg ...string) {\n\tx.Code = int32(code)\n\tif len(msg) == 0 {\n\t\tx.Msg = code.String()\n\t} else {\n\t\tx.Msg = strings.Join(msg, \"，\")\n\t}\n}\n")
        }
        if len(items) >= 5 {
            i.body.WriteString("\nfunc (x *" + items[2] + ") SetBody(code status.Code, " + items[4] + " " + items[3] + ") {\n\tx.Code = int32(code)\n\tx.Msg = code.String()\n\tx." + strings.Title(items[4]) + " = " + items[4] + "\n}\n")
        }
        i.response = ""
    }
}

func (i *Inject) process(line []byte) {
    temp := string(line)
    if strings.Contains(temp, "import (") {
        i.headerEnd = true
        i.header.Write(line)
        i.header.WriteString("\n")
    } else if strings.Contains(temp, "// @inject_tag") {
        i.tag = temp

        i.body.Write(line)
        i.body.WriteString("\n")
    } else if strings.Contains(temp, "// @inject_response") {
        i.exist = true
        i.response = temp

        i.body.Write(line)
        i.body.WriteString("\n")
    } else {
        if len(i.tag) > 0 {
            i.processTag(line)
        } else if len(i.response) > 0 {
            i.processResponse(line)
        } else {
            if i.headerEnd {
                i.body.Write(line)
                i.body.WriteString("\n")
            } else {
                i.header.Write(line)
                i.header.WriteString("\n")
            }
        }
    }
}
