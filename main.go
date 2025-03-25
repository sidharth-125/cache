package main

import (
	"container/heap"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sidharth-125/token_heap/entity"
)

var data []entity.Token = []entity.Token{
	{
		Id:        uuid.New().String(),
		UserEmail: "test1@yopmail.com",
		Expiry:    5,
	},

	{
		Id:        uuid.New().String(),
		UserEmail: "test2@yopmail.com",
		Expiry:    7,
	},

	{
		Id:        uuid.New().String(),
		UserEmail: "test3@yopmail.com",
		Expiry:    6,
	},
}

func Setup() entity.TokenPool {
	tokenPool := make(entity.TokenPool, 0)
	for _, item := range data {
		tokenPool = append(tokenPool, &item)
	}

	heap.Init(&tokenPool)
	return tokenPool
}

func main() {
	q := Setup()
	starter := time.NewTicker(time.Second * 2)

	for range starter.C {
		fmt.Println("ticking")
		remainingItems := make([]entity.Token, 0)
		for q.Len() > 0 {
			tokenData := q.Pop().(*entity.Token)
			if tokenData.Expiry > 3 {
				tokenData.Expiry -= 1
				remainingItems = append(remainingItems, *tokenData)
			} else {
				fmt.Printf("user : %s expired, expiry :%v\n", tokenData.UserEmail, tokenData.Expiry)
			}
		}

		if len(remainingItems) == 0 {
			fmt.Println("no items left")
			return
		}

		q.Update(remainingItems)
	}
}
