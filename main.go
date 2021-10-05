package main

import (
	"bufio"
	"bytes"
	"encoding/json"
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

type BasicInfo struct {
	Token   string `json:"token"`
	FarmUrl string `json:"farmUrl"`
}

func main() {
	//token = os.Args[1]
	scanner := bufio.NewScanner(os.Stdin)
	fileBasicInfo, err := os.Open("BasicInfo.json")
	if err != nil {
		fmt.Println("First Login, please complete the parameters.")
		updateParameters()
	} else {
		fmt.Println("Parameters already exists, do you want to change them? (y/n)")
		if scanner.Scan() {
			updateParametersYesNo := scanner.Text()
			if updateParametersYesNo == "y" {
				updateParameters()
			} else {
				basicInfoJson, _ := ioutil.ReadAll(fileBasicInfo)
				var loginData BasicInfo
				json.Unmarshal(basicInfoJson, &loginData)
				token = loginData.Token
				farmUrl = loginData.FarmUrl
			}
		}
	}

	for {
		fmt.Println("Checking " + time.Now().String())
		mainLogic()
		rand.Seed(time.Now().UnixNano())
		n := utils.RandFloats(5, 20)
		s := fmt.Sprintf("Waiting %f minutes to check again", n)

		fmt.Println(s)
		time.Sleep(time.Duration(n) * time.Minute)
	}
}
func updateParameters() bool {
	fmt.Println("Ingrese su token por favor ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		token = scanner.Text()
	}
	fmt.Println("Ingrese el numero de url q se esta utilizando")
	fmt.Println("1  https://backend-farm.plantvsundead.com")
	fmt.Println("2  https://backend-farm-stg.plantvsundead.com")
	var farmUrlId int
	fmt.Scanln(&farmUrlId)
	if farmUrlId == 2 {
		farmUrl = "https://backend-farm-stg.plantvsundead.com"
	} else {
		farmUrl = "https://backend-farm.plantvsundead.com"
	}
	BasicInfo := BasicInfo{
		Token:   token,
		FarmUrl: farmUrl,
	}
	file, _ := json.MarshalIndent(BasicInfo, "", " ")
	_ = ioutil.WriteFile("BasicInfo.json", file, 0644)
	return true
}

func mainLogic() bool {
	farmStatus := farmStatus()
	isMyTurn := gjson.Get(farmStatus, "data.status")
	if isMyTurn.String() == "1" {
		myFarm := farms("")
		//myFarm = testvars.TestFarms
		plantIds := gjson.Get(myFarm, "data.#._id")
		var countPlants int
		plantIds.ForEach(func(key, value gjson.Result) bool {
			plantId := value.String()
			stage := gjson.Get(myFarm, "data."+strconv.Itoa(countPlants)+".stage").String()
			totalHarvest := gjson.Get(myFarm, "data."+strconv.Itoa(countPlants)+".totalHarvest").Int()
			needWater := gjson.Get(myFarm, "data."+strconv.Itoa(countPlants)+".needWater").String()
			if stage == "new" {
				fmt.Println("Plant " + plantId + " is new and a new Pot needs to be added")
				utils.AddRandomSleep(7, 23)
				applyToolSmallPot(plantId)
			}
			fixWater(plantId, needWater, stage)
			fixWater(plantId, needWater, stage)
			hasCrow := gjson.Get(myFarm, "data."+strconv.Itoa(countPlants)+".hasCrow").String()
			fixCrow(plantId, hasCrow)
			isTempPlant := gjson.Get(myFarm, "data."+strconv.Itoa(countPlants)+".isTempPlant").Bool()
			//if stage == "cancelled" && totalHarvest != 0 {
			if totalHarvest != 0 {
				fmt.Println("Plant " + plantId + " needs to be harvested")
				utils.AddRandomSleep(7, 23)
				harvestPlant := harvestPlant(plantId)
				if harvestPlant && stage == "cancelled" && isTempPlant != true {
					fmt.Println("Plant " + plantId + " needs to be removed")
					utils.AddRandomSleep(7, 23)
					removePlant(plantId)
				}
			}
			if stage == "cancelled" && totalHarvest == 0 {
				removePlant(plantId)
			}
			countPlants += 1
			return true // keep iterating
		})
		checkFreeSpotsAndAddNewPlants()
		doWorldTree()
		return true
	} else {
		const (
			RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
		)
		t, _ := time.Parse(RFC3339Nano, gjson.Get(farmStatus, "data.nextGroup").String())
		fmt.Println(farmStatus)
		//fmt.Println("Not your turn, turn at " + t.String())
		turnTime := t.Add(time.Hour * -3).String()
		fmt.Println("Not your turn, turn at " + utils.Substr(turnTime, 0, len(turnTime)-13))
		return true
	}
}

func checkFreeSpotsAndAddNewPlants() bool {
	urlLands := farmUrl + "/my-lands"
	//fmt.Println(urlHarvest)
	limit := []string{"limit", "10"}
	offset := []string{"offset", "0"}
	header := [][]string{limit, offset}
	lands := api(urlLands, "GET", token, "", header)
	//addNewPlant()
	//return lands
	status := gjson.Get(lands, "status").Int()
	landIds := gjson.Get(lands, "data.#._id")
	var countLands int
	if status == 0 {
		landIds.ForEach(func(key, value gjson.Result) bool {
			landId := gjson.Get(lands, "data."+strconv.Itoa(countLands)+".land.landId").String()
			capacityPlant := gjson.Get(lands, "data."+strconv.Itoa(countLands)+".land.capacity.plant").String()
			capacityMotherTree := gjson.Get(lands, "data."+strconv.Itoa(countLands)+".land.capacity.motherTree").String()
			totalFarmingPlant := gjson.Get(lands, "data."+strconv.Itoa(countLands)+".totalFarming.plant").String()
			totalFarmingMotherTree := gjson.Get(lands, "data."+strconv.Itoa(countLands)+".totalFarming.motherTree").String()
			countLands += 1
			if capacityMotherTree > totalFarmingMotherTree {
				addNewPlant(landId, 2)
				fmt.Println("A new mama has been added to land " + landId)
			}
			if capacityPlant > totalFarmingPlant {
				if addNewPlant(landId, 1) {
					fmt.Println("A new sappling has been added to land " + landId)
				}
			}
			//fmt.Println(landId + " - " + capacityPlant + " - " + capacityMotherTree + " - " + totalFarmingPlant + " - " + totalFarmingMotherTree + " - ")
			return true // keep iterating
		})
	}
	return true
}

func addNewPlant(landId string, sunflowerId int) bool {
	//https://backend-farm.plantvsundead.com/farms
	//{"landId": 0,"sunflowerId": 1}
	sunflowerAvailability := getSunfloweAvailability(sunflowerId)
	if !sunflowerAvailability {
		switch sunflowerId {
		case 1:
			buySunflowerSapling(1)
		case 2:
			buySunflowerMama(1)
		}
	}
	urlAddPlant := farmUrl + "/farms"
	payload := `{"landId":` + landId + `,"sunflowerId":` + strconv.Itoa(sunflowerId) + `}`
	addPlant := api(urlAddPlant, "POST", token, payload, nil)
	if gjson.Get(addPlant, "status").Int() == 0 {
		fmt.Println("A new plant has been planted")
		return true
	} else {
		fmt.Println("There was an error planting")
		return false
	}
}

func getSunfloweAvailability(sunflowerId int) bool {
	mySunflowers := getMySunflowers()
	mySunfloweCount := 0
	if int(gjson.Get(mySunflowers, "data.0.sunflowerId").Int()) == sunflowerId {
		mySunfloweCount = int(gjson.Get(mySunflowers, "data.0.usages").Int())
	} else {
		mySunfloweCount = int(gjson.Get(mySunflowers, "data.1.usages").Int())
	}
	if mySunfloweCount >= 1 {
		return true
	} else {
		return false
	}
}

func harvestPlant(plantId string) bool {
	urlHarvest := farmUrl + "/farms/" + plantId + "/harvest"
	//fmt.Println(urlHarvest)
	harvest := api(urlHarvest, "POST", token, "", nil)
	//{"status":0,"data":{"amount":250}}
	status := gjson.Get(harvest, "status").Int()
	amount := gjson.Get(harvest, "data.amount").String()
	if status == 0 {
		fmt.Println("Plant " + plantId + " has been harvested and you get: " + amount + " LE")
		return true
	} else {
		fmt.Println("There was an error harvesting plant " + plantId)
		return false
	}
}

func removePlant(plantId string) bool {
	urlHarvest := farmUrl + "/farms/" + plantId + "/deactivate"
	harvest := api(urlHarvest, "POST", token, "", nil)
	//{"status":0,"data":{"amount":250}}
	status := gjson.Get(harvest, "status").Int()
	if status == 0 {
		fmt.Println("Plant " + plantId + " has been removed.")
	} else {
		fmt.Println("There was an error removing plant " + plantId)
	}
	return true
}

func fixWater(plantId string, needWater string, stage string) bool {
	var message string
	if needWater == "true" && stage != "cancelled" {
		fmt.Println("Plant " + plantId + " needs water")
		utils.AddRandomSleep(7, 23)
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
		utils.AddRandomSleep(8, 23)
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

func hasTool(toolIdCheck int) bool {
	myTools := myTools()
	var tool int64
	myToolsId := gjson.Get(myTools, "data.#.toolId")
	countTool := 0
	myToolsId.ForEach(func(key, value gjson.Result) bool {
		toolId := gjson.Get(myTools, "data."+strconv.Itoa(countTool)+".toolId").Int()
		if int(toolId) == toolIdCheck {
			tool = gjson.Get(myTools, "data."+strconv.Itoa(countTool)+".usages").Int()
		}
		countTool += 1
		return true
	})
	if tool > 0 {
		return true
	} else {
		return false
	}
}

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
	return string(farms)
}

var counter int

func applyTool(farmId string, toolId int, captchaInfo string) bool {
	if counter > 5 {
		counter = 1
	}
	urlApplyTool := farmUrl + "/farms/apply-tool"
	limit := []string{"limit", "10"}
	offset := []string{"offset", "0"}
	header := [][]string{limit, offset}
	tokenCaptcha := `{"challenge":"default","seccode":"default","validate":"default"}`
	if captchaInfo != "" {
		tokenCaptcha = captchaInfo
		fmt.Println("A captcha has been solved")
	}
	//fmt.Println(tokenCaptcha)
	appliedTool := api(urlApplyTool, "POST", token, `{"farmId":"`+farmId+`","toolId":`+strconv.Itoa(toolId)+`,"token":`+tokenCaptcha+`}`, header)
	state := gjson.Get(appliedTool, "status").Int()
	if state == 0 {
		return true
	} else if state == 20 {
		fmt.Println("Plant already has 2 waters")
		return true
	} else {
		if state == 556 {
			captchaSolved := fixCaptcha()
			applyTool(farmId, toolId, captchaSolved)
			//fmt.Println("Solving Captcha")
			return true
		}
		if counter == 5 {
			return false
		} else {
			counter++
			return applyTool(farmId, toolId, "")
		}
	}
}
func fixCaptcha() string {
	var captchaSolution string
	key2Captcha := "0865e8ccc6d546acc46ecff9f1b17ffc"
	urlNewCaptcha := farmUrl + "/captcha/register"
	newCaptcha := api(urlNewCaptcha, "GET", token, "", nil)
	if gjson.Get(newCaptcha, "status").Int() == 0 {
		gt := gjson.Get(newCaptcha, "data.gt").String()
		challenge := gjson.Get(newCaptcha, "data.challenge").String()
		//fmt.Println(gt + " - " + challenge)

		url2Captcha := "http://2captcha.com/in.php?key=" + key2Captcha + "&method=geetest&json=1&gt=" + gt + "&challenge=" + challenge + "&pageurl=" + farmUrl
		solutionId2Captcha := api(url2Captcha, "POST", "", "", nil)
		if gjson.Get(solutionId2Captcha, "status").Int() == 1 {
			requestId := gjson.Get(solutionId2Captcha, "request").String()
			utils.AddRandomSleep(15, 20)
			urlGet2CaptchaSolution := "http://2captcha.com/res.php?key=" + key2Captcha + "&action=get&json=1&id=" + requestId
			get2CaptchaSolution := Get2CaptchaSolution(urlGet2CaptchaSolution)
			challenge = gjson.Get(get2CaptchaSolution, "request.geetest_challenge").String()
			seccode := gjson.Get(get2CaptchaSolution, "request.geetest_seccode").String()
			validate := gjson.Get(get2CaptchaSolution, "request.geetest_validate").String()
			captchaSolution = `{"challenge":"` + challenge + `","seccode":"` + seccode + `","validate":"` + validate + `"}`
			urlValidateCaptcha := farmUrl + "/captcha/validate"
			valitateCaptcha := api(urlValidateCaptcha, "POST", token, captchaSolution, nil)
			utils.AddRandomSleep(0, 1)
			statusValdiateCaptcha := gjson.Get(valitateCaptcha, "status").Int()
			if statusValdiateCaptcha == 0 {
				return captchaSolution
			}
		}
	}

	return captchaSolution
}
func Get2CaptchaSolution(urlGet2CaptchaSolution string) string {
	get2CaptchaSolution := api(urlGet2CaptchaSolution, "POST", "", "", nil)
	request := gjson.Get(get2CaptchaSolution, "request").String()
	if request == "CAPCHA_NOT_READY" {
		utils.AddRandomSleep(4, 7)
		get2CaptchaSolution = Get2CaptchaSolution(urlGet2CaptchaSolution)
	}
	return get2CaptchaSolution
}

func applyToolWater(plantId string) bool {
	hasTool := hasTool(3)
	if !hasTool {
		fmt.Print("Buying Water")
		utils.AddRandomSleep(3, 7)
		buyWater(1)
	}
	utils.AddRandomSleep(3, 7)
	if applyTool(plantId, 3, "") == true {
		fmt.Println("The plant " + plantId + " has been watered")
		return true
	} else {
		fmt.Println("The plant " + plantId + " has't been watered.")
		return false
	}
}

func applyToolScareCrow(plantId string) bool {
	hasTool := hasTool(4)
	if !hasTool {
		fmt.Print("Buying ScareCrow")
		utils.AddRandomSleep(3, 7)
		buyScareCrow(1)
	}
	utils.AddRandomSleep(3, 7)
	if applyTool(plantId, 4, "") == true {
		fmt.Println("The Crow in plant" + plantId + " has been scared")
		return true
	} else {
		fmt.Println("The crown in plant " + plantId + " has't been scared.")
		return false
	}
}

func applyToolSmallPot(plantId string) bool {
	hasTool := hasTool(1)
	if !hasTool {
		fmt.Print("Buying a Small Pot")
		utils.AddRandomSleep(3, 7)
		buySmallPot(1)
	}
	utils.AddRandomSleep(3, 7)
	if applyTool(plantId, 1, "") == true {
		fmt.Println("The small pot has been added to plant" + plantId)
		return true
	} else {
		fmt.Println("There was an error adding the small pot")
		return false
	}
}

func applyToolBigPot(plantId string) bool {
	if applyTool(plantId, 2, "") == true {
		fmt.Println("The big pot has been added to plant" + plantId)
		return true
	} else {
		fmt.Println("There was an error adding the big pot")
		return false
	}
}

func buyTools(toolId int, cant int) string {
	urlBuyTools := farmUrl + "/buy-tools"
	buyTools := api(urlBuyTools, "POST", token, `{"amount":`+strconv.Itoa(cant)+`,"toolId":`+strconv.Itoa(toolId)+`}`, nil)
	//fmt.Println(string(buyTools))
	return string(buyTools)
}
func buyWater(cant int) bool {
	leWallet := getLeWallet()
	if leWallet >= 50*cant {
		buyTools(3, cant)
		return true
	} else {
		return false
	}
}
func buyScareCrow(cant int) bool {
	leWallet := getLeWallet()
	if leWallet >= 20*cant {
		buyTools(4, cant)
		return true
	} else {
		return false
	}
}
func buySmallPot(cant int) bool {
	leWallet := getLeWallet()
	if leWallet >= 100*cant {
		buyTools(1, cant)
		return true
	} else {
		return false
	}
}

func getFarmingStats() string {
	//{"status":0,"data":{"totalHarvestable":0,"pvuToFarm":3000000,"seedsToFarm":22000,"pvuMyFarmed":61.21,"seedsMyFarmed":0,"leWallet":960,"usagesSunflower":85}}
	urlFarmingStats := farmUrl + "/farming-stats"
	farmingStats := api(urlFarmingStats, "POST", token, "", nil)
	//fmt.Println(string(buyTools))
	return string(farmingStats)
}

func getMySunflowers() string {
	urlFarmingStats := farmUrl + "/my-sunflowers"
	farmingStats := api(urlFarmingStats, "POST", token, "", nil)
	//fmt.Println(string(buyTools))
	return string(farmingStats)
}

func getLeWallet() int {
	farmingStats := getFarmingStats()
	leWallet := gjson.Get(farmingStats, "data.leWallet").Int()
	return int(leWallet)
}

func buySunflowers(toolId int, cant int) string {
	urlBuyTools := farmUrl + "/buy-sunflowers"
	buyTools := api(urlBuyTools, "POST", token, `{"amount":`+strconv.Itoa(cant)+`,"toolId":`+strconv.Itoa(toolId)+`}`, nil)
	return string(buyTools)
}

func buySunBox(cant int) bool {
	leWallet := getLeWallet()
	if leWallet >= 100*cant {
		buySunflowers(3, cant)
		return true
	} else {
		return false
	}
}
func buySunflowerMama(cant int) bool {
	leWallet := getLeWallet()
	if leWallet >= 200*cant {
		buySunflowers(2, cant)
		return true
	} else {
		return false
	}
}
func buySunflowerSapling(cant int) bool {
	leWallet := getLeWallet()
	if leWallet >= 100*cant {
		buySunflowers(1, cant)
		return true
	} else {
		return false
	}
}
func myTools() string {
	urlMyTools := farmUrl + "/my-tools"
	myTools := api(urlMyTools, "GET", token, "", nil)
	return string(myTools)
}

func getWorldTreeReward(n int) string {
	urlGetWorldTreeReward := farmUrl + "/world-tree/claim-reward"
	utils.AddRandomSleep(3, 12)
	worldTreeReward := api(urlGetWorldTreeReward, "POST", token, `{"type":`+strconv.Itoa(n)+`}`, nil)
	fmt.Println("Reward " + strconv.Itoa(n) + " has been taken")
	return string(worldTreeReward)
}

func getWorldTreeData() string {
	urlGgetWorldTreeData := farmUrl + "/world-tree/datas"
	worldTreeData := api(urlGgetWorldTreeData, "GET", token, "", nil)
	return string(worldTreeData)
}

func giveWatersWorldTree(n int) bool {
	if hasWatter() {
		urlGiveWatersWorldTree := farmUrl + "/world-tree/give-waters"
		utils.AddRandomSleep(3, 12)
		api(urlGiveWatersWorldTree, "POST", token, `{"amount":`+strconv.Itoa(n)+`}`, nil)
		fmt.Println(strconv.Itoa(n) + " waters were given to the World Tree")
		return true
	} else {
		fmt.Println("Couldn`t water the World Tree, no waters availables")
		return false
	}

}

func getWorldTreeYesterdayReward() string {
	urlWorldTreeYesterdayReward := farmUrl + "/world-tree/claim-yesterday-reward"
	utils.AddRandomSleep(7, 23)
	worldTreeYesterdayReward := api(urlWorldTreeYesterdayReward, "POST", token, "", nil)
	fmt.Println("Reward from yesterday has been taken")
	return string(worldTreeYesterdayReward)
}

func doWorldTree() {
	wolrdTreeData := getWorldTreeData()
	yesterdayReward := gjson.Get(wolrdTreeData, "data.yesterdayReward").Bool()
	if yesterdayReward {
		getWorldTreeYesterdayReward()
	}
	myWater := gjson.Get(wolrdTreeData, "data.myWater").Int()
	if myWater < 20 {
		giveWatersWorldTree(20)
		wolrdTreeData = getWorldTreeData()
	}
	rewardAvailable := gjson.Get(wolrdTreeData, "data.rewardAvailable").Bool()
	totalWatersNow := gjson.Get(wolrdTreeData, "data.totalWater").String()
	if rewardAvailable {
		rewardIds := gjson.Get(wolrdTreeData, "data.reward.#.type")
		rewardIds.ForEach(func(key, value gjson.Result) bool {
			rewardStatus := gjson.Get(wolrdTreeData, "data.reward."+strconv.Itoa(int(value.Int())-1)+".status").String()
			targetWaters := gjson.Get(wolrdTreeData, "data.reward."+strconv.Itoa(int(value.Int())-1)+".target").String()
			if rewardStatus == "finish" {
				getWorldTreeReward(int(value.Int()))
			} else if rewardStatus == "notfinish" {
				fmt.Println("Reward " + value.String() + " has not been finished yet. " + totalWatersNow + "/" + targetWaters)
				return false
			}
			return true
		})

	}
	//gjson.Get(wolrdTreeData, "data.totalWater").String()
	//gjson.Get(wolrdTreeData, "data."+key.String()+".toolId").Int()
}

func api(url string, method string, token string, rawBody string, headers [][]string) string {
	var jsonData = []byte(rawBody)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if len(headers) != 0 {
		for _, element := range headers {
			request.Header.Set(element[0], element[1])
		}
	}
	if token != "" {
		request.Header.Set("Authorization", token)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("host", "https://marketplace.plantvsundead.com")
	request.Header.Set("referer", "https://marketplace.plantvsundead.com")
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("pragma", "no-cache")
	request.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")
	client := &http.Client{}
	response, err := client.Do(request)
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
