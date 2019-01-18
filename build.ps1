$env:GOOS="windows"
go build -ldflags "-s -w" -o TmdTextEpub.exe main.go

$env:GOOS="linux"
go build -ldflags "-s -w" -o TmdTextEpub-linux main.go

$env:GOOS="darwin"
go build -ldflags "-s -w" -o TmdTextEpub-darwin main.go


echo "done!"
Start-Sleep -Seconds 20 main.go