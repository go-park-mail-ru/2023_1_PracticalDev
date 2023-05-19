#!/bin/bash

registry_id="crp0bgh5kbj19plq4i85"
targets=("auth" "images" "pickpin" "search" "shortener")

for target in "${targets[@]}"
do
	docker push cr.yandex/"$registry_id"/"$target"
done
