#!/bin/bash

containers=("scylla-1" "scylla-2" "scylla-3")

# Loop through each container
echo "Removing container:"
for container in "${containers[@]}"; do
    # Check if the container exists
    if docker ps -a --format '{{.Names}}' | grep -Eq "^${container}$"; then
        docker container rm -f $container
    fi
    echo
done
