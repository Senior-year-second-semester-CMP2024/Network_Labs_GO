#!/bin/bash

# Number of instances to run
NUM_INSTANCES=3

# Build the Go program
go build -o ./data_keeper/data_keeper.bin ./data_keeper

# Run instances
for ((i=0; i<NUM_INSTANCES; i++)); do
    # Calculate base port for this instance
    BASE_PORT=$((6000 + (i * 100)))

    # Run instance with calculated base port
    ./data_keeper/data_keeper.bin -name="node$i" -ports="$BASE_PORT $((BASE_PORT+1)) $((BASE_PORT+2))" -master_port=8080 &
done
# to run one node: 
# ./data_keeper/data_keeper.bin -name="node0" -ports="6000 6001 6002" -master_port=8080
# ./data_keeper/data_keeper.bin -name="node1" -ports="6100 6101 6102" -master_port=8080
# ./data_keeper/data_keeper.bin -name="node2" -ports="6200 6201 6202" -master_port=8080
# ./data_keeper/data_keeper.bin -name="node3" -ports="6300 6301 6302" -master_port=8080
# ./data_keeper/data_keeper.bin -name="node4" -ports="6400 6401 6402" -master_port=8080
