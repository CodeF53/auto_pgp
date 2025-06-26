# auto_pgp
Small tool to to automatically get Abacus 2fa login code.

xsel, xclip, or wl-clipboard needed on linux. If you use xsel on wayland, you cannot do ctrl+c to copy message, you need to do `xsel -b`

## how to use
1. copy PGP message from abacus's "Two-factor Authentication" page
2. run the program `autopgp` or `autopgp.exe` depending on your os
3. paste your new clipboard contents into the "Security code" field and proceed

## How can I trust this won't do something stupid?
Read the source code, `./main.go` is less than 100 lines

All this does is:
1. Read clipboard content
2. Read your pgp key
3. Decode clipboard using pgp key
4. Use regex to find login key
5. Send that key to your clipboard

## how compile?!?!?
read [the release workflow](https://github.com/CodeF53/auto_pgp/blob/master/.github/workflows/deploy.yml)
