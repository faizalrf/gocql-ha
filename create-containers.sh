#!/bin/bash

# Pull the new image if available
echo "Pulling the latest ScyllaDB Image from scylladb/scylla:latest"
docker pull scylladb/scylla > /dev/null 2>&1

commands=(
    'scylla-1:docker run --name scylla-1 -d -p 9002:9042 scylladb/scylla:latest --reactor-backend=epoll --overprovisioned 1 --smp 1'
    'scylla-2:docker run --name scylla-2 -d -p 9003:9042 scylladb/scylla:latest --reactor-backend=epoll --seeds="$(docker inspect --format="{{.NetworkSettings.IPAddress}}" scylla-1)" --overprovisioned 1 --smp 1'
    'scylla-3:docker run --name scylla-3 -d -p 9004:9042 scylladb/scylla:latest --reactor-backend=epoll --seeds="$(docker inspect --format="{{.NetworkSettings.IPAddress}}" scylla-1)" --overprovisioned 1 --smp 1'
)

# Function to check nodetool status
check_nodetool_status() {
    local container=$1
    retries=30
    sleep_time=3

    for ((i=1; i<=retries; i++)); do
        echo
        echo "Checking nodetool status for container: \"$container\" (Attempt $i)"
        docker exec -it $container nodetool status > /dev/null 2>&1
        if [ $? -eq 0 ]; then
            echo "Container \"$container\" is running successfully."
            docker exec -it $container nodetool status
            echo
            return 0
        else
            echo "nodetool status check failed for \"$container\". Retrying in $sleep_time seconds..."
            sleep $sleep_time
        fi
    done
    echo "Failed to start container \"$container\" after $retries attempts."
    echo
    return 1
}

for cmd in "${commands[@]}"; do
    # Split the string into container name and command
    container=$(echo $cmd | cut -d ':' -f 1)
    command=$(echo $cmd | cut -d ':' -f 2-)
    echo "Starting container: $container"
    eval $command > /dev/null 2>&1
    ret=$?
    sleep 10
    if [ $ret -eq 0 ]; then
        echo "Container \"$container\" started successfully, waiting for ScyllaDB service to start"
        if ! check_nodetool_status $container; then
            echo "Failed to verify ScyllaDB \"$container\" status. Exiting."
            exit 1
        fi
    else
        echo "Failed to start container \"$container\". Exiting."
        exit 1
    fi
    docker cp schema.cql $container:/tmp/schema.cql    
done

echo "All containers started successfully!"
