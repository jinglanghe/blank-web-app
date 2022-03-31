package utils

import (
	uuid "github.com/apulis/sdk/go-utils/uuid"
)

var (
	generator = uuid.NewMist()
)

func NextID() int64 {
	return generator.Generate()
}
