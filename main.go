package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sort"
)

// All numbers are float64 so maths is simpler
type item struct {
	Name              string
	BuyPrice          float64
	SellPrice         float64
	BuyOrders         float64
	SellOffers        float64
	DemandSupplyRatio float64
	Margin            float64
	Roi               float64
	Score             float64
}

func main() {
	// Request data
	response, err := http.Get("https://crossoutdb.com/api/v1/items")
	if err != nil {
		log.Panic(err)
	}

	// Read data
	rawData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panic(err)
	}

	// Unmarshal data
	var items []item
	err = json.Unmarshal(rawData, &items)
	if err != nil {
		log.Panic(err)
	}

	// Correct item values
	for i := range items {
		item := &items[i]
		item.BuyPrice /= 100
		item.SellPrice /= 100
		item.Margin /= 100
	}

	// Filter useless items
	for i := len(items) - 1; i >= 0; i-- {
		item := &items[i]
		if item.Margin < 0 ||
			item.BuyPrice > 10000 ||
			item.DemandSupplyRatio == 0 {
			length := len(items)
			*item = items[length-1]
			items = items[:length-1]
		}
	}

	// Score items
	for i := range items {
		item := &items[i]
		item.Score = math.Pow(item.Margin, 2) *
			item.Roi *
			(1 / item.DemandSupplyRatio) *
			item.BuyOrders *
			item.SellOffers
		// item.Score = 1 / item.DemandSupplyRatio * item.Margin
	}

	// Sort by score
	sort.SliceStable(items, func(a, b int) bool {
		return items[a].Score < items[b].Score
	})

	// Print
	for _, item := range items {
		json, err := json.MarshalIndent(item, "", "\t")
		if err != nil {
			log.Panic(err)
		}
		log.Print(string(json))
	}
}
