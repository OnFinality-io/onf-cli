package service

import (
	"fmt"
	"github.com/OnFinality-io/onf-cli/pkg/watcher"
	"github.com/ethereum/go-ethereum/rpc"
)

const MaxRetryCount = 15

type SessionKey struct {
	Scheme    string `json:"scheme"`
	Phrase    string `json:"phrase"`
	PublicKey string `json:"publicKey"`
}

func InsertSessionKey(url string, key *SessionKey) error {
	c, _ := rpc.Dial(url)
	defer c.Close()
	err := c.Call(nil, "author_insertKey", key.Scheme, key.Phrase, key.PublicKey)
	if err != nil {
		retry := 0
		watcher := watcher.Watcher{Second: 2}
		watcher.Run(func(done chan bool) {
			err = c.Call(nil, "author_insertKey", key.Scheme, key.Phrase, key.PublicKey)
			//if err != nil {
			//	fmt.Printf("author_insertKey url:%s, err %v \n", url,err)
			//}
			if err == nil {
				done <- true
			}
			retry++
			if retry > MaxRetryCount {
				fmt.Printf("Author_insertKey url:%s, err %v \n", url, err)
				done <- true
			}
		})
	}
	return err
}
