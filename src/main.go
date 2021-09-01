package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tidwall/gjson"
)

var token string
var farmUrl string

func main() {
	token = os.Args[1]
	farmUrl = "https://backend-farm-stg.plantvsundead.com"
	farmStatus := farmStatus()
	isMyTurn := gjson.Get(farmStatus, "data.status")
	if isMyTurn.String() != "1" {
		myFarm := farms("")
		//myFarm = testvars.TestFarms
		plantIds := gjson.Get(myFarm, "data.#._id")
		var countPlants int
		plantIds.ForEach(func(key, value gjson.Result) bool {
			plantId := value.String()
			needWater := gjson.Get(myFarm, "data."+strconv.Itoa(countPlants)+".needWater").String()
			if needWater == "true" {
				fmt.Println("Plant " + plantId + " needs water")
				rand.Seed(time.Now().UnixNano())
				n := rand.Intn(10)
				time.Sleep(time.Duration(n) * time.Second)
				fmt.Printf("Waiting %d seconds...\n", n)
				applyTool := applyTool(plantId)
				if applyTool != true {
					fmt.Println("Water has been applied to " + plantId)
				}
			} else {
				fmt.Println("Plant " + plantId + " doesn´t need water")
			}
			countPlants += 1
			return true // keep iterating

		})

		//idPlant := len(gjson.Get(myFarm, "data.0._id"))
		//fmt.Println(idPlant.String())
	}

	//applyTool("612a41891cd86b000992c675")
}

func hasWatter() bool {
	myTools := myTools()
	//myTools = testvars.TestTools
	var waters int64
	myToolsId := gjson.Get(myTools, "data.#.toolId")
	myToolsId.ForEach(func(key, value gjson.Result) bool {
		toolId := gjson.Get(myTools, "data."+key.String()+".toolId").Int()
		if toolId == 3 {
			waters = gjson.Get(myTools, "data."+key.String()+".usages").Int()
		}
		return true
	})
	if waters > 0 {
		return true
	} else {
		return false
	}
}

//https://backend-farm.plantvsundead.com/farms/6127c50f1a8e8c001cb68ba9
func farmStatus() string {
	urlFarms := farmUrl + "/farm-status"
	farms := api(urlFarms, "GET", token, "", nil)
	return string(farms)
}
func farms(farmId string) string {
	urlFarms := farmUrl + "/farms"
	if farmId != "" {
		urlFarms += "/other/" + farmId
	}
	limit := []string{"limit", "10"}
	offset := []string{"offset", "0"}
	header := [][]string{limit, offset}
	farms := api(urlFarms, "GET", token, "", header)
	//fmt.Println(string(myTools))
	return string(farms)
}

func applyTool(farmId string) bool {
	urlApplyTool := farmUrl + "/farms/apply-tool"
	limit := []string{"limit", "10"}
	offset := []string{"offset", "0"}
	header := [][]string{limit, offset}
	applyTool := api(urlApplyTool, "POST", token, `{"farmId":"`+farmId+`","toolId":3,"token":{"challenge":"default","seccode":"default","validate":"default"}}`, header)
	//fmt.Println(string(applyTool))
	//return string(applyTool)
	state := gjson.Get(applyTool, "status").Int()
	if state == 200 {
		return true
	} else {
		return false
	}

}

func buyTools(toolId int, cant int) string {
	urlBuyTools := farmUrl + "/farms/buy-tools"
	buyTools := api(urlBuyTools, "POST", token, `{"amount":`+strconv.Itoa(cant)+`,"toolId":`+strconv.Itoa(toolId)+`}`, nil)
	//fmt.Println(string(buyTools))
	return string(buyTools)
}

func buySunflowers(toolId int, cant int) string {
	urlBuyTools := farmUrl + "/farms/buy-sunflowers"
	buyTools := api(urlBuyTools, "POST", token, `{"amount":`+strconv.Itoa(cant)+`,"toolId":`+strconv.Itoa(toolId)+`}`, nil)
	return string(buyTools)
}

func myTools() string {
	urlMyTools := farmUrl + "/farms/my-tools"
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
	//request.Header.Set("limit", "10")
	//request.Header.Set("offset", "0")
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
