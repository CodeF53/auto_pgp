import { file } from 'bun'
import clipboard from 'clipboardy'
import * as openpgp from 'openpgp'

const armoredMessage = clipboard.readSync()

const message = await openpgp.readMessage({ armoredMessage })
  .catch(() => { console.error('Invalid PGP Message in clipboard'); process.exit() })

const armoredKey = await file('./private_key.asc')
  .text()
  .catch(() => { console.error('Failed to read private_key.asc, does it exist?'); process.exit() })
const decryptionKeys = await openpgp.readPrivateKey({ armoredKey })
  .catch(() => { console.error('Invalid PGP Private Key'); process.exit() })

const { data: decryptedMessage } = await openpgp.decrypt({ message, decryptionKeys })
  .catch(() => { console.error('PGP Decryption Error'); process.exit() })

const key = /[0-9a-f]{56}/.exec(decryptedMessage)![0]

clipboard.writeSync(key)
console.log('\u001B[32mKey copied to clipboard!\u001B[0m')
