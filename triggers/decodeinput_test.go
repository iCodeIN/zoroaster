package trigger

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestDecodeInputData(t *testing.T) {

	const data = "0xef34358800000000000000000000000000000000000000000000000011f7372757499810000000000000000000000000000000000000000000001ba270fd1e13e6d83f8000000000000000000000000000000000000000000000000000000000000186a00000000000000000000000000000000000000000000000000000000000002f7300000000000000000000000000000000000000000000000011f73727574998100000000000000000000000000000000000000000000000000000000000002bd100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002cc99991d58d500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000198f46f520f33cd4329bd4be380a25a90536cd500000000000000000000000093eeb4782cf0fd477ff7a969224cb6caf30aeb8c000000000000000000000000ac801efa162c5c5a931b86c799528d9463db370f000000000000000000000000000000000000000000000000000000000000001b000000000000000000000000000000000000000000000000000000000000001c11569bd0758cf5a325278fea75e81daf20756bc7316d73732f1bc47dc1dba05a431895a099dedb1fdf91b059c04f909b986b3ce8832b0df094ff0097780ea01c3631cd50e33ed246b2e1063d7ae8b49e211d76592af8da2dedf337f928fb804f51de6b19ed5ab0e9c589eb775edb45eed5851cbd125feb359d276bd24206f58f"
	const abi = `[{"constant":false,"inputs":[{"name":"assertion","type":"bool"}],"name":"assert","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"setOwner","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"token","type":"address"},{"name":"amount","type":"uint256"},{"name":"user","type":"address"},{"name":"nonce","type":"uint256"},{"name":"v","type":"uint8"},{"name":"r","type":"bytes32"},{"name":"s","type":"bytes32"},{"name":"feeWithdrawal","type":"uint256"}],"name":"adminWithdraw","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"lastActiveTransaction","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"token","type":"address"},{"name":"amount","type":"uint256"}],"name":"depositToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"withdrawn","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"admins","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"admin","type":"address"},{"name":"isAdmin","type":"bool"}],"name":"setAdmin","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"tokens","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"feeAccount","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"invalidOrder","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"getOwner","outputs":[{"name":"out","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"a","type":"uint256"},{"name":"b","type":"uint256"}],"name":"safeSub","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"user","type":"address"},{"name":"nonce","type":"uint256"}],"name":"invalidateOrdersBefore","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"a","type":"uint256"},{"name":"b","type":"uint256"}],"name":"safeMul","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"deposit","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"traded","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"expiry","type":"uint256"}],"name":"setInactivityReleasePeriod","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"a","type":"uint256"},{"name":"b","type":"uint256"}],"name":"safeAdd","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"tradeValues","type":"uint256[8]"},{"name":"tradeAddresses","type":"address[4]"},{"name":"v","type":"uint8[2]"},{"name":"rs","type":"bytes32[4]"}],"name":"trade","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"inactivityReleasePeriod","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"token","type":"address"},{"name":"amount","type":"uint256"}],"name":"withdraw","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"bytes32"}],"name":"orderFills","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"token","type":"address"},{"name":"user","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"feeAccount_","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":false,"stateMutability":"nonpayable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":true,"name":"previousOwner","type":"address"},{"indexed":true,"name":"newOwner","type":"address"}],"name":"SetOwner","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"tokenBuy","type":"address"},{"indexed":false,"name":"amountBuy","type":"uint256"},{"indexed":false,"name":"tokenSell","type":"address"},{"indexed":false,"name":"amountSell","type":"uint256"},{"indexed":false,"name":"expires","type":"uint256"},{"indexed":false,"name":"nonce","type":"uint256"},{"indexed":false,"name":"user","type":"address"},{"indexed":false,"name":"v","type":"uint8"},{"indexed":false,"name":"r","type":"bytes32"},{"indexed":false,"name":"s","type":"bytes32"}],"name":"Order","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"tokenBuy","type":"address"},{"indexed":false,"name":"amountBuy","type":"uint256"},{"indexed":false,"name":"tokenSell","type":"address"},{"indexed":false,"name":"amountSell","type":"uint256"},{"indexed":false,"name":"expires","type":"uint256"},{"indexed":false,"name":"nonce","type":"uint256"},{"indexed":false,"name":"user","type":"address"},{"indexed":false,"name":"v","type":"uint8"},{"indexed":false,"name":"r","type":"bytes32"},{"indexed":false,"name":"s","type":"bytes32"}],"name":"Cancel","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"tokenBuy","type":"address"},{"indexed":false,"name":"amountBuy","type":"uint256"},{"indexed":false,"name":"tokenSell","type":"address"},{"indexed":false,"name":"amountSell","type":"uint256"},{"indexed":false,"name":"get","type":"address"},{"indexed":false,"name":"give","type":"address"}],"name":"Trade","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"token","type":"address"},{"indexed":false,"name":"user","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"balance","type":"uint256"}],"name":"Deposit","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"token","type":"address"},{"indexed":false,"name":"user","type":"address"},{"indexed":false,"name":"amount","type":"uint256"},{"indexed":false,"name":"balance","type":"uint256"}],"name":"Withdraw","type":"event"}]`

	res, err := DecodeInputData(data, abi)
	assert.Nil(t, err)

	tradeValues := res["tradeValues"].([8]*big.Int)

	x := new(big.Int)
	x.SetString("130500409274193550000000", 10)

	assert.Equal(t, tradeValues[1].Cmp(x), 0)

	tradeAddresses := res["tradeAddresses"].([4]common.Address)
	assert.Equal(t, tradeAddresses[1].String(), "0x0198f46f520F33cd4329bd4bE380a25a90536CD5")

	v := res["v"].([2]uint8)
	y := [2]uint8{27, 28}
	assert.Equal(t, v, y)

	rs := res["rs"].([4][32]uint8) // bytes32[4]
	z := [32]uint8{17, 86, 155, 208, 117, 140, 245, 163, 37, 39, 143, 234, 117, 232, 29, 175, 32, 117, 107, 199, 49, 109, 115, 115, 47, 27, 196, 125, 193, 219, 160, 90}
	assert.Equal(t, rs[0], z)
}
