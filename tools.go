package kafcli

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func parseLang(lang string) string {
	if lang == "" {
		return "en"
	}
	var langs = "en,de,fr,it,es,zh,ja,pt,ru,nl"
	if strings.Contains(langs, lang) {
		return lang
	}
	return "en"
}

func run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func lookKindlegen() string {
	command := "kindlegen"
	if runtime.GOOS == "windows" {
		command = "kindlegen.exe"
	}
	kindlegen, err := exec.LookPath(command)
	if err != nil {
		currentDir, err := os.Executable()
		if err != nil {
			return ""
		}
		kindlegen = filepath.Join(filepath.Dir(currentDir), command)
		if exist, _ := isExists(kindlegen); !exist {
			return ""
		}
		fmt.Println("kindlegen: ", kindlegen)
	}
	return kindlegen
}

func converToMobi(bookname, lang string) {
	command := lookKindlegen()
	fmt.Printf("\n检测到Kindle格式转换器: %s，正在把书籍转换成Kindle格式...\n", command)
	fmt.Println("转换mobi比较花时间, 大约耗时1-10分钟, 请等待...")
	start := time.Now()
	run(command, "-dont_append_source", "-locale", lang, "-c1", bookname)
	// 计算耗时
	end := time.Now().Sub(start)
	fmt.Println("转换为mobi格式耗时:", end)
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

func getClientID() string {
	clientID := fmt.Sprintf("%d", rand.Uint32())
	config, err := os.UserConfigDir()
	if err != nil {
		return clientID
	}
	filepath := fmt.Sprintf("%s/kaf-wifi/config", config)
	if exist, _ := isExists(filepath); exist {
		bs, err := ioutil.ReadFile(filepath)
		if err != nil {
			return clientID
		}
		clientID = string(bs)
	} else {
		err := os.MkdirAll(fmt.Sprintf("%s/kaf-wifi", config), 0700)
		if err != nil {
			return clientID
		}
		_ = os.WriteFile(filepath, []byte(clientID), 0700)
	}
	return clientID
}

func GetEnv(key, defaultvalue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultvalue
}

var colors = []string{
	"#61005e",
	"#70706d",
	"#890029",
	"#c4000e",
	"#6d001d",
	"#6a00bd",
	"#f10000",
	"#0071b1",
	"#f9bc00",
	"#2c0077",
	"#ba009a",
	"#009047",
	"#009d9e",
	"#222e85",
	"#bd002e",
	"#009d1a",
	"#75a500",
}

func ParseInt(v string) int {
	v = strings.ReplaceAll(v, ",", "")
	i, err := strconv.ParseInt(v, 0, 32)
	if err != nil {
		return 0
	}
	return int(i)
}

func GenCover(title, author, color string, img int) (string, error) {
	query := url.Values{}
	query.Add("title", title)
	query.Add("author", author)
	query.Add("g_loc", "BR")
	query.Add("top_text", "kaf")
	query.Add("g_text", "")
	if img >= 0 && img <= 41 {
		query.Add("img_id", fmt.Sprintf("%d", img))
	} else {
		query.Add("img_id", fmt.Sprintf("%d", rand.Intn(41)))
	}
	if strings.HasPrefix(color, "#") {
		query.Add("color", strings.TrimLeft(color, "#"))
	} else {
		i := ParseInt(color)
		if i == 0 {
			i = rand.Intn(17)
		}
		color := colors[i]
		query.Add("color", strings.TrimLeft(color, "#"))
	}

	uri := fmt.Sprintf("https://orly.nanmu.me/api/generate?%s", query.Encode())
	res, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	tempDir, err := os.MkdirTemp("", "kaf-cli")
	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	coverfile := filepath.Join(tempDir, fmt.Sprintf("%s.jpg", title))
	err = os.WriteFile(coverfile, bs, 0666)
	if err != nil {
		return "", err
	}
	return coverfile, nil
}

type Number interface {
	~int | ~uint
}

func defaultString(src, dst string) string {
	if src == "" {
		return dst
	}
	return src
}
func defalutInt[T Number](src, dst T) T {
	if src == 0 {
		return dst
	}
	return src
}
func defaultBool(src, dst bool) bool {
	if src {
		return src
	}
	return dst
}
