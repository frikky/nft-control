package main

import (
	"context"
	"crypto/ecdsa"
	//"fmt"
	"log"
	"math/big"

	store "./contracts" // for demo

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	//ethClient := "https://rinkeby.infura.io"
	ethClient := "https://ropsten.infura.io/v3/85b826fbec8f4c5d8bcff63e1c6f48f0"
	log.Printf("[DEBUG] Checking %s", ethClient)
	client, err := ethclient.Dial(ethClient)
	if err != nil {
		log.Fatalf("Client error: %s", err)
	}

	log.Printf("[DEBUG] Setting up with private key")
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatalf("Key error: %s", err)
	}

	publicKey := privateKey.Public()
	log.Printf("[DEBUG] Setting up with private key (2)")
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatalf("PubKey error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Printf("[DEBUG] Setting up with private key (3)")
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Nonce error: %s", err)
	}

	log.Printf("[DEBUG] Setting up with private key (4)")
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Gas error: %s", err)
	}

	log.Printf("[DEBUG] Setting up with private key. Gas price: %d", gasPrice)
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	input := "1.0"
	_ = input
	address, tx, instance, err := store.DeployStore(auth, client, input)
	if err != nil {
		log.Fatalf("Address error: %s", err)
	}

	log.Println(address.Hex())   // 0x147B8eb97fD247D06C4006D269c90C1908Fb5D54
	log.Println(tx.Hash().Hex()) // 0xdae8ba5444eefdc99f4d45cd0c4f24056cba6a02cefbf78066ef9f4188ff7dc0

	_ = instance
}
