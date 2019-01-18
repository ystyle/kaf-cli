## TmdTextEpub

> 把txt文本转成epub电子书的命令行工具

### Usage
```$xslt
Usage of TmdTextEpub.exe:
  -author string
        作者 (default "YSTYLE")
  -filename string
        txt 文件名
  -match string
        匹配标题的正则表达式, 例: -match 第.{1,8}章 表示第和章字之间可以有1-8个任意文字 (default "第.{1,8}章")

```

### 示例
1. [点击下载](https://github.com/ystyle/TmdTextEpub/releases)
1. 把小说和`TmdTextEpub.exe`放到`D`盘
1. 按以下其中一种方法打开命令行
    - 按`win + r` 输入 `cmd` 然后输入以下命令
    - 按`win + x + i` 输入以下命令


```shell
cd d:/
d:/TmdTextEpub.exe -author 乱 -filename d:/全职法师.txt -match 第.{1,4}章
```

### 手工构建
```$xslt
go build -ldflags "-s -w" -o TmdTextEpub.exe main.go
```