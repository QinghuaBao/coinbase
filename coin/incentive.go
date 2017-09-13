package coin

import (
	"github.com/golang/protobuf/proto"
	"encoding/base64"
)

const (
	//coinbase 100 phcoin need txout
	INCENT_T0          float64 = 2
	INCENT_ALPHA0      float64 = 0.7
	//pre 100 phcoin adjust
	INCENT_THREADSHOLD int64   = 100*100000
)

//Dynamic adjustment
func updatePovSession(store Store, session *HydruscoinInfo_POVSession) {
	logger.Debug("Adjusting POV session parameters")
	session.CurrentAlpha = float32((INCENT_ALPHA0 * float64(session.TxCount)) / INCENT_T0)
	if session.CurrentAlpha > 0.7 {
		session.CurrentAlpha = 0.7
	}
	session.CurrentTotalIncentive = 0
	session.TxCount = 0
}

//timestamp make txout cant repeat
func doIncentive(store Store, incentives *map[string]*TX_TXOUT, version uint64, coin *Hydruscoin, timestamp int64) ([]byte, error) {
	logger.Debug("Doing Incentives")

	tx := &TX{Txout: make([]*TX_TXOUT, len(*incentives)), Version: version, Timestamp: timestamp, Founder: "blockchain"}

	i := 0
	for _, val := range *incentives {
		tx.Txout[i] = val
		i++
	}

	arg, err := proto.Marshal(tx)
	if err != nil {
		return nil, err
	}
	base64 := base64.StdEncoding.EncodeToString(arg)
	return coin.coinbase(store, []string{base64})
}
