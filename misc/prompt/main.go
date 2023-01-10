package main

import (
  "bufio"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "net/url"
  "os"
  "strings"

  "github.com/urfave/cli/v2"
)

func SendPrompt(requestURL string, prompt string) string {
  // create new request
  req, err := http.NewRequest(http.MethodGet, requestURL, nil)
  if err != nil {
    log.Fatalf("client: could not create request: %s\n", err)
  }

  // add query string values to request and encode them
  req.URL.RawQuery = url.Values{"prompt": {prompt}}.Encode()

  // send request and store response
  resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error making http request: %s\n", err)
	}

  // get actual chatgpt generated text
  bodyText, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatalf("error trying to read body of response: %s\n", err)
  }

  // make sure to close after function closes
  defer resp.Body.Close()

  // finally return chatgpt response
  return string(bodyText)
}

func main() {
  // setup vars to hold flag values
  var port string
  var address string

  // setup cli app config
  app := &cli.App{
    Name:  "prompt",
    Usage: "send prompts to cligpt server",
    Flags: []cli.Flag{
            &cli.StringFlag{
                Name:  "port",
                Value: "8080",
                Usage: "port of cligpt server",
                Destination: &port,
            },
            &cli.StringFlag{
                Name:  "address",
                Value: "0.0.0.0",
                Usage: "address of cligpt server",
                Destination: &address,
            },
        },
    Action: func(*cli.Context) error {
      // setup chatgpt url
      server := "http://" + address + ":" + port + "/chatgpt"

      // show server setup
      fmt.Printf("Sending request to: %s\n", server)

      // setup stdin reader
      reader := bufio.NewReader(os.Stdin)

      // loop forever on input
      for {
        // print prompt character
        fmt.Print("> ")

        // read the keyboad input.
        input, err := reader.ReadString('\n')
        if err != nil {
          fmt.Fprintln(os.Stderr, err)
        }

        // clean off newline
        input = strings.TrimSuffix(input, "\n")

        // send the prompt text to cligpt server.
        fmt.Println(SendPrompt(server, input))
      }
      return nil
    },
  }

  // begin cli app
  if err := app.Run(os.Args); err != nil {
    log.Fatal(err)
  }
}
