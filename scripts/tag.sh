#!/bin/bash

tag="${1}"

while (($#)); do
    case "$2" in

    -a)
        git tag -a "${tag}" -m ""
        exit 0
        ;;

    -p)
        git push origin "${tag}"
        exit 0
        ;;

    -x)
        git tag -a "${tag}" -m ""
        git push origin "${tag}"
        exit 0
        ;;

    -d)
        git tag -d "${tag}"
        git push --delete origin "${tag}"
        exit 0
        ;;

    esac
done
