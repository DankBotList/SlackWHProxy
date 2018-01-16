package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

type Configuration struct {
	ForwardTo map[string][]string `json:"forward_to"`
}

func main() {

	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		data, err := json.MarshalIndent(Configuration{ForwardTo:map[string][]string{"key": {"https://discordapp.com/meme", "https://discordapp.com/dank"}}}, "", "    ")
		if err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile("config.json", data, os.FileMode(0644)); err != nil {
			panic(err)
		}
		fmt.Println("DEFAULT CONFIG SAVED, PLEASE EDIT!")
		fmt.Println("DEFAULT CONFIG SAVED, PLEASE EDIT!")
		os.Exit(1)
		return
	}

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	config := &Configuration{}
	if err := json.Unmarshal(data, config); err != nil {
		panic(err)
	}

	engine := gin.Default()

	for k, v := range config.ForwardTo {
		engine.POST(k, CreateHandler(v))
	}

	engine.Run(":8080")

}

func CreateHandler(recepients []string) func(*gin.Context) {
	return func(ctx *gin.Context) {

		client := http.Client{}

		errorOcurred := false
		var lastResponse *http.Response

		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.Error(err)
			return
		}

		for _, v := range recepients {
			buf := new(bytes.Buffer)
			buf.Write(body)
			req, err := http.NewRequest("POST", v, buf)
			if err != nil {
				fmt.Println(err)
				continue
			}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if resp.StatusCode < 200 || resp.StatusCode > 209 {
				errorOcurred = true
				lastResponse = resp
			} else if !errorOcurred {
				lastResponse = resp
			}

		}

		ctx.Status(lastResponse.StatusCode)
		for headerName, vals := range lastResponse.Header {
			for _, val := range vals {
				ctx.Header(headerName, val)
			}
		}

		contentType := lastResponse.Header.Get("Content-type")
		data, err := ioutil.ReadAll(lastResponse.Body)
		if err != nil {
			ctx.Error(err)
			return
		}
		ctx.Data(lastResponse.StatusCode, contentType, data)

	}
}