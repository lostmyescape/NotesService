package quotes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type QuoteResponse struct {
	Quote struct {
		Body string `json:"body"`
	} `json:"quote"`
}

// GetQuote получает ответ от внешнего апи
func GetQuote() (string, error) {
	// получаем ответ
	resp, err := http.Get("https://favqs.com/api/qotd")
	if err != nil {
		return "", fmt.Errorf("failed to fetch quote: %w", err)
	}
	defer resp.Body.Close()

	// проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result QuoteResponse

	// считываем данные из json и записываем их в result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode JSON: %w", err)
	}

	return result.Quote.Body, nil
}
