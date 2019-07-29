package deployers

import "github.com/fiveateooate/deployinator/internal/vaulthelpers"

// VaultDeployer type to hold policy info
type VaultDeployer struct {
	Name  string
	Rules string
}

// NewVaultDeployer return a new vaultpolicydeployer type
func NewVaultDeployer() *VaultDeployer {
	vd := new(VaultDeployer)
	return vd
}

// Deploy - make the vault policy
func (vd *VaultDeployer) Deploy() error {
	vp := vaulthelpers.NewVaultPolicy()
	vp.CreatePolicy()
	return nil
}
