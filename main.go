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
		log.Fatal(runProgram(keys, vals))
	}
}

func getToken() (listOfKeys, listOfValues []string) {
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
		key, value := keyValueSplitter(scanner)
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

func runProgram(listOfKeys, listOfValues []string) error {
	if findValueByKey(listOfKeys, "openai") == -1 {
		log.Fatal("empty openai key")
	}
	config := &model.GptConfig{
		Temperature: 0.5,
	}
	idx := findValueByKey(listOfKeys, "openai")
	defaultEnv := &model.Env{
		GptToken: listOfValues[idx],
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
					case "q":
						os.Exit(1)
					default:
						return fmt.Errorf("err: %s", "invalid command")
					}
				} else {
					switch contains := !strings.ContainsRune(q, '@'); contains {
					case true:
						body := program.MakeRequest(q, config, defaultEnv, nil)
						histories = append(histories, body)

						printBody(body)
					case false:
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
				continue l1
			}
		}
	}
}

func printBody(b model.QnaBody) {
	fmt.Printf("-----\nQ:%s\nA:%s\n-----\n", b.Q, b.A)
}
