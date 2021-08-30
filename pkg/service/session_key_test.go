package service

import "testing"

func TestInsertSessionKey(t *testing.T) {
	url := "https://node-6841254544723111936.jm.onfinality.io/rpc?apikey=5a5bc2dc-eb7e-4121-90d0-14ae8dd6d820"
	InsertSessionKey(url, &SessionKey{Scheme: "gran1", PublicKey: "publicKey", Phrase: "phrase"})
}
