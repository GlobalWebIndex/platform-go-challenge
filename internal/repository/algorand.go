package repository

import (
	"log"
	"ownify_api/internal/domain"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	_ "github.com/go-sql-driver/mysql"

	"github.com/spf13/viper"
)

type AlgoHandler interface {
	NewWalletQuery() WalletQuery
}

type algoHandler struct {
	client      *algod.Client
	Indexer     *indexer.Client
	testClient  *algod.Client
	testIndexer *indexer.Client
}

// NewWalletQuery implements AlgoHandler
func (*algoHandler) NewWalletQuery() WalletQuery {
	return &walletQuery{}
}

var Client *algod.Client

func NewAlgoHandler(client *algod.Client, indexer *indexer.Client, testClient *algod.Client, testIndexer *indexer.Client) AlgoHandler {
	return &algoHandler{client, indexer, testClient, testIndexer}
}

func NewAlgoClient(net string) (*algod.Client, *indexer.Client, error) {

	viper.AddConfigPath("../config")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("cannot read from a config")
	}
	algodAddress := viper.Get("algod.client").(string)
	indexerAddress := viper.Get("algod.indexer").(string)
	if net == domain.TestNet {
		algodAddress = viper.Get("algod.client.test").(string)
		indexerAddress = viper.Get("algod.indexer.test").(string)
	}

	// create algorand client
	algodClient, err := algod.MakeClient(algodAddress, "")
	if err != nil {
		return nil, nil, err
	}

	algodIndexer, err := indexer.MakeClient(indexerAddress, "")
	if err != nil {
		return nil, nil, err
	}
	return algodClient, algodIndexer, nil
}

func (d *dbHandler) NewWalletQuery() WalletQuery {
	return &walletQuery{}
}
