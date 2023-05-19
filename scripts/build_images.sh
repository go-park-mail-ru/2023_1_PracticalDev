#!/bin/bash

registry_id="crp0bgh5kbj19plq4i85"
targets=("auth" "images" "pickpin" "search" "shortener")

for target in "${targets[@]}"
do
	DOCKER_BUILDKIT=1 docker build -f cmd/"$target"/Dockerfile -t cr.yandex/"$registry_id"/"$target" .
done
