# raft-grpc-example

This is some example code for how to use [Hashicorp's Raft implementation](https://github.com/hashicorp/raft) with gRPC.

## Start your own cluster

bash start\_cluster.sh

You start up three nodes, and bootstrap one of them. Then you tell the bootstrapped node where to find peers. Those peers sync up to the state of the bootstrapped node and become members of the cluster. Once your cluster is running, you never need to pass `--raft_bootstrap` again.

[raftadmin](https://github.com/Jille/raftadmin) is used to communicate with the cluster and add the other nodes.

## Start the client

bash start\_client.sh

## Set fault injection and start a cluster with flaw

For crash based failure, we manually kill the process.

For CPU&Memory limitation failure, we run 'bash set\_cgroup.sh' to configure the cgroup. Then, we run 'bash start\_cluster\_fault.sh' to start a cluster with CPU or memory failure.
