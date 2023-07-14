package program

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/JUNGHUNKIM-7/cli_gpt/model"
	"github.com/fatih/color"
)

const (
	mod = "gpt-4"
	url = "https://api.openai.com/v1/chat/completions"
)

func GetCompletion(sysRole, cliRole model.Roles, config *model.GptConfig, envConfig *model.Env) (result string) {
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	var requestBody model.ReqBody

	if config != nil {
		requestBody = model.ReqBody{
			Model:            mod,
			Messages:         *RolesNew(sysRole, cliRole),
			Temperature:      config.Temperature,
			TopP:             config.TopP,
			N:                config.N,
			FrequencyPenalty: config.FrequencyPenalty,
			PresencePenalty:  config.PresencePenalty,
			MaxTokens:        config.MaxTokens,
		}
	} else {
		requestBody = model.ReqBody{
			Model:    mod,
			Messages: *RolesNew(sysRole, cliRole),
		}
	}

	jsonBody, err := json.Marshal(requestBody)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	if envConfig != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", envConfig.GptToken))
	} else {
		log.Fatal("env not initalized")
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	var resbody model.ResBody
	json.Unmarshal(body, &resbody)

	result = resbody.Choices[0].Message.Content
	return
}

func RolesNew(sysRole, cliRole model.Roles) *[]model.Roles {
	return &[]model.Roles{
		{
			Role:    sysRole.Role,
			Content: sysRole.Content,
		},
		{
			Role:    cliRole.Role,
			Content: cliRole.Content,
		},
	}
}

func MakeRequest(qs string, defaultConfig *model.GptConfig, defaultEnv *model.Env, wg *sync.WaitGroup) model.QnaBody {
	if wg != nil {
		defer wg.Done()
	}
	a := GetCompletion(
		model.Roles{
			Role:    "system",
			Content: "You are a helpful assistant.",
		}, model.Roles{
			Role:    "user",
			Content: any(qs).(string),
		},
		defaultConfig,
		defaultEnv,
	)
	return model.QnaBody{
		Q: qs,
		A: a,
	}
}

func PrintBody(i int, b model.QnaBody, answerFunc func(format string, a ...interface{})) {
	fmt.Printf("Q[%s]: %s\n", color.CyanString(strconv.Itoa(i)), color.CyanString(b.Q))
	answerFunc("A: %s\n", b.A)
}
