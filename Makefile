GOOS_LINUX=linux
GOOS_WINDOWS=windows
GOARCH_AMD64=amd64
GOARCH_ARM64=arm64
GOARCH_ARM=arm
BUILD_DIR=./build
BIN_NAME=netiso
MAIN_GO=cmd/netiso/netiso.go

run:
	@echo "Running NetISO..."
	go run ${MAIN_GO}

clean:
	@echo "Cleaning build directory..."
	rm -rf ${BUILD_DIR}/${BIN_NAME}*

all: build-all

build-all: build-linux-amd64 build-linux-arm64 build-linux-arm build-windows-amd64 build-windows-arm64 build-windows-arm
	@echo "Build for all platforms and architectures completed."

build-linux-amd64:
	@echo "Building for Linux AMD64..."
	GOOS=${GOOS_LINUX} GOARCH=${GOARCH_AMD64} go build -o ${BUILD_DIR}/${BIN_NAME}-linux-amd64 ${MAIN_GO}

build-linux-arm64:
	@echo "Building for Linux ARM64..."
	GOOS=${GOOS_LINUX} GOARCH=${GOARCH_ARM64} go build -o ${BUILD_DIR}/${BIN_NAME}-linux-arm64 ${MAIN_GO}

build-linux-arm:
	@echo "Building for Linux ARM..."
	GOOS=${GOOS_LINUX} GOARCH=${GOARCH_ARM} go build -o ${BUILD_DIR}/${BIN_NAME}-linux-arm ${MAIN_GO}

build-windows-amd64:
	@echo "Building for Windows AMD64..."
	GOOS=${GOOS_WINDOWS} GOARCH=${GOARCH_AMD64} go build -o ${BUILD_DIR}/${BIN_NAME}.exe ${MAIN_GO}

build-windows-arm64:
	@echo "Building for Windows ARM64..."
	GOOS=${GOOS_WINDOWS} GOARCH=${GOARCH_ARM64} go build -o ${BUILD_DIR}/${BIN_NAME}-arm64.exe ${MAIN_GO}

build-windows-arm:
	@echo "Building for Windows ARM..."
	GOOS=${GOOS_WINDOWS} GOARCH=${GOARCH_ARM} go build -o ${BUILD_DIR}/${BIN_NAME}-arm.exe ${MAIN_GO}
