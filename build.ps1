$env:GOOS="linux"
go build -ldflags "-s -w" -o kaf-cli-linux main.go

$env:GOOS="darwin"
go build -ldflags "-s -w" -o kaf-cli-darwin main.go


$env:GOOS="windows"
go build -ldflags "-s -w" -o kaf-cli.exe main.go

$env:GOOS="windows"
$env:GOARCH="386"
go build -ldflags "-s -w" -o kaf-cli_32.exe main.go

echo "done!"
Start-Sleep -Seconds 20 main.go