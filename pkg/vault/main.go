package vault

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"

	hashiVault "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
)

func configDefaults() error {
	viper.SetDefault("vault_auth_method", "token")
	viper.SetDefault("vault_approle_mountpoint", "/v1/auth/approle/login")
	viper.SetDefault("vault_k8s_mountpoint", "auth/kubernets/login")
	viper.SetDefault("vault_k8s_token_path", "/var/run/secrets/kubernetes.io/serviceaccount/token")
	return nil
}

// Configure allows to create a Vault client connected to the vault server to easily read secrets
func Configure() (Vault, error) {
	var v Vault
	var err error

	err = configDefaults()
	if err != nil {
		return v, err
	}

	v.Client, err = hashiVault.NewClient(&hashiVault.Config{
		Address: viper.GetString("vault_url"),
	})
	if err != nil {
		return v, err
	}

	switch authMethod := viper.GetString("vault_auth_method"); authMethod {
	case "token":
		err := v.loginToken()
		if err != nil {
			return v, err
		}
	case "approle":
		err = v.loginApprole()
		if err != nil {
			return v, err
		}
	case "kubernetes":
		err = v.loginK8s()
		if err != nil {
			return v, err
		}
	default:
		errorMessage := fmt.Sprintf("error login into vault, %s auth method not supported. Valid values are: token, approle and kubernetes", authMethod)
		return v, errors.New(errorMessage)
	}

	if err := v.readAllConfig(); err != nil {
		return v, err
	}

	return v, nil
}

func (v Vault) loginApprole() error {
	r := v.Client.NewRequest("POST", viper.GetString("vault_approle_mountpoint"))
	raw := map[string]interface{}{
		"role_id":   viper.GetString("vault_role_id"),
		"secret_id": viper.GetString("vault_secret_id"),
	}
	ctx := context.Background()
	err := r.SetJSONBody(raw)
	if err != nil {
		return err
	}
	resp, err := v.Client.RawRequestWithContext(ctx, r)
	if err != nil {
		return err
	}
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Errorf("Error while closing body %s", err)
			}
		}()
	}
	token, parseErr := hashiVault.ParseSecret(resp.Body)
	if parseErr != nil {
		log.Errorf("could not parse secret token")
		return err
	}
	v.Client.SetToken(token.Auth.ClientToken)
	return nil
}

func (v Vault) loginToken() error {
	if viper.GetString("vault_token") != "" {
		v.Client.SetToken(viper.GetString("vault_token"))
		return nil
	}
	if os.Getenv("VAULT_TOKEN") != "" {
		v.Client.SetToken(os.Getenv("VAULT_TOKEN"))
		return nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.New(fmt.Sprintf("unable to find user home dir %s", err))
	}
	path := fmt.Sprintf("%s/%s", home, ".vault-token")
	clientToken, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New(fmt.Sprintf("unable to find vault token in user home dir %s", err))
	}
	v.Client.SetToken(string(clientToken))
	return nil
}

func (v Vault) loginK8s() error {
	var jwtToken string
	var err error
	jwtToken = viper.GetString("VaultK8sToken")
	if viper.GetString("VaultK8sTokenPath") != "" {
		jwtTokenByte, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
		jwtToken = string(jwtTokenByte)
		if err != nil {
			return errors.New("the vault JWT is not available so the vault secrets cannot be set")
		}
	}

	resp, err := v.Client.Logical().Write(
		viper.GetString("vault_k8s_mountpoint"),
		map[string]interface{}{
			"jwt":  jwtToken,
			"role": viper.GetString("VaultK8sRole"),
		},
	)

	if err != nil {
		return errors.New(fmt.Sprintf("cluster authentication failed for vault: %s", err))
	}

	if resp.Auth.ClientToken == "" {
		return errors.New("expected a client token")
	}

	v.Client.SetToken(resp.Auth.ClientToken)
	return nil
}

// ReadSecret read the secret at the given path
func (v Vault) ReadSecret(path string) (secret *hashiVault.Secret, err error) {
	client := v.Client.Logical()
	secret, err = client.Read(path)
	if err != nil {
		return
	}
	return
}

// ReadKey read the secret at the given path
func (v Vault) ReadKey(path string, key string) string {
	secret, err := v.ReadSecret(path)
	if err != nil {
		return ""
	}
	val, ok := secret.Data[key].(string)
	if !ok {
		return ""
	}
	return val
}

func (v Vault) readAllConfig() error {
	allSettings := viper.AllSettings()
	for key := range allSettings {
		val := viper.GetString(key)
		if vaultSecret := strings.Split(val, "VAULT::"); len(vaultSecret) > 1 {
			secretPath := strings.Split(vaultSecret[1], ":")
			if len(secretPath) < 2 {
				return errors.New("should specify path and key")
			}
			viper.Set(key, v.ReadKey(secretPath[0], secretPath[1]))
		}
	}
	return nil
}
