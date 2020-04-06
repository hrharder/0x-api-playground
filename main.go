package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"

	"github.com/hrharder/0x-api-playground/client"
)

func main() {
	base := flag.String("raw-url", "https://kovan.api.0x.org", "set the base 0x API URL")
	flag.Parse()

	cl, err := client.New(*base)
	if err != nil {
		log.Fatal(err)
	}

	res, err := cl.Quote("0x1FcAf05ABa8c7062D6F08E25c77Bf3746fCe5433", "0x48178164eB4769BB919414Adc980b659a634703E", big.NewInt(100000), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.Price)
}
