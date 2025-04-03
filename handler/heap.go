package handler

import (
	"container/heap"
	"log"
	"sync"
	"time"

	"github.com/sidharth-125/token_heap/entity"
)

func SetupTokenPool() entity.TokenPool {
	tokenPool := make(entity.TokenPool, 0)
	heap.Init(&tokenPool)
	return tokenPool
}

func TokenManager(pool *entity.TokenPool, mu *sync.RWMutex, tokenMap map[string]*entity.Token) {
	mu.Lock()
	defer mu.Unlock()

	if pool.Len() == 0 {
		log.Println("token pool is empty")
		return
	}

	remainingItems := make([]entity.Token, 0)
	for pool.Len() > 0 {
		tokenData := pool.Pop().(*entity.Token)
		if tokenData.Expiry.Before(time.Now()) {
			log.Printf("token  %s expired", tokenData.Id)
			delete(tokenMap, tokenData.Id)

		} else {
			remainingItems = append(remainingItems, *tokenData)
		}
	}

	if len(remainingItems) > 0 {
		pool.Update(remainingItems)
	}

}
