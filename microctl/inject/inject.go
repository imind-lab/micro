package inject

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

func Process(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var header bytes.Buffer
	var body bytes.Buffer
	var content bytes.Buffer

	exist := false
	headerEnd := false
	tag := ""
	reponse := ""

	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		temp := string(line)
		if strings.Contains(temp, "import (") {
			headerEnd = true
			header.Write(line)
			header.WriteString("\n")
		} else if strings.Contains(temp, "// @inject_tag") {
			tag = temp

			body.Write(line)
			body.WriteString("\n")
		} else if strings.Contains(temp, "// @inject_response") {
			exist = true
			reponse = temp

			body.Write(line)
			body.WriteString("\n")
		} else {
			if len(tag) > 0 {
				items := strings.Split(tag, " ")
				if len(items) >= 3 {
					new := "omitempty\" " + items[2]
					tmp := strings.Replace(temp, "omitempty\"", new, 1)
					body.WriteString(tmp)
				} else {
					body.Write(line)
				}
				body.WriteString("\n")
				tag = ""
			} else if len(reponse) > 0 {
				body.Write(line)
				body.WriteString("\n")

				if temp == "}" {
					body.WriteString("\n")
					body.WriteString(reponse)

					items := strings.Split(reponse, " ")
					if len(items) >= 3 {
						body.WriteString("\nfunc (x *" + items[2] + ") SetCode(code status.Code, message ...string) {\n\tx.Code = int32(code)\n\tif len(message) == 0 {\n\t\tx.Message = code.String()\n\t} else {\n\t\tx.Message = message\n\t}\n}\n")
					}
					if len(items) >= 5 {
						body.WriteString("\nfunc (x *" + items[2] + ") SetBody(code status.Code, " + items[4] + " " + items[3] + ") {\n\tx.Code = int32(code)\n\tx.Message = code.String()\n\tx." + strings.Title(items[4]) + " = " + items[4] + "\n}\n")
					}
					reponse = ""
				}
			} else {
				if headerEnd {
					body.Write(line)
					body.WriteString("\n")
				} else {
					header.Write(line)
					header.WriteString("\n")
				}
			}
		}
	}

	content.Write(header.Bytes())
	if exist {
		content.WriteString("\t\"strings\"\n")
		content.WriteString("\t\"github.com/imind-lab/micro/status\"\n")
	}
	content.Write(body.Bytes())

	return content.Bytes(), nil
}

func Save(path string, content []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(content)
	if err != nil {
		return err
	}
	writer.Flush()

	return nil
}
