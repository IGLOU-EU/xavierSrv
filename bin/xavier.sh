#!/bin/bash
trap "exit 0" 2 3
time=60

cerebro () {
    outS=""

    while read line ; do
        [ "${line:0:1}" == "#" ] && continue
        [ "${line:0:1}" == "%" ] && cmdR+=("${line:1}") && continue

        url=${line#*:}
        status=${line%%:*}

        rq=$(curl -I -s "$url" -o /dev/null -w "%{http_code}\n")
        [ $rq -eq $status ] || outS="${outS}${rq} ${url#*//}\n"
    done < "${urls}"

    if [ -z "$outS" ]; then
        outC="OK!\n$(date +"%d/%m/%Y %T")"
    else
        if [ "$outS" != "$outC" ]; then
            xPsy
        fi

        outC="$outS"
    fi

    echo -e "$outC"
    echo -e "$outC" > "$html"
}

xPsy () {
    local _buff=""

    for ((i = 0; i < ${#cmdR[@]}; i++))
    do
        _buff="${cmdR[$1]}"
        _buff="${_buff/~outS~/$outS}"
        _buff="$(eval  "$_buff" 2>&1)"

        echo "$_buff"
    done
}

root="$(cd "$(dirname "$0")"; pwd)"
html="${1:-${root}/status.html}"
urls="${2:-${root}/url.list}"
outS=""
outC=""
cmdR=()

if [ "$1" == "-h" ]; then
    echo -e "usage : $(basename "$0") [<html output path> [<url list path>]]\n"
    exit 0
fi

if [ ! -w "${html%/*}" ]; then
    echo "Unable to write ${html}"
    exit 1
fi

if [ ! -r "${urls}" ]; then
    echo "Unable to read ${urls}"
    exit 1
fi

while true; do
    cerebro
    sleep $time
done
