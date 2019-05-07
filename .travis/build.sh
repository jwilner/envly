#!/usr/bin/env bash

# a way over complicated build script because it was fun

set -e

# this cray https://stackoverflow.com/a/32641190/1567452
function _powerset {
    local items=("${@}")
    local n="${#items[@]}"
    local powersize=$((1 << "${n}"))

    local i=0
    while [[ $i -lt "${powersize}" ]] ; do
       local subset=()
       local j=0
       while [[ $j -lt $n ]]
       do
           if [[ $(((1 << $j) & $i)) -gt 0 ]]
           then
               subset+=("${items[$j]}")
           fi
           j=$(($j + 1))
       done
       echo "${subset[@]}"
       i=$(($i + 1))
    done
}

function _build {
    local tags="${@}"
    local target="target/envly-$(echo ${tags} | tr ' ' '-')-${GOOS}-amd64"

    go build \
        -o "${target}" \
        -tags "'${tags}'" || (
        echo "failed" >> "${FAIL_FILE}" && return 1
    )

    echo "built ${target}" >&2
}

function main {
    export GOARCH=amd64 CGO_ENABLED=0

    mkdir -p target

    local all=$(echo *_backend.go | sed 's!_backend\.go!!g')

    export FAIL_FILE=$(mktemp)
    trap "rm ${FAIL_FILE}" exit

    _powerset ${all} | while read tags; do
        if [[ -n "${tags}" ]]; then
            for goos in darwin linux; do
                GOOS="${goos}" _build "${tags}" &
            done
        fi
    done
    wait
}

main
