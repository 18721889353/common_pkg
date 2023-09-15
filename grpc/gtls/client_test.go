package gtls

import (
	"testing"

	"github.com/18721889353/common_pkg/grpc/gtls/certfile"

	"github.com/stretchr/testify/assert"
)

func TestGetClientTLSCredentials(t *testing.T) {
	credentials, err := GetClientTLSCredentials("192.168.209.8", certfile.Path("one-way/server.crt"))
	assert.NoError(t, err)
	assert.NotNil(t, credentials)

	_, err = GetClientTLSCredentials("192.168.209.8", certfile.Path("one-way/notfound.crt"))
	assert.Error(t, err)
}

func TestGetClientTLSCredentialsByCA(t *testing.T) {
	credentials, err := GetClientTLSCredentialsByCA(
		"192.168.209.8",
		certfile.Path("two-way/ca.pem"),
		certfile.Path("two-way/client/client.pem"),
		certfile.Path("two-way/client/client.key"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, credentials)

	_, err = GetClientTLSCredentialsByCA(
		"192.168.209.8",
		certfile.Path("two-way/ca.pem"),
		certfile.Path("two-way/client/notfound.pem"),
		certfile.Path("two-way/client/notfound.key"),
	)
	assert.Error(t, err)

	_, err = GetClientTLSCredentialsByCA(
		"192.168.209.8",
		certfile.Path("two-way/notfound.pem"),
		certfile.Path("two-way/client/client.pem"),
		certfile.Path("two-way/client/client.key"),
	)
	assert.Error(t, err)
}
