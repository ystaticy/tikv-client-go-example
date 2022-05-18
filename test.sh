

port=$1
ip_addr="172.16.5.32"
addr="$ip_addr:$port"


for i in {1..100}
do


./scan-demo --pd $addr
sleep 1s
tiup ctl:v6.0.0 tikv --pd $addr compact-cluster

sleep 5m

done