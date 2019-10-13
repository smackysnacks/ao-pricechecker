package marketdata

import (
	"encoding/json"
	"net/http"
	"strings"
)

const pricesBaseURL = "https://www.albion-online-data.com/api/v2/stats/Prices/"

// Query returns a list of `MarketResponse` prices
func Query(items []string) (res []*MarketResponse, err error) {
	url := pricesBaseURL + strings.Join(items, ",")

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&res)
	return res, err
}
