yockey
======

yockey is a simple command-line tool for working with Ethereum keyfiles.


# Usage

### `yockey generate`

Generate a new keyfile.
If you want to use an existing private key to use in the keyfile, it can be 
specified by setting `--privatekey` with the location of the file containing the 
private key.


### `yockey inspect <keyfile>`

Print various information about the keyfile.
Private key information can be printed by using the `--private` flag;
make sure to use this feature with great caution!


### `yockey sign <keyfile> <message/file>`

Sign the message with a keyfile.
It is possible to refer to a file containing the message.


### `yockey verify <address> <signature> <message/file>`

Verify the signature of the message.
It is possible to refer to a file containing the message.


## Passphrases

For every command that uses a keyfile, you will be prompted to provide the 
passphrase for decrypting the keyfile.  To avoid this message, it is possible
to pass the passphrase by using the `--passphrase` flag pointing to a file that
contains the passphrase.
