package snow_flake

import (
	"common_pkg/snow_flake/snowflake_interf"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

type snowflake struct {
	node *sf.Node
}

func NewSnowflake(machineID int64) (snowflake_interf.InterfaceSnowFlake, error) {
	var st time.Time
	st, err := time.Parse("2006-01-02", "2023-01-01")
	if err != nil {
		return nil, err
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err := sf.NewNode(machineID)
	return &snowflake{
		node: node,
	}, nil
}

func (s *snowflake) GetId() int64 {
	return s.node.Generate().Int64()
}
