package service

import (
	"github.com/ethereum/go-ethereum/rpc"
)

type SessionKey struct {
	Scheme    string
	Phrase    string
	PublicKey string
}

func InsertSessionKey(url string, key *SessionKey) error {
	c, _ := rpc.Dial(url)
	defer c.Close()
	return c.Call(nil, "author_insertKey", key.Scheme, key.Phrase, key.PublicKey)
}
