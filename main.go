// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Test backup with exceeding GC safe point.

package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:42977"},
		DialTimeout: 5 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "/keyspaces/tidb/", clientv3.WithPrefix())
	defer cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}

	fmt.Printf("aaaaa")

	for _, ev := range resp.Kvs {
		fmt.Printf("get all (/sample_key/name/) - %s : %s\n", ev.Key, ev.Value)
	}

}
