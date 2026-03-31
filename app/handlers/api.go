package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type API struct {
	DB       *sql.DB
	Contract *bind.BoundContract
	Client   *ethclient.Client
}

type SetRequest struct {
	Value uint64 `json:"value"`
}

func (api *API) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "ERRO - método", http.StatusMethodNotAllowed)
		return
	}

	var output []interface{}
	err := api.Contract.Call(&bind.CallOpts{Context: r.Context()}, &output, "get")
	if err != nil {
		http.Error(w, fmt.Sprintf("ERRO - consulta: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "\nValor no contrato: %v\n", output[0])
}

func (api *API) SetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERRO - método", http.StatusMethodNotAllowed)
		return
	}

	var requisition SetRequest
	if err := json.NewDecoder(r.Body).Decode(&requisition); err != nil {
		http.Error(w, "ERRO - requisição", http.StatusBadRequest)
		return
	}
	// TODO - verificar ???
	pkey := os.Getenv("PRIVATE_KEY")
	if len(pkey) > 2 && pkey[:2] == "0x" {
		pkey = pkey[2:]
	}
	privateKey, err := crypto.HexToECDSA(pkey)
	if err != nil {
		http.Error(w, "ERRO - chave privada", http.StatusInternalServerError)
		return
	}

	chainID, err := api.Client.ChainID(r.Context())
	if err != nil {
		http.Error(w, "ERRO - ChainID", http.StatusInternalServerError)
		return
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		http.Error(w, "ERRO - transação", http.StatusInternalServerError)
		return
	}

	tx, err := api.Contract.Transact(auth, "set", big.NewInt(int64(requisition.Value)))
	if err != nil {
		http.Error(w, fmt.Sprintf("ERRO - envio transação: %v", err), http.StatusInternalServerError)
		return
	}

	receipt, err := bind.WaitMined(context.Background(), api.Client, tx)
	if err != nil {
		http.Error(w, "ERRO - mineração", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "\nTransação no bloco: %v\nHash: %s\n", receipt.BlockNumber, tx.Hash().Hex())
}

func (api *API) SyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERRO - Método", http.StatusMethodNotAllowed)
		return
	}

	var output []interface{}
	err := api.Contract.Call(&bind.CallOpts{Context: r.Context()}, &output, "get")
	if err != nil {
		http.Error(w, fmt.Sprintf("ERRO - leitura blockchain: %v", err), http.StatusInternalServerError)
		return
	}
	value := output[0].(*big.Int).String()

	query := `INSERT INTO contract_state (id, contract_value)
		VALUES (1, $1)
		ON CONFLICT (id) DO UPDATE SET contract_value = EXCLUDED.contract_value;`
	_, err = api.DB.Exec(query, value)
	if err != nil {
		http.Error(w, fmt.Sprintf("ERRO - salvar db: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "\nValor %s salvo (db).\n", value)
}
