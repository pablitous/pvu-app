package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println("Calling API...")
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://backend-farm.plantvsundead.com/farms?limit=10&offset=0", nil)
	if err != nil {
		fmt.Print(err.Error())
	}
	req.Header.Add("Authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwdWJsaWNBZGRyZXNzIjoiMHhiOGNjMjQ1M2FhMmE0ZjdhYjYwMzc4NTRhNGRmY2IxMjMxYzE0NmJjIiwibG9naW5UaW1lIjoxNjMwMDExNDA1NjYyLCJjcmVhdGVEYXRlIjoiMjAyMS0wOC0xMyAxNTowODoxNiIsImlhdCI6MTYzMDAxMTQwNX0.R-8CXgdxZpdj6U-GqXE0gY4dqm_bpPY5qZFfzjIiJwU")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Printf(string(body))
}
