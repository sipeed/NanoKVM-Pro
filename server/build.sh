#!/bin/bash

set -e

# Cleanup function to remove temporary files on exit
cleanup() {
    if [ -f "dl_lib/libkvm.so" ] || [ -f "dl_lib/libkvm.so.0" ]; then
        echo "[INFO] Cleaning up temporary files..."
        rm -f dl_lib/libkvm.so*
    fi
}

# Register cleanup function to be called on exit
trap cleanup EXIT

# Print colored log messages
log_info() {
    echo -e "\033[32m[INFO]\033[0m $1"
}

log_error() {
    echo -e "\033[31m[ERROR]\033[0m $1" >&2
}

log_warn() {
    echo -e "\033[33m[WARN]\033[0m $1" >&2
}

# Check if required dependencies are available
check_dependencies() {
    log_info "Checking dependencies..."

    # Check for Go
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed or not in PATH"
        return 1
    fi

    # Check for patchelf
    if ! command -v patchelf &> /dev/null; then
        log_error "patchelf is not installed or not in PATH"
        return 1
    fi

    # Check for Python (required for getconfig.py)
    if ! command -v python3 &> /dev/null && ! command -v python &> /dev/null; then
        log_error "Python is not installed or not in PATH"
        return 1
    fi

    log_info "All dependencies are available"
    return 0
}

# Check if current directory is the server project root
is_prj_rootpath() {
    local REQUIRED_DIRS=("common" "config" "dl_lib" "include" "logger" "middleware" "proto" "router" "service" "utils")
    local all_dirs_exist=true
    local pwd=$(pwd)

    for dir in "${REQUIRED_DIRS[@]}"; do
        if [ ! -d "$dir" ]; then
            log_error "Missing required project directory: $dir"
            all_dirs_exist=false
        fi
    done

    if [ "$all_dirs_exist" = false ]; then
        log_error "Current directory: ${pwd} is not the <server> project root directory"
        log_error "This script must be executed from the <server> project root directory."
        return 1
    fi

    log_info "Project<server> root directory verified successfully"
    return 0
}

# Find ARM64 cross compiler toolchain
find_pro_cc() {
    local INI_FILE="../support/toolchains/toolchain.ini"
    if [ ! -f "$INI_FILE" ]; then
        log_error "Toolchain configuration file not found: $INI_FILE"
        return 1
    fi

    if [ ! -f "../support/scripts/getconfig.py" ]; then
        log_error "Toolchain configuration script not found: ../support/scripts/getconfig.py"
        return 1
    fi

    # Try python3 first, then python
    local python_cmd=""
    if command -v python3 &> /dev/null; then
        python_cmd="python3"
    elif command -v python &> /dev/null; then
        python_cmd="python"
    fi

    cc=$($python_cmd ../support/scripts/getconfig.py "$INI_FILE" "toolchain" "cc")
    if [ -z "$cc" ]; then
        log_error "Failed to get cross compiler from toolchain configuration"
        return 1
    fi

    # Verify the compiler exists
    if ! command -v "$cc" &> /dev/null; then
        log_error "Cross compiler not found: $cc"
        return 1
    fi

    echo "$cc"
    return 0
}

# Build the project for ARM64 platform
build_project() {
    log_info "Starting project build..."
    local cc="$1"
    local arch="$2"

    log_info "Using ARCH: $arch"
    log_info "Using CC: $cc"

    # Calculate sysroot path from compiler path
    local toolchain_dir=$(dirname $(dirname "$cc"))
    local sysroot="${toolchain_dir}/aarch64-none-linux-gnu/libc"
    log_info "Using sysroot: $sysroot"

    # Build shared library
    log_info "Building shared library..."
    if ! "$cc" -fPIC -shared -o dl_lib/libkvm.so dl_lib/kvmimpl.cpp; then
        log_error "Failed to build shared library"
        return 1
    fi

    if ! cp dl_lib/libkvm.so dl_lib/libkvm.so.0; then
        log_error "Failed to copy shared library"
        return 1
    fi

    # Get version information
    VERSION=$(cat version.txt 2>/dev/null || echo "dev")
    BUILD_TIME=$(date +"%Y-%m-%d %H:%M:%S")
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    GO_VERSION=$(go version | awk '{print $3}')
    BUILD_USER=$(whoami)@$(hostname)
    GIT_BRANCH=$(git branch --show-current 2>/dev/null || echo "unknown")

    log_info "Building version: $VERSION"
    log_info "Commit: $COMMIT"
    log_info "Build time: $BUILD_TIME"

    # CGO flags for cross-compilation with opus
    local cgo_cflags="-I${sysroot}/usr/include"
    local cgo_ldflags="-L${sysroot}/usr/lib"
    
    log_info "CGO_CFLAGS: $cgo_cflags"
    log_info "CGO_LDFLAGS: $cgo_ldflags"

    # Build Go application
    log_info "Building Go application..."
    if ! CGO_ENABLED=1 GOOS=linux GOARCH=${arch} CC=${cc} \
        CGO_CFLAGS="$cgo_cflags" CGO_LDFLAGS="$cgo_ldflags" \
        GOEXPERIMENT=boringcrypto go build -ldflags "\
        -X 'main.Version=$VERSION' \
        -X 'main.BuildTime=$BUILD_TIME' \
        -X 'main.Commit=$COMMIT' \
        -X 'main.GoVersion=$GO_VERSION' \
        -X 'main.GitBranch=$GIT_BRANCH' \
        -X 'main.BuildUser=$BUILD_USER'"; then
        log_error "Go build failed"
        return 1
    fi

    log_info "Project built successfully"

    # Add runtime library path
    log_info "Adding runtime library path to executable..."
    if ! patchelf --add-rpath \$ORIGIN/dl_lib NanoKVM-Server; then
        log_error "Failed to add runtime library path"
        return 1
    fi

    log_info "Runtime library path added successfully"
    return 0
}

# Main function - optimized for ARM64 only
main() {
    log_info "Starting ARM64 build process for NanoKVM-Pro..."

    # Check dependencies
    check_dependencies || return 1

    # Verify project structure
    is_prj_rootpath || return 1

    # Find cross compiler
    local cc=$(find_pro_cc || return 1)
    local arch="arm64"

    # Build the project
    build_project "$cc" "$arch" || return 1

    log_info "Build process completed successfully"
    log_info "Executable: NanoKVM-Server"
}

main "$@"
