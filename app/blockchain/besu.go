package blockchain

import (
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	ContractABI = `[{"inputs":[],"name":"get","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"x","type":"uint256"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
)

func InitContract(nodeURL string, contractAddressHex string) (*ethclient.Client, *bind.BoundContract) {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Fatalf("FALHA - Besu: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(ContractABI))
	if err != nil {
		log.Fatalf("FALHA - ABI: %v", err)
	}

	address := common.HexToAddress(contractAddressHex)
	boundContract := bind.NewBoundContract(address, parsedABI, client, client, client)

	return client, boundContract
}
