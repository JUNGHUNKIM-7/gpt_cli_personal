package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/JUNGHUNKIM-7/cli_gpt/model"
	"github.com/JUNGHUNKIM-7/cli_gpt/program"
	"github.com/fatih/color"
)

func main() {
	if keys, vals := program.GetToken(); keys == nil || vals == nil {
		log.Fatal("keys||vals are nil")
	} else {
		log.Fatal(runProgram(keys, vals))
	}
}

func runProgram(listOfKeys, listOfValues []string) error {
	if program.FindValueByKey(listOfKeys, listOfValues, "openai") == "" {
		log.Fatal("empty openai key")
	}
	config := &model.GptConfig{
		Temperature: 0.5,
	}
	token := program.FindValueByKey(listOfKeys, listOfValues, "openai")
	defaultEnv := &model.Env{
		GptToken: token,
	}
	wg := sync.WaitGroup{}
	histories := make([]model.QnaBody, 0)

	//color theme
	y := color.New(color.FgYellow)
	e := color.New(color.FgRed).Add(color.Italic)
	g := color.New(color.FgGreen)
	eln := e.PrintlnFunc()
	errf := e.PrintfFunc()
	answerf := g.PrintfFunc()

l1:
	for {
		y.Printf("Message? |> ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			q := scanner.Text()

			switch {
			case len(q) > 1:
				if strings.ContainsRune(q, '-') {
					switch flag := strings.Split(q, " "); flag[0][1:] {
					case "h":
						if len(histories) < 1 {
							eln("Err: No Histories")
							continue l1
						}
						for i, v := range histories {
							program.PrintBody(i, v, answerf)
						}
					case "r":
						param := strings.Split(q, " ")[1]
						idx, err := strconv.Atoi(param)
						if err != nil {
							errf("%s can't be converted to index\n", param)
							continue l1
						}
						if len(histories) < 1 {
							eln("Err: No Histories")
							continue l1
						}
						histories = append(histories[:idx], histories[idx+1:]...)
						for i, v := range histories {
							program.PrintBody(i, v, answerf)
						}
					case "q":
						os.Exit(1)
					default:
						errf("Err: %s\n", "Invalid Command")
					}
				} else {
					switch contains := !strings.ContainsRune(q, '@'); contains {
					case true:
						body := program.MakeRequest(q, config, defaultEnv, nil)
						histories = append(histories, body)

						program.PrintBody(-1, body, answerf)
					case false:
						qs := strings.Split(q, "@")

						wg.Add(len(qs))
						for i, v := range qs {
							go func(idx int, str string) {
								body := program.MakeRequest(str, config, defaultEnv, &wg)
								histories = append(histories, body)

								program.PrintBody(idx, body, answerf)
							}(i, v)
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
