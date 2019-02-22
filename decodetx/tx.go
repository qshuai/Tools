package decodetx

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcutil"
)

func tx(rawtx string) (*btcutil.Tx, error) {
	bs, err := hex.DecodeString(rawtx)
	if err != nil {
		return nil, err
	}

	tx, err := btcutil.NewTxFromReader(bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}

	return tx, nil
}
