#!/bin/bash

installPath=$1
tranPath=""

if [ "$installPath" != "" ]; then
    tranPath=$installPath
else
    tranPath=/usr/local/bin
fi

UNAME=$(uname)
ARCH=$(uname -m)

rmOldFiles() {
    if [ -f $tranPath/tran ]; then
        sudo rm -rf $tranPath/tran*
    fi
}

v=$(curl --silent "https://api.github.com/repos/abdfnx/tran/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

releases_api_url=https://github.com/scmn-dev/tran/releases/download

successInstall() {
    echo "🙏 Thanks for installing Tran! If this is your first time using the CLI, be sure to run `tran --help` first."
}

mainCheck() {
    echo "Installing tran version $v"
    name=""

    if [ "$UNAME" == "Linux" ]; then
        if [ $ARCH = "x86_64" ]; then
            name="tran_linux_${v}_amd64"
        elif [ $ARCH = "i686" ]; then
            name="tran_linux_${v}_386"
        elif [ $ARCH = "i386" ]; then
            name="tran_linux_${v}_386"
        elif [ $ARCH = "arm64" ]; then
            name="tran_linux_${v}_arm64"
        elif [ $ARCH = "arm" ]; then
            name="tran_linux_${v}_arm"
        fi

        tranURL=$releases_api_url/$v/$name.zip

        wget $tranURL
        sudo chmod 755 $name.zip
        unzip $name.zip
        rm $name.zip

        # tran
        sudo mv $name/bin/tran $tranPath

        rm -rf $name

    elif [ "$UNAME" == "Darwin" ]; then
        if [ $ARCH = "x86_64" ]; then
            name="tran_macos_${v}_amd64"
        elif [ $ARCH = "arm64" ]; then
            name="tran_macos_${v}_arm64"
        fi

        tranURL=$releases_api_url/$v/$name.zip

        wget $tranURL
        sudo chmod 755 $name.zip
        unzip $name.zip
        rm $name.zip

        # tran
        sudo mv $name/bin/tran $tranPath

        rm -rf $name

    elif [ "$UNAME" == "FreeBSD" ]; then
        if [ $ARCH = "x86_64" ]; then
            name="tran_freebsd_${v}_amd64"
        elif [ $ARCH = "i386" ]; then
            name="tran_freebsd_${v}_386"
        elif [ $ARCH = "i686" ]; then
            name="tran_freebsd_${v}_386"
        elif [ $ARCH = "arm64" ]; then
            name="tran_freebsd_${v}_arm64"
        elif [ $ARCH = "arm" ]; then
            name="tran_freebsd_${v}_arm"
        fi

        tranURL=$releases_api_url/$v/$name.zip

        wget $tranURL
        sudo chmod 755 $name.zip
        unzip $name.zip
        rm $name.zip

        # tran
        sudo mv $name/bin/tran $tranPath

        rm -rf $name
    fi

    # chmod
    sudo chmod 755 $tranPath/tran
}

rmOldFiles
mainCheck

if [ -x "$(command -v tran)" ]; then
    successInstall
else
    echo "Download failed 😔"
    echo "Please try again."
fi
