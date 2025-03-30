package entity

type Token struct {
	Expiry    int64  `json:"-"`
	Id        string `json:"token_id"`
	UserEmail string `json:"email"`
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
	return t[i].Expiry < t[j].Expiry
}

func (t *TokenPool) Update(items []Token) {
	(*t) = (*t)[:0]
	for _, data := range items {
		t.Push(data)
	}
}
