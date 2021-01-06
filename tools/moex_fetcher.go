package tools

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"time"
)

// sort desc to get fresh prices on the top
const (
	MoexUrl          = "https://iss.moex.com//iss/history/engines/stock/markets/index/securities/SBSPB.json?sort_order=desc&iss.only=history&history.columns=TRADEDATE,CLOSE,OPEN"
	MoexDateLayout   = "2006-01-02"
	BrokerDateLayout = "02.01.2006"
)

type MoexFetcherInterface interface {
	Fetch(from string) error
}

type MoexFetcher struct {
	ticker       string
	TickerPrices map[string]float64
	MoexFetcherInterface
}

func NewMoexFetcher(ticker string) *MoexFetcher {
	return &MoexFetcher{
		ticker:       ticker,
		TickerPrices: make(map[string]float64),
	}
}

func (fetcher *MoexFetcher) Fetch(from string) error {
	dateBroker, err := time.Parse(BrokerDateLayout, from)
	if err != nil {
		return err
	}
	//fmt.Println(dateBroker)

	page := 0
	for {
		fetchedEnough, err := fetcher.checkFetchedEnough(dateBroker)
		if err != nil {
			return err
		}
		if fetchedEnough {
			return nil
		}

		url := fmt.Sprintf("%s&start=%v", MoexUrl, page*100)
		page++
		err = fetcher.fetchHistoryPrices(url)
		if err != nil {
			return err
		}

		if page > 6 {
			fmt.Println("=============================================")
			fmt.Println("Return because i is too much")
			fmt.Println("=============================================")
		}
	}

}

type ParsedJSON struct {
	History History `json:"history"`
}

type History struct {
	Prices [][]interface{} `json:"data"`
}

func (fetcher *MoexFetcher) fetchHistoryPrices(moexUrl string) error {
	var parsedJson ParsedJSON
	fmt.Println(moexUrl)
	body, err := getRequest(moexUrl)
	if err != nil {
		return err
	}
	json.Unmarshal(body, &parsedJson)

	for i := 0; i < len(parsedJson.History.Prices); i++ {
		fmt.Print(".")
		row := parsedJson.History.Prices[i]
		//fmt.Println(row)
		avg := (row[1].(float64) + row[2].(float64)) / 2
		avg = math.Round(avg*1000) / 1000
		//fmt.Println(avg)

		moexDate := fmt.Sprintf("%v", row[0])
		fetcher.TickerPrices[moexDate] = avg
	}
	return nil
}

func (fetcher *MoexFetcher) checkFetchedEnough(from time.Time) (bool, error) {
	if len(fetcher.TickerPrices) == 0 {
		return false, nil
	}

	keys := make([]string, 0, len(fetcher.TickerPrices))
	for k := range fetcher.TickerPrices {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	lastFetchedDateMoex, err := time.Parse(MoexDateLayout, keys[0])
	if err != nil {
		return false, err
	}
	fmt.Println("\n Last fetched date", lastFetchedDateMoex)

	return from.After(lastFetchedDateMoex), nil
}

func getRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
