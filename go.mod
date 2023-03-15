module github.com/ystaticy/tikv-client-go-example

go 1.17

require (
	github.com/pingcap/log v1.1.1-0.20221110025148-ca232912c9f3
	github.com/tikv/client-go/v2 v2.0.1
	github.com/tikv/pd/client v0.0.0-20230209034200-6d23a31c24be
	go.etcd.io/etcd/client/v3 v3.5.2
	go.uber.org/zap v1.24.0
)

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pingcap/errors v0.11.5-0.20211224045212-9687c2b0f87c // indirect
	github.com/pingcap/failpoint v0.0.0-20220801062533-2eaa32854a6c // indirect
	github.com/pingcap/kvproto v0.0.0-20230206112125-0561adc37543 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.39.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	go.etcd.io/etcd/api/v3 v3.5.2 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.2 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/net v0.5.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20230202175211-008b39050e57 // indirect
	google.golang.org/grpc v1.52.3 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect

)

// Use github.com/tidbcloud/pd-cse/client
replace github.com/tikv/pd/client => github.com/tidbcloud/pd-cse/client v0.0.0-20230314083701-071f222e87f4

replace github.com/tikv/client-go/v2 => github.com/ystaticy/client-go/v2 v2.0.5-0.20230314084751-1cd49730c1e0

replace github.com/pingcap/kvproto => github.com/tidbcloud/kvproto v0.0.0-20230314071356-ede8ef250f4f
