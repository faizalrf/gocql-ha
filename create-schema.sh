#!/bin/bash

docker container exec -i scylla-1 cqlsh -f /tmp/schema.cql