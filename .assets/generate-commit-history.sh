#!/bin/bash
# This script regenerate a complete git history from the eldest year found

set_commit_date() {
    GIT_COMMITTER_DATE="$1" GIT_AUTHOR_DATE="$1" git commit --amend --no-edit --date "$1"
}

is_number='^[0-9]+$'

# count number of changes and additions
# changed=$(git ls-files . --modified)
# added=$(git ls-files . --exclude-standard --others)

# reset soft to the root
git update-ref -d HEAD
git reset

# get years contained
init=()
utils=()
years=()
for entry in $(ls --almost-all .); do
    if [[ $entry =~ $is_number ]]; then
        years+=($entry)
    elif [[ $entry = "utils" ]]; then
        utils+=($entry)
    else
        init+=($entry)
    fi
done

# commit first the init content the 01/01/years[0]
for init in "${init[@]}"; do
    git add $init
done
git commit --message "Initial commit"
set_commit_date "01 Jan ${years[0]} 00:00:00"

# commit then the utility stuff, also the 01/01/years[0]
for utils in "${utils[@]}"; do
    git add $utils
done
git commit --message "Add utility stuff"
set_commit_date "01 Jan ${years[0]} 00:00:00"

# then for each year, and each day of december, commit to the correct date
for year in "${years[@]}"; do
    for day in $(ls --almost-all $year); do
        if [[ $day =~ $is_number ]]; then
            git add $year/$day
            git commit --message "${year} day ${day}"
            set_commit_date "${day} Dec ${year} 00:00:00"
        fi
    done
done

# finish by pushing branch if asked
if [[ $1 == "--push" ]]; then
    git push --force
else
    echo "skipping branch push (use --push)"
fi
