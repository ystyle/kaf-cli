module github.com/ystyle/kaf-cli

require (
	github.com/766b/mobi v0.0.0-20200528201125-c87aa9e3c890
	github.com/bmaupin/go-epub v0.11.0
	github.com/leotaku/mobi v0.0.0-20220405163106-82e29bde7964
	github.com/ystyle/google-analytics v0.0.0-20210425064301-a7f754dd0649
	golang.org/x/net v0.0.0-20210505024714-0287a6fb4125
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da
	golang.org/x/text v0.3.7
)

require (
	github.com/gabriel-vasile/mimetype v1.3.1 // indirect
	github.com/gofrs/uuid v3.1.0+incompatible // indirect
	github.com/vincent-petithory/dataurl v0.0.0-20191104211930-d1553a71de50 // indirect
)

replace golang.org/x/text => github.com/golang/text v0.3.2

go 1.19
