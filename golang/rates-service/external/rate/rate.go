package rate

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeremyseow/rates-service/limiter"
)

type RatesResp struct {
	Data struct {
		BaseCurr string            `json:"currency"`
		Rates    map[string]string `json:"rates"`
	} `json:"data"`
}

type AllFiatResp struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

type AllCryptoResp struct {
	Data []struct {
		Code string `json:"code"`
	} `json:"data"`
}

type RatesClient struct {
	limiter *limiter.RateLimiter
}

func NewRatesClient() *RatesClient {
	return &RatesClient{
		limiter: limiter.NewRateLimiter(),
	}
}

func (r *RatesClient) GetAllFiat() (*AllFiatResp, error) {
	url := "https://api.coinbase.com/v2/currencies"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	payload := &AllFiatResp{}
	err = json.NewDecoder(resp.Body).Decode(payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (r *RatesClient) GetAllCrypto() (*AllCryptoResp, error) {
	url := "https://api.coinbase.com/v2/currencies/crypto"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	payload := &AllCryptoResp{}
	err = json.NewDecoder(resp.Body).Decode(payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (r *RatesClient) GetRates(baseCurr string) (*RatesResp, error) {
	if !r.limiter.Allow() {
		fmt.Println("rate limit exceeded")
		return &RatesResp{}, nil
	}

	url := fmt.Sprintf("https://api.coinbase.com/v2/exchange-rates?currency=%s", baseCurr)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	payload := &RatesResp{}
	err = json.NewDecoder(resp.Body).Decode(payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
