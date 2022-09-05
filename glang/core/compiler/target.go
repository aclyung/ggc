package compiler

import (
	"fmt"
	"os"
	"strings"
)

type Target struct {
	Arch   Arch
	Vendor Vendor
	Sys    Sys
}

type Arch string

var archs = map[string]Arch{
	"i386":    I386,
	"x86_64":  X86_64,
	"amd64":   AMD64,
	"arm":     ARM,
	"arm64":   ARM64,
	"aarch64": AARCH64,
}

const (
	I386    Arch = "i386"
	X86_64  Arch = "x86_64"
	AMD64   Arch = X86_64
	ARM     Arch = "arm"
	ARM64   Arch = "arm64"
	AARCH64 Arch = "aarch64"
)

type Vendor string

var vendors = map[string]Vendor{
	"pc":    PC,
	"apple": APPLE,
}

const (
	PC    Vendor = "pc"
	APPLE Vendor = "apple"
)

type Sys string

var systems = map[string]Sys{
	"linux":  LINUX,
	"darwin": DARWIN,
	"win32":  WIN32,
}

const (
	LINUX  Sys = "linux"
	DARWIN Sys = "darwin"
	WIN32  Sys = "win32"
)

func TargetFromTriple(t string) Target {
	if t == "" {
		return NewTarget(ARM, APPLE, DARWIN)
	}

	triple := strings.Split(t, "-")
	if len(triple) != 3 {
		fmt.Println("Invalid target triple '" + t + "'")
		os.Exit(1)
	}
	arch, archValid := archs[triple[0]]
	vendor, vendorValid := vendors[triple[1]]
	sys, sysValid := systems[triple[2]]
	if !archValid || !vendorValid || !sysValid {
		fmt.Println("Invalid target triple '" + t + "'")
		os.Exit(1)
	}
	return NewTarget(arch, vendor, sys)
}

func NewTarget(a Arch, v Vendor, os Sys) Target {
	return Target{a, v, os}
}

func (t Target) String() string {
	return string("--target="+t.Arch+"-") + string(t.Sys) + "-" + string(t.Vendor)
}
