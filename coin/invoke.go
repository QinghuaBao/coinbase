package coin

func (coin *Hydruscoin) registerAccount(store Store, args []string) ([]byte, error) {
	addr := args[0]

	if tmpaccount, err := store.GetAccount(addr); err == nil && tmpaccount != nil && tmpaccount.Addr == addr {
		logger.Warningf("account(%s) already registered.", addr)
		return nil, ErrAlreadyRegisterd
	}

	account := &Account{
		Addr:    addr,
		Balance: 0,
		Txouts:  make(map[string]*TX_TXOUT),
	}
	if err := store.PutAccount(account); err != nil {
		logger.Errorf("store.PutAccount(%#v) return error: %v", account, err)
		return nil, err
	}

	return nil, nil
}
