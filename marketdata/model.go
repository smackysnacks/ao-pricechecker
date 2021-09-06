package marketdata

import (
	"bytes"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Time wraps a `time.Time` structure for dates returned by the albion online
// data project
type Time struct {
	time time.Time
}

// UnmarshalJSON parses a `Time` from an arbitrary byte sequence. Error is
// always nil, defaulting a `time.Time` when parsing fails
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

// MarketResponses represents a slice of `MarketResponse` that provides a
// `String()` for formatting results in a nice table format
type MarketResponses []*MarketResponse

func (r *MarketResponses) Table() string {
	var data [][]string
	var buf bytes.Buffer

	for _, mr := range *r {
		data = append(data, []string{mr.City, strconv.FormatInt(int64(mr.Quality), 10),
			strconv.FormatInt(mr.SellPriceMin, 10), strconv.FormatInt(mr.SellPriceMax, 10),
			strconv.FormatInt(mr.BuyPriceMin, 10), strconv.FormatInt(mr.BuyPriceMax, 10)})
	}

	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{"City", "Quality", "Sell Min", "Sell Max", "Buy Min", "Buy Max"})
	table.AppendBulk(data)
	table.Render()

	return buf.String()
}

func (r *MarketResponses) Tables(maxlength int) []string {
	panic("unimplemented")
}

func (r *MarketResponses) ShortSummary() string {
	panic("unimplemented")
}
