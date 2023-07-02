package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type CoinMarket struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	CurrentPrice float64 `json:"current_price"`
}

func main() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	updateTimer := time.NewTimer(15 * time.Minute)

	for {
		select {
		case <-interrupt:
			fmt.Println("Программа остановлена")
			return
		case <-updateTimer.C:
			fmt.Println("Обновление курсов...")
			markets, err := getCoinMarkets(client)
			if err != nil {
				fmt.Println("Не удалось получить данные:", err)
				return
			}
			fmt.Println("Курсы обновлены")
			displayCoinMarkets(markets)
			updateTimer.Reset(15 * time.Minute)
		default:
			fmt.Println("Введите символ криптовалюты (например, btc):")
			var symbol string
			fmt.Scanln(&symbol)
			market, err := getCoinMarket(client, symbol)
			if err != nil {
				fmt.Println("Не удалось получить данные:", err)
				continue
			}
			if market == nil {
				fmt.Println("Криптовалюта не найдена")
				continue
			}
			displayCoinMarket(*market)
		}
	}
}

func getCoinMarkets(client *http.Client) ([]CoinMarket, error) {
	response, err := client.Get("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1")
	if err != nil {
		return nil, err
	}
	fmt.Println(response)
	defer response.Body.Close()

	var markets []CoinMarket
	err = json.NewDecoder(response.Body).Decode(&markets)
	if err != nil {
		return nil, err
	}

	return markets, nil
}

func getCoinMarket(client *http.Client, symbol string) (*CoinMarket, error) {
	markets, err := getCoinMarkets(client)
	if err != nil {
		return nil, err
	}

	for _, market := range markets {
		if market.Symbol == symbol {
			return &market, nil
		}
	}

	return nil, nil
}

func displayCoinMarkets(markets []CoinMarket) {
	fmt.Println("Список криптовалют:")
	for _, market := range markets {
		fmt.Printf("Название: %s, Символ: %s, Цена: %.2f\n", market.Name, market.Symbol, market.CurrentPrice)
	}
}

func displayCoinMarket(market CoinMarket) {
	fmt.Printf("Название: %s, Символ: %s, Цена: %.2f\n", market.Name, market.Symbol, market.CurrentPrice)
}
