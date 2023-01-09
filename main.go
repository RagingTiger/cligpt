package main

import (
  "encoding/json"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "strings"

  "github.com/RagingTiger/cligpt/cli"
)

// setting up cli args
var server = flag.Bool("server", false, "run ChatGPT CLI server")

// custom data types for sending to ChatGPT API
type Data struct {
  ID      string                 `json:"id"`
  Object  string                 `json:"object"`
  Created int                    `json:"created"`
  Model   string                 `json:"model"`
  Choices []TextCompletionChoice `json:"choices"`
  Usage   TextCompletionUsage    `json:"usage"`
}

type TextCompletionChoice struct {
  Text         string  `json:"text"`
  Index        int     `json:"index"`
  LogProbs     *string `json:"logprobs"`
  FinishReason string  `json:"finish_reason"`
}

type TextCompletionUsage struct {
  PromptTokens     int `json:"prompt_tokens"`
  CompletionTokens int `json:"completion_tokens"`
  TotalTokens      int `json:"total_tokens"`
}

func HealthStatus(w http.ResponseWriter, req *http.Request) {
  // set status to 200
  w.WriteHeader(http.StatusOK)

  // write response
  w.Write([]byte("Status 200: Server Accessible\n"))
}

func CallAPI(prompt string, auth string, model string, max_tokens string) string {
  // setting up http client to send request to API
  client := &http.Client{}

  // setting up payload that will be sent via POST
  var data = strings.NewReader(`{
      "model": "` + model + `",
      "prompt": "` + prompt + `",
      "temperature": 0.7,
      "max_tokens": ` + max_tokens + `,
      "top_p": 1,
      "frequency_penalty": 0,
      "presence_penalty": 0
    }`)

  // create POST request with payload
  req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", data)
  if err != nil {
    fmt.Println(err, req)
  }
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Authorization", `Bearer `+auth+``)

  // send request
  resp, err := client.Do(req)
  if err != nil {
    log.Fatal("The token is valid?", err)
  }

  // check if the token is valid.
  if resp.StatusCode == 401 {
    fmt.Println("The token is invalid")
    os.Exit(0)
  }

  // get actual ChatGPT response
  defer resp.Body.Close()
  bodyText, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Println(err)
  }
  var response Data
  json.Unmarshal(bodyText, &response)
  if err != nil {
    log.Println(err)
  }

  // select first response and return its text
  choice := response.Choices[0]
  text := choice.Text
  return text
}

func main() {
  // parse cli
  flag.Parse()

  // create config dir
  var err error
  var config map[string]string
  if err = cli.CreateConfigDirectory(); err != nil {
    log.Fatal(err)
  }

  // get config file
  if config, err = cli.ReadYml(); err != nil {
    log.Fatal(err)
  }

  // check token is correct length
  if len(config["auth"]) < 51 {
    log.Fatal("Ensure to insert a valid token in cligpt.yml file.")
  }

  // decide execution mode
  if *server {
    // checking for addtional arguments that will be ignored
    if len(flag.Args()) > 0 {
      log.Printf("Additional positional arguments ignored in server mode.")
    }

    // notify of server startup
    log.Printf("Starting up ChatGPT server.")

    // setting up handler
    ChatGPTInterface := func(w http.ResponseWriter, req *http.Request) {
      // parse form for query args
      if err = req.ParseForm(); err != nil {
        log.Printf("Error parsing form: %s", err)
      }

      // get query arg for prompt
      var prompt = req.Form.Get("prompt")

      // call ChatGPT API with prompt and config vars
      text := CallAPI(
        prompt,
        config["auth"],
        config["model"],
        config["max_tokens"],
      )

      // log prompt and response text
      log.Printf("Prompt: %s", prompt)
      log.Printf("Response: %s\n\n", text)
      log.Printf("End")

      // respond with ChatGPT output
      fmt.Fprintf(w, "%s\n\n", text)
    }

    // register handlers
    http.HandleFunc("/", HealthStatus)
    http.HandleFunc("/health", HealthStatus)
    http.HandleFunc("/chatgpt", ChatGPTInterface)

    // start up server
    log.Fatal(http.ListenAndServe(":80", nil))

  } else {
    // check prompt is submitted
    if len(flag.Args()) == 0 {
      log.Fatal("No prompt submitted.")
    }

    // warn if more than one CLI arg submitted
    if len(flag.Args()) > 1 {
      log.Printf("Ignoring additional arguments: %s", flag.Args()[1:])
    }

    // call ChatGPT API with prompt and config vars
    text := CallAPI(
      flag.Args()[0],
      config["auth"],
      config["model"],
      config["max_tokens"],
    )

    // print out response text
    fmt.Println(text)
  }
}
