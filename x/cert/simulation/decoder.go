package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/certikfoundation/shentu/x/cert/types"
)

// DecodeCertID converts the binary representation of cert id into an uint64 value.
func DecodeCertID(bz []byte) uint64 {
	if bz == nil {
		return 0
	}
	id := binary.LittleEndian.Uint64(bz)
	return id
}

// DecodeCertIDs converts the binary representation of cert ids into []uint64.
func DecodeCertIDs(bz []byte) []uint64 {
	data := make([]uint64, len(bz)/8)
	for i := range data {
		data[i] = binary.LittleEndian.Uint64(bz[8*i:])
	}
	return data
}

// NewDecodeStore unmarshals the KVPair's Value to the corresponding type of cert module.
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.CertifiersStoreKey()):
			var certifierA, certifierB types.Certifier
			cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &certifierA)
			cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &certifierB)
			return fmt.Sprintf("%v\n%v", certifierA, certifierB)

		case bytes.Equal(kvA.Key[:1], types.ValidatorsStoreKey()):
			var validatorA, validatorB types.Validator
			cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &validatorA)
			cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &validatorB)
			return fmt.Sprintf("%v\n%v", validatorA, validatorB)

		case bytes.Equal(kvA.Key[:1], types.PlatformsStoreKey()):
			var platformA, platformB types.Platform
			cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &platformA)
			cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &platformB)
			return fmt.Sprintf("%v\n%v", platformA, platformB)

		case bytes.Equal(kvA.Key[:1], types.CertificatesStoreKey()):
			var certA, certB types.Certificate
			err := cdc.UnmarshalInterface(kvA.Value, &certA)
			if err != nil {
				panic(err)
			}
			err = cdc.UnmarshalInterface(kvB.Value, &certB)
			if err != nil {
				panic(err)
			}
			return fmt.Sprintf("%v\n%v", certA, certB)

		case bytes.Equal(kvA.Key[:1], types.LibrariesStoreKey()):
			var libraryA, libraryB types.Library
			cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &libraryA)
			cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &libraryB)
			return fmt.Sprintf("%v\n%v", libraryA, libraryB)

		case bytes.Equal(kvA.Key[:1], types.CertifierAliasesStoreKey()):
			var certifierA, certifierB types.Certifier
			cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &certifierA)
			cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &certifierB)
			return fmt.Sprintf("%v\n%v", certifierA, certifierB)

		case bytes.Equal(kvA.Key[:1], types.NextCertificateIDKey()), bytes.Equal(kvA.Key[:1], types.ContentCertIDStoreKeyPrefix):
			idA := DecodeCertID(kvA.Value)
			idB := DecodeCertID(kvB.Value)
			return fmt.Sprintf("%v\n%v", idA, idB)

		case bytes.Equal(kvA.Key[:1], types.CertifierCertIDsStoreKeyPrefix):
			idsA := DecodeCertIDs(kvA.Value)
			idsB := DecodeCertIDs(kvB.Value)
			return fmt.Sprintf("%v\n%v", idsA, idsB)

		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
