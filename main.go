package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func FetchApi(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Cound not create request: %s", err.Error())
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("Could not make request: %s", err.Error())
		return nil, err
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Could not read response body: %s", err.Error())
		return nil, err
	}
	defer response.Body.Close()

	return responseBody, nil
}

func FetchBrasilApi(cep string, ch chan<- string) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	responseBody, err := FetchApi(url)
	if err != nil {
		return
	}

	address := string(responseBody)

	ch <- address
}

func FetchViaCepApi(cep string, ch chan<- string) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	responseBody, err := FetchApi(url)
	if err != nil {
		return
	}

	address := string(responseBody)

	ch <- address
}

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	cep := "01153000"

	go FetchBrasilApi(cep, ch1)
	go FetchViaCepApi(cep, ch2)

	select {
	case address := <-ch1:
		fmt.Printf("Received response from Brasil API: %s\n", address)
	case address := <-ch2:
		fmt.Printf("Received response from Via CEP: %s\n", address)
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}
}
