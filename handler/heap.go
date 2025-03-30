package handler

import (
	"container/heap"

	"github.com/sidharth-125/token_heap/entity"
)

func SetupTokenPool() entity.TokenPool {
	tokenPool := make(entity.TokenPool, 0)
	heap.Init(&tokenPool)
	return tokenPool
}
