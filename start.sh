#! /bin/bash

export DEBIAN_FRONTEND=noninteractive

aptitude update
aptitude install -y golang git
git clone https://github.com/icco/numbers
cd numbers; go install; ./numbers -p 80
