#!/bin/bash 

set -e
project_dir=$(pwd)
if [ -f $project_dir/bin/fresh ]; then
    exit 0
fi

printf "🔨 Installing fresh\n" 

mkdir -p bin
tmp_dir=$(mktemp -d -t ci-XXXXXXXXXX)
cd $tmp_dir
GOPATH=$tmp_dir go install github.com/pilu/fresh@latest
chmod -R u+rw $tmp_dir
cp $tmp_dir/bin/fresh $project_dir/bin/fresh
rm -rf $tmp_dir

printf "👍 Done\n"
