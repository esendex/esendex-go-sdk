package main

import (
	"encoding/json"
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

const pageSize = 20

func pageOpts(page int) xesende.Option {
	startIndex := (page - 1) * pageSize

	return xesende.Page(startIndex, pageSize)
}

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
		SentCmd(client),
		MessageCmd(client),
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

			data, _ := json.MarshalIndent(resp.Messages, "", "  ")
			fmt.Printf("%s\r\n", data)
		},
	}

	cmd.Flag.IntVar(&page, "page", 0, "")

	return cmd
}

func SentCmd(client *xesende.Client) *hadfield.Command {
	var page int

	cmd := &hadfield.Command{
		Usage: "sent [options]",
		Short: "lists sent messages",
		Long: `
  Sent displays a list of sent messages.

    --page <num>    # Display given page
`,
		Run: func(cmd *hadfield.Command, args []string) {
			resp, err := client.Sent(pageOpts(page))
			if err != nil {
				log.Fatal(err)
			}

			data, _ := json.MarshalIndent(resp.Messages, "", "  ")
			fmt.Printf("%s\r\n", data)
		},
	}

	cmd.Flag.IntVar(&page, "page", 1, "")

	return cmd
}

func MessageCmd(client *xesende.Client) *hadfield.Command {
	return &hadfield.Command{
		Usage: "message MESSAGEID",
		Short: "displays a messag",
		Long: `
  Message displays the details for a message.
`,
		Run: func(cmd *hadfield.Command, args []string) {
			if len(args) < 1 {
				log.Fatal("Require MESSAGEID parameter")
			}

			resp, err := client.Message(args[0])
			if err != nil {
				log.Fatal(err)
			}

			data, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Printf("%s\r\n", data)
		},
	}
}
