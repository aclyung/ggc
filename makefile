NAME = glang
OS = ${MAKEOS}
ARCH = ${MAKEARCH}
OUTDIR = ./build
OUTNAME = name

# .PHONY: build
all: arm_all amd_all
# OS configurations
darwin linux:
	$(eval OS=$@)

windows:
	$(eval OS=$@)
	$(eval EXT=.exe)

# ARCH configurations
arm64 amd64 386:
	$(eval ARCH=$@)

arm_all: 
	make arm64 darwin build
	make arm64 windows build
	make arm64 linux build

amd_all: 
	make amd64 darwin build
	make amd64 windows build
	make amd64 linux build

cross: arm_all amd_all

name:
	$(eval OUTNAME = ${OUTDIR}/${NAME}_${OS}_${ARCH}${EXT})

test:
	$(eval OUTNAME = ${NAME}_${OS}_${ARCH}${EXT})
	env GOOS=${OS} GOARCH=${ARCH} go build -o ${OUTNAME}

rm_test:
	rm ./glang_*

build: name
	env GOOS=${OS} GOARCH=${ARCH} go build -o ${OUTNAME}



