# gocql Test

This little repository uses GO to connect to a Scylla cluster that is running on GCP/AWS. There are a few pre-requisites to running this code.

To start witn, provision 4 nodes on GCP/AWS within the same VPC and create and install ScyllaDB cluster. To make it simple, we can also use the ScyllaDB Ansible Role and Terraform to do this quickly. Refer to <https://github.com/faizalrf/scylla-terraform-ansible> for the ready made code to use on GCP for this.

Once the nodes are ready, connect to one of the nodes using `cqlsh` and create the schema from `schema.cql` 

Now, connect to the 4th node which does not have ScyllaDB installed, we will use this for the GO, do the following on this node.

- Install and configure `go` from <https://golang.org/dl/>
    - Download this gocql-ha repo from git and `cd gocql-ha` folder
    - Execute `go mod init gocql-ha` to init `go`
    - Execute `go mod tidy` to add module requirements and sums
    - Get the `gocql` drivers `go get github.com/gocql/gocql` 

Then build then execute the `main.go` Easiest and fastest way is to build the code by `go build` and it will generate the executable for you. Then simply execute that file to get the test going.

Now, once the code starting executing 

```
❯ go run main.go
Automatic SimpleRetry GOCQL Policy Test!
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


Whie the code executes and writes data to the table, stop one of the ScyllaDB nodes one by one until only one remains. The code will not fail and continue to execute without fail. Try different sequence of shutdown of the ScyllaDB nodes but ensure, at least one node is alive. The code will not fail and continue to execute as per normal.

```
.
.
Inserted Row Count: 1271
Inserted Row Count: 1272
Inserted Row Count: 1273
Inserted Row Count: 1274
Inserted Row Count: 1275
Inserted Row Count: 1276
Inserted Row Count: 1277
Inserted Row Count: 1278
.
.
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

Without the Simpe Retry policy, the app would have failed with something like the following error.

```
2024/06/20 02:43:06 gocql: unable to dial control conn 127.0.0.1:9003: EOF
```

This app uses `CL=LOCAL_ONE` to show the highest availability.

## Thanks!
