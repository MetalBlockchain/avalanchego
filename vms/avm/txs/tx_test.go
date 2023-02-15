// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"errors"
	"testing"

	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/codec/linearcodec"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/utils/crypto"
	"github.com/MetalBlockchain/metalgo/utils/units"
	"github.com/MetalBlockchain/metalgo/utils/wrappers"
	"github.com/MetalBlockchain/metalgo/vms/avm/fxs"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
)

var (
	networkID       uint32 = 10
	chainID                = ids.ID{5, 4, 3, 2, 1}
	platformChainID        = ids.Empty.Prefix(0)

	keys = crypto.BuildTestKeys()

	assetID = ids.ID{1, 2, 3}
)

func setupCodec() codec.Manager {
	parser, err := NewParser([]fxs.Fx{
		&secp256k1fx.Fx{},
	})
	if err != nil {
		panic(err)
	}
	return parser.Codec()
}

func NewContext(tb testing.TB) *snow.Context {
	ctx := snow.DefaultContextTest()
	ctx.NetworkID = networkID
	ctx.ChainID = chainID
	avaxAssetID, err := ids.FromString("2XGxUr7VF7j1iwUp2aiGe4b6Ue2yyNghNS1SuNTNmZ77dPpXFZ")
	if err != nil {
		tb.Fatal(err)
	}
	ctx.AVAXAssetID = avaxAssetID
	ctx.XChainID = ids.Empty.Prefix(0)
	ctx.CChainID = ids.Empty.Prefix(1)
	aliaser := ctx.BCLookup.(ids.Aliaser)

	errs := wrappers.Errs{}
	errs.Add(
		aliaser.Alias(chainID, "X"),
		aliaser.Alias(chainID, chainID.String()),
		aliaser.Alias(platformChainID, "P"),
		aliaser.Alias(platformChainID, platformChainID.String()),
	)
	if errs.Errored() {
		tb.Fatal(errs.Err)
	}
	return ctx
}

func TestTxNil(t *testing.T) {
	ctx := NewContext(t)
	c := linearcodec.NewDefault()
	m := codec.NewDefaultManager()
	if err := m.RegisterCodec(CodecVersion, c); err != nil {
		t.Fatal(err)
	}

	tx := (*Tx)(nil)
	if err := tx.SyntacticVerify(ctx, m, ids.Empty, 0, 0, 1); err == nil {
		t.Fatalf("Should have erred due to nil tx")
	}
}

func TestTxEmpty(t *testing.T) {
	ctx := NewContext(t)
	c := setupCodec()
	tx := &Tx{}
	if err := tx.SyntacticVerify(ctx, c, ids.Empty, 0, 0, 1); err == nil {
		t.Fatalf("Should have erred due to nil tx")
	}
}

func TestTxInvalidCredential(t *testing.T) {
	ctx := NewContext(t)
	c := setupCodec()

	tx := &Tx{
		Unsigned: &BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*avax.TransferableInput{{
				UTXOID: avax.UTXOID{
					TxID:        ids.Empty,
					OutputIndex: 0,
				},
				Asset: avax.Asset{ID: assetID},
				In: &secp256k1fx.TransferInput{
					Amt: 20 * units.KiloAvax,
					Input: secp256k1fx.Input{
						SigIndices: []uint32{
							0,
						},
					},
				},
			}},
		}},
		Creds: []*fxs.FxCredential{{Verifiable: &avax.TestVerifiable{Err: errors.New("")}}},
	}
	tx.SetBytes(nil, nil)

	if err := tx.SyntacticVerify(ctx, c, ids.Empty, 0, 0, 1); err == nil {
		t.Fatalf("Tx should have failed due to an invalid credential")
	}
}

func TestTxInvalidUnsignedTx(t *testing.T) {
	ctx := NewContext(t)
	c := setupCodec()

	tx := &Tx{
		Unsigned: &BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*avax.TransferableInput{
				{
					UTXOID: avax.UTXOID{
						TxID:        ids.Empty,
						OutputIndex: 0,
					},
					Asset: avax.Asset{ID: assetID},
					In: &secp256k1fx.TransferInput{
						Amt: 20 * units.KiloAvax,
						Input: secp256k1fx.Input{
							SigIndices: []uint32{
								0,
							},
						},
					},
				},
				{
					UTXOID: avax.UTXOID{
						TxID:        ids.Empty,
						OutputIndex: 0,
					},
					Asset: avax.Asset{ID: assetID},
					In: &secp256k1fx.TransferInput{
						Amt: 20 * units.KiloAvax,
						Input: secp256k1fx.Input{
							SigIndices: []uint32{
								0,
							},
						},
					},
				},
			},
		}},
		Creds: []*fxs.FxCredential{
			{Verifiable: &avax.TestVerifiable{}},
			{Verifiable: &avax.TestVerifiable{}},
		},
	}
	tx.SetBytes(nil, nil)

	if err := tx.SyntacticVerify(ctx, c, ids.Empty, 0, 0, 1); err == nil {
		t.Fatalf("Tx should have failed due to an invalid unsigned tx")
	}
}

func TestTxInvalidNumberOfCredentials(t *testing.T) {
	ctx := NewContext(t)
	c := setupCodec()

	tx := &Tx{
		Unsigned: &BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Ins: []*avax.TransferableInput{
				{
					UTXOID: avax.UTXOID{TxID: ids.Empty, OutputIndex: 0},
					Asset:  avax.Asset{ID: assetID},
					In: &secp256k1fx.TransferInput{
						Amt: 20 * units.KiloAvax,
						Input: secp256k1fx.Input{
							SigIndices: []uint32{
								0,
							},
						},
					},
				},
				{
					UTXOID: avax.UTXOID{TxID: ids.Empty, OutputIndex: 1},
					Asset:  avax.Asset{ID: assetID},
					In: &secp256k1fx.TransferInput{
						Amt: 20 * units.KiloAvax,
						Input: secp256k1fx.Input{
							SigIndices: []uint32{
								0,
							},
						},
					},
				},
			},
		}},
		Creds: []*fxs.FxCredential{{Verifiable: &avax.TestVerifiable{}}},
	}
	tx.SetBytes(nil, nil)

	if err := tx.SyntacticVerify(ctx, c, ids.Empty, 0, 0, 1); err == nil {
		t.Fatalf("Tx should have failed due to an invalid number of credentials")
	}
}
