// Authored and revised by YOC team, 2017-2018
// License placeholder #1

package main

import (
	"fmt"
	"os"

	"github.com/Yocoin15/Yocoin_Sources/cmd/utils"
	"gopkg.in/urfave/cli.v1"
)

const (
	defaultKeyfileName = "keyfile.json"
)

// Git SHA1 commit hash of the release (set via linker flags)
var gitCommit = ""

var app *cli.App

func init() {
	app = utils.NewApp(gitCommit, "an YoCoin key manager")
	app.Commands = []cli.Command{
		commandGenerate,
		commandInspect,
		commandSignMessage,
		commandVerifyMessage,
	}
}

// Commonly used command line flags.
var (
	passphraseFlag = cli.StringFlag{
		Name:  "passwordfile",
		Usage: "the file that contains the passphrase for the keyfile",
	}
	jsonFlag = cli.BoolFlag{
		Name:  "json",
		Usage: "output JSON instead of human-readable format",
	}
	messageFlag = cli.StringFlag{
		Name:  "message",
		Usage: "the file that contains the message to sign/verify",
	}
)

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
