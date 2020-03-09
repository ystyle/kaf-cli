module github.com/ystyle/TmdTextEpub

require (
	github.com/bmaupin/go-epub v0.5.3
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a
	golang.org/x/text v0.3.0
)

replace golang.org/x/text => github.com/golang/text v0.3.2

go 1.13
