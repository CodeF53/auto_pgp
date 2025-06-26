package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/atotto/clipboard"
)

var pathToKey string
var toStdout bool

func init() {
	flag.StringVar(&pathToKey, "key", "private_key.asc", "Path to private key")
	flag.BoolVar(&toStdout, "o", false, "Print output to stdout instead of copying to clipboard")
}

func main() {
	flag.Parse()

	keyReader, err := os.Open(pathToKey)
	if err != nil { errHandle(fmt.Sprintf("Failed to read %s, does it exist?", pathToKey)) }
	defer keyReader.Close()

	keyring, err := openpgp.ReadArmoredKeyRing(keyReader)
	if err != nil { errHandle(fmt.Sprintf("No valid keys found in %s", pathToKey)) }

	if decKey := keyring.DecryptionKeys()[0].PrivateKey; decKey.Encrypted{
		for {
			fmt.Print("\nInput password for private key: ")
			stdinScanner := bufio.NewScanner(os.Stdin)
			stdinScanner.Scan()
			input := stdinScanner.Bytes()
			if err := decKey.Decrypt(input); err == nil {
				fmt.Println("")
				break
			}
			fmt.Println("\033[31mIncorrect password\033[0m")
		}
	}

	clipContents, err := clipboard.ReadAll()
	if err != nil { errHandle("Failed to read clipboard") }

	armoredMessage, err := armor.Decode(strings.NewReader(clipContents))
	if err != nil { errHandle("Invalid message in clipboard. Make sure ASCII armor (-----BEGIN/END PGP MESSAGE-----) is correct") }

	messageReader, err := openpgp.ReadMessage(armoredMessage.Body, keyring, nil, nil)
	if err != nil { errHandle("Failed to read message, is the key correct?") }

	message, _ := io.ReadAll(messageReader.UnverifiedBody)

	regex := regexp.MustCompile(`[0-9a-f]{56}`)
	final := regex.Find(message)
	if final == nil { errHandle("Message did not contain an abacus verification key") }

	clipboard.WriteAll(string(final))
	fmt.Println("\033[32mKey copied to clipboard!\033[0m")
	time.Sleep(2 * time.Second)
	os.Exit(0)
}

func errHandle(message string) {
	fmt.Printf("\033[31mError: %s\033[0m\n", message)
	time.Sleep(3 * time.Second)
	os.Exit(1)
}
