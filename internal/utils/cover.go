package utils

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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

func ParseInt(v string) int {
	v = strings.ReplaceAll(v, ",", "")
	i, err := strconv.ParseInt(v, 0, 32)
	if err != nil {
		return 0
	}
	return int(i)
}
