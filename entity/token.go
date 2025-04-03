package entity

import "time"

type Token struct {
	Expiry    time.Time `json:"expiry"`
	Id        string    `json:"token_id"`
	UserEmail string    `json:"email"`
}

type TokenPool []*Token

func (t *TokenPool) Push(x any) {
	data := x.(Token)
	*t = append(*t, &data)
}

func (t *TokenPool) Pop() any {
	length := len(*t)
	token := (*t)[length-1]
	*t = (*t)[:length-1]
	return token
}

func (t TokenPool) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TokenPool) Len() int {
	return len(t)
}

func (t TokenPool) Less(i, j int) bool {
	first := t[i].Expiry
	second := t[j].Expiry
	return first.Before(second)
}

func (t *TokenPool) Update(items []Token) {
	(*t) = (*t)[:0]
	for _, data := range items {
		t.Push(data)
	}
}
