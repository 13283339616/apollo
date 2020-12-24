package apollo

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func dealError(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadProperties(path string) map[string]string {
	fmt.Print(path)
	file, err := os.Open(path)
	dealError(err)
	defer file.Close()
	reader := bufio.NewReader(file)
	configMap := make(map[string]string, 128)
	for {
		readString, err := reader.ReadString('\n')
		if err == io.EOF {
			return configMap
		}
		dealError(err)
		readString = strings.TrimSpace(readString)
		if strings.HasPrefix(readString, "#") || readString == "" {
			continue
		} else {
			if !strings.Contains(readString, "=") {
				dealError(errors.New("配置文件有误"))
			}
			if strings.HasPrefix(readString, "=") {
				dealError(errors.New("配置文件有误"))
			}
			index := strings.Index(readString, "=")
			configMap[readString[0:index]] = readString[index+1 : len(readString)]
		}

	}

}
