#!/usr/bin/env bash
set -euo pipefail

# Function to check if an executable exists and is executable
check_executable() {
    local exe_path=$1
    local exe_name=$(basename "$exe_path")

    if [[ ! -x "$exe_path" ]]; then
        echo "Error: Executable file does not exist or is not executable [$exe_name]"
        return 1
    fi

    # Try to get version information with a timeout protection
    if ! timeout 5s "$exe_path" --version &>/dev/null; then
        echo "Error: Unable to get version information [$exe_name]"
        return 1
    fi

    echo "Verification passed: $exe_name $("$exe_path" --version | head -n1)"
    return 0
}

# Function to validate the toolchain
validate_toolchain() {
    local target_dir=$1
    local cc cxx ld

    # Build full paths
    cc="${target_dir}/${cc_path}"
    cxx="${target_dir}/${cxx_path}"
    ld="${target_dir}/${ld_path}"

    # Parallel verification of the three executable files
    local check_results=0
    check_executable "$cc" || check_results=$((check_results | 1))
    check_executable "$cxx" || check_results=$((check_results | 2))
    check_executable "$ld" || check_results=$((check_results | 4))

    if [[ $check_results -ne 0 ]]; then
        echo "Key component verification failed, error code: $check_results"
        return 1
    fi
    return 0
}

# Function to generate the toolchain configuration file
gen_toolchain_path() {
    local target_dir=$1
    local cc cxx ld

    # Build full paths
    cc="${target_dir}/${cc_path}"
    cxx="${target_dir}/${cxx_path}"
    ld="${target_dir}/${ld_path}"

    # Create toolchain.ini file
    cat <<EOF >../toolchains/toolchain.ini
[toolchain]
cc = $cc
cxx = $cxx
ld = $ld
EOF

    return 0
}

# Function to generate the Conan profile file
gen_conan_profile() {
    local target_dir=$1
    local cc cxx ld

    # Build full paths
    cc="${target_dir}/${cc_path}"
    cxx="${target_dir}/${cxx_path}"
    ld="${target_dir}/${ld_path}"

    cat <<EOF >./NanoKVM-Pro
[settings]
os=Linux
arch=armv8
compiler=gcc
build_type=Release
compiler.cppstd=gnu23
compiler.libcxx=libstdc++
compiler.version=11
[buildenv]
CC=$cc
CXX=$cxx
LD=$ld
EOF
}

# Function to check if an installation exists
is_installed() {
    local target_dir=$1
    [[ -d "${target_dir}/bin" ]] && return 0 || return 1
}

# Function to prompt the user for reinstallation options
prompt_reinstall() {
    local target_dir=$1
    echo "Detected installed version: $(basename "$target_dir")"
    PS3='Please select an option: '
    select opt in "Skip installation" "Reinstall" "Exit"; do
        case $opt in
        "Skip installation")
            echo "Skipped installation"
            exit 0
            ;;
        "Reinstall")
            echo "Preparing to reinstall..."
            rm -rf "$target_dir" # Delete old version
            mkdir -p "$target_dir"
            return 0
            ;;
        "Exit")
            echo "Operation cancelled"
            exit 1
            ;;
        *)
            echo "Invalid option, please select again"
            ;;
        esac
    done
}

# Main installation process
main() {
    local config_file="./config.ini"

    # Parse configuration (unchanged)
    local section="toolchain"
    local name url sha256
    name=$(./getconfig.py "$config_file" "$section" "name")
    url=$(./getconfig.py "$config_file" "$section" "url")
    sha256=$(./getconfig.py "$config_file" "$section" "sha256")

    local cc_path cxx_path ld_path
    cc_path=$(./getconfig.py "$config_file" "$section" "cc")
    cxx_path=$(./getconfig.py "$config_file" "$section" "cxx")
    ld_path=$(./getconfig.py "$config_file" "$section" "ld")

    # Validate key fields (unchanged)
    if [[ -z "$name" || -z "$url" || -z "$sha256" ]]; then
        echo "Error: Missing necessary fields in configuration file"
        echo "Parsed results:"
        echo "Name: $name"
        echo "URL: $url"
        echo "SHA256: $sha256"
        exit 1
    fi
    if [[ -z "$cc_path" || -z "$cxx_path" || -z "$ld_path" ]]; then
        echo "Error: Missing compiler path configuration"
        exit 1
    fi

    # Prepare installation directory (added installation detection)
    local script_dir target_dir
    script_dir=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd -P)
    target_dir="${script_dir}/../toolchains/${name}"

    # Added installation detection and user interaction
    if is_installed "$target_dir"; then
        prompt_reinstall "$target_dir"
    else
        mkdir -p "$target_dir"
    fi

    # Download file (supports resuming download, unchanged)
    local temp_file
    temp_file=$(mktemp "${TMPDIR:-/tmp}/toolchain.XXXXXX")
    echo "Downloading toolchain: ${url}"
    if ! curl -#L -o "$temp_file" "$url"; then
        echo "Download failed, please check:"
        echo "1. URL validity"
        echo "2. Network connection"
        echo "3. Disk space"
        rm -f "$temp_file"
        exit 1
    fi

    # Hash verification (unchanged)
    echo "Verifying file integrity..."
    local computed_sha256
    computed_sha256=$(sha256sum "$temp_file" | awk '{print $1}')
    if [[ "$computed_sha256" != "$sha256" ]]; then
        echo "Security verification failed!"
        echo "Expected value: $sha256"
        echo "Actual value: $computed_sha256"
        rm -f "$temp_file"
        exit 1
    fi

    # Extract installation (unchanged)
    echo "Installing to: ${target_dir}"
    if ! tar -xJf "$temp_file" -C "$target_dir" --strip-components=1; then
        echo "Extraction failed, possible reasons:"
        echo "1. File corruption (please re-download)"
        echo "2. Insufficient disk space"
        echo "3. Permission issues"
        rm -f "$temp_file"
        exit 1
    fi
    rm -f "$temp_file"

    if ! validate_toolchain "$target_dir"; then
        echo "Detected damaged installation, triggering reinstallation..."
        rm -rf "$target_dir"
        # Re-execute installation process
        main "$@"
        return
    fi

    if ! gen_toolchain_path "$target_dir"; then
        return
    fi

    if ! gen_conan_profile "$target_dir"; then
        return
    fi
}

main
