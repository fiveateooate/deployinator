package vaulthelpers

import (
	vault "github.com/hashicorp/vault/api"
)

// VaultPolicy type to hold policy info
type VaultPolicy struct {
	Path  string
	Rules string
}

// NewVaultPolicy return a new vaultpolicydeployer type
func NewVaultPolicy() *VaultPolicy {
	vp := new(VaultPolicy)
	return vp
}

// CreatePolicy - make the vault policy
func (vp *VaultPolicy) CreatePolicy() error {
	vs := vault.Sys{}
	if err := vs.PutPolicy(vp.Path, vp.Rules); err != nil {
		return err
	}
	return nil
}
