//go:generate go run marketdata/gen.go

package main

import (
	"flag"
	"fmt"
	"os"
	"pricechecker/bot"
)

var token string

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	// for _, s := range os.Args[1:] {
	// 	item := marketdata.Closest(s)
	// 	fmt.Println(item)
	// }
	b, err := bot.New(token)
	if err != nil {
		fmt.Printf("failed to create bot: %v\n", err)
		os.Exit(1)
	}

	b.Run()
}
