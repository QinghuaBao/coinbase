package coin

import (
	"github.com/golang/protobuf/proto"
	"encoding/base64"
)

const (
	INCENT_T0          float64 = 1000
	INCENT_ALPHA0      float64 = 0.1
	INCENT_THREADSHOLD int64   = 1000000
)

//Dynamic adjustment
func updatePovSession(store Store, session *HydruscoinInfo_POVSession) {
	logger.Debug("Adjusting POV session parameters")
	session.CurrentAlpha = float32((INCENT_ALPHA0 * float64(session.TxCount)) / INCENT_T0)
	session.CurrentTotalIncentive = 0
	session.TxCount = 0
}

func doIncentive(store Store, incentives map[string]*TX_TXOUT, version uint64, coin *Hydruscoin) ([]byte, error) {
	logger.Debug("Doing Incentives")

	tx := &TX{Txout: make([]*TX_TXOUT, len(*incentives)), Version: version, Timestamp: 1060762481, Founder: "blockchain"}

	for _, val := range incentives {
		tx.Txout = append(tx.Txout, val)
	}

	arg, err := proto.Marshal(tx)
	if err != nil {
		return nil, err
	}
	base64 := base64.StdEncoding.EncodeToString(arg)
	return coin.coinbase(store, []string{base64})
}
