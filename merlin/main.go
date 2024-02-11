package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Code int `json:"code"`
	Data struct {
		List []struct {
			Address    string `json:"address"`
			CntBitmaps string `json:"cnt_bitmaps"`
			Percent    string `json:"percent"`
			Rank       string `json:"rank"`
		} `json:"list"`
	} `json:"data"`
}

type ResponseOne struct {
	Data struct {
		DomainAssets []struct {
			SuffixInfo struct {
				Name string `json:"name"`
			} `json:"suffixInfo"`
			Count int `json:"count"`
		} `json:"domainAssets"`
	} `json:"data"`
}

type ResponseTwo struct {
	ChainStats struct {
		FundedTxoSum int `json:"funded_txo_sum"`
	} `json:"chain_stats"`
	MempoolStats struct {
		FundedTxoSum int `json:"funded_txo_sum"`
	} `json:"mempool_stats"`
}

type ResponseThree struct {
	Data struct {
		Brc420Assets []Brc420Asset `json:"brc420Assets"`
	} `json:"data"`
}

type Brc420Asset struct {
	TokenInfo struct {
		Name string `json:"name"`
	} `json:"tokenInfo"`
	Count int `json:"count"`
}

func getDomainAssets(address string) (int, error) {
	url := fmt.Sprintf("https://search.idclub.io/assets/address/%s", address)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("请求失败：%s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败：%s", err)
	}

	var result ResponseOne
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, fmt.Errorf("解析bitmap响应失败：%s", err)
	}

	if len(result.Data.DomainAssets) > 0 {
		return result.Data.DomainAssets[0].Count, nil
	}

	return 0, nil
}

func getBTCAmount(address string) (float64, error) {
	url := fmt.Sprintf("https://mempool.space/api/address/%s", address)

	// sleep for 3 seconds to avoid frequent request
	time.Sleep(3 * time.Second)

	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("请求失败：%s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败：%s", err)
	}

	var data ResponseTwo
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, fmt.Errorf("解析btc响应失败：%s", err)
	}

	totalFundedTxoSum := data.ChainStats.FundedTxoSum + data.MempoolStats.FundedTxoSum
	btcAmount := float64(totalFundedTxoSum) / 100000000

	return btcAmount, nil
}

func getBrc420Assets(address string) ([]Brc420Asset, error) {
	url := fmt.Sprintf("https://search.idclub.io/assets/address/%s", address)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败：%s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败：%s", err)
	}

	var response ResponseThree
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("解析420响应失败：%s", err)
	}

	return response.Data.Brc420Assets, nil
}

func GetCntBitmaps() (string, error) {
	url := "https://www.geniidata.com/api/dashboard/chart/public/data?chartId=126020&pageSize=3&page=1&searchKey=&searchValue="

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var responseObject Response
	err = json.Unmarshal(body, &responseObject)
	if err != nil {
		return "", err
	}

	for _, item := range responseObject.Data.List {
		if item.Rank == "1" {
			return item.CntBitmaps, nil
		}
	}

	return "", fmt.Errorf("cnt_bitmaps not found")
}

func main() {
	for {
		// Create a ticker that ticks every 10 seconds
		ticker := time.NewTicker(10 * time.Second)

		// Create a channel to receive a signal when the program should stop
		stop := make(chan bool)

		go func() {
			for {
				select {
				case <-ticker.C:
					// Call the three functions and get their results
					brc420Result, err := brc420Address()
					if err != nil {
						fmt.Println(err)
						continue
					}

					bitmapResult, err := bitmapAddress()
					if err != nil {
						fmt.Println(err)
						continue
					}

					btcResult, err := btcAddress()
					if err != nil {
						fmt.Println(err)
						continue
					}

					// Create a string to hold the combined results
					combinedResult := fmt.Sprintf("%s\n%s\n%s\n", bitmapResult, brc420Result, btcResult)
					combinedResult = strings.TrimSuffix(combinedResult, "\n")

					// Write the combined result to a file
					err = ioutil.WriteFile("/app/merlin/merlinall.txt", []byte(combinedResult), 0644)
					if err != nil {
						fmt.Println(err)
						continue
					}

					fmt.Println("Results written to merlinall.txt")
				case <-stop:
					ticker.Stop()
					return
				}
			}
		}()

		// Wait for a signal to stop the program
		<-stop
	}
}

func brc420Address() (string, error) {
	address := "bc1q4gfsheqz7ll2wdgfwjh2l5hhr45ytc4ekgxaex"
	results, err := getBrc420Assets(address)
	if err != nil {
		return "", err
	}

	resultString := ""
	for _, result := range results {
		resultString += fmt.Sprintf("%s: %d\n", result.TokenInfo.Name, result.Count)
	}

	return resultString, nil
}

func bitmapAddress() (string, error) {
	Bitmap := "Bitmap"
	//address := "bc1qptgujmlkez7e6744yctzjgztu0st372mxs6702"
	//count, err := getDomainAssets(address)
	count, err := GetCntBitmaps()
	fmt.Println(count)
	if err != nil {
		return "", err
	}

	resultString := fmt.Sprintf("%s: %s\n", Bitmap, count)
	//resultString := fmt.Sprintf("%s: %d\n", Bitmap, count)
	return resultString, nil
}

func btcAddress() (string, error) {
	addresses := []string{
		"bc1qua5y9yhknpysslxypd4dahagj9jamf90x4v90x",
		"15zVuow5e9Zwj4nTrxSH3Rvupk32wiKEsr",
		"bc1qq3c6kehun66sdek3q0wmu540n3vg0hgrekkjce",
		"124SzTv3bBXZVPz2Li9ADs9oz4zCfT3VmM",
		"bc1qtu66zfqxj6pam6e0zunwnggh87f5pjr7vdr5cd",
		"bc1qyqt9zs42qmyf373k7yvy0t3askxd927v304xlv",
		"16LDby5cWxzQqTFJrA1DDmbwABumCQHteG",
		"1EEU18ZvWrbMxdXEuqdii6goDKbAbaXiA1",
	}

	totalAmount := 0.0
	for _, address := range addresses {
		amount, err := getBTCAmount(address)
		if err != nil {
			return "", err
		}
		totalAmount += amount
	}

	resultString := fmt.Sprintf("Total BTC amount: %.8f\n", totalAmount)
	return resultString, nil
}
