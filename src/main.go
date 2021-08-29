package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testvars"
)

var token string

func main() {
	token = os.Args[1]
	urlFarms := "https://backend-farm.plantvsundead.com/farms"
	farms := api(urlFarms, token, "")
	farms = testvars.TestFarms
	//json.Unmarshal([]byte(farms), &myStoredVariable)
	var result map[string]interface{}
	json.Unmarshal([]byte(farms), &result)

	fmt.Println(string(farms))
	//applyTool("612a41891cd86b000992c675")
}

func applyTool(farmId string) int {
	urlApplyTool := "https://backend-farm-stg.plantvsundead.com/farms/apply-tool"
	applyTool := api(urlApplyTool, token, `{"farmId":"`+farmId+`","toolId":3,"token":{"challenge":"default","seccode":"default","validate":"default"}}`)
	fmt.Println(string(applyTool))
	return 200
}

func api(url string, token string, rawBody string) string {
	var jsonData = []byte(rawBody)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	request.Header.Set("Authorization", token)
	request.Header.Set("limit", "10")
	request.Header.Set("offset", "0")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("host", "https://marketplace.plantvsundead.com")
	//request.Header.Set("referer", "https://marketplace.plantvsundead.com")
	request.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")
	client := &http.Client{}
	response, err := client.Do(request)
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
