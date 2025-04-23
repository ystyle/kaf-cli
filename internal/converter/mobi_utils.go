package converter

import (
	"fmt"
	"time"

	"github.com/ystyle/kaf-cli/internal/utils"
)

func ConverToMobi(bookname, lang string) {
	command := utils.LookKindlegen()
	fmt.Printf("\n检测到Kindle格式转换器: %s，正在把书籍转换成Kindle格式...\n", command)
	fmt.Println("转换mobi比较花时间, 大约耗时1-10分钟, 请等待...")
	start := time.Now()
	utils.Run(command, "-dont_append_source", "-locale", lang, "-c1", bookname)
	// 计算耗时
	end := time.Now().Sub(start)
	fmt.Println("转换为mobi格式耗时:", end)
}
