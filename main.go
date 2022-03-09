package main

import (
	"encoding/csv"
	"io/ioutil"
	"os"
	"strings"
)

var wordlist = []string{
	"don't",
	"dont",
	"doesn't",
	"sorry",
	"no",
	"not",
	"please",
	"wrong",
	"incorrect",
	"invalid",
	"success",
	"fail",
	"failure",
	"unsuccessful",
	"success",
	"successful",
	"attempt",
	"valid",
	"cannot",
	"can't",
	"couldn't",
	"could not",
	"cant",
	"wont",
	"won't",
	"will not",
	"approved",
	"declined",
	"rejected",
	"accepted",
	"must",
	"mustn't",
	"must not",
	"updated",
	"correct",
	"try again",
	"never",
	"error",
	"permissions",
	"denied",
	"unable",
	"disallowed",
	"congratulations",
	"approved",
}

var worldListSecondary = map[string]string{
	"account":      "Account",
	"purchase":     "Checkout",
	"sku":          "Catalog",
	"wish list":    "Wish List",
	"login":        "Log in",
	"log in":       "Log in",
	"register":     "Log in",
	"signup":       "Log in",
	"sign up":      "Log in",
	"newsletter":   "Newsletter",
	"news":         "Newsletter",
	"payment":      "Checkout",
	"checkout":     "Checkout",
	"product":      "Catalog",
	"option":       "Catalog",
	"paypal":       "Checkout",
	"braintree":    "Checkout",
	"brain tree":   "Checkout",
	"credit":       "Checkout",
	"gateway":      "Checkout",
	"cart":         "Checkout",
	"customer":     "Checkout",
	"address":      "Checkout",
	"password":     "Account",
	"email":        "Account",
	"subscription": "Account",
	"user":         "Account",
	"subscriber":   "Account",
	"stripe":       "Checkout",
	"recaptcha":    "Log in",
}

var blackList = []string{
	"index",
	"indexer",
	"parent",
	"child",
	"image",
	"schema",
	"simple product",
	"configurable product",
	"entity",
	"model",
	"magento",
	"not exist",
	"load",
	"tokens",
	"template",
	"block",
	"widget",
	"site url",
	"store url",
	"website url",
	"creditmemo",
	"memo",
	"exception",
	"grid",
	"attribute",
	"token",
	"api",
	"configuration",
	"cron",
	"system",
	"module",
}

func main() {
	// loop through all files in input dir
	fileInfo, err := ioutil.ReadDir("./inputs")

	if err != nil {
		panic(err)
	}

	for _, file := range fileInfo {
		fname := file.Name()

		csvFile, err := os.Open("./inputs/" + fname)

		if err != nil {
			panic(err)
		}

		reader := csv.NewReader(csvFile)
		reader.LazyQuotes = true

		records, err := reader.ReadAll()
		if err != nil {
			panic(err)
		}

		createSortedFiles(records, fname)
	}
}

func createSortedFiles(records [][]string, filePrefix string) {
	hitCache := make([][]string, 0, len(records))
	missCache := make([][]string, 0, len(records))

	for _, line := range records {
		isHit := false
		isBlackListed := false
		firstColumn := line[0]
		for _, word := range blackList {
			if strings.Contains(strings.ToLower(firstColumn), word) {
				isBlackListed = true
				break
			}
		}
		if isBlackListed {
			continue
		}
		for _, word := range wordlist {
			if strings.Contains(strings.ToLower(firstColumn), word) {
				isHit = true
				break
			}
		}
		if isHit {
			for word, associatedWord := range worldListSecondary {
				isHit = false
				nl := append([]string{associatedWord}, line...)
				if strings.Contains(strings.ToLower(firstColumn), word) {
					hitCache = append(hitCache, nl)
					isHit = true
					break
				}
			}
		}
		if !isHit {
			missCache = append(missCache, line)
		}

	}

	hitResultFile, err := os.Create("./hits/" + filePrefix + "_hits.csv")
	if err != nil {
		panic(err)
	}

	defer hitResultFile.Close()

	missResultFile, err := os.Create("./misses/" + filePrefix + "_misses.csv")
	if err != nil {
		panic(err)
	}

	defer missResultFile.Close()

	hitWriter := csv.NewWriter(hitResultFile)
	missWriter := csv.NewWriter(missResultFile)
	if err = hitWriter.WriteAll(hitCache); err != nil {
		panic(err)
	}
	if err = missWriter.WriteAll(missCache); err != nil {
		panic(err)
	}
}
