package vault

import hashiVault "github.com/hashicorp/vault/api"

// Vault structure of the wrapper around Vault client library
type Vault struct {
	Client *hashiVault.Client `json:"client"`
}
