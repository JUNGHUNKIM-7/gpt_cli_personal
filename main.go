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
	runProgram(getToken())
}

func getToken() string {
	var (
		token       string
		// temperature string
	)

	f, err := os.Open("./env.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		key, value := keyValueSplitter(scanner)
		if key != "openai" {
			log.Fatal("not found key \"openai\"")
		}
		token = value
	}
	return token
}

func keyValueSplitter(scanner *bufio.Scanner) (string, string) {
	vals := strings.Split(scanner.Text(), "=")
	key, value := vals[0], vals[len(vals)-1]
	return key, value
}

func runProgram(gptToken string) {
	defaultConfig := &model.GptConfig{
		Temperature: 0.8,
	}

	defaultEnv := &model.Env{
		GptToken: gptToken,
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
						body := program.MakeRequest(q, defaultConfig, defaultEnv, nil)
						histories = append(histories, body)
						printBody(body)
					} else {
						qs := strings.Split(q, "@")
						wg.Add(len(qs))

						for _, v := range qs {
							go func(str string) {
								body := program.MakeRequest(str, defaultConfig, defaultEnv, &wg)
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
