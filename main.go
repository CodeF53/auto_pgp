package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"regexp"

	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

func main() {
	keyReader, err := os.Open("private_key.asc")
	if err != nil { errHandle("Failed to read private_key.asc, does it exist?") }
	defer keyReader.Close()

	keyring, err := openpgp.ReadArmoredKeyRing(keyReader)
	if err != nil { errHandle("Invalid key in priv.asc") }

	if keyring[0].PrivateKey.Encrypted {
		for {
			fmt.Print("Input password for private key:")
			var input string
			fmt.Scanln(input)
			if err := keyring[0].PrivateKey.Decrypt([]byte(input)); err == nil {
				fmt.Println("")
				break
			}
			fmt.Println("\033[31mIncorrect password\033[0m")
		}
	}

	clip := clipboard.New()
	clipContents, err := clip.PasteText()
	if err != nil { errHandle("Failed to read clipboard") }

	armoredMessage, err := armor.Decode(strings.NewReader(clipContents))
	if err != nil { errHandle(err.Error()) }

	messageReader, err := openpgp.ReadMessage(armoredMessage.Body, keyring, nil, nil)
	if err != nil { errHandle(err.Error()) }

	message, _ := io.ReadAll(messageReader.UnverifiedBody)

	regex := regexp.MustCompile(`[0-9a-f]{56}`)
	final := regex.Find(message)
	if final == nil { errHandle("Message did not contain an abacus verification key") }

	clip.CopyText(string(final))
	fmt.Println("\033[32mKey copied to clipboard!\033[0m")
	time.Sleep(2 * time.Second)
	os.Exit(0)
}

func errHandle(message string) {
	fmt.Println("\033[31mError: ", message, "\033[0m")
	time.Sleep(3 * time.Second)
	os.Exit(1)
}
