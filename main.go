package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/JUNGHUNKIM-7/cli_gpt/model"
	"github.com/JUNGHUNKIM-7/cli_gpt/program"
)

func main() {
	if keys, vals := getToken(); keys == nil || vals == nil {
		log.Fatal("keys||vals are nil")
	} else {
		runProgram(keys, vals)
	}
}

func getToken() (ks, vs []string) {
	var (
		keys []string
		vals []string
	)
	f, err := os.Open("./env.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		key, value := keyValueSplitter(scanner)
		keys = append(keys, key)
		vals = append(vals, value)
	}

	for _, v := range keys {
		if v == "openai" {
			ks = keys
			vs = vals
			return
		}
	}
	ks = nil
	vs = nil
	return
}

func keyValueSplitter(scanner *bufio.Scanner) (string, string) {
	vals := strings.Split(scanner.Text(), "=")
	key, value := vals[0], vals[len(vals)-1]
	return key, value
}

func findValueByKey(keys []string, key string) int {
	for i, v := range keys {
		if v == key {
			return i
		}
	}
	return -1
}

func runProgram(keys, vals []string) {
	if findValueByKey(keys, "openai") == -1 {
		log.Fatal("empty openai key")
	}

	config := &model.GptConfig{
		Temperature: 0.5,
	}
	idx := findValueByKey(keys, "openai")
	defaultEnv := &model.Env{
		GptToken: vals[idx],
	}
	wg := sync.WaitGroup{}
	histories := make([]model.QnaBody, 0)

l1:
	for {
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			q := scanner.Text()

			switch {
			case len(q) > 1:
				if strings.ContainsRune(q, '-') {
					switch flag := strings.Split(q, " "); flag[0][1:] {
					case "h":
						for _, v := range histories {
							printBody(v)
						}
					case "v":
						fmt.Println("get v")
					default:
					}
				} else {
					if !strings.ContainsRune(q, '@') {
						body := program.MakeRequest(q, config, defaultEnv, nil)
						histories = append(histories, body)

						printBody(body)
					} else {
						qs := strings.Split(q, "@")

						wg.Add(len(qs))
						for _, v := range qs {
							go func(str string) {
								body := program.MakeRequest(str, config, defaultEnv, &wg)
								histories = append(histories, body)

								printBody(body)
							}(v)
						}
						wg.Wait()
					}
				}
			case len(q) == 0:
				break l1
			}
		}
	}
}

func printBody(b model.QnaBody) {
	fmt.Printf("Q:%s\nA:%s\n\n", b.Q, b.A)
}
