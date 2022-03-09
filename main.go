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

var worldListSecondary = []string{
	"account",
	"login",
	"log in",
	"register",
	"signup",
	"sign up",
	"newsletter",
	"news",
	"payment",
	"checkout",
	"product",
	"option",
	"checkout",
	"paypal",
	"braintree",
	"brain tree",
	"credit",
	"gateway",
	"cart",
	"customer",
	"address",
	"password",
	"email",
	"subscription",
	"user",
	"subscriber",
	"stripe",
	"recaptcha",
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
		firstColumn := line[0]
		isHit := false
		for _, word := range wordlist {
			if strings.Contains(strings.ToLower(firstColumn), word) {
				isHit = true
				break
			}
		}
		if isHit {
			for _, word := range worldListSecondary {
				isHit = false
				if strings.Contains(strings.ToLower(firstColumn), word) {
					hitCache = append(hitCache, line)
					isHit = true
					break
				}
			}
		}
		if !isHit {
			missCache = append(missCache, line)
		}

	}

	hitResultFile, err := os.Create("./results/" + filePrefix + "_hits.csv")
	if err != nil {
		panic(err)
	}

	defer hitResultFile.Close()

	missResultFile, err := os.Create("./results/" + filePrefix + "_misses.csv")
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
