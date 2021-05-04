package hdwallet

import (
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/ethereum/go-ethereum/crypto/secp256k1"
)

func GenerateAddress(xPubKStr string, id uint32, index uint32) (string, error) {
	xpub, err := hdkeychain.NewKeyFromString(xPubKStr)
	if err != nil {
		return "", err
	}
	child, err := xpub.Child(id)
	if err != nil {
		return "", err
	}

	indexChild, err := child.Child(index)
	if err == nil {
		childEC, _ := indexChild.ECPubKey()
		address := crypto.PubkeyToAddress(*childEC.ToECDSA()).Hex()
		return address, nil
	}

	return "", err
}

func GetPrivateKey(xPvKStr string, id uint32, index uint32) (string, error) {
	xprivate, err := hdkeychain.NewKeyFromString(xPvKStr)
	if err != nil {
		return "", err
	}
	child, err := xprivate.Child(id)
	if err != nil {
		return "", err
	}

	indexChild, err := child.Child(index)
	if err == nil {
		child, _ := indexChild.ECPrivKey()
		privatekey := crypto.FromECDSA(child.ToECDSA())

		return hexutil.Encode(privatekey), nil
	}

	return "", err
}
