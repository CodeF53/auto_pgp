package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"regexp"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

func main() {
	keyFile, err := os.Open("private_key.asc")
	if err != nil { errHandle("Failed to read private_key.asc, does it exist?") }
	defer keyFile.Close()

	clip := clipboard.New()
	pgp := crypto.PGP()

	clipContents, err := clip.PasteText()
	if err != nil { log.Panic("Failed to read clipboard") }

	privKey, err := crypto.NewKeyFromReader(keyFile)
	if err != nil { errHandle("Invalid key in file private_key.asc") }
	decHandler, err := pgp.Decryption().DecryptionKey(privKey).New()
	if err != nil { errHandle("yeah this wasn't my fault (i think). something went really fucking wrong") }
	decrypted, err := decHandler.Decrypt([]byte(clipContents), crypto.Auto)
	if err != nil { errHandle("Decryption failed, is the message in your clipboard correct?") }

	regex := regexp.MustCompile(`[0-9a-f]{56}`)
	final := regex.FindString(string(decrypted.Bytes()[:]))

	clip.CopyText(final)
	fmt.Println("\033[32mKey copied to clipboard!\033[0m")
	time.Sleep(2 * time.Second)
	os.Exit(0)
}

func errHandle(message string) {
	fmt.Println("\033[31mError: ", message, "\033[0m")
	time.Sleep(3 * time.Second)
	os.Exit(1)
}
