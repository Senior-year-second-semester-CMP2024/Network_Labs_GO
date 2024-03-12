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
