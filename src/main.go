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
	"utils"

	"github.com/tidwall/gjson"
)

var token string
var farmUrl string

func main() {
	token = os.Args[1]
	farmUrl = "https://backend-farm-stg.plantvsundead.com"
	farmStatus := farmStatus()
	isMyTurn := gjson.Get(farmStatus, "data.status")
	if isMyTurn.String() == "1" {

		myFarm := farms("")
		//myFarm = testvars.TestFarms
		plantIds := gjson.Get(myFarm, "data.#._id")
		var countPlants int
		plantIds.ForEach(func(key, value gjson.Result) bool {
			plantId := value.String()
			needWater := gjson.Get(myFarm, "data."+strconv.Itoa(countPlants)+".needWater").String()
			fixWater(plantId, needWater)
			hasCrow := gjson.Get(myFarm, "data."+strconv.Itoa(countPlants)+".hasCrow").String()
			fixCrow(plantId, hasCrow)
			countPlants += 1
			return true // keep iterating
		})
		//idPlant := len(gjson.Get(myFarm, "data.0._id"))
		//fmt.Println(idPlant.String())
	} else {
		const (
			RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
		)
		t, _ := time.Parse(RFC3339Nano, gjson.Get(farmStatus, "data.nextGroup").String())
		turnTime := t.Add(time.Hour * -3).String()
		fmt.Println("Not your turn, turn at " + utils.Substr(turnTime, 0, len(turnTime)-14))
	}

	//applyTool("612a41891cd86b000992c675")
}

func fixWater(plantId string, needWater string) bool {
	var message string
	if needWater == "true" {
		fmt.Println("Plant " + plantId + " needs water")
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(10)
		time.Sleep(time.Duration(n) * time.Second)
		fmt.Printf("Waiting %d seconds...\n", n)
		applyToolWater := applyToolWater(plantId)
		if applyToolWater != true {
			message = "Water has been applied to " + plantId
		}
	} else {
		message = "Plant " + plantId + " doesn´t need water"
	}
	fmt.Println(message)
	return true
}
func fixCrow(plantId string, hasCrow string) bool {
	var message string
	if hasCrow == "true" {
		fmt.Println("Plant " + plantId + " has a crow and needs to be scared")
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(10)
		time.Sleep(time.Duration(n) * time.Second)
		fmt.Printf("Waiting %d seconds...\n", n)
		applyToolScarecrow := applyToolScareCrow(plantId)
		if applyToolScarecrow != true {
			message = "Crow has been scared in " + plantId
		}
	} else {
		message = "Plant " + plantId + " doesn´t have a crow"
	}
	fmt.Println(message)
	return true
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

func applyTool(farmId string, toolId int) bool {
	counter := 1
	urlApplyTool := farmUrl + "/farms/apply-tool"
	limit := []string{"limit", "10"}
	offset := []string{"offset", "0"}
	header := [][]string{limit, offset}
	appliedTool := api(urlApplyTool, "POST", token, `{"farmId":"`+farmId+`","toolId":`+strconv.Itoa(toolId)+`,"token":{"challenge":"default","seccode":"default","validate":"default"}}`, header)
	state := gjson.Get(appliedTool, "status").Int()
	//fmt.Println(state)
	if state == 0 {
		return true
	} else {
		if counter == 5 {
			return false
		} else {
			counter++
			return applyTool(farmId, toolId)
		}

	}

}

func applyToolWater(plantId string) bool {
	if applyTool(plantId, 3) == true {
		fmt.Println("The plant " + plantId + " has been watered")
		return true
	} else {
		fmt.Println("The plant " + plantId + " has't been watered.")
		return false
	}
}

func applyToolScareCrow(plantId string) bool {
	if applyTool(plantId, 4) == true {
		fmt.Println("The Crow in plant" + plantId + " has been scared")
		return true
	} else {
		fmt.Println("The crown in plant " + plantId + " has't been scared.")
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
