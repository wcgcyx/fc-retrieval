#!/bin/bash

./lotus daemon --genesis=devgen.car --bootstrap=false &
sleep 5
./lotus-miner run --nosync