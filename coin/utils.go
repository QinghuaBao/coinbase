package coin

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"

	"github.com/golang/protobuf/proto"
)

// ParseTXBytes unmarshal txData into TX object
func ParseTXBytes(txData []byte) (*TX, error) {
	tx := new(TX)
	err := proto.Unmarshal(txData, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// TxHash generates the Hash for the transaction.
func TxHash(tx *TX) string {
	txBytes, err := proto.Marshal(tx)
	if err != nil {
		return ""
	}

	fHash := sha256.Sum256(txBytes)
	lHash := sha256.Sum256(fHash[:])
	return hex.EncodeToString(lHash[:])
}

// ParseHydruscoinInfoBytes unmarshal infoBytes into HydruscoinInfo
func ParseHydruscoinInfoBytes(infoBytes []byte) (*HydruscoinInfo, error) {
	info := new(HydruscoinInfo)
	if err := proto.Unmarshal(infoBytes, info); err != nil {
		return nil, err
	}

	return info, nil
}

func ParseTestBytes(infoBytes []byte) (*Test, error) {
	info := new(Test)
	if err := proto.Unmarshal(infoBytes, info); err != nil {
		return nil, err
	}

	return info, nil
}

func ParseIncentiveBytes(infoBytes []byte) (*Incentive, error) {
	info := new(Incentive)
	if err := proto.Unmarshal(infoBytes, info); err != nil {
		return nil, err
	}

	return info, nil
}

func MapToSlice(txout map[string]*TX_TXOUT, txoutmapslice []*TxoutMap) {
	temp := make([]string, len(txout))

	i := 0
	for key, _ := range txout {
		temp[i] = key
		i++
	}

	sort.Strings(temp)

	//txoutmapslice = make([]*coin.TxoutMap, len(txout))
	for i = 0; i < len(txout); i++ {
		//fmt.Println(temp[i])

		//txoutmapslice[i] = {Key: temp[i], Txouts: txout[temp[i]]}
		//txoutmapslice[i] = new(coin.TxoutMap)
		txoutmapslice[i].Key = temp[i]
		txoutmapslice[i].Txouts = txout[temp[i]]
		//fmt.Println(txoutmapslice[i])
	}
}

func GenerateAccount(account *Account) *AccountSlice {
	acountslice := new(AccountSlice)
	acountslice.Addr = account.Addr
	acountslice.Balance = account.Balance

	acountslice.Txoutmap = make([]*TxoutMap, len(account.Txouts))
	for i := 0; i < len(account.Txouts); i++ {
		acountslice.Txoutmap[i] = new(TxoutMap)
	}

	MapToSlice(account.Txouts, acountslice.Txoutmap)

	return acountslice
}

func ParseAccount(accountslice *AccountSlice) *Account {
	account := new(Account)
	account.Addr = accountslice.Addr
	account.Balance = accountslice.Balance
	account.Txouts = make(map[string]*TX_TXOUT)

	for _, txoutmap := range accountslice.Txoutmap {
		account.Txouts[txoutmap.Key] = txoutmap.Txouts
	}

	return account
}
