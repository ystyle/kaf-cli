## kaf-cli

> 把txt文本转成epub电子书的命令行工具

### 功能
- 傻瓜操作模式(把txt文件拖到`kaf-cli.exe`上面自动转换)
- 自定义封面
- 支持生成Orly风格的书籍封面
- 自动识别书名和章节(示例中所有用法都会自动识别)
- 自动识别字符编码(自动解决中文乱码)
- 自定义章节标题识别规则
- 自定义卷的标题识别规则
- 自动给章节正文生成加粗居中的标题
- 卷名使用精美的双线边框样式
- 章节序号和标题分离显示（可选）
- 自定义 CSS 样式支持
- 智能 HTML 标签处理（保留 epub 支持的标签，转义其他标签）
- 自定义标题对齐方式
- 段落自动识别
- 段落自动缩进
- 自定义段落缩进字数
- 自定义段落间距
- 自定义行间距
- 自定义书籍语言
- epub格式支持嵌入字体
- 知轩藏书格式文件名会自动提取书名和作者, 例: `《希灵帝国》（校对版全本）作者：远瞳.txt` 或 `《希灵帝国》作者：远瞳.txt`
- 超快速(epub格式生成300章/s以上速度, 4000章15s不到)
- 自动转为mobi格式

### 下载
- 电脑版kaf-cli: [Github下载](https://github.com/ystyle/kaf-cli/releases/latest)
- 手机版kaf: [Github下载](https://github.com/ystyle/kaf-cli/releases/tag/android)
- 电脑版wifi传书kaf-wifi: [Github下载](https://github.com/ystyle/kaf-wifi/releases/latest)
- 全部软件 [百度网盘下载 `https://pan.baidu.com/s/1EPkLJ7WIJYdYtRHBEMqw0w?pwd=h4np`](https://pan.baidu.com/s/1EPkLJ7WIJYdYtRHBEMqw0w?pwd=h4np)
- 包管理器安装:
  - Archlinux 可以在aur上安装 [`yay -S kaf-cli kaf-wifi`](https://aur.archlinux.org/packages/kaf-cli/) (包括龙芯loongarch64版本)
  - windows使用winget: `winget install kaf-cli`

### 使用方法
1. 解压, 把小说直接拖到 `kaf-cli.exe` 文件上面
1. 等转换完，目录下会生成epub、azw3、mobi文件
   - mobi格式需要有kindlegen才会生成(windows、mac版本已经自带)
1. 自定义封面功能
   在拖拽模式下, 如果目录下有`cover.png`文件会自动添加为封面、支持jpg、png格式， 如果需要指定其它文件或jpg格式需要使用命令行模式   
1. 其它自定义功能请用命令行模式

### 效果
![效果图片](assets/2021-06-20_12-13-34.png)
![效果图片](assets/2020-01-21_12-02.png)

### 命令行模式参数

```text
Usage of kaf-cli:
```shell
  -align string
        标题对齐方式: left、center、righ。环境变量KAF_CLI_ALIGN可修改默认值 (default "center")
  -author string
        作者 (default "YSTYLE")
  -bookname string
        书名: 默认为txt文件名
  -bottom string
        段落间距(单位可以为em、px) (default "1em")
  -cover string
        封面图片可为: 本地图片, 和orly。 设置为orly时生成orly风格的封面, 需要连接网络。 (default "cover.png")
  -cover-orly-color string
        orly封面的主题色, 可以为1-16和hex格式的颜色代码, 不填时随机
  -cover-orly-idx int
        orly封面的动物, 可以为0-41, 不填时随机, 具体图案可以查看: https://orly.nanmu.me (default -1)
  -custom-css-file string
        自定义 CSS 文件路径，用于覆盖默认样式
  -exclude string
        排除无效章节/卷的正则表达式 (default "^第[0-9一二三四五六七八九十零〇百千两 ]+(部门|部队|部属|部分|部件|部落|部.*：$)")
  -filename string
        txt 文件名
  -font string
        嵌入字体, 之后epub的正文都将使用该字体
  -format string
        书籍格式: all、epub、mobi、azw3。环境变量KAF_CLI_FORMAT可修改默认值 (default "epub")
  -indent uint
        段落缩进字数 (default 2)
  -lang string
        设置语言: en,de,fr,it,es,zh,ja,pt,ru,nl。环境变量KAF_CLI_LANG可修改默认值 (default "en")
  -line-height string
        行高(用于设置行间距, 默认为1.5rem)
  -match string
        匹配标题的正则表达式, 不写可以自动识别, 如果没生成章节就参考教程。例: -match 第.{1,8}章 表示第和章字之间可以有1-8个任意文字
  -max uint
        标题最大字数 (default 35)
  -out string
        输出文件名，不需要包含格式后缀
  -separate-chapter-number
        是否分离章节序号和标题样式（序号单独一行显示）
  -tips
        添加本软件教程 (default true)
  -unknow-title string
        未知章节默认名称 (default "章节正文")
  -volume-match string
        卷匹配规则,设置为false可以禁用卷识别 (default "^第[0-9一二三四五六七八九十零〇百千两 ]+[卷部]")
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

### 自定义 CSS 样式

kaf-cli 支持通过 CSS 文件自定义样式，可以覆盖默认的样式设置。

#### 可用的 CSS 类名

- `h2.volume` - 卷名样式
- `h3.title` - 章节标题样式
- `h3.title span.chapter-number` - 章节序号样式
- `.content` - 正文段落样式
- `body` - 整体样式

#### 使用示例

1. 创建一个 CSS 文件（如 `mystyle.css`）：
```css
/* 修改卷名样式 */
h2.volume {
    border-top: 5px solid red;
    border-bottom: 5px solid red;
    color: darkred;
}

/* 修改章节标题样式 */
h3.title {
    color: blue;
    border-bottom: 3px solid blue;
}

/* 修改章节序号样式 */
h3.title span.chapter-number {
    color: green;
}
```

2. 使用 `--custom-css-file` 参数：
```shell
kaf-cli -filename 小说.txt --custom-css-file mystyle.css
```

#### 高级用法

- 调整行间距：`.content { line-height: 2; }`
- 修改字体大小：`h3.title { font-size: 1.5em; }`
- 修改缩进：`.content { text-indent: 4em; }`
- 修改边距：`body { margin: 2em; }`

### 自定义章节匹配规则
>以下全部示例都可以自动识别，不需要自己设定标题格式了， 一般用上用上面的例子就行了

>规则支持[正则表达式](http://deerchao.net/tutorials/regex/regex.htm)， 要自定义标题格式参考以下几个例子, 以下例子小说都在D盘

快速入门
>注意`.`是指任意字，`{1,8}`里边的数值代表1-8个字, `\d+`是指纯数字，`\w`是指中英文的字



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

### MCP 版本使用说明
>以cherry studio为例
- 1. 点击下载[kaf-mcp](https://github.com/ystyle/kaf-cli/releases?q=mcp-), 并解压, 复制kaf-mcp的路径
- 2. 配置mcp
   - 命令这里填上面的kaf-mcp命令的路径
   - 在环境变量里添加`KAF_DIR=xxx`, xxx换成txt小说所在的目录，转换后的小说也会放这里
   ![](assets/mcp-cherry-studio-config.png)
- 3. 在ai对话里转换小说
   ![](assets/mcp-cherry-studio-use.png)

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

