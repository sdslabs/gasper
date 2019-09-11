#!/bin/sh

project_dir=$(pwd)
if [ -f $project_dir/bin/fresh ]; then
    exit 0
fi
echo "  >  Installing fresh..."
mkdir -p bin
tmp_dir=$(mktemp -d -t ci-XXXXXXXXXX)
cd $tmp_dir
GOPATH=$tmp_dir go get github.com/pilu/fresh
cp $tmp_dir/bin/fresh $project_dir/bin/fresh
rm -rf $tmp_dir
