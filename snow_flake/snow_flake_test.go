package snow_flake

import (
	"testing"
)

func TestGetId(t *testing.T) {
	newSnowflake, err := NewSnowflake(1)
	t.Log(err)
	id := newSnowflake.GetId()
	t.Log(id)
}
