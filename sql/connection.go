package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

const (
	Dbtype   string = "mysql"
	Dbname   string = "yby"
	Dbhostip string = "61.183.76.102:20011"
	User     string = "root"
	Password string = "244486"
)

type transction struct {
	id             uint64
	key            string
	time           string
	txin_addr      string
	txout_addr     string
	coin           float32
	version        int
	tran_founder   string
	change         float32
	transctionType string
}

func QueryTran(strsql string) (map[uint64]transction, error) {
	return query(strsql)
}

func InsertTran(key string, timeunix int64, txin_addr string, txout_addr string, coin float32, version uint64, tran_founder string, change float32, transctionType string) error {
	return insert(key, timeunix, txin_addr, txout_addr, coin, version, tran_founder, change, transctionType)
}

func query(strsql string) (map[uint64]transction, error) {
	openconf := mysql.Config{User: User, Passwd: Password, Addr: Dbhostip, DBName: Dbname, Net: "tcp"}
	fmt.Println(openconf.FormatDSN())
	db, err := sql.Open(Dbtype, openconf.FormatDSN())
	defer db.Close()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(strsql)
	if err != nil {
		return nil, err
	}

	tran := make(map[uint64]transction)
	for rows.Next() {
		var id uint64
		var key string
		var time string
		var txin_addr string
		var txout_addr string
		var coin float32
		var version int
		var tran_founder string
		var change float32
		var transctionType string

		rows.Columns()

		err = rows.Scan(&id, &key, &time, &txin_addr, &txout_addr, &coin, &version, &tran_founder, &change, &transctionType)
		tran[id] = transction{id, key, time, txin_addr, txout_addr, coin, version, tran_founder, change, transctionType}
		if err != nil {
			return nil, err
		}
	}
	return tran, err
	//fmt.Println(tran)
}

func Haskey(key string) (bool, error) {
	tran := make(map[uint64]transction)
	sqlstr := fmt.Sprintf("SELECT * FROM yby.transfers where transfers.key = '%v'", key)
	tran, err := query(sqlstr)
	if err != nil {
		return false, err
	}
	if len(tran) != 0 {
		return true, nil
	} else {
		return false, nil
	}
	//	else {
	//		return true, nil
	//	}
}

func insert(key string, timeunix int64, txin_addr string, txout_addr string, coin float32, version uint64, tran_founder string, change float32, utxoType string) error {
	open := fmt.Sprintf("%v:%v@tcp(%v)/%v", User, Password, Dbhostip, Dbname)
	db, err := sql.Open(Dbtype, open)
	defer db.Close()
	if err != nil {
		return err
	}
	//INSERT INTO `coinbase_transction_key`.`transction_key` (`key`, `transction_time`, `txin_addr`, `txout_addr`, `coin`, `version`, `transction_founder`) VALUES ('hh', 'h', '123', '456', '77', '1', 'foam');
	sqlstr := fmt.Sprintf("INSERT INTO `transfers` (`key`, `transction_time`, `txin_addr`, `txout_addr`, `coin`, `version`, `transction_founder`, `change`, `transction_type`) VALUES ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v')", key, time.Unix(timeunix, 0).Format("2006-01-02 15:04:05"), txin_addr, txout_addr, coin, version, tran_founder, change, utxoType)
	stmt, err := db.Prepare(sqlstr)
	fmt.Println(sqlstr)
	if err != nil {
		return err
	}

	res, err := stmt.Exec()
	if err != nil {
		return err
	}
	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
