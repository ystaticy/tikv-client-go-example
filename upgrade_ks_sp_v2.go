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
	"math"
	"time"

	"github.com/pingcap/kvproto/pkg/keyspacepb"
	"github.com/pingcap/log"
	"github.com/tikv/client-go/v2/oracle"
	pd "github.com/tikv/pd/client"
	"go.uber.org/zap"
)

var (
	ca         = flag.String("ca", "", "CA certificate path for TLS connection")
	cert       = flag.String("cert", "", "certificate path for TLS connection")
	key        = flag.String("key", "", "private key path for TLS connection")
	pdAddr     = flag.String("pd", "127.0.0.1:43277", "PD address")
	opType     = flag.String("op", "", "optype")
	serviceID  = flag.String("serviceid", "test-service", "serviceid")
	keyspaceID = flag.Uint64("keyspaceid", math.MaxUint64, "serviceid")
)

func main() {
	flag.Parse()
	if *pdAddr == "" {
		log.Panic("pd address is empty")
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

	if *opType == "upgradeAll" {
		upgradeAllToGCV2(ctx, pdclient)
	} else if *opType == "upgradeKeyspace" {
		upgradeKeyspaceToGCV2(ctx, pdclient, uint32(*keyspaceID))
	} else if *opType == "updateserviceV1" {
		updateServiceV1(ctx, pdclient)
	} else if *opType == "updategcv1" {
		updateGCV1(ctx, pdclient)
	} else {
		log.Info("please use: -op xxx")
	}
}

func upgradeAllToGCV2(ctx context.Context, pdclient pd.Client) {
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
		if gcSafePointV2 == gcSafePointV1 {
			log.Info("[gc upgrade] update gc safe point v2 success.", zap.Uint32("KeyspaceID", keyspaceMeta.Id), zap.Uint64("gcSafePointV2", gcSafePointV2))
		} else {
			log.Error("[gc upgrade] update gc safe point v2 error, because safe point v2 is not newest.", zap.Uint32("KeyspaceID", keyspaceMeta.Id))
		}
	}
}

func upgradeKeyspaceToGCV2(ctx context.Context, pdclient pd.Client, keyspaceID uint32) {
	// Get gc safe point v1.
	gcSafePointV1, err := pdclient.UpdateGCSafePoint(ctx, 0)
	if err != nil {
		log.Panic("get gc safe point v1 from pd client failed", zap.Error(err))
	}
	log.Info("[gc upgrade] start gc upgrade. Get gc safe point v1 from pd client.", zap.Uint64("gcSafePointV1", gcSafePointV1), zap.Uint32("keyspaceID", keyspaceID))
	// update keyspace gc safe point v2.

	gcSafePointV2, err := pdclient.UpdateGCSafePointV2(ctx, keyspaceID, gcSafePointV1)
	if err != nil {
		log.Error("[gc upgrade] update gc safe point v2 error", zap.Uint32("KeyspaceID", keyspaceID), zap.Error(err))
	}

	serviceSafePointV2, err := pdclient.UpdateServiceSafePointV2(ctx, keyspaceID, "gcworker", 90, gcSafePointV1)
	if err != nil {
		log.Error("[gc upgrade] update gc safe point v2 error", zap.Uint32("KeyspaceID", keyspaceID), zap.Error(err))
	}
	if gcSafePointV2 == gcSafePointV1 && serviceSafePointV2 == gcSafePointV1 {
		log.Info("[gc upgrade] update gc safe point v2 success.", zap.Uint32("KeyspaceID", keyspaceID))
	} else {
		log.Error("[gc upgrade] update gc safe point v2 error, because safe point v2 is not newest.", zap.Uint32("KeyspaceID", keyspaceID), zap.Uint64("serviceSafePointV2", serviceSafePointV2), zap.Uint64("gcSafePointV2", gcSafePointV2))
	}
}

func updateServiceV1(ctx context.Context, pdclient pd.Client) {
	p, l, err := pdclient.GetTS(ctx)
	if err != nil {
		log.Panic("get ts failed", zap.Error(err))
	}
	now := oracle.ComposeTS(p, l)
	log.Info("get now ts to update service safe point v1.", zap.Uint64("now", now))
	// update all keyspace gc safe point v2.
	getServiceV1, err := pdclient.UpdateServiceGCSafePoint(ctx, *serviceID, 10, now)
	if err != nil {
		log.Panic("[gc upgrade] update service safe point v1 error", zap.String("serviceID", *serviceID), zap.Error(err))
	} else {
		log.Info("[gc upgrade] update service safe point v1 succuss", zap.Uint64("getServiceV1", getServiceV1))
	}
}

func updateGCV1(ctx context.Context, pdclient pd.Client) {
	// Get gc safe point v1.
	gcSafePointV1, err := pdclient.UpdateGCSafePoint(ctx, 0)
	if err != nil {
		log.Panic("get gc safe point v1 from pd client failed", zap.Error(err))
	}
	log.Info("get gc safe point v1 from pd client.", zap.Uint64("gcSafePointV1", gcSafePointV1))
}

func getAllKeyspace(pdclient pd.Client) []*keyspacepb.KeyspaceMeta {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watchChan, err := pdclient.WatchKeyspaces(ctx)
	if err != nil {
		log.Panic("WatchKeyspaces error")
	}
	initialLoaded := <-watchChan
	return initialLoaded
}
