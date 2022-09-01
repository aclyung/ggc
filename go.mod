module main

// legacy code is imcomplete and highly unstable
// using core library is recommended

require (
	almeng.com/glang v1.0.0
)

require (
	github.com/almenglee/general v0.1.0 // indirect
	github.com/llir/ll v0.0.0-20210719001141-246f2b6b1fa9 // indirect
	github.com/llir/llvm v0.3.4 // indirect
	github.com/mewmew/float v0.0.0-20211212214546-4fe539893335 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace almeng.com/glang v1.0.0 => ./glang

go 1.18
