package storage

import "sync"

type RatesStorage struct {
	FiatMap   map[string]bool
	CryptoMap map[string]bool
	RatesMap  sync.Map
}

type RatesMap struct {
	RetrieveTime int64
	Rates        map[string]string
}

func NewRatesStorage() *RatesStorage {
	// can populate this with rates client
	fiatMap := map[string]bool{
		"USD": true,
		"EUR": true,
		"SGD": true,
	}
	cryptoMap := map[string]bool{
		"BTC": true,
		"ETH": true,
		"XRP": true,
	}
	return &RatesStorage{
		FiatMap:   fiatMap,
		CryptoMap: cryptoMap,
	}
}
