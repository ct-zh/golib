#!/bin/bash
# server1 ip: 192.168.199.198
mkdir $(pwd)/redisVolume1
mkdir $(pwd)/redisVolume2
mkdir $(pwd)/redisVolume3

mv $(pwd)/sentinel1.conf $(pwd)/redisVolume1/sentinel.conf
mv $(pwd)/sentinel2.conf $(pwd)/redisVolume2/sentinel.conf
mv $(pwd)/sentinel3.conf $(pwd)/redisVolume3/sentinel.conf

docker pull redis:latest
docker-compose -f sentinel.yml up -d

