package nacoscli

import (
	"context"
	"shop/sun_utils/pkg/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	ipAddr      = "192.168.132.128"
	port        = 8848
	namespaceID = "a3143e55-cd06-4379-9875-4a106bdd527f"
)

func TestParse(t *testing.T) {
	conf := new(map[string]interface{})
	params := &Params{
		Scheme:      "http",
		IPAddr:      ipAddr,
		Port:        uint64(port),
		NamespaceID: namespaceID,
		Group:       "dev",
		DataID:      "platform_admin_user_srv.json",
		Format:      "json",
		ContextPath: "/nacos",
	}
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		err := Init(conf, params)
		t.Log(err, conf)
	})

	for {
		select {}
	}
	//conf = new(map[string]interface{})
	//params = &Params{
	//	Group:  "dev",
	//	DataID: "platform_admin_user_srv.json",
	//	Format: "json",
	//}
	//clientConfig := &constant.ClientConfig{
	//	NamespaceId:         namespaceID,
	//	TimeoutMs:           1000,
	//	NotLoadCacheAtStart: true,
	//	LogDir:              "./tmp/nacos/log",
	//	CacheDir:            "./tmp/nacos/cache",
	//	LogLevel:            "debug",
	//}
	//serverConfigs := []constant.ServerConfig{
	//	{
	//		IpAddr: ipAddr,
	//		Port:   uint64(port),
	//	},
	//}
	//
	//utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
	//	err := Init(conf, params,
	//		WithClientConfig(clientConfig),
	//		WithServerConfigs(serverConfigs),
	//	)
	//	t.Log(err, conf)
	//})
}

func TestNewNamingClient(t *testing.T) {
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		namingClient, err := NewNamingClient(ipAddr, port, namespaceID)
		t.Log(err, namingClient)
	})
}

func TestError(t *testing.T) {
	p := &Params{}
	p.Group = ""
	err := p.valid()
	assert.Error(t, err)

	p.Group = "group"
	p.DataID = ""
	err = p.valid()
	assert.Error(t, err)

	p.Group = "group"
	p.DataID = "id"
	p.Format = ""
	err = p.valid()
	assert.Error(t, err)

	p.Group = "group"
	p.DataID = "id"
	p.Format = "yml"
	err = p.valid()
	assert.NoError(t, err)

	p.Group = "group"
	p.DataID = "id"
	p.Format = "unknown"
	err = p.valid()
	assert.Error(t, err)

	err = Init(nil, p)
	assert.Error(t, err)
}
