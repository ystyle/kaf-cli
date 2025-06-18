module github.com/ystyle/kaf-cli

require (
	github.com/766b/mobi v0.0.0-20200528201125-c87aa9e3c890
	github.com/go-shiori/go-epub v1.2.1
	github.com/leotaku/mobi v0.5.0
	github.com/ystyle/google-analytics v0.0.0-20210425064301-a7f754dd0649
	golang.org/x/net v0.39.0
	golang.org/x/text v0.24.0
)

require (
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/gofrs/uuid/v5 v5.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mark3labs/mcp-go v0.27.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/vincent-petithory/dataurl v1.0.0 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
)

replace github.com/go-shiori/go-epub v1.2.1 => github.com/ystyle/go-epub v0.0.0-20250425133851-dba4e6a949ec

go 1.23.1
