// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/crypto"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/components/verify"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
)

func TestImportTxSerialization(t *testing.T) {
	expected := []byte{
		// Codec version
		0x00, 0x00,
		// txID:
		0x00, 0x00, 0x00, 0x03,
		// networkID:
		0x00, 0x00, 0x00, 0x02,
		// blockchainID:
		0xff, 0xff, 0xff, 0xff, 0xee, 0xee, 0xee, 0xee,
		0xdd, 0xdd, 0xdd, 0xdd, 0xcc, 0xcc, 0xcc, 0xcc,
		0xbb, 0xbb, 0xbb, 0xbb, 0xaa, 0xaa, 0xaa, 0xaa,
		0x99, 0x99, 0x99, 0x99, 0x88, 0x88, 0x88, 0x88,
		// number of base outs:
		0x00, 0x00, 0x00, 0x00,
		// number of base inputs:
		0x00, 0x00, 0x00, 0x00,
		// Memo length:
		0x00, 0x00, 0x00, 0x04,
		// Memo:
		0x00, 0x01, 0x02, 0x03,
		// Source Chain ID:
		0x1f, 0x8f, 0x9f, 0x0f, 0x1e, 0x8e, 0x9e, 0x0e,
		0x2d, 0x7d, 0xad, 0xfd, 0x2c, 0x7c, 0xac, 0xfc,
		0x3b, 0x6b, 0xbb, 0xeb, 0x3a, 0x6a, 0xba, 0xea,
		0x49, 0x59, 0xc9, 0xd9, 0x48, 0x58, 0xc8, 0xd8,
		// number of inputs:
		0x00, 0x00, 0x00, 0x01,
		// utxoID:
		0x0f, 0x2f, 0x4f, 0x6f, 0x8e, 0xae, 0xce, 0xee,
		0x0d, 0x2d, 0x4d, 0x6d, 0x8c, 0xac, 0xcc, 0xec,
		0x0b, 0x2b, 0x4b, 0x6b, 0x8a, 0xaa, 0xca, 0xea,
		0x09, 0x29, 0x49, 0x69, 0x88, 0xa8, 0xc8, 0xe8,
		// output index
		0x00, 0x00, 0x00, 0x00,
		// assetID:
		0x1f, 0x3f, 0x5f, 0x7f, 0x9e, 0xbe, 0xde, 0xfe,
		0x1d, 0x3d, 0x5d, 0x7d, 0x9c, 0xbc, 0xdc, 0xfc,
		0x1b, 0x3b, 0x5b, 0x7b, 0x9a, 0xba, 0xda, 0xfa,
		0x19, 0x39, 0x59, 0x79, 0x98, 0xb8, 0xd8, 0xf8,
		// input:
		// input ID:
		0x00, 0x00, 0x00, 0x05,
		// amount:
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0xe8,
		// num sig indices:
		0x00, 0x00, 0x00, 0x01,
		// sig index[0]:
		0x00, 0x00, 0x00, 0x00,
		// number of credentials:
		0x00, 0x00, 0x00, 0x00,
	}

	tx := &Tx{Unsigned: &ImportTx{
		BaseTx: BaseTx{BaseTx: avax.BaseTx{
			NetworkID: 2,
			BlockchainID: ids.ID{
				0xff, 0xff, 0xff, 0xff, 0xee, 0xee, 0xee, 0xee,
				0xdd, 0xdd, 0xdd, 0xdd, 0xcc, 0xcc, 0xcc, 0xcc,
				0xbb, 0xbb, 0xbb, 0xbb, 0xaa, 0xaa, 0xaa, 0xaa,
				0x99, 0x99, 0x99, 0x99, 0x88, 0x88, 0x88, 0x88,
			},
			Memo: []byte{0x00, 0x01, 0x02, 0x03},
		}},
		SourceChain: ids.ID{
			0x1f, 0x8f, 0x9f, 0x0f, 0x1e, 0x8e, 0x9e, 0x0e,
			0x2d, 0x7d, 0xad, 0xfd, 0x2c, 0x7c, 0xac, 0xfc,
			0x3b, 0x6b, 0xbb, 0xeb, 0x3a, 0x6a, 0xba, 0xea,
			0x49, 0x59, 0xc9, 0xd9, 0x48, 0x58, 0xc8, 0xd8,
		},
		ImportedIns: []*avax.TransferableInput{{
			UTXOID: avax.UTXOID{TxID: ids.ID{
				0x0f, 0x2f, 0x4f, 0x6f, 0x8e, 0xae, 0xce, 0xee,
				0x0d, 0x2d, 0x4d, 0x6d, 0x8c, 0xac, 0xcc, 0xec,
				0x0b, 0x2b, 0x4b, 0x6b, 0x8a, 0xaa, 0xca, 0xea,
				0x09, 0x29, 0x49, 0x69, 0x88, 0xa8, 0xc8, 0xe8,
			}},
			Asset: avax.Asset{ID: ids.ID{
				0x1f, 0x3f, 0x5f, 0x7f, 0x9e, 0xbe, 0xde, 0xfe,
				0x1d, 0x3d, 0x5d, 0x7d, 0x9c, 0xbc, 0xdc, 0xfc,
				0x1b, 0x3b, 0x5b, 0x7b, 0x9a, 0xba, 0xda, 0xfa,
				0x19, 0x39, 0x59, 0x79, 0x98, 0xb8, 0xd8, 0xf8,
			}},
			In: &secp256k1fx.TransferInput{
				Amt:   1000,
				Input: secp256k1fx.Input{SigIndices: []uint32{0}},
			},
		}},
	}}

	c := setupCodec()
	if err := tx.SignSECP256K1Fx(c, nil); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, tx.ID().String(), "9wdPb5rsThXYLX4WxkNeyYrNMfDE5cuWLgifSjxKiA2dCmgCZ")
	result := tx.Bytes()
	if !bytes.Equal(expected, result) {
		t.Fatalf("\nExpected: 0x%x\nResult:   0x%x", expected, result)
	}

	credBytes := []byte{
		// type id
		0x00, 0x00, 0x00, 0x09,

		// there are two signers (thus two signatures)
		0x00, 0x00, 0x00, 0x02,

		// 65 bytes
		0x8c, 0xc7, 0xdc, 0x8c, 0x11, 0xd3, 0x75, 0x9e, 0x16, 0xa5,
		0x9f, 0xd2, 0x9c, 0x64, 0xd7, 0x1f, 0x9b, 0xad, 0x1a, 0x62,
		0x33, 0x98, 0xc7, 0xaf, 0x67, 0x02, 0xc5, 0xe0, 0x75, 0x8e,
		0x62, 0xcf, 0x15, 0x6d, 0x99, 0xf5, 0x4e, 0x71, 0xb8, 0xf4,
		0x8b, 0x5b, 0xbf, 0x0c, 0x59, 0x62, 0x79, 0x34, 0x97, 0x1a,
		0x1f, 0x49, 0x9b, 0x0a, 0x4f, 0xbf, 0x95, 0xfc, 0x31, 0x39,
		0x46, 0x4e, 0xa1, 0xaf, 0x00,

		// 65 bytes
		0x8c, 0xc7, 0xdc, 0x8c, 0x11, 0xd3, 0x75, 0x9e, 0x16, 0xa5,
		0x9f, 0xd2, 0x9c, 0x64, 0xd7, 0x1f, 0x9b, 0xad, 0x1a, 0x62,
		0x33, 0x98, 0xc7, 0xaf, 0x67, 0x02, 0xc5, 0xe0, 0x75, 0x8e,
		0x62, 0xcf, 0x15, 0x6d, 0x99, 0xf5, 0x4e, 0x71, 0xb8, 0xf4,
		0x8b, 0x5b, 0xbf, 0x0c, 0x59, 0x62, 0x79, 0x34, 0x97, 0x1a,
		0x1f, 0x49, 0x9b, 0x0a, 0x4f, 0xbf, 0x95, 0xfc, 0x31, 0x39,
		0x46, 0x4e, 0xa1, 0xaf, 0x00,

		// type id
		0x00, 0x00, 0x00, 0x09,

		// there are two signers (thus two signatures)
		0x00, 0x00, 0x00, 0x02,

		// 65 bytes
		0x8c, 0xc7, 0xdc, 0x8c, 0x11, 0xd3, 0x75, 0x9e, 0x16, 0xa5,
		0x9f, 0xd2, 0x9c, 0x64, 0xd7, 0x1f, 0x9b, 0xad, 0x1a, 0x62,
		0x33, 0x98, 0xc7, 0xaf, 0x67, 0x02, 0xc5, 0xe0, 0x75, 0x8e,
		0x62, 0xcf, 0x15, 0x6d, 0x99, 0xf5, 0x4e, 0x71, 0xb8, 0xf4,
		0x8b, 0x5b, 0xbf, 0x0c, 0x59, 0x62, 0x79, 0x34, 0x97, 0x1a,
		0x1f, 0x49, 0x9b, 0x0a, 0x4f, 0xbf, 0x95, 0xfc, 0x31, 0x39,
		0x46, 0x4e, 0xa1, 0xaf, 0x00,

		// 65 bytes
		0x8c, 0xc7, 0xdc, 0x8c, 0x11, 0xd3, 0x75, 0x9e, 0x16, 0xa5,
		0x9f, 0xd2, 0x9c, 0x64, 0xd7, 0x1f, 0x9b, 0xad, 0x1a, 0x62,
		0x33, 0x98, 0xc7, 0xaf, 0x67, 0x02, 0xc5, 0xe0, 0x75, 0x8e,
		0x62, 0xcf, 0x15, 0x6d, 0x99, 0xf5, 0x4e, 0x71, 0xb8, 0xf4,
		0x8b, 0x5b, 0xbf, 0x0c, 0x59, 0x62, 0x79, 0x34, 0x97, 0x1a,
		0x1f, 0x49, 0x9b, 0x0a, 0x4f, 0xbf, 0x95, 0xfc, 0x31, 0x39,
		0x46, 0x4e, 0xa1, 0xaf, 0x00,
	}
	if err := tx.SignSECP256K1Fx(c, [][]*crypto.PrivateKeySECP256K1R{{keys[0], keys[0]}, {keys[0], keys[0]}}); err != nil {
		t.Fatal(err)
	}
	require.Equal(t, tx.ID().String(), "pCW7sVBytzdZ1WrqzGY1DvA2S9UaMr72xpUMxVyx1QHBARNYx")
	result = tx.Bytes()

	// there are two credentials
	expected[len(expected)-1] = 0x02
	expected = append(expected, credBytes...)
	if !bytes.Equal(expected, result) {
		t.Fatalf("\nExpected: 0x%x\nResult:   0x%x", expected, result)
	}
}

func TestImportTxSyntacticVerify(t *testing.T) {
	ctx := NewContext(t)
	c := setupCodec()

	tx := &ImportTx{
		BaseTx: BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Outs: []*avax.TransferableOutput{{
				Asset: avax.Asset{ID: assetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 12345,
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
					},
				},
			}},
		}},
		SourceChain: platformChainID,
		ImportedIns: []*avax.TransferableInput{{
			UTXOID: avax.UTXOID{
				TxID: ids.ID{
					0xff, 0xfe, 0xfd, 0xfc, 0xfb, 0xfa, 0xf9, 0xf8,
					0xf7, 0xf6, 0xf5, 0xf4, 0xf3, 0xf2, 0xf1, 0xf0,
					0xef, 0xee, 0xed, 0xec, 0xeb, 0xea, 0xe9, 0xe8,
					0xe7, 0xe6, 0xe5, 0xe4, 0xe3, 0xe2, 0xe1, 0xe0,
				},
				OutputIndex: 0,
			},
			Asset: avax.Asset{ID: assetID},
			In: &secp256k1fx.TransferInput{
				Amt: 54321,
				Input: secp256k1fx.Input{
					SigIndices: []uint32{2},
				},
			},
		}},
	}
	tx.Initialize(nil)

	if err := tx.SyntacticVerify(ctx, c, ids.Empty, 0, 0, 0); err != nil {
		t.Fatal(err)
	}
}

func TestImportTxSyntacticVerifyInvalidMemo(t *testing.T) {
	ctx := NewContext(t)
	c := setupCodec()

	tx := &ImportTx{
		BaseTx: BaseTx{BaseTx: avax.BaseTx{
			NetworkID:    networkID,
			BlockchainID: chainID,
			Outs: []*avax.TransferableOutput{{
				Asset: avax.Asset{ID: assetID},
				Out: &secp256k1fx.TransferOutput{
					Amt: 12345,
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{keys[0].PublicKey().Address()},
					},
				},
			}},
			Memo: make([]byte, avax.MaxMemoSize+1),
		}},
		SourceChain: platformChainID,
		ImportedIns: []*avax.TransferableInput{{
			UTXOID: avax.UTXOID{
				TxID: ids.ID{
					0xff, 0xfe, 0xfd, 0xfc, 0xfb, 0xfa, 0xf9, 0xf8,
					0xf7, 0xf6, 0xf5, 0xf4, 0xf3, 0xf2, 0xf1, 0xf0,
					0xef, 0xee, 0xed, 0xec, 0xeb, 0xea, 0xe9, 0xe8,
					0xe7, 0xe6, 0xe5, 0xe4, 0xe3, 0xe2, 0xe1, 0xe0,
				},
				OutputIndex: 0,
			},
			Asset: avax.Asset{ID: assetID},
			In: &secp256k1fx.TransferInput{
				Amt: 54321,
				Input: secp256k1fx.Input{
					SigIndices: []uint32{2},
				},
			},
		}},
	}
	tx.Initialize(nil)

	if err := tx.SyntacticVerify(ctx, c, ids.Empty, 0, 0, 0); err == nil {
		t.Fatalf("should have erred due to memo field being too long")
	}
}

func TestImportTxNotState(t *testing.T) {
	intf := interface{}(&ImportTx{})
	if _, ok := intf.(verify.State); ok {
		t.Fatalf("shouldn't be marked as state")
	}
}
