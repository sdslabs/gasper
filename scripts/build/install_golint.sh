#!/bin/bash 

set -e
project_dir=$(pwd)
if [ -f $project_dir/bin/golint ]; then
    exit 0
fi

printf "🔨 Installing golint\n" 

mkdir -p bin
tmp_dir=$(mktemp -d -t ci-XXXXXXXXXX)
cd $tmp_dir
GOPATH=$tmp_dir go install golang.org/x/lint/golint@latest
chmod -R u+rw $tmp_dir
cp $tmp_dir/bin/golint $project_dir/bin/golint
rm -rf $tmp_dir

printf "👍 Done\n"
