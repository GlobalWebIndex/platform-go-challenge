package repository

import (
	"log"
	"ownify_api/internal/domain"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/common"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/indexer"
	_ "github.com/go-sql-driver/mysql"

	"github.com/spf13/viper"
)

type AlgoHandler interface {
	NewWalletQuery() WalletQuery
	NewProductQuery() ProductQuery
}

type algoHandler struct{}

func NewAlgoHandler() AlgoHandler {
	return &algoHandler{}
}
func NewClient(net string) (*algod.Client, *indexer.Client, error) {
	viper.AddConfigPath("../config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("cannot read from a config")
	}
	algodAddress := viper.Get("algod.client").(string)
	indexerAddress := viper.Get("algod.indexer").(string)
	const algodToken = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	
	if net == domain.TestNet {
		algodAddress = viper.Get("algod.client.test").(string)
		indexerAddress = viper.Get("algod.indexer.test").(string)
	}

	// create algorand client
	commonClient, err := common.MakeClient(algodAddress, "X-API-Key", "") //algod.MakeClient(algodAddress, algodToken)
	algodClient := (*algod.Client)(commonClient)
	if err != nil {
		return nil, nil, err
	}

	algodIndexer, err := indexer.MakeClient(indexerAddress, algodToken)
	if err != nil {
		return nil, nil, err
	}

	return algodClient, algodIndexer, nil
}

func (a *algoHandler) NewWalletQuery() WalletQuery {
	return &walletQuery{}
}

func (a *algoHandler) NewProductQuery() ProductQuery {
	return &productQuery{}
}
