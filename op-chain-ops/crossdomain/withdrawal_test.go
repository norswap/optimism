package crossdomain_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-chain-ops/crossdomain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

// ok get some legacy messages
func TestWithdrawalHash(t *testing.T) {
	sender := common.HexToAddress("0x4200000000000000000000000000000000000010")
	target := common.HexToAddress("0x99c9fc46f92e8a1c0dec1b1747d010903e884be1")

	w := crossdomain.NewWithdrawal(
		big.NewInt(110525),
		&sender,
		&target,
		big.NewInt(0),
		big.NewInt(0),
		hexutil.MustDecode("0xA9F9E675000000000000000000000000514910771AF9CA656AF840DFF83E8264ECF986CA000000000000000000000000350A791BFC2C21F9ED5D10980DAD2E2638FFA7F60000000000000000000000001D29B76B4BEB9954EAE4EFC0B817FA3F083FD9240000000000000000000000001D29B76B4BEB9954EAE4EFC0B817FA3F083FD924000000000000000000000000000000000000000000000022340B2DB1AC76534500000000000000000000000000000000000000000000000000000000000000C00000000000000000000000000000000000000000000000000000000000000000"),
	)

	hash, err := w.LegacyHash()
	require.Nil(t, err)

	fmt.Println(hash)
}
