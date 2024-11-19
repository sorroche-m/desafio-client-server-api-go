package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type USDBRL struct {
	Code       string `json: "code"`
	Codein     string `json: "codein"`
	Name       string `json: "name"`
	High       string `json: "high"`
	Low        string `json: "low"`
	VarBid     string `json: "varBid"`
	PctChange  string `json: "pctChange"`
	Bid        string `json: "bid"`
	Ask        string `json: "ask"`
	Timestamp  string `json: "timestamp"`
	CreateDate string `json: "create_date"`
}

type Response struct {
	USDBRL USDBRL `json: "USDBRL"`
}

type CotacaoDolar struct {
	Dolar float64 `json: "dolar"`
}

func main() {
	http.HandleFunc("/cotacao", cotacaoHandler)
	println("Servidor escutando na porta 8080...")
	http.ListenAndServe(":8080", nil)
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	timeout := 200 * time.Millisecond
	w.Header().Set("Content-Type", "application/json")
	log.Println("Requisição iniciada")
	cotacao, err := cotacaoDolarComTimeout(timeout)
	if err != nil {
		http.Error(w, "Erro ao buscar a cotação do dólar", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf(`{"dolar": %s}`, cotacao.Bid)))
	defer log.Println("Requisição finalizada")
}

func cotacaoDolarComTimeout(timeout time.Duration) (USDBRL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		fmt.Println("Erro ao fazer a requisição:", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Timeout ao fazer a requisição:", err)
			return USDBRL{}, err

		}
		fmt.Println("Erro ao fazer a requisição:", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Erro ao ler a resposta:", err)
	}

	// Parse do JSON para a struct Response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Erro ao fazer unmarshal:", err)
	}

	fmt.Println("Cotação do Dólar:", response.USDBRL.Bid)

	return response.USDBRL, nil

}
