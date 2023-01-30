package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"ownify_api/internal/utils"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

type WalletQuery interface {
	AddNewAccount(role string, email string) (*string, error)
	GetMyAccounts(role string, email string) ([]string, error)
	MintOwnify(email string, pubKey string, products []dto.BriefProduct, net string) ([]uint64, error)
	UpdatePinCode(role string, email string, newPinCode string) error

	MakeTransaction(role string, userId string, pubKey string, rawTx []byte, net string) (*string, error)
}

type walletQuery struct{}

func (w *walletQuery) AddNewAccount(
	role string,
	email string,
) (*string, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, err
	}
	tableName := domain.BusinessTableName
	if role == domain.PersonalWallet {
		tableName = domain.PersonTableName
	}
	//get user pin code hash from wallets table
	var pin string
	sqlBuilder := utils.NewSqlBuilder()
	sql, err := sqlBuilder.Select(tableName, []string{"pin"}, []utils.Tuple{{Key: "email", Val: email}}, "OR")
	if err != nil {
		return nil, err
	}

	err = DB.QueryRow(*sql).Scan(&pin)
	if err != nil {
		return nil, err
	}
	//create new EOA in algorand.
	newAcc := crypto.GenerateAccount()
	mnemonic, err := mnemonic.FromPrivateKey(newAcc.PrivateKey)
	if err != nil {
		return nil, err
	}
	pubKey := newAcc.Address.String()

	//encrypt mnemonic.
	cipher, err := utils.Encrypt(mnemonic, pin)
	if err != nil {
		return nil, err
	}

	//inset to wallet table.
	cols := []string{"chain_id", "pub_addr", "email", "user_role", "seed_cipher"}
	values := []interface{}{0, pubKey, email, role, cipher}

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
	email string,
) ([]string, error) {
	var accounts = []string{}
	sqlBuilder := utils.NewSqlBuilder()
	sql, err := sqlBuilder.Select(domain.BusinessTableName, []string{
		"pub_addr",
	}, []utils.Tuple{{Key: "email", Val: email}}, "OR")
	if err != nil {
		return []string{}, err
	}
	err = DB.QueryRow(*sql).Scan(&accounts)
	if err != nil {
		return []string{}, err
	}
	return accounts, nil
}

func (w *walletQuery) MintOwnify(
	email string,
	pubKey string,
	products []dto.BriefProduct,
	net string,
) ([]uint64, error) {

	//get seed from wallet table.
	cipherR := make(chan domain.Result[string])
	pinR := make(chan domain.Result[string])
	go func() {
		var cipher string
		sqlBuilder := utils.NewSqlBuilder()
		sql, err := sqlBuilder.Select(domain.WalletTableName, []string{
			"seed_cipher",
		}, []utils.Tuple{{Key: "email", Val: email}, {Key: "pub_addr", Val: pubKey}}, "AND")
		if err != nil {
			cipherR <- domain.Result[string]{Err: err}
			return
		}
		err = DB.QueryRow(*sql).Scan(&cipher)
		if err != nil {
			cipherR <- domain.Result[string]{Err: err}
			return
		}
		seed := ""
		err = DB.QueryRow(*sql).Scan(&seed)
		if seed == "" {
			cipherR <- domain.Result[string]{Err: err}
			return
		}
		cipherR <- domain.Result[string]{Val: seed}
	}()

	//get pin code from business table.
	go func() {
		pin := ""
		sqlBuilder := utils.NewSqlBuilder()
		sql, err := sqlBuilder.Select(domain.BusinessTableName, []string{
			"pin",
		}, []utils.Tuple{{Key: "email", Val: email}}, "AND")
		if err != nil {
			pinR <- domain.Result[string]{Err: err}
			return
		}
		err = DB.QueryRow(*sql).Scan(&pin)
		if err != nil {
			pinR <- domain.Result[string]{Err: err}
			return
		}
		pinR <- domain.Result[string]{Val: pin}
	}()

	pin := <-pinR
	cipher := <-cipherR

	if !pin.Ok() {
		return nil, pin.Err
	}
	if !cipher.Ok() {
		return nil, pin.Err
	}

	// decrypt cipher for recover account.
	seed, err := utils.Decrypt(cipher.Val, pin.Val)
	if err != nil {
		return nil, pin.Err
	}
	prv, err := mnemonic.ToPrivateKey(seed)
	if err != nil {
		return nil, pin.Err
	}

	//algorand client initialize
	client, _, err := NewClient(net)
	if err != nil {
		return nil, pin.Err
	}

	//build ASA transaction
	txns := []types.Transaction{}
	for _, product := range products {
		note, err := json.Marshal(product)
		if err != nil {
			return nil, err
		}
		txParams, err := client.SuggestedParams().Do(context.Background())
		if err != nil {
			return nil, err
		}
		metaHash := utils.Hash(fmt.Sprintf("%v", note))

		txn, err := transaction.MakeAssetCreateTxn(pubKey,
			note,
			txParams, 1, 0,
			false, pubKey, pubKey, pubKey, pubKey,
			domain.OwnifyAssetName, domain.OwnifyAssetUnit, domain.OwnifyAssetMetaUrl, metaHash)

		if err != nil {
			return nil, err
		}
		txns = append(txns, txn)
	}

	groupedTxs, err := transaction.AssignGroupID(txns, pubKey)
	if err != nil {
		return nil, pin.Err
	}
	var stxs []byte
	for _, txn := range groupedTxs {
		_, stx, _ := crypto.SignTransaction(prv, txn)
		stxs = append(stxs, stx...)
	}

	if err != nil {
		return nil, err
	}

	pendingTxID, err := client.SendRawTransaction(stxs).Do(context.Background())

	if err != nil {
		return nil, err
	}
	confirmedTx, err := transaction.WaitForConfirmation(client, pendingTxID, 4, context.Background())
	if err != nil {
		return nil, err
	}
	endIndex := confirmedTx.AssetIndex + uint64(len(products))

	//add product to db.

	return utils.MakeRange(confirmedTx.AssetIndex, endIndex), nil
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
