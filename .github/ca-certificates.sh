#!/usr/bin/env bash

set -e
[ $(uname) != Linux ] && {
    echo "This operating system is not supported."
    exit 1
}

Green_font_prefix="\033[32m"
Red_font_prefix="\033[31m"
Green_background_prefix="\033[42;37m"
Red_background_prefix="\033[41;37m"
Font_color_suffix="\033[0m"
INFO="[${Green_font_prefix}INFO${Font_color_suffix}]"
ERROR="[${Red_font_prefix}ERROR${Font_color_suffix}]"

echo -e "${INFO} Updating ca-certificates.crt ..."
[ $EUID != 0 ] && {
    SUDO=sudo
    echo -e "${INFO} You may need to enter a password to authorize."
}
$SUDO mkdir -vp /etc/ssl/certs
$SUDO mv -vf ca-certificates.crt /etc/ssl/certs/ca-certificates.crt &&
    echo -e "${INFO} ca-certificates.crt update completed !" ||
    {
        echo -e "${ERROR} ca-certificates.crt update failed !"
        exit 1
    }
exit 0
