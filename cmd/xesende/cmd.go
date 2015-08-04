package main

import (
	"flag"
	"fmt"
	"log"

	"hawx.me/code/hadfield"
	"hawx.me/code/xesende"
)

var (
	accountReference = flag.String("account-reference", "", "")
	username         = flag.String("username", "", "")
	password         = flag.String("password", "", "")
)

var templates = hadfield.Templates{
	Help: `usage: example [command] [arguments]

  This is an example.

  Commands: {{range .}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}
`,
	Command: `usage: example {{.Usage}}
{{.Long}}
`,
}

func main() {
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("Require --username and --password options")
	}

	client := xesende.New(*username, *password)

	commands := hadfield.Commands{
		ReceivedCmd(client),
	}

	hadfield.Run(commands, templates)
}

func ReceivedCmd(client *xesende.Client) *hadfield.Command {
	var page int

	cmd := &hadfield.Command{
		Usage: "received [options]",
		Short: "lists received messages",
		Long: `
  Received displays a list of received messages.

    --page <num>    # Display given page
`,
		Run: func(cmd *hadfield.Command, args []string) {
			resp, err := client.Received()
			if err != nil {
				log.Fatal(err)
			}

			for _, message := range resp.Messages {
				fmt.Printf("At: %s \r\nFrom: %s \r\nBody: %s\r\n", message.ReceivedAt, message.From, message.BodyURI)
			}
		},
	}

	cmd.Flag.IntVar(&page, "page", 0, "")

	return cmd
}
