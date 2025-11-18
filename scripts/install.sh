#!/bin/bash

set -e

if [ -z "${BIN_DIR}" ]; then
	BIN_DIR=$(pwd)
fi

THE_ARCH_VER="v0.1.0-alpha.3"
THE_ARCH_BIN=""
DEST=${BIN_DIR}/dnocs-${THE_ARCH_VER}.tar.gz

OS=$(uname -s)
ARCH=$(uname -m)

if type "tput" >/dev/null 2>&1; then
	bold=$(tput bold || true)
	italic=$(tput sitm || true)
	normal=$(tput sgr0 || true)
fi

case ${OS} in
Linux*)
	case ${ARCH} in
	aarch64)
		THE_ARCH_BIN="dnocs-${THE_ARCH_VER}-linux-arm64.tar.gz"
		;;
	x86_64)
		THE_ARCH_BIN="dnocs-${THE_ARCH_VER}-linux-amd64.tar.gz"
		;;
	*)
		THE_ARCH_BIN=""
		;;
	esac
	;;
# Darwin*)
# 	case ${ARCH} in
# 	arm64)
# 		THE_ARCH_BIN="marijan-${THE_ARCH_VER}-darwin-arm64.tar.gz"
# 		;;
# 	*)
# 		THE_ARCH_BIN="marijan-${THE_ARCH_VER}-darwin-amd64.tar.gz"
# 		;;
# 	esac
# 	;;
# Windows | MINGW64_NT*)
# 	echo "❗ Use WSL to run Marijan on Windows: https://learn.microsoft.com/windows/wsl/"
# 	exit 1
# 	;;
*)
	THE_ARCH_BIN=""
	;;
esac

if [ -z "${THE_ARCH_BIN}" ]; then
	echo "❗ DNOCS is not supported on ${OS} and ${ARCH}"
	exit 1
fi

SUDO=""

echo "📦 Downloading ${bold}DNOCS${normal} for ${OS} (${ARCH}):"

# check if $DEST is writable and suppress an error message
touch "${DEST}" 2>/dev/null

# we need sudo powers to write to DEST
if [ $? -eq 1 ]; then
	echo "❗ You do not have permission to write to ${italic}${DEST}${normal}, enter your password to grant sudo powers"
	SUDO="sudo"
fi

if type "curl" >/dev/null 2>&1; then
	curl -L --progress-bar "https://github.com/devetek/d-panel-cli/releases/download/${THE_ARCH_VER}/${THE_ARCH_BIN}" -o "${DEST}"
elif type "wget" >/dev/null 2>&1; then
	${SUDO} wget "https://github.com/devetek/d-panel-cli/releases/download/${THE_ARCH_VER}/${THE_ARCH_BIN}" -O "${DEST}"
else
	echo "❗ Please install ${italic}curl${normal} or ${italic}wget${normal} to download DNOCS"
	exit 1
fi

# extract the tar.gz file, and delete the tar.gz file
tar -xzvf "${DEST}"
rm "${DEST}"

echo
echo "🥳 DNOCS downloaded successfully to ${italic}${DEST}${normal}"
echo "🔧 Move the binary to ${italic}/usr/local/bin/${normal} or another directory in your ${italic}PATH${normal} to use it globally:"
echo "   ${bold}sudo mv ${DEST} /usr/local/bin/${normal}"
echo
echo "⭐ If you like DNOCS, please give it a star on GitHub: ${italic}https://github.com/devetek/d-panel-cli${normal}"