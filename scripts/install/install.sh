#!/bin/sh
set -e

reset="\033[0m"
cyan="\033[36m"
green="\033[32m"
red="\033[31m"
VERSION=$(curl -sS https://raw.githubusercontent.com/OnFinality-io/onf-cli/feature/networkV2/VERSION)

#Check OS
if [ "$(uname)" = "Darwin" ]; then
  LOCATION="/usr/local/bin"
  URL="https://github.com/OnFinality-io/onf-cli/releases/download/v$VERSION/onf-darwin-amd64-v$VERSION"
  SYSTEM="MACOS"
elif [ "$(uname -s)" = "Linux" ]; then
  URL="https://github.com/OnFinality-io/onf-cli/releases/download/v$VERSION/onf-linux-amd64-v$VERSION"
  LOCATION="/usr/local/bin"
  SYSTEM="LINUX"
elif [ "$(uname -s)" = "MINGW64_NT"  ]; then
  URL="https://github.com/OnFinality-io/onf-cli/releases/download/v$VERSION/onf-windows-amd64-v$VERSION.exe"
  LOCATION="%WINDIR%\system32"
  SYSTEM="WINDOWS"
fi

#Download binary
printf %s"$cyan> Downloading ...$reset\n"
cd "$LOCATION"
if curl -L "$URL" --output onf ;then
  if [ "$SYSTEM" != "WINDOWS" ];then
    chmod 773 onf
  fi
  printf %s"$green > onf command v$VERSION is ready. $reset\n"
else
  printf %s"$red> Failed to download $URL.$reset\n"
  exit 1;
fi
