#!/usr/bin/env bash
#=================================================
# https://github.com/P3TERX/ca-certificates.crt
# Description: Update ca-certificates.crt
# Version: 1.1
# Lisence: MIT
# Author: P3TERX
# Blog: https://p3terx.com
#=================================================

function __wget() {
    : ${DEBUG:=0}
    local URL=$1
    local tag="Connection: close"
    local mark=0

    if [ -z "${URL}" ]; then
        printf "Usage: %s \"URL\" [e.g.: %s http://www.google.com/]" \
               "${FUNCNAME[0]}" "${FUNCNAME[0]}"
        return 1;
    fi
    read proto server path <<<$(echo ${URL//// })
    DOC=/${path// //}
    HOST=${server//:*}
    PORT=${server//*:}
    [[ x"${HOST}" == x"${PORT}" ]] && PORT=80
    [[ $DEBUG -eq 1 ]] && echo "HOST=$HOST"
    [[ $DEBUG -eq 1 ]] && echo "PORT=$PORT"
    [[ $DEBUG -eq 1 ]] && echo "DOC =$DOC"

    exec 3<>/dev/tcp/${HOST}/$PORT
    echo -en "GET ${DOC} HTTP/1.1\r\nHost: ${HOST}\r\n${tag}\r\n\r\n" >&3
    while read line; do
        [[ $mark -eq 1 ]] && echo $line
        if [[ "${line}" =~ "${tag}" ]]; then
            mark=1
        fi
    done <&3
    exec 3>&-
}

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

echo && echo -e "$INFO Check downloader ..."
if [ $(command -v wget) ]; then
    DOWNLOADER='wget -nv -N'
elif [ $(command -v curl) ]; then
    DOWNLOADER='curl -fsSLO'
else
    echo -e "$ERROR curl or wget is not installed."
    DOWNLOADER='__wget'
fi

echo -e "${INFO} Download ca-certificates.crt ..."
${DOWNLOADER} https://raw.githubusercontent.com/P3TERX/ca-certificates.crt/download/ca-certificates.crt ||
    ${DOWNLOADER} https://cdn.jsdelivr.net/gh/P3TERX/ca-certificates.crt@download/ca-certificates.crt ||
    ${DOWNLOADER} https://gh.p3terx.workers.dev/ca-certificates.crt/download/ca-certificates.crt

[ -s ca-certificates.crt ] && echo -e "${INFO} ca-certificates.crt Download completed !" || {
    echo -e "${ERROR} Unable to download ca-certificates.crt, network failure or API error."
    rm -f ca-certificates.crt
    exit 1
}

exit 0
