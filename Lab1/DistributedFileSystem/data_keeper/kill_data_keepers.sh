#!/bin/bash

# Kill instances
for PID in $(pgrep -f "data_keeper.bin"); do
    kill $PID
done