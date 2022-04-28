#sudo mkdir /sys/fs/cgroup/cpu/assign2
#sudo bash -c "echo 200000 > /sys/fs/cgroup/cpu/assign2/cpu.cfs_quota_us"
#sudo bash -c "echo 1000000 > /sys/fs/cgroup/cpu/assign2/cpu.cfs_period_us"

#sudo mkdir /sys/fs/cgroup/memory/assign2
#sudo bash -c "echo 10M > /sys/fs/cgroup/memory/assign2/memory.limit_in_bytes"
