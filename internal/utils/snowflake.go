package utils

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

func Snowflake() int64 {
	nodeID := int64(1)

	// 创建一个雪花节点
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("Failed to create snowflake node: %v", err)
	}

	// 生成一个唯一 ID
	id := node.Generate()
	return id.Int64()
}
