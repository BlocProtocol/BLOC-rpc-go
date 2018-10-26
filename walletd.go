/*

Copyright (c) 2018 Rashed Mohammed, The TurtleCoin Developers, The BLOC Developers

Please see the included LICENSE file for more information

*/

package blocrpc

import (
	"bytes"
	"errors"
)

// Walletd structure contains the URL and Port
// info of the node and RPC Password for RPC calls
type Walletd struct {
	URL         string
	Port        int
	RPCPassword string
}

func (wallet *Walletd) check() error {
	if wallet.URL == "" {
		wallet.URL = "127.0.0.1"
	}
	if wallet.Port == 0 {
		wallet.Port = 8070
	}
	if wallet.RPCPassword == "" {
		return errors.New("RPCPassword not specified")
	}
	return nil
}

/*
Save method saves the wallet without closing it.
*/
func (wallet *Walletd) Save() (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	return wallet.makePostRequest("save", params), nil
}

/*
Reset method resyncs the wallet if no viewSecretKey is given.
If viewSecretKey is given then it replaces the existing wallet with a new one
corresponding to the viewSecretKey
*/
func (wallet *Walletd) Reset(scanHeight int) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["scanHeight"] = scanHeight
	return wallet.makePostRequest("reset", params), nil
}

/*
CreateAddress method creates a new address inside the container along with old addresses
*/
func (wallet *Walletd) CreateAddress(
	spendSecretKey string,
	spendPublicKey string,
	scanHeight int,
	newAddress bool) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}

	params := make(map[string]interface{})

	if spendSecretKey != "" && spendPublicKey != "" {
		return nil, errors.New("Cannot specify both spend keys.. either or both should be empty")
	} else if spendPublicKey == "" {
		params["spendSecretKey"] = spendSecretKey
	} else {
		params["spendPublicKey"] = spendPublicKey
	}

	if newAddress {
		params["newAddress"] = newAddress
	} else {
		params["scanHeight"] = scanHeight
	}

	return wallet.makePostRequest("createAddress", params), nil
}

/*
DeleteAddress method deletes the specified address from the container
*/
func (wallet *Walletd) DeleteAddress(address string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["address"] = address
	return wallet.makePostRequest("deleteAddress", params), nil
}

/*
GetSpendKeys method returns the spendPublicKey and spendSecretKey corresponding
the given input wallet address
*/
func (wallet *Walletd) GetSpendKeys(address string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["address"] = address
	return wallet.makePostRequest("getSpendKeys", params), nil
}

/*
GetBalance method returns the balance present in the specified address
If the address is empty then returns the balance present in the container
*/
func (wallet *Walletd) GetBalance(address string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["address"] = address
	return wallet.makePostRequest("getBalance", params), nil
}

/*
GetBlockHashes method returns array of hashes starting from specified blockIndex upto blockCount
*/
func (wallet *Walletd) GetBlockHashes(firstBlockIndex int, blockCount int) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["firstBlockIndex"] = firstBlockIndex
	params["blockCount"] = blockCount
	return wallet.makePostRequest("getBlockHashes", params), nil
}

/*
GetTransactionHashes method returns array of objects containing block and transaction hashes
of the specified address
*/
func (wallet *Walletd) GetTransactionHashes(
	addresses []string,
	blockHash string,
	firstBlockIndex int,
	blockCount int,
	paymentID string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})

	if blockHash != "" {
		params["blockHash"] = blockHash
	} else {
		params["firstBlockIndex"] = firstBlockIndex
	}

	params["addresses"] = addresses
	params["blockCount"] = blockCount
	params["paymentId"] = paymentID
	return wallet.makePostRequest("getTransactionHashes", params), nil
}

/*
GetTransactions method returns array of objects containing block and transaction details
of the specified address
*/
func (wallet *Walletd) GetTransactions(
	addresses []string,
	blockHash string,
	firstBlockIndex int,
	blockCount int,
	paymentID string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})

	if blockHash != "" {
		params["blockHash"] = blockHash
	} else {
		params["firstBlockIndex"] = firstBlockIndex
	}

	params["addresses"] = addresses
	params["blockCount"] = blockCount
	params["paymentId"] = paymentID
	return wallet.makePostRequest("getTransactions", params), nil
}

/*
GetUnconfirmedTransactionHashes method returns array of hashes of unconfirmed transactions of the specified address
*/
func (wallet *Walletd) GetUnconfirmedTransactionHashes(addresses []string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["addresses"] = addresses
	return wallet.makePostRequest("getUnconfirmedTransactionHashes", params), nil
}

/*
GetTransaction method returns the transaction details of a particular specified transaction hash
*/
func (wallet *Walletd) GetTransaction(transactionHash string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})

	if transactionHash == "" {
		return nil, errors.New("transactionHash cannot be empty.. please specify a valid hash")
	}

	params["transactionHash"] = transactionHash
	return wallet.makePostRequest("getTransaction", params), nil
}

/*
SendTransaction method sends specified transactions
*/
func (wallet *Walletd) SendTransaction(
	addresses []string,
	transfers []map[string]interface{},
	fee int,
	unlockTime int,
	extra string,
	paymentID string,
	changeAddress string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["addresses"] = addresses
	params["transfers"] = transfers
	params["fee"] = fee
	params["unlockTime"] = unlockTime
	params["changeAddress"] = changeAddress

	if extra != "" && paymentID != "" {
		return nil, errors.New("Can't set paymentID and extra together.. either or both should be empty")
	} else if extra != "" {
		params["extra"] = extra
	} else {
		params["paymentId"] = paymentID
	}

	return wallet.makePostRequest("sendTransaction", params), nil
}

/*
CreateDelayedTransaction method allows you to create a delayed transaction
Such transactions are not sent into the network automatically and should be pushed
using SendDelayedTransaction method
*/
func (wallet *Walletd) CreateDelayedTransaction(
	addresses []string,
	transfers []map[string]interface{},
	fee int,
	unlockTime int,
	extra string,
	paymentID string,
	changeAddress string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["addresses"] = addresses
	params["transfers"] = transfers
	params["fee"] = fee
	params["unlockTime"] = unlockTime
	params["changeAddress"] = changeAddress

	if extra != "" && paymentID != "" {
		return nil, errors.New("Can't set paymentID and extra together.. either or both should be empty")
	} else if extra != "" {
		params["extra"] = extra
	} else {
		params["paymentId"] = paymentID
	}

	return wallet.makePostRequest("createDelayedTransaction", params), nil
}

/*
GetDelayedTransactionHashes method returns array of delayedTransactionHashes
*/
func (wallet *Walletd) GetDelayedTransactionHashes() (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	return wallet.makePostRequest("getDelayedTransactionHashes", params), nil
}

/*
DeleteDelayedTransaction method deletes the specified delayedTransactionHash
*/
func (wallet *Walletd) DeleteDelayedTransaction(transactionHash string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["transactionHash"] = transactionHash
	return wallet.makePostRequest("deleteDelayedTransaction", params), nil
}

/*
SendDelayedTransaction method sends the delayedTransaction created using CreateDelayedTransaction
method into the network
*/
func (wallet *Walletd) SendDelayedTransaction(transactionHash string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["transactionHash"] = transactionHash
	return wallet.makePostRequest("sendDelayedTransaction", params), nil
}

/*
GetViewKey method returns the viewSecretKey of the wallet
*/
func (wallet *Walletd) GetViewKey() (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	return wallet.makePostRequest("getViewKey", params), nil
}

/*
GetMnemonicSeed method returns the 25 word random seed corresponding to
the given input wallet address
*/
func (wallet *Walletd) GetMnemonicSeed(address string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["address"] = address
	return wallet.makePostRequest("getMnemonicSeed", params), nil
}

/*
GetStatus method returns the sync state of the wallet and known top block height
*/
func (wallet *Walletd) GetStatus() (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	return wallet.makePostRequest("getStatus", params), nil
}

/*
GetAddresses method returns an array of addresses present in the container
*/
func (wallet *Walletd) GetAddresses() (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	return wallet.makePostRequest("getAddresses", params), nil
}

/*
SendFusionTransaction method allows you to send a fusion transaction from selected address to destination
address. If there aren't any outputs that can be optimized it returns an error.
*/
func (wallet *Walletd) SendFusionTransaction(
	threshold int,
	addresses []string,
	destinationAddress string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["threshold"] = threshold
	params["addresses"] = addresses
	params["destinationAddress"] = destinationAddress
	return wallet.makePostRequest("sendFusionTransaction", params), nil
}

/*
EstimateFusion method returns the number of outputs that can be optimized
This is helpful for sending fusion transactions
*/
func (wallet *Walletd) EstimateFusion(threshold int, addresses []string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["threshold"] = threshold
	params["addresses"] = addresses
	return wallet.makePostRequest("estimateFusion", params), nil
}

/*
CreateIntegratedAddress method creates a unique 236 char long address which corresponds to
the specified address with paymentID
*/
func (wallet *Walletd) CreateIntegratedAddress(address string, paymentID string) (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["address"] = address
	params["paymentId"] = paymentID
	return wallet.makePostRequest("createIntegratedAddress", params), nil
}

/*
GetFeeInfo method returns the fee information that the service picks up from the
connected daemon
*/
func (wallet *Walletd) GetFeeInfo() (*bytes.Buffer, error) {
	err := wallet.check()
	if err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	return wallet.makePostRequest("getFeeInfo", params), nil
}
