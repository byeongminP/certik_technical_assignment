package types

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(constantFee sdk.Coin, startingCertificateID uint64) GenesisState {
	return GenesisState{}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		NextCertificateID: uint64(1),
	}
}

// ValidateGenesis - validate crisis genesis data
func ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	if err := ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", ModuleName, err)
	}
	if data.NextCertificateID < 1 {
		return fmt.Errorf("failed to validate %s genesis state: NextCertificateID must be positive", ModuleName)
	}
	return nil
}

// GetGenesisStateFromAppState returns cert module GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.Marshaler, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisState) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for i := 0; i < len(g.Certificates); i++ {
		var cert Certificate
		err := unpacker.UnpackAny(g.Certificates[i], &cert)
		if err != nil {
			return err
		}
	}

	for _, platform := range g.Platforms {
		err := platform.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}

	for _, validator := range g.Validators {
		err := validator.UnpackInterfaces(unpacker)
		if err != nil {
			return err
		}
	}
	return nil
}
