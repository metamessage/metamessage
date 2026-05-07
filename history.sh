#!/bin/bash

VERSION="$1"
TARGET_HEADER="## ${VERSION}"
in_changelog=false

while IFS= read -r line; do
    if [[ "$line" == "# "* ]]; then
        continue
    fi

    if [[ $in_changelog == false ]] && [[ "$line" == "${TARGET_HEADER}"* ]]; then
        in_changelog=true
        continue
    fi

    if [[ $in_changelog == true ]] && [[ "$line" == "## "* ]]; then
        break
    fi

    echo "$line"
done < HISTORY.md
