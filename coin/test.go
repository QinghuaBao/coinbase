package coin

import "encoding/base64"

func (coin *Hydruscoin) test(store Store, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, ErrInvalidArgs
	}

	txDataBase64 := args[0]
	txData, err := base64.StdEncoding.DecodeString(txDataBase64)
	if err != nil {
		logger.Errorf("Decoding base64 error: %v\n", err)
		return nil, err
	}

	test, err := ParseTestBytes(txData)
	if err != nil {
		logger.Errorf("Unmarshal tx bytes error: %v\n", err)
		return nil, err
	}
	logger.Debugf("test: %v", test)

	INCENT_T0 = test.INCENT_T0
	INCENT_ALPHA0 = test.INCENT_ALPHA0
	//pre 100 phcoin adjust
	INCENT_THREADSHOLD = test.GetINCENT_THREADSHOLD()
	return nil, nil
}
