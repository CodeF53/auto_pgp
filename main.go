package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

func main() {
	keyFile, err := os.Open("private_key.asc")
	if err != nil { log.Fatal("Failed to read private_key.asc, does it exist?") }
	defer keyFile.Close()

	clip := clipboard.New()
	pgp := crypto.PGP()

	clipContents, err := clip.PasteText()
	if err != nil { log.Fatal("Failed to read clipboard") }

	privKey, err := crypto.NewKeyFromReader(keyFile)
	if err != nil { log.Fatal("Invalid key in file private_key.asc") }
	decHandler, err := pgp.Decryption().DecryptionKey(privKey).New()
	if err != nil { log.Fatal(err) }
	decrypted, err := decHandler.Decrypt([]byte(clipContents), crypto.Auto)
	if err != nil { log.Fatal("Decryption failed, are you using the correct key?") }

	regex := regexp.MustCompile(`[0-9a-f]{56}`)
	final := regex.FindString(string(decrypted.Bytes()[:]))

	clip.CopyText(final)
	fmt.Println("\u001B[32mKey copied to clipboard!\u001B[0m")
}
