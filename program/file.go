package program

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func GetToken() (listOfKeys, listOfValues []string) {
	var (
		keys []string
		vals []string
	)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := fmt.Sprintf("%s\\env.txt", dir)

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		key, value := KeyValueSplitter(scanner)
		keys = append(keys, key)
		vals = append(vals, value)
	}

	for _, v := range keys {
		if v == "openai" {
			listOfKeys = keys
			listOfValues = vals
			return
		}
	}
	listOfKeys = nil
	listOfValues = nil
	return
}

func KeyValueSplitter(scanner *bufio.Scanner) (string, string) {
	vals := strings.Split(scanner.Text(), "=")
	key, value := vals[0], vals[len(vals)-1]
	return key, value
}

func FindValueByKey(keys, vals []string, key string) string {
	for i, v := range keys {
		if v == key {
			return vals[i]
		}
	}
	return ""
}
