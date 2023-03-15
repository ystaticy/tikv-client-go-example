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
	"os"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func main4() {

	num := 1
	t7 := time.Now().Nanosecond()
	for num <= 200000 {
		os.Getenv("PATH")
		num++
	}
	t8 := time.Now().Nanosecond()
	t9 := (t8 - t7)
	fmt.Println("time key  :", zap.Int("t9", t9))

}

func main3() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.0.100:2379"},
		DialTimeout: 5 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "/", clientv3.WithPrefix())
	defer cancel()
	if err != nil {
		fmt.Println("get failed, err:", err)
		return
	}

	fmt.Printf("aaaaa")

	for _, ev := range resp.Kvs {
		fmt.Printf("get all (path) - %s : %s\n", ev.Key, ev.Value)
	}

	defer cancel()
}

func main2() {
	//cli, err := clientv3.New(clientv3.Config{
	//	Endpoints:   []string{"172.16.5.32:2379"},
	//	DialTimeout: 5 * time.Second,
	//})
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//resp, err := cli.Get(ctx, "/", clientv3.WithPrefix())
	//defer cancel()
	//if err != nil {
	//	fmt.Println("get failed, err:", err)
	//	return
	//}
	//
	//fmt.Printf("aaaaa")
	//
	//for _, ev := range resp.Kvs {
	//	fmt.Printf("get all (/sample_key/name/) - %s : %s\n", ev.Key, ev.Value)
	//}

	/**

	78000001748000FFFFFFFFFFFFFE5F7200FF0000000000000000FE
	78000001748000FFFFFFFFFFFFFE5F72FFFFFFFFFFFFFFFFFF00FF0000000000000000F7


	7800000374800000FF00000000445F7280FF000000000F3E5800FE
	*/
	//a := "78000001748000FFFFFFFFFFFFFE5F7200FF0000000000000000FE"
	////b := "78000001748000FFFFFFFFFFFFFE5F72FFFFFFFFFFFFFFFFFF00FF0000000000000000F7"
	//
	//c := "7800000374800000FF00000000445F7280FF000000000F3E5800FE"
	//d := strings.Compare(a, c)
	//
	//fmt.Printf("res:", d)

	//timeout := time.Second * 10
	//ctx, cancel := context.WithTimeout(context.Background(), timeout)
	//pdclient, err := pd.NewClientWithContext(ctx, []string{*pdAddr}, pd.SecurityOption{
	//	CAPath:   *ca,
	//	CertPath: *cert,
	//	KeyPath:  *key,
	//})
	//defer cancel()
	//if err != nil {
	//
	//}
	//
	//now := oracle.ComposeTS(p, l)
	//nowMinusOffset := oracle.GetTimeFromTS(now).Add(-*gcOffset)
	//newSP := oracle.ComposeTS(oracle.GetPhysical(nowMinusOffset), 0)
	//
	//pdclient.UpdateServiceSafePointV2(ctx, 1, "test_aaa", 100000, newSP)

	//meta, err := pdclient.LoadKeyspace(ctx, "ks1")
	//
	//fmt.Println("meta.Id  :", meta.Id)
	//
	//meta2, err := pdclient.LoadKeyspace(ctx, "ks2")
	//
	//fmt.Println("meta2.Id  :", meta2.Id)
	////regions, err := pdclient.ScanRegions(ctx, []byte{}, []byte{}, 100000)
	////
	////for i := range regions {
	////	startkey := regions[i].Meta.StartKey
	////	endkey := regions[i].Meta.EndKey
	////	if len(startkey) > 3 {
	////		prefix := startkey[0:4]
	////		out := Prefix(prefix)
	////		if out {
	////			fmt.Println("startkey  :", zap.Binary("startkey", startkey), zap.Binary("endkey", endkey))
	////		}
	////	}
	////}
	//
	//str := "780000026D44444CFF4A6F624869FF7374FF6F7279000000FC00FF0000000000006800FF00000000000003FFFF0000000000000000FFF700000000000000F8"
	//test, _ := hex.DecodeString(str)
	//fmt.Println("key  :", zap.Binary("test", test))
	//region, err := pdclient.GetRegion(ctx, test)
	//startkey := region.Meta.StartKey
	//engkey := region.Meta.EndKey
	//
	//start := hex.EncodeToString(startkey)
	//end := hex.EncodeToString(engkey)
	//fmt.Println("region key  :", zap.String("startkey", start), zap.String("endkey", end))

}

func Prefix(bytes []byte) bool {

	if bytes[0] == 120 && bytes[1] == 0 && bytes[2] == 0 && bytes[3] == 1 {
		return true
	}
	return false
}
