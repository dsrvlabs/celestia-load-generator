package main

import (
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

type keyInfo struct {
	Mnemonic string `json:"mnemonic"`
	Address  string `json:"address"`
}

func createKey() (*keyInfo, error) {
	kr := keyring.NewInMemory()
	algoList := keyring.SigningAlgoList{hd.Secp256k1}
	signAlgo, err := keyring.NewSigningAlgoFromString("secp256k1", algoList)

	bip44 := hd.CreateHDPath(118, 0, 0)
	key, mnemonic, err := kr.NewMnemonic("archway", keyring.English, bip44.String(), keyring.DefaultBIP39Passphrase, signAlgo)
	if err != nil {
		return nil, err
	}

	keyOut, err := keyring.MkAccKeyOutput(key)
	if err != nil {
		return nil, err
	}

	return &keyInfo{
		Mnemonic: mnemonic,
		Address:  keyOut.Address,
	}, nil
}

//func createKeyCreateTask(keyChan chan<- keyInfo) {
//	for {
//		newKey, err := createKey()
//		if err != nil {
//			time.Sleep(1 * time.Second)
//			continue
//		}
//
//		log.Println("New Account Created", newKey.Address)
//		keyChan <- *newKey
//	}
//}
