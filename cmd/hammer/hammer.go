// Binary hammer sends requests to your Raft cluster as fast as it can.
// It sends the written out version of the Dutch numbers up to 2000.
// In the end it asks the Raft cluster what the longest three words were.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	_ "github.com/Jille/grpc-multi-resolver"
	pb "github.com/Jille/raft-grpc-example/proto"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/health"
)

var n_operations int = 1000
var n_milliseconds int = 20000

func main() {
	serviceConfig := `{"healthCheckConfig": {"serviceName": "Example"}, "loadBalancingConfig": [ { "round_robin": {} } ]}`
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}
	conn, err := grpc.Dial("multi:///localhost:50051,localhost:50052,localhost:50053",
		grpc.WithDefaultServiceConfig(serviceConfig), grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))
	if err != nil {
		log.Fatalf("dialing failed: %v", err)
	}
	defer conn.Close()
	c := pb.NewExampleClient(conn)

	//ch := generateWords()
	workload, _ := strconv.ParseFloat(os.Args[1], 64)
	ch := genWordsRate(workload)

	var wg sync.WaitGroup

	n_goroutine := 200
	delays := make([][]int, n_goroutine)
	for i := 0; i < n_goroutine; i++ {
		delays[i] = make([]int, 1)
	}
	global_ts_begin := time.Now().UnixMilli()

	for i := 0; n_goroutine > i; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			word_id := 0
			for w := range ch {
				ts_enqueue := time.Now().UnixMicro()
				_, err := c.AddWord(context.Background(), &pb.AddWordRequest{Word: w})
				if err != nil {
					log.Fatalf("AddWord RPC failed: %v", err)
				}
				word_id += 1
				ts_dequeue := time.Now().UnixMicro()
				//ts_enqueue, _ := strconv.ParseInt(w, 10, 64)
				//fmt.Printf("latency: %d\n", ts_dequeue-ts_enqueue)
				delays[i] = append(delays[i], int(ts_dequeue-ts_enqueue))
				word_id = 0
			}
		}(i)
	}
	wg.Wait()

	global_ts_end := time.Now().UnixMilli()

	var delay_s []int
	for i := 0; i < n_goroutine; i++ {
		delay_s = append(delay_s, delays[i]...)
	}
	sort.Ints(delay_s)
	fmt.Printf("workload: %f all_requests: %d thorughput: %d p50_delay: %d\n",
		workload,
		len(delay_s),
		1000*len(delay_s)/int(global_ts_end-global_ts_begin),
		delay_s[len(delay_s)/2])
	//fmt.Printf("ewma_latency: %f\n", avg_latency)

	resp, err := c.GetWords(context.Background(), &pb.GetWordsRequest{})
	if err != nil {
		log.Fatalf("GetWords RPC failed: %v", err)
	}
	fmt.Println(resp)
}

func generateWords() <-chan string {
	ch := make(chan string, 1)
	go func() {
		//for i := 1; 2000 > i; i++ {
		for i := 1; n_operations > i; i++ {
			ts_enqueue := time.Now().UnixMicro()
			ch <- strconv.FormatInt(ts_enqueue, 10)
		}
		close(ch)
	}()
	return ch
}

// rate (rate kops/sec)
// generate reqs per 10 milliseconds
func genWordsRate(rate float64) <-chan string {
	ch := make(chan string, 10000)
	go func() {
		for j := 0; j < n_milliseconds/10; j++ {
			fmt.Printf("millisecond: %d\n", j)
			begin_ts := time.Now().UnixMicro()
			time.Sleep(time.Millisecond * 10)
			end_ts := time.Now().UnixMicro()
			n_requests := int64(rate * float64(end_ts-begin_ts) / 1000.0)
			for i := int64(0); i < n_requests; i++ {
				ts_enqueue := time.Now().UnixMicro()
				ch <- strconv.FormatInt(ts_enqueue, 10)
			}
		}
		close(ch)
	}()
	return ch
}
