// Authored and revised by YOC team, 2017-2018
// License placeholder #1

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestMessageSignVerify(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "yockey-test")
	if err != nil {
		t.Fatal("Can't create temporary directory:", err)
	}
	defer os.RemoveAll(tmpdir)

	keyfile := filepath.Join(tmpdir, "the-keyfile")
	message := "test message"

	// Create the key.
	generate := runYOCKey(t, "generate", keyfile)
	generate.Expect(`
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
Repeat passphrase: {{.InputLine "foobar"}}
`)
	_, matches := generate.ExpectRegexp(`Address: (0x[0-9a-fA-F]{40})\n`)
	address := matches[1]
	generate.ExpectExit()

	// Sign a message.
	sign := runYOCKey(t, "signmessage", keyfile, message)
	sign.Expect(`
!! Unsupported terminal, password will be echoed.
Passphrase: {{.InputLine "foobar"}}
`)
	_, matches = sign.ExpectRegexp(`Signature: ([0-9a-f]+)\n`)
	signature := matches[1]
	sign.ExpectExit()

	// Verify the message.
	verify := runYOCKey(t, "verifymessage", address, signature, message)
	_, matches = verify.ExpectRegexp(`
Signature verification successful!
Recovered public key: [0-9a-f]+
Recovered address: (0x[0-9a-fA-F]{40})
`)
	recovered := matches[1]
	verify.ExpectExit()

	if recovered != address {
		t.Error("recovered address doesn't match generated key")
	}
}
