package marketdata

import (
	"bytes"
	"time"
)

type Time struct {
	time time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	t.time, _ = time.Parse("2006-01-02T15:04:05", string(b))
	return nil
}

func (t Time) String() string {
	return t.time.Format(time.RFC3339)
}

// MarketResponse contains price information for an item within some City
type MarketResponse struct {
	ItemID           string `json:"item_id"`
	City             string `json:"city"`
	Quality          int32  `json:"quality"`
	SellPriceMin     int64  `json:"sell_price_min"`
	SellPriceMinDate Time   `json:"sell_price_min_date"`
	SellPriceMax     int64  `json:"sell_price_max"`
	SellPriceMaxDate Time   `json:"sell_price_max_date"`
	BuyPriceMin      int64  `json:"buy_price_min"`
	BuyPriceMinDate  Time   `json:"buy_price_min_date"`
	BuyPriceMax      int64  `json:"buy_price_max"`
	BuyPriceMaxDate  Time   `json:"buy_price_max_date"`
}
