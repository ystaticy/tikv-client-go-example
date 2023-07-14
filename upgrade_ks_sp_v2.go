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
	"flag"
	"time"

	"github.com/pingcap/kvproto/pkg/keyspacepb"
	"github.com/pingcap/log"
	pd "github.com/tikv/pd/client"
	"go.uber.org/zap"
)

var (
	ca       = flag.String("ca", "", "CA certificate path for TLS connection")
	cert     = flag.String("cert", "", "certificate path for TLS connection")
	key      = flag.String("key", "", "private key path for TLS connection")
	pdAddr   = flag.String("pd", "127.0.0.1:43277", "PD address")
	gcOffset = flag.Duration("gc-offset", time.Second*1,
		"Set GC safe point to current time - gc-offset, default: 10s")
	updateService = flag.Bool("update-service", false, "use new service to update min SafePoint")
)

func main() {
	flag.Parse()
	if *pdAddr == "" {
		log.Panic("pd address is empty")
	}
	if *gcOffset == time.Duration(0) {
		log.Panic("zero gc-offset is not allowed")
	}

	timeout := time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	defer cancel()
	pdclient, err := pd.NewClientWithContext(ctx, []string{*pdAddr}, pd.SecurityOption{
		CAPath:   *ca,
		CertPath: *cert,
		KeyPath:  *key,
	})
	if err != nil {
		log.Panic("create pd client failed", zap.Error(err))
	}

	// Get gc safe point v1.
	gcSafePointV1, err := pdclient.UpdateGCSafePoint(ctx, 0)
	if err != nil {
		log.Panic("get gc safe point v1 from pd client failed", zap.Error(err))
	}
	log.Info("get gc safe point v1 from pd client.", zap.Uint64("gcSafePointV1", gcSafePointV1))
	// update all keyspace gc safe point v2.

	// get all keyspace
	keyspaces := getAllKeyspace(pdclient)
	for i := range keyspaces {
		keyspaceMeta := keyspaces[i]
		if keyspaceMeta.State != keyspacepb.KeyspaceState_ENABLED {
			continue
		}
		log.Info("[gc upgrade] start gc upgrade", zap.Uint32("KeyspaceID", keyspaceMeta.Id))

		// ------ do pdclient.updateGCsafepointv2(ksid,safepoint)
		gcSafePointV2, err := pdclient.UpdateGCSafePointV2(ctx, keyspaceMeta.Id, gcSafePointV1)
		if err != nil {
			log.Error("[gc upgrade] update gc safe point v2 error", zap.Uint32("KeyspaceID", keyspaceMeta.Id), zap.Error(err))
		}
		if gcSafePointV2 != gcSafePointV1 {
			log.Error("[gc upgrade] update gc safe point v2 error, because safe point v2 is not newest.", zap.Uint32("KeyspaceID", keyspaceMeta.Id))
		}
	}

}

func getAllKeyspace(pdclient pd.Client) []*keyspacepb.KeyspaceMeta {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watchChan, err := pdclient.WatchKeyspaces(ctx)
	if err != nil {
		log.Error("WatchKeyspaces error")
	}
	initialLoaded := <-watchChan
	return initialLoaded
}
