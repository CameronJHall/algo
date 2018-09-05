package robinhood

import (
	"encoding/json"
	"net/http"
)

type InstrumentList struct {
	List     []Instrument `json:"results"`
	NextPage string       `json:"next"`
}

type Instrument struct {
	Symbol string `json:"symbol"`
	ID     string `json:"id"`
}

func GetInstrumentID(symbol string, client *http.Client) (id string, err error) {
	instrumentList := InstrumentList{}

	// Create the new request
	req, err := http.NewRequest("GET", "https://api.robinhood.com/instruments/", nil)
	if err != nil {
		return "", err
	}

	query := req.URL.Query()
	query.Add("symbol", symbol)
	req.URL.RawQuery = query.Encode()

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// Decode the body using the initialized quotes struct
	err = json.NewDecoder(resp.Body).Decode(&instrumentList)
	if err != nil {
		return "", err
	}

	return instrumentList.List[0].ID, nil
}
