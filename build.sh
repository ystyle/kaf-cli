GOOS=windows go build -ldflags "-s -w" -o TmdTextEpub.exe main.go

GOOS=linux go build -ldflags "-s -w" -o TmdTextEpub-linux main.go

GOOS=darwin go build -ldflags "-s -w" -o TmdTextEpub-darwin main.go

echo "done!"
