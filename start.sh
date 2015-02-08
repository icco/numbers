#! /bin/sh

sudo aptitude update
sudo aptitude install -y golang git
git clone https://github.com/icco/numbers
cd numbers; go install; sudo ./numbers -p 80
