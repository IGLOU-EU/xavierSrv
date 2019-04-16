#!/bin/sh
trap "exit 0" 2 3
time=60

cerebro () {
    while read line ; do
        [ "${line:0:1}" == "#" ] && continue

        url=${line#*:}
        status=${line%%:*}

        rq=$(curl -I -s "$url" -o /dev/null -w "%{http_code}\n")
        [ $rq -eq $status ] || outS="${outS}${rq} ${url#*//}\n"
    done < "${urls}"

    if [ -z "$outS" ]; then
        out="OK!\n$(date +"%d/%m/%Y %T")"
    else
        out="$outS"
    fi

    echo -e "$out"
    echo -e "$out" > "$html"
}

root="$(cd "$(dirname "$0")"; pwd)"
html="${1:-${root}/status.html}"
urls="${2:-${root}/url.list}"
outS=""

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
