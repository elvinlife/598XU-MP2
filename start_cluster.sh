#!/bin/bash
sudo rm -rf /tmp/my-raft-cluster
mkdir /tmp/my-raft-cluster
mkdir /tmp/my-raft-cluster/nodeA
mkdir /tmp/my-raft-cluster/nodeB
mkdir /tmp/my-raft-cluster/nodeC

sudo killall raft-grpc-example
./raft-grpc-example --raft_bootstrap --raft_id=nodeA --address=localhost:50051 --raft_data_dir /tmp/my-raft-cluster &> nodeA.log &
sudo taskset -cp 0 $!
./raft-grpc-example --raft_id=nodeB --address=localhost:50052 --raft_data_dir /tmp/my-raft-cluster &> nodeB.log &
sudo taskset -cp 1 $!
./raft-grpc-example --raft_id=nodeC --address=localhost:50053 --raft_data_dir /tmp/my-raft-cluster &> nodeC.log &
sudo taskset -cp 2 $!

sleep 3
../raftadmin/cmd/raftadmin/raftadmin localhost:50051 add_voter nodeB localhost:50052 0 &
sleep 1
../raftadmin/cmd/raftadmin/raftadmin --leader multi:///localhost:50051,localhost:50052 add_voter nodeC localhost:50053 0 &
