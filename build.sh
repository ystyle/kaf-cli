export GOPROXY=https://goproxy.cn,direct
GOOS=windows go build -ldflags "-s -w" -o TmdTextEpub.exe main.go
GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o TmdTextEpub_32.exe main.go
GOOS=linux go build -ldflags "-s -w" -o TmdTextEpub-linux main.go
GOOS=darwin go build -ldflags "-s -w" -o TmdTextEpub-darwin main.go
echo "done!"
