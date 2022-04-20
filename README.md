## kaf-cli

> 把txt文本转成epub电子书的命令行工具

### 功能
- 傻瓜操作模式(把txt文件拖到`kaf-cli.exe`上面自动转换)
- 自定义封面
- 自动识别书名和章节(示例中所有用法都会自动识别)
- 自动识别字符编码(自动解决中文乱码)
- 自定义章节标题识别规则
- 自动给章节正文生成加粗居中的标题
- 自定义标题对齐方式
- 段落自动识别
- 段落自动缩进
- 自定义段落缩进字数
- 自定义段落间距
- 自定义书籍语言
- 知轩藏书格式文件名会自动提取书名和作者, 例: `《希灵帝国》（校对版全本）作者：远瞳.txt`
- 超快速(130章/s以上速度, 4000章30s不到)
- 自动转为mobi格式

### 下载
- 电脑版kaf-cli: [Github下载](https://github.com/ystyle/kaf-cli/releases/latest)
- 手机版kaf: [Github下载](https://github.com/ystyle/kaf-cli/releases/tag/android)
- 电脑版wifi传书kaf-wifi: [Github下载](https://github.com/ystyle/kaf-wifi/releases/latest)
- 全部软件 [百度网盘下载 `https://pan.baidu.com/s/1EPkLJ7WIJYdYtRHBEMqw0w?pwd=h4np`](https://pan.baidu.com/s/1EPkLJ7WIJYdYtRHBEMqw0w?pwd=h4np)
- Archlinux 可以在aur上安装 [`yay -S kaf-cli kaf-wifi`](https://aur.archlinux.org/packages/kaf-cli/)

### 使用方法
1. 解压, 把小说直接拖到 `kaf-cli.exe` 文件上面
1. 等转换完，目录下会生成epub、azw3、mobi文件
   - mobi格式需要有kindlegen才会生成(windows、mac版本已经自带)
1. 自定义封面功能
   在拖拽模式下, 如果目录下有`cover.png`文件会自动添加为封面、支持jpg、png格式， 如果需要指定其它文件或jpg格式需要使用命令行模式   
1. 其它自定义功能请用命令行模式

### 效果
![效果图片](2021-06-20_12-13-34.png)
![效果图片](2020-01-21_12-02.png)

### 命令行模式参数

```text
Usage of kaf-cli.exe:
  -align string
        标题对齐方式: left、center、righ (default "center")
  -author string
        作者 (default "YSTYLE")
  -bookname string
        书名: 默认为txt文件名
  -bottom string
        段落间距(单位可以为em、px) (default "1em")
  -cover string
        封面图片 (default "cover.png")
  -filename string
        txt 文件名
  -format string
        书籍格式: all、epub、mobi、azw3 (default "all")
  -indent uint
        段落缩进字数 (default 2)
  -lang string
        设置语言: en,de,fr,it,es,zh,ja,pt,ru,nl。 环境变量KAF_CLI_LANG可修改默认值 (default "zh")
  -match string
        匹配标题的正则表达式, 不写可以自动识别, 如果没生成章节就参考教程。例: -match 第.{1,8}章 表示第和章字之间可以有1-8个任意文字 (default "自动匹配,可自定义")
  -max uint
        标题最大字数 (default 35)
  -out string
        输出文件名，不需要包含格式后缀
  -tips
        添加本软件教程 (default true)

```

>PS: 在darwin(mac、osx)上`-tips`参数要设置为false的方法 `kaf-cli -filename 小说.txt -tips=0`

### 命令行模式说明

转换`全职法师.txt`, 并设置作者名为`乱`
```shell
# windows 10: win + s 搜索powershell 
cd d:/
d:/kaf-cli.exe -author 乱 -filename d:/全职法师.txt

# linux / mac下, 把kaf-cli-linux/kaf-cli-darwin重命名为kaf-cli, 放到用户目录
# 把小说和kaf-cli放到用户目录下,  打开终端执行
cd ~
./kaf-cli -author 乱 -filename ./全职法师.txt

# 如果kaf-cli放到path里了, 或者在aur安装的可以执行:
kaf-cli -author 乱 -filename ~/全职法师.txt

# 命令行的简单模式（功能和拖拽模式一样）
kaf-cli ~/全职法师.txt
```

### 自定义章节匹配规则
>以下全部示例都可以自动识别，不需要自己设定标题格式了， 一般用上用上面的例子就行了

>规则支持[正则表达式](http://deerchao.net/tutorials/regex/regex.htm)， 要自定义标题格式参考以下几个例子, 以下例子小说都在D盘



自定义章节匹配, 章节格式为`第x节`: 
```shell
d:/kaf-cli.exe -filename d:/ebbok.txt -match "第.{1,8}节"
```

自定义章节匹配, 章节格式为`Section 1` ~ `Section 100`: 
```shell
d:/kaf-cli.exe -filename d:/ebbok.txt -match "Section \d+"
```

自定义章节匹配, 章节格式为`Chapter xxx`: 
```shell
d:/kaf-cli.exe -filename d:/ebbok.txt -match "Chapter .{1,8}"
```


### 在任意位置执行命令
- windows
  - 把`kaf-cli.exe` 和 `kindlegen.exe` 放`c:/windows/`下边
  - 以后可以把小说放任意目录，都可以很简单执行转换， 第一步只需要做一次， 以下为每次转换小说的操作，
    - 打开小说在的文件夹, 按住`Shift键`不放，鼠标右击文件夹空白位置
    - 在右键菜单选择 `用命令行打开` 或 `以PowerShell打开`
    - 以上命令可以改为 `kaf-cli.exe -filename 全职法师.txt`,  现在可以不用写盘符了
- linux(理论上mac也可以是这样的)
  - 软件可以放任意地方, 比如`~/application/kaf-cli`，在`~/.bashrc` 或 `~/.zshrc` 最后一行添加 `export PATH=$HOME/application:$PATH`
  - 打开终端, 执行命令为: `kaf-cli -filename ~/全职法师.txt`


### 手动把书转为kindle的mobi格式
>新版如果检测到有kindlegen程序时会自动转为mobi

1. 下载[kindlegen](https://github.com/ystyle/kaf-cli/releases/kindlegen/) (github备份，官网已经不提供下载)
2. 同样放到`d:`盘根目录下， 把epub拖拽到kindlegen.exe上面， 或执行以下命令转换
  ```shell
  cd d:/
  d:/kindlegen.exe d:/全职法师.epub
  ```
3. 在d盘就能找到mobi文件，复制到kindle的documents目录下，打开kindle就能看到小说了

