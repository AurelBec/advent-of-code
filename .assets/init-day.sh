#!/bin/bash
# This script initialize a day directory, if not exists

# read date from arguments
IFS='/' read -ra date <<<$1

year=${date[0]}
day=${date[1]}
if [ -z "$day" ]; then
    day=$2
fi
day=$(printf "%02d" $day)
URL=https://adventofcode.com/$year/day/$(echo $day | sed 's/^0*//')
title=$(curl --silent --url $URL | sed -nr 's/.*--- (.*) ---.*/\1/p')

# ensure year and day are valid
if [ -z "$year" ] || [ -z "$day" ]; then
    echo "usage: '$0 \$year/\$day' or '$0 \$year \$day'"
    exit 1
fi

if [ -d "$year/$day" ]; then
    echo "$year/$day already exists, abort"
    exit 1
fi

# create directory and basic files
mkdir -p $year/$day
touch $year/$day/example.txt

# create go file with correct header and day title
cp $(dirname "$0")/template.go $year/$day/main.go
sed -i "s|https://adventofcode.com/|$URL|g" $year/$day/main.go
sed -i "s|--- ---|--- $year $title ---|g" $year/$day/main.go

echo "repository '$year/$day' ready ($title)"

# get input using token is present
if [ -z "$AOC_TOKEN" ]; then
    echo "skip input retrieval, token is not set"
    exit
fi

echo "retrieve input using token '$AOC_TOKEN'..."

response=$(
    curl \
        --cookie "session=$AOC_TOKEN" \
        --url $URL/input \
        --write-out "\n%{http_code}" \
        --silent
)
code=$(tail -n1 <<<"$response")
input=$(sed '$ d' <<<"$response")

if [ "$code" -ne 200 ]; then
    echo "received unexpected response: $code | $input"
    rm --recursive --force $year/$day
    exit 1
fi

echo "$input" >$year/$day/input.txt
echo "user input successully fetched"
