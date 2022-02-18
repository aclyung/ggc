module main


// legacy code is imcomplete and highly unstable
// using core library is recommended

require (
	almeng.com/glang v1.0.0
	github.com/c-bata/go-prompt v0.2.6
)

require (
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mattn/go-tty v0.0.3 // indirect
	github.com/pkg/term v1.2.0-beta.2 // indirect
	golang.org/x/sys v0.0.0-20200918174421-af09f7315aff // indirect
)

replace almeng.com/glang v1.0.0 => ./glang

go 1.17
