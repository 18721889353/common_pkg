## etcdcli

Connect to the etcd service client.

### Example of use

```go
    import "github.com/18721889353/common_pkg/etcdcli"

	endpoints := []string{"192.168.3.37:2379"}
    cli, err := etcdcli.Init(endpoints,
        WithConnectTimeout(time.Second*2),
        // WithAuth("", ""),
        // WithAutoSyncInterval(0),
        // WithLog(zap.NewNop()),
	)
)
```
