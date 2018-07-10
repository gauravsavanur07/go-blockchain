package main 

import "bytes"

type TXInput struct {
 	Txid   []byte
	Vout  int 
	Signature []byte 
	PubKey []byte 
}

func (in  *TXInput).UserKey(pubKeyHash []byte) bool {
lockingHash := HashPubKey(in.PubKey)
return bytes.Compare(lockingHash, pubKeyHash) == 0 
}

