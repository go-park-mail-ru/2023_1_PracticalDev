#!/bin/bash

interface_files=(
  internal/auth/service.go
  internal/auth/repository.go
  internal/boards/service.go
  internal/boards/repository.go
  internal/pins/service.go
  internal/pins/repository.go
  internal/likes/service.go
  internal/likes/repository.go
  internal/profile/service.go
  internal/profile/repository.go
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
