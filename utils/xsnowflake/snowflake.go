package xsnowflake

import "github.com/bwmarrin/snowflake"

var node *snowflake.Node

func init() {
	node, _ = snowflake.NewNode(0)
}

func Init(machineId int64) error {
	var err error
	node, err = snowflake.NewNode(machineId)
	if err != nil {
		return err
	}
	return nil
}

func GetNode() *snowflake.Node {
	return node
}

func GenId() int64 {
	if node == nil {
		return 0
	}

	return node.Generate().Int64()
}
