package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceInstance(t *testing.T) {
	s := NewServiceInstance("foo", "bar", []string{"grpc://192.168.209.8:8282"},
		WithVersion("v1.0.0"),
		WithMetadata(map[string]string{"foo": "bar"}),
	)
	assert.NotNil(t, s)
}
