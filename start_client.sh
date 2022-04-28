#!/bin/bash
# the request generation rate
rateArray=("1" "4" "7" "9" "10" "11")
for rate in ${rateArray[@]}; do
  go run cmd/hammer/hammer.go $rate > hammer${rate}.log
  taskset -cp 3 $!
done
