package sql

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestQueryTran(t *testing.T) {
	x, err := QueryTran("SELECT * FROM yby.transfer")
	fmt.Println(x, err)
}

//InsertTran(id uint64, key string, timeunix int64, txin_addr string, txout_addr string, coin int, version int, tran_founder string)
func TestInsertTran(t *testing.T) {
	fmt.Println(InsertTran("bqh1", time.Now().UTC().Unix(), "123", "456", 20, 1, "foam", 11, "transfer"))
}

func TestHaskey(t *testing.T) {
	fmt.Println(Haskey("bqh1"))
}
