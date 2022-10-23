module main

// legacy code is imcomplete and highly unstable
// using core library is recommended

require almeng.com/glang v1.0.0

require (
	almeng.com/glang-vm v0.0.0-00010101000000-000000000000 // indirect
	github.com/almenglee/general v0.1.0 // indirect
	github.com/inhies/go-bytesize v0.0.0-20220417184213-4913239db9cf // indirect
	github.com/llir/llvm v0.3.4 // indirect
	github.com/mewmew/float v0.0.0-20211212214546-4fe539893335 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/tools v0.1.12 // indirect
)

replace almeng.com/glang v1.0.0 => ./glang

replace almeng.com/glang-vm => ./vm

go 1.18
