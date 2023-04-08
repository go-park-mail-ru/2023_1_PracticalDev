#!/bin/bash

interface_files=(
  internal/boards/service.go
  internal/boards/repository.go
  internal/images/service.go
  internal/images/repository.go
)

echo "Generating mocks..."
for file in ${interface_files[@]}; do
  out_file=$(dirname $file)
  out_file+="/mocks/"
  out_file+=$(basename $file)

  echo -e Generate $out_file
  mockgen -source=$file -destination=$out_file -package=mocks
done
echo "Mocks were generated."
