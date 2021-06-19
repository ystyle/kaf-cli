package kafcli

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"
)

func parseLang(lang string) string {
	var langs = "en,de,fr,it,es,zh,ja,pt,ru,nl"
	if strings.Contains(langs, lang) {
		return lang
	}
	return "en"
}

func run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func lookKindlegen() string {
	command := "kindlegen"
	if runtime.GOOS == "windows" {
		command = "kindlegen.exe"
	}
	kindlegen, _ := exec.LookPath(command)
	if kindlegen == "" {
		currentDir := path.Dir(os.Args[0])
		kindlegen = path.Join(currentDir, command)
		if exist, _ := isExists(kindlegen); !exist {
			return kindlegen
		}
	}
	return kindlegen
}

func converToMobi(bookname string) {
	command := lookKindlegen()
	fmt.Printf("\n检测到Kindle格式转换器: %s，正在把书籍转换成Kindle格式...\n", command)
	fmt.Println("转换mobi比较花时间, 大约耗时1-10分钟, 请等待...")
	start := time.Now()
	run(command, "-dont_append_source", "-locale", lang, "-c1", bookname)
	// 计算耗时
	end := time.Now().Sub(start)
	fmt.Println("转换为Kindle格式耗时:", end)
}

func isExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
