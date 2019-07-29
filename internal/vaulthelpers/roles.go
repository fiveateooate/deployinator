package vaulthelpers

import (
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

// VaultRole type to hold policy info
type VaultRole struct {
	Name  string
	Rules string
}

// NewVaultRole return a new vaultpolicydeployer type
func NewVaultRole() *VaultRole {
	vr := new(VaultRole)
	return vr
}

// CreateRole - make the vault policy
func (vr *VaultRole) CreateRole() error {
	var (
		config   vault.Config
		response *vault.Response
		err      error
	)
	vcli, _ := vault.NewClient(&config)
	request := vcli.NewRequest("PUT", "someurl")
	if response, err = vcli.RawRequest(request); err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}
