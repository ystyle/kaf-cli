### 环境准备
- 安装`wasmtime`运行时， 或其它支持go 1.21版本生成的wasip1 wasm的运行时
  - 经测试wasmer不支持运行gov1.21.0生成的wasm

### 运行
```shell
wasmtime run --mapdir /tmp::/tmp --mapdir /txt::.  ./kaf-cli.wasm -- -filename /txt/《诡秘之主》（精校版全本）作者：爱潜水的乌贼.txt -format epub
```
- `--mapdir` 需要把运行时里的`/tmp`和txt存放位置映射到物理机的目录
- 运行时的参数和`kaf-cli`的参数需要用`--`分隔开
- `kaf-cli`的小说路径参数，需要指定运行时里映射的路径

### 源码构建
1. 需要提前安装[`go编译器`](https://go.dev)
2. 下载：https://github.com/ystyle/kaf-cli
3. 编译`wasm/wasi`版本: `OARCH=wasm GOOS=wasip1 go build -o kaf-cli.wasm cmd/cli.go`