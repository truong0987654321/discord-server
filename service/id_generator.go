package service

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func init() {
	const nodeID int64 = 1
	var err error
	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("Failed tp init snowflake node: %v", err)
	}
}

func GenerateId() string {
	id := node.Generate()

	return id.String()
}
