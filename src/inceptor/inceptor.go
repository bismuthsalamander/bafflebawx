package inceptor

import (
	"crypto/rand"
	"math/big"
)

var max64 *big.Int

func init() {
	max64 = big.NewInt(2)
	max64.Exp(big.NewInt(2), big.NewInt(64), nil)
	max64.Sub(max64, big.NewInt(1))
}

func Uint64() (uint64, error) {
	nBig, err := rand.Int(rand.Reader, max64)
	if err != nil {
		return 0, err
	}
	return nBig.Uint64(), err
}
