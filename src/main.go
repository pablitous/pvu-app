package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var token string

func main() {
	token = os.Args[1]


	myFarm := farms("")
	someoneFarm := farms("0xe404a13f95d805f2f2158f2b7fcbde042d2167b5")
	//farms = testvars.TestFarms
	//json.Unmarshal([]byte(farms), &myStoredVariable)
	var result map[string]interface{}
	json.Unmarshal([]byte(farms), &result)
	fmt.Println(string(farms))
	//applyTool("612a41891cd86b000992c675")
}

https://backend-farm.plantvsundead.com/farms/6127c50f1a8e8c001cb68ba9
func farms(farmId string) string {
	urlFarms := "https://backend-farm.plantvsundead.com/farms"
	if farmId == ""{
		urlFarms += "other/"+farmId
	}	
	limit := []string{"limit", "10"}
	offset := []string{"offset", "0"}
	header := [][]string{limit, offset}
	farms := api(urlFarms, "GET", token, "", nil, header)
	//fmt.Println(string(myTools))
	return string(farms)
}

func applyTool(farmId string) string {
	urlApplyTool := "https://backend-farm-stg.plantvsundead.com/farms/apply-tool"
	limit := []string{"limit", "10"}
	offset := []string{"offset", "0"}
	header := [][]string{limit, offset}
	applyTool := api(urlApplyTool, "POST", token, `{"farmId":"`+farmId+`","toolId":3,"token":{"challenge":"default","seccode":"default","validate":"default"}}`, header)
	//fmt.Println(string(applyTool))
	return string(applyTool)
}

func buyTools(toolId int, cant int) string {
	urlBuyTools := "https://backend-farm-stg.plantvsundead.com/farms/buy-tools"
	buyTools := api(urlBuyTools, "POST", token, `{"amount":`+strconv.Itoa(cant)+`,"toolId":`+strconv.Itoa(toolId)+`}`, nil)
	//fmt.Println(string(buyTools))
	return string(buyTools)
}

func buySunflowers(toolId int, cant int) string {
	urlBuyTools := "https://backend-farm-stg.plantvsundead.com/farms/buy-sunflowers"
	buyTools := api(urlBuyTools, "POST", token, `{"amount":`+strconv.Itoa(cant)+`,"toolId":`+strconv.Itoa(toolId)+`}`, nil)
	return string(buyTools)
}

func myTools(toolId int, cant int) string {
	urlMyTools := "https://backend-farm-stg.plantvsundead.com/farms/my-tools"
	myTools := api(urlMyTools, "GET", token, "", nil)
	//fmt.Println(string(myTools))
	return string(myTools)
}

func api(url string, method string, token string, rawBody string, headers [][]string) string {

	var jsonData = []byte(rawBody)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if len(headers) != 0 {
		for _, element := range headers {
			request.Header.Set(element[0], element[1])
		}

	}
	request.Header.Set("Authorization", token)
	request.Header.Set("limit", "10")
	request.Header.Set("offset", "0")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("host", "https://marketplace.plantvsundead.com")
	request.Header.Set("referer", "https://marketplace.plantvsundead.com")
	request.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")
	client := &http.Client{}
	response, err := client.Do(request)
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
