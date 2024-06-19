# gocql Test

This little repository uses GO to connect to a Scylla cluster that is running locally using doker containers. There are a few pre-requisites to running this code.

- Install and configure `go` from <https://golang.org/dl/>
    - Download this gocql-ha repo from git and `cd gocql-ha` folder
    - Execute `go mod init gocql-ha` to init `go`
    - Execute `go mod tidy` to add module requirements and sums
    - Get the `gocql` drivers `go get github.com/gocql/gocql` 
- Install docker

Once ready with the above two, execute the scripts in this order

```
./create-containers.sh
./creeate-schema.sh
```

Then build then execute the `main.go` Easiest and fastest way is to build the code by `go build` and it will generate the executable for you. Then simply execute that file to get the test going.

Now, once the code starting executing 

```
❯ go build
❯ ./gocql-ha
Test code for gocql!
Connected to ScyllaDB!
Inserted Row Count: 1
Inserted Row Count: 2
Inserted Row Count: 3
Inserted Row Count: 4
Inserted Row Count: 5
Inserted Row Count: 6
Inserted Row Count: 7
Inserted Row Count: 8
Inserted Row Count: 9
Inserted Row Count: 10
.
.
.
```


Whie running stop the containersuntil only one ScyllaDB node remains. You may facew thee error sooner but app will reconnect automatically and continue writing from where it left off.

```
Inserted Row Count: 1271
Inserted Row Count: 1272
Inserted Row Count: 1273
2024/06/20 02:43:06 gocql: unable to dial control conn 127.0.0.1:9003: EOF
Error inserting data, attempt 1: gocql: no hosts available in the pool
2024/06/20 02:43:06 Attempting to reconnect to the ScyllaDB...
2024/06/20 02:43:09 error: failed to connect to "[HostInfo hostname=\"127.0.0.1\" connectAddress=\"127.0.0.1\" peer=\"<nil>\" rpc_address=\"172.17.0.3\" broadcast_address=\"172.17.0.3\" preferred_ip=\"<nil>\" connect_addr=\"127.0.0.1\" connect_addr_source=\"connect_address\" port=9003 data_centre=\"datacenter1\" rack=\"rack1\" host_id=\"83082671-05cc-4803-ba73-dd7705c5740c\" version=\"v3.0.8\" state=UP num_tokens=256]" due to error: EOF
2024/06/20 02:43:16 gocql: unable to dial control conn 172.17.0.4:9042: dial tcp 172.17.0.4:9042: operation was canceled
2024/06/20 02:43:16 gocql: unable to dial control conn 172.17.0.2:9042: dial tcp 172.17.0.2:9042: operation was canceled
2024/06/20 02:43:16 gocql: unable to connect to any ring node: dial tcp 172.17.0.2:9042: operation was canceled
2024/06/20 02:43:16 gocql: control falling back to initial contact points.
2024/06/20 02:43:16 gocql: unable to dial control conn 127.0.0.1:9002: dial tcp 127.0.0.1:9002: operation was canceled
2024/06/20 02:43:16 gocql: unable to dial control conn 127.0.0.1:9003: dial tcp 127.0.0.1:9003: operation was canceled
2024/06/20 02:43:16 gocql: unable to dial control conn 127.0.0.1:9004: dial tcp 127.0.0.1:9004: operation was canceled
2024/06/20 02:43:16 gocql: unable to reconnect control connection: dial tcp 127.0.0.1:9004: operation was canceled
2024/06/20 02:43:16 gocql: unable to dial control conn 127.0.0.1:9002: dial tcp 127.0.0.1:9002: connect: connection refused
2024/06/20 02:43:16 gocql: unable to dial control conn 127.0.0.1:9003: dial tcp 127.0.0.1:9003: connect: connection refused
Inserted Row Count: 1274
Inserted Row Count: 1275
Inserted Row Count: 1276
Inserted Row Count: 1277
Inserted Row Count: 1278
```

At this point only 1 node was running

```
❯ docker container exec -it scylla-3 nodetool status
Datacenter: datacenter1
=======================
Status=Up/Down
|/ State=Normal/Leaving/Joining/Moving
-- Address    Load      Tokens Owns Host ID                              Rack
DN 172.17.0.2 391.55 KB 256    ?    d90d3103-e0c4-4503-a609-7364fd6eea2c rack1
DN 172.17.0.3 379.70 KB 256    ?    83082671-05cc-4803-ba73-dd7705c5740c rack1
UN 172.17.0.4 409.22 KB 256    ?    fdf89c1b-3046-4839-971c-b147e28bad4f rack1
```

Without the retry bit in the code, the app would have failed at the line 

```
2024/06/20 02:43:06 gocql: unable to dial control conn 127.0.0.1:9003: EOF
```

This app uses `CL=LOCAL_ONE` to show the highest availability.

## Thanks!
