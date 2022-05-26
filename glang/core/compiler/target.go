package compiler

type Target struct {
	Arch   Arch
	Vendor Vendor
	Sys    Sys
}

type Arch string

const (
	I386    Arch = "i386"
	X86_64  Arch = "x86_64"
	ARM     Arch = "arm"
	ARM64   Arch = "arm64"
	AARCH64 Arch = "aarch64"
)

type Vendor string

const (
	PC    Vendor = "pc"
	APPLE Vendor = "apple"
)

type Sys string

const (
	LINUX  Sys = "linux"
	DARWIN Sys = "darwin"
	WIN32  Sys = "win32"
)

func NewTarget(a Arch, v Vendor, os Sys) Target {
	return Target{a, v, os}
}

func (t Target) String() string {
	return string("--target="+t.Arch+"-") + string(t.Sys) + "-" + string(t.Vendor)
}
