package utils

import (
	"bytes"
	"strings"
)

const (
	htmlPStart = `<p class="content">`
	htmlPEnd   = "</p>"
)

func AddPart(buff *bytes.Buffer, content string) {
	if strings.HasSuffix(content, "==") ||
		strings.HasSuffix(content, "**") ||
		strings.HasSuffix(content, "--") ||
		strings.HasSuffix(content, "//") {
		buff.WriteString(content)
		return
	}
	buff.WriteString(htmlPStart)
	buff.WriteString(content)
	buff.WriteString(htmlPEnd)
}
