package main

import (
	"flag"
	"log"

	"github.com/esendex/esendex-go-sdk"
	"github.com/gobs/pretty"
	"hawx.me/code/hadfield"
)

var (
	accountReference = flag.String("account-reference", "", "")
	username         = flag.String("username", "", "")
	password         = flag.String("password", "", "")
)

const pageSize = 20

func pageOpts(page int) esendex.Option {
	startIndex := (page - 1) * pageSize

	return esendex.Page(startIndex, pageSize)
}

var templates = hadfield.Templates{
	Help: `usage: esendex [command] [arguments]

  A command line client for the Esendex REST API.

  Options:
    --username USER    # Username to authenticate with
    --password PASS    # Password to authenticate with
    --help             # Display this message

  Commands: {{range .}}
    {{.Name | printf "%-15s"}} # {{.Short}}{{end}}
`,
	Command: `usage: esendex {{.Usage}}
{{.Long}}
`,
}

func main() {
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("Both --username and --password options are required.")
	}

	client := esendex.New(*username, *password)

	commands := hadfield.Commands{
		receivedCmd(client),
		sentCmd(client),
		messageCmd(client),
		accountsCmd(client),
	}

	hadfield.Run(commands, templates)
}

func receivedCmd(client *esendex.Client) *hadfield.Command {
	var page int

	cmd := &hadfield.Command{
		Usage: "received [options]",
		Short: "lists received messages",
		Long: `
  Received displays a list of received messages.

    --page NUM       # Display given page
`,
		Run: func(cmd *hadfield.Command, args []string) {
			resp, err := client.Received()
			if err != nil {
				log.Fatal(err)
			}

			pretty.PrettyPrint(resp.Messages)
		},
	}

	cmd.Flag.IntVar(&page, "page", 0, "")

	return cmd
}

func sentCmd(client *esendex.Client) *hadfield.Command {
	var page int

	cmd := &hadfield.Command{
		Usage: "sent [options]",
		Short: "lists sent messages",
		Long: `
  Sent displays a list of sent messages.

    --page NUM       # Display given page
`,
		Run: func(cmd *hadfield.Command, args []string) {
			resp, err := client.Sent(pageOpts(page))
			if err != nil {
				log.Fatal(err)
			}

			pretty.PrettyPrint(resp.Messages)
		},
	}

	cmd.Flag.IntVar(&page, "page", 1, "")

	return cmd
}

func messageCmd(client *esendex.Client) *hadfield.Command {
	return &hadfield.Command{
		Usage: "message MESSAGEID",
		Short: "displays a message",
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

			pretty.PrettyPrint(resp)
		},
	}
}

func accountsCmd(client *esendex.Client) *hadfield.Command {
	return &hadfield.Command{
		Usage: "accounts",
		Short: "list accounts",
		Long: `
  List accounts available to the user.
`,
		Run: func(cmd *hadfield.Command, args []string) {
			resp, err := client.Accounts()
			if err != nil {
				log.Fatal(err)
			}

			pretty.PrettyPrint(resp.Accounts)
		},
	}
}
