package switchAccount

import (
	"flag"
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"
	acmeclient "github.com/cpu/acmeshell/acme/client"
	"github.com/cpu/acmeshell/shell/commands"
)

type switchAccountCmd struct {
	commands.BaseCmd
}

var SwitchAccountCommand = switchAccountCmd{
	commands.BaseCmd{
		Cmd: &ishell.Cmd{
			Name:     "switchAccount",
			Aliases:  []string{"switch", "switchAcct", "switchReg", "switchRegistration"},
			Func:     switchAccountHandler,
			Help:     "Switch the active ACME account",
			LongHelp: `TODO(@cpu): Write this!`,
		},
	},
}

func (a switchAccountCmd) Setup(client *acmeclient.Client) (*ishell.Cmd, error) {
	return SwitchAccountCommand.Cmd, nil
}

func switchAccountHandler(c *ishell.Context) {
	switchAccountFlags := flag.NewFlagSet("switchAccount", flag.ContinueOnError)

	accountIndex := switchAccountFlags.Int("account", -1, "account number to switch to. "+
		"leave blank to pick interactively")

	err := switchAccountFlags.Parse(c.Args)
	if err != nil && err != flag.ErrHelp {
		c.Printf("switchAccount: error parsing input flags: %s\n", err.Error())
		return
	} else if err == flag.ErrHelp {
		return
	}

	client := commands.GetClient(c)

	if *accountIndex >= 0 {
		if *accountIndex >= len(client.Accounts) {
			c.Printf("switchAccount: provided account index (%d) "+
				"is larger than number of accounts (%d)\n",
				*accountIndex, len(client.Accounts))
			return
		}

		client.ActiveAccount = client.Accounts[*accountIndex]
		c.Printf("Active account is now #%d - %q\n", *accountIndex, client.ActiveAccount.ID)
		return
	}

	accountList := make([]string, len(client.Accounts))
	for i, acct := range client.Accounts {
		line := fmt.Sprintf("%3d)", i)

		contacts := "none"
		if len(acct.Contact) > 0 {
			contacts = strings.Join(acct.Contact, ", ")
		}
		line += fmt.Sprintf(" %s", contacts)
		line += fmt.Sprintf(" %q", acct.ID)
		accountList[i] = line
	}

	choice := c.MultiChoice(accountList, "Which account would you like to switch to?")

	client.ActiveAccount = client.Accounts[choice]
	c.Printf("Active account is now #%d - %q\n", choice, client.ActiveAccount.ID)
}