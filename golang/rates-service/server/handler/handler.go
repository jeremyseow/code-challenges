package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeremyseow/rates-service/dto"
	"github.com/jeremyseow/rates-service/external/rate"
	"github.com/jeremyseow/rates-service/storage"
	"github.com/jeremyseow/rates-service/util"
	"github.com/jeremyseow/rates-service/webhook"
)

type Handlers struct {
	ratesClient  *rate.RatesClient
	ratesStorage *storage.RatesStorage
}

func NewHandlers(ratesClient *rate.RatesClient, ratesStorage *storage.RatesStorage) *Handlers {
	return &Handlers{
		ratesClient:  ratesClient,
		ratesStorage: ratesStorage,
	}
}

func (h *Handlers) GetUserAgent(c *gin.Context) {
	ua := c.Request.Header["User-Agent"]
	c.Header("Test-Header", "abc")
	c.String(http.StatusOK, "%s", ua)
}

func (h *Handlers) PostFile(c *gin.Context) {
	content := c.Param("content")
	err := util.WriteFile("./assets/", "hello.txt", []byte(content))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, gin.H{"results": "success"})
}

func (h *Handlers) PostWebHook(c *gin.Context) {
	webHookReq := &dto.WebHookRequest{}
	err := c.BindJSON(webHookReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	err = webhook.SendWebHook(webHookReq.URL, webHookReq.ID, webHookReq.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, webHookReq)
}

func (h *Handlers) GetRates(c *gin.Context) {
	currType := c.Query("base")
	if currType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "missing parameter: base",
		})
		return
	}

	defaultCurr := "USD"
	desiredCurrMap, targetCurrMap := h.ratesStorage.FiatMap, h.ratesStorage.CryptoMap
	if currType != "fiat" {
		desiredCurrMap, targetCurrMap = h.ratesStorage.CryptoMap, h.ratesStorage.FiatMap
	}

	if val, ok := h.ratesStorage.RatesMap.Load(defaultCurr); ok {
		rateMap, ok := val.(*storage.RatesMap)
		if ok && time.Now().Unix()-rateMap.RetrieveTime < 30 && len(rateMap.Rates) > 0 {
			fmt.Println("cache hit")
			results, err := h.populateMap(desiredCurrMap, targetCurrMap, rateMap.Rates)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, results)
			return
		}
	}

	fmt.Println("cache miss")

	rate, err := h.ratesClient.GetRates(defaultCurr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	rateMap := &storage.RatesMap{
		RetrieveTime: time.Now().Unix(),
		Rates:        rate.Data.Rates,
	}

	if len(rate.Data.Rates) != 0 {
		h.ratesStorage.RatesMap.Store(defaultCurr, rateMap)
	}

	results, err := h.populateMap(desiredCurrMap, targetCurrMap, rateMap.Rates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *Handlers) populateMap(desiredCurrMap, targetCurrMap map[string]bool, ratesMap map[string]string) (*dto.RateDTO, error) {
	results := dto.RateDTO{}
	for desiredCurr := range desiredCurrMap {
		convertedRates := map[string]string{}
		for targetCurr := range targetCurrMap {
			if ratesMap[desiredCurr] == "" || ratesMap[targetCurr] == "" {
				fmt.Printf("missing rate for %s and %s\n", desiredCurr, targetCurr)
				continue
			}

			desiredRate, err := strconv.ParseFloat(ratesMap[desiredCurr], 64)
			if err != nil {
				return nil, err
			}
			targetRate, err := strconv.ParseFloat(ratesMap[targetCurr], 64)
			if err != nil {
				return nil, err
			}

			convertedRates[targetCurr] = fmt.Sprintf("%f", (1/desiredRate)*targetRate)
		}
		results[desiredCurr] = convertedRates
	}

	return &results, nil
}
