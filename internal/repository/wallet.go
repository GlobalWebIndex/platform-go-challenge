package repository

import (
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"

	sq "github.com/Masterminds/squirrel"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
)

type WalletQuery interface {
	AddNewAccount(role string, userId string) (*string, error)
	GetMyAccounts(role string, userId string, net string) ([]string, error)
	MintOwnify(pubKey string, products []dto.BriefProduct, net string) ([]string, error)
	UpdatePinCode(role string, userId string, newPinCode string) error

	MakeTransaction(role string, userId string, pubKey string, rawTx []byte, net string) (*string, error)
}

type walletQuery struct{}

// func NewClient() (*algod.Client, *indexer.Client, error) {

// 	viper.AddConfigPath("../config")
// 	viper.SetConfigName("config")
// 	err := viper.ReadInConfig()
// 	if err != nil {
// 		log.Fatalln("cannot read from a config")
// 	}

// 	algodAddress := viper.Get("algod.client").(string)
// 	indexerAddress := "algod.indexer"

// 	// create algorand client
// 	algodClient, err := algod.MakeClient(algodAddress, "")
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	algodIndexer, err := indexer.MakeClient(indexerAddress, "")
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return algodClient, algodIndexer, nil
// }

func (w *walletQuery) AddNewAccount(
	role string,
	userId string,
) (*string, error) {

	//create new EOA in algorand.
	newAcc := crypto.GenerateAccount()
	mnemonic, err := mnemonic.FromPrivateKey(newAcc.PrivateKey)
	if err != nil {
		return nil, err
	}
	pubKey := newAcc.Address.String()

	//get user pin code hash from wallets table
	var pin string
	err = pgQb().Select("pin").Where(sq.Eq{"email": userId}).From("ownify.wallets").QueryRow().Scan(&pin)
	if err != nil {
		return nil, err
	}

	//encrypt mnemonic.
	cipher, err := utils.Encrypt(mnemonic, pin)
	if err != nil {
		return nil, err
	}

	//inset to wallet table.
	cols := []string{"chain_id", "pub_addr", "email", "user_role", "seed_cipher"}
	values := []interface{}{0, pubKey, userId, role, cipher}
	sqlBuilder := utils.NewSqlBuilder()
	query, err := sqlBuilder.Insert("wallets", cols, values)
	if err != nil {
		return nil, err
	}
	_, err = DB.Exec(*query)
	if err != nil {
		return nil, err
	}
	return &pubKey, nil
}

func (w *walletQuery) GetMyAccounts(
	role string,
	userId string,
	net string,
) ([]string, error) {
	var accounts = []string{}
	err := pgQb().Select("pub_addr").Where(sq.Eq{"email": userId}).From("ownify.wallets").QueryRow().Scan(&accounts)
	if err != nil {
		return []string{}, err
	}
	return accounts, nil
}

func (w *walletQuery) MintOwnify(
	pubKey string,
	products []dto.BriefProduct,
	net string,
) ([]string, error) {
	var cipher string
	var email string
	err := pgQb().Select("pub_addr", "email").Where(sq.Eq{"pub_addr": pubKey}).From("ownify.wallets").QueryRow().Scan(&email, &cipher)
	if err != nil {
		return []string{}, err
	}

	return []string{}, nil
}

func (w *walletQuery) UpdatePinCode(role string, userId string, newPinCode string) error {
	//var cipher string
	///var email string
	// err := pgQb().Select("pub_addr", "email").Where(sq.Eq{"pub_addr": pubKey}).From("ownify.wallets").QueryRow().Scan(&email, &cipher)
	// if err != nil {
	// 	return []string{}, err
	// }

	return nil
}

func (w *walletQuery) MakeTransaction(role string, userId string, pubKey string, rawTx []byte, net string) (*string, error) {
	txId := ""
	return &txId, nil
}
