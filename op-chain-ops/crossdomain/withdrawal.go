package crossdomain

import (
	"math/big"
	"strings"

	"github.com/ethereum-optimism/optimism/op-bindings/predeploys"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	// Standard ABI types
	Uint256Type, _ = abi.NewType("uint256", "", nil)
	BytesType, _   = abi.NewType("bytes", "", nil)
	AddressType, _ = abi.NewType("address", "", nil)
	// messagePasserABI is the JSON ABI for the OVM_L2ToL1MessagePasser method
	// `passMessageToL1`
	messagePasserABI = "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_message\",\"type\":\"bytes\"}],\"name\":\"passMessageToL1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	// messagePasser represents the ABI type for `passMessageToL1`
	messagePasser abi.ABI
)

// Create the required ABI
func init() {
	var err error
	messagePasser, err = abi.JSON(strings.NewReader(messagePasserABI))
	if err != nil {
		panic(err)
	}
}

// Withdrawal represents a withdrawal transaction on L2
type Withdrawal struct {
	Nonce    *big.Int
	Sender   *common.Address
	Target   *common.Address
	Value    *big.Int
	GasLimit *big.Int
	Data     []byte
}

// Encode will serialize the Withdrawal so that it is suitable for hashing.
func (w *Withdrawal) Encode() ([]byte, error) {
	args := abi.Arguments{
		{Name: "nonce", Type: Uint256Type},
		{Name: "sender", Type: AddressType},
		{Name: "target", Type: AddressType},
		{Name: "value", Type: Uint256Type},
		{Name: "gasLimit", Type: Uint256Type},
		{Name: "data", Type: BytesType},
	}
	enc, err := args.Pack(w.Nonce, w.Sender, w.Target, w.Value, w.GasLimit, w.Data)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// EncodeLegacy will serialze the Withdrawal in the legacy format so that it
// is suitable for hashing.
func (w *Withdrawal) EncodeLegacy() ([]byte, error) {
	msg, err := EncodeCrossDomainMessageV0(w.Target, w.Sender, w.Data, w.Nonce)
	if err != nil {
		return nil, err
	}

	enc, err := messagePasser.Pack("passMessageToL1", msg)
	if err != nil {
		return nil, err
	}

	out := make([]byte, len(enc)+len(predeploys.L2ToL1MessagePasserAddr))
	copy(out, enc)
	copy(out[len(enc):], predeploys.L2ToL1MessagePasserAddr.Bytes())

	return out, nil
}

// Hash will hash the Withdrawal. This is the hash that is computed in
// the L2ToL1MessagePasser. The encoding is the same as the v1 cross domain
// message encoding without the 4byte selector prepended.
func (w *Withdrawal) Hash() (common.Hash, error) {
	encoded, err := w.Encode()
	if err != nil {
		return common.Hash{}, err
	}
	hash := crypto.Keccak256(encoded)
	return common.BytesToHash(hash), nil
}

// LegacyHash will compute the legacy style hash that is computed in the
// OVM_L2ToL1MessagePasser.
func (w *Withdrawal) LegacyHash() (common.Hash, error) {
	encoded, err := w.EncodeLegacy()
	if err != nil {
		return common.Hash{}, nil
	}
	hash := crypto.Keccak256(encoded)
	return common.BytesToHash(hash), nil
}
