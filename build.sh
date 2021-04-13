export GOPROXY=https://goproxy.cn,direct
GOOS=windows go build -ldflags "-s -w" -o kaf-cli.exe main.go
GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o kaf-cli_32.exe main.go
GOOS=linux go build -ldflags "-s -w" -o kaf-cli-linux main.go
GOOS=darwin go build -ldflags "-s -w" -o kaf-cli-darwin main.go
echo "done!"
