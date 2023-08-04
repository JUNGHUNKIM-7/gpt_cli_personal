package main

import (
	"bufio"
	"context"
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
		uri := program.FindValueByKey(keys, vals, "mongouri")
		if uri == "" {
			log.Fatal("mongouri isn't provided")
		}
		program.Initialize(uri)
		log.Fatal(runProgram(keys, vals))
	}
}

func runProgram(listOfKeys, listOfValues []string) error {
	defer func() {
		if err := program.Client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if program.FindValueByKey(listOfKeys, listOfValues, "openai") == "" {
		log.Fatal("empty openai key")
	}
	config := &model.GptConfig{
		Temperature: 0.7,
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
	c := color.New(color.FgCyan)
	cln := c.PrintlnFunc()
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
						if isNoHistories(histories, cln) {
							continue l1
						}
						for i, v := range histories {
							program.PrintBody(i, v, answerf)
						}
					case "r":
						if searchParam, ok := checkParam(q, eln); !ok {
							continue l1
						} else {
							idx, err := convertSearchParam(searchParam, errf)
							if err != nil {
								continue l1
							}
							if isNoHistories(histories, cln) {
								continue l1
							}
							if checkNegativeValue(idx, histories, eln) {
								continue l1
							}
							histories = append(histories[:idx], histories[idx+1:]...)
							for i, v := range histories {
								program.PrintBody(i, v, answerf)
							}
						}
					case "g":
						param := strings.Split(q, " ")[1:]
						if len(param) == 0 {
							eln("Search param hasn't not provided")
							continue l1
						}
						searchParam := strings.Join(param, " ")
						historyFrom := program.GetData(searchParam)
						for i, v := range historyFrom {
							program.PrintBody(i, *v, answerf)
						}
					case "s":
						if searchParam, ok := checkParam(q, eln); !ok {
							continue l1
						} else {
							if _, err := strconv.Atoi(searchParam); err != nil {
								eln("Can't be converted to index")
								continue l1
							}
							idx, err := convertSearchParam(searchParam, errf)
							if err != nil {
								continue l1
							}
							if isNoHistories(histories, cln) {
								continue l1
							}
							if checkNegativeValue(idx, histories, eln) {
								continue l1
							}
							body := histories[idx]
							program.SetData(&body)
							histories = append(histories[:idx], histories[idx+1:]...)
						}
					case "a":
						if isNoHistories(histories, cln) {
							continue l1
						}
						program.SetAll(histories)
						histories = make([]model.QnaBody, 0)
					case "q":
						os.Exit(1)
					default:
						eln("Invalid Command")
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

func convertSearchParam(searchParam string, errf func(format string, a ...interface{})) (idx int, err error) {
	idx, err = strconv.Atoi(searchParam)
	if err != nil {
		errf("%s can't be converted to index\n", searchParam)
		return
	}
	return
}

func isNoHistories(histories []model.QnaBody, cln func(a ...interface{})) bool {
	if len(histories) < 1 {
		cln("No Histories")
		return true
	}
	return false
}

func checkParam(q string, eln func(a ...interface{})) (string, bool) {
	param := strings.Split(q, " ")[1:]
	if len(param) == 0 {
		eln("Search param hasn't not provided")
		return "", false
	}
	searchParam := param[0]
	if searchParam == "" {
		eln("Index hasn't not provided")
		return "", false
	}
	return searchParam, true
}

func checkNegativeValue(s int, histories []model.QnaBody, eln func(a ...interface{})) bool {
	if s < 0 || len(histories) <= s {
		eln("Index can't be negative or exceeded than histories")
		return true
	}
	return false
}
