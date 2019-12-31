## TmdTextEpub

> 把txt文本转成epub电子书的命令行工具

### 功能
- 自动识别书名和章节(以下示例的章节名都可以自动匹配)
- 自定义章节匹配
- 自动给章节正文生成加粗居中的标题
- 段落自动识别
- 段落自动缩进

### 使用方法
```shell
Usage of D:\TmdTextEpub.exe:
  -author string
        作者 (default "YSTYLE")
  -bookname string
        书名: 默认为txt文件名
  -filename string
        txt 文件名
  -match string
        匹配标题的正则表达式, 不写可以自动识别, 如果没生成章节就参考教程。例: -match 第.{1,8}章 表示第和章字之间可以有1-8个任意文字 (default "自动匹配,可自定义")
  -tips
        添加本软件教程 (default true)
```

### 示例
1. [点击下载](https://github.com/ystyle/TmdTextEpub/releases/latest)
1. 把小说和`TmdTextEpub.exe`放到`D`盘
1. 按以下其中一种方法打开命令行
    - 按`win + r` 输入 `cmd` 然后输入以下命令
    - 按`win + x + i` 输入以下命令

把`全职法师.txt`生成epub, 并设置作者名为`乱`
```shell
cd d:/
d:/TmdTextEpub.exe -author 乱 -filename d:/全职法师.txt
```

自定义章节匹配, 章节格式为`第x节`: 
```shell
cd d:/
d:/TmdTextEpub.exe -filename d:/ebbok.txt -match "第.{1,8}节"
```

自定义章节匹配, 章节格式为`Section 1` ~ `Section 100`: 
```shell
cd d:/
d:/TmdTextEpub.exe -filename d:/ebbok.txt -match "Section \d+"
```

自定义章节匹配, 章节格式为`Chapter xxx`: 
```shell
cd d:/
d:/TmdTextEpub.exe -filename d:/ebbok.txt -match "Chapter .{1,8}"
```

### 把书转为kindle的mobi格式
1. 在官网下载[kindlegen](https://www.amazon.com/gp/feature.html?ie=UTF8&docId=1000765211)
2. 同样放到`d:`盘根目录下， 执行以下命令转换
  ```shell
  cd d:/
  d:/kindlegen.exe d:/全职法师.mobi
  ```
3. 在d盘就能找到mobi文件了，复制到kindle的documents目录下，打开kindle就能看到小说了

### 手工构建
```$xslt
go build -ldflags "-s -w" -o TmdTextEpub.exe main.go
```
快捷方法:
- windows: 在`build.ps1`上右键选择`以 Power Shell 运行`
- linux: 在本项目源码目录执行 `./build.sh`
