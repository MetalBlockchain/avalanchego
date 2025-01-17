// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	_ "embed"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow/snowtest"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/units"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/stakeable"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
	"github.com/MetalBlockchain/metalgo/vms/types"
)

var (
	//go:embed convert_subnet_tx_test_simple.json
	convertSubnetTxSimpleJSON []byte
	//go:embed convert_subnet_tx_test_complex.json
	convertSubnetTxComplexJSON []byte
)

func TestConvertSubnetTxSerialization(t *testing.T) {
	var (
		addr = ids.ShortID{
			0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb,
			0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb,
			0x44, 0x55, 0x66, 0x77,
		}
		avaxAssetID = ids.ID{
			0x21, 0xe6, 0x73, 0x17, 0xcb, 0xc4, 0xbe, 0x2a,
			0xeb, 0x00, 0x67, 0x7a, 0xd6, 0x46, 0x27, 0x78,
			0xa8, 0xf5, 0x22, 0x74, 0xb9, 0xd6, 0x05, 0xdf,
			0x25, 0x91, 0xb2, 0x30, 0x27, 0xa8, 0x7d, 0xff,
		}
		customAssetID = ids.ID{
			0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
			0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
			0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
			0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
		}
		txID = ids.ID{
			0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
			0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
			0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
			0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
		}
		subnetID = ids.ID{
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
			0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
			0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
		}
		managerChainID = ids.ID{
			0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
			0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
			0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		}
		managerAddress = []byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0xde, 0xad,
		}
	)

	tests := []struct {
		name          string
		tx            *ConvertSubnetTx
		expectedBytes []byte
		expectedJSON  []byte
	}{
		{
			name: "simple",
			tx: &ConvertSubnetTx{
				BaseTx: BaseTx{
					BaseTx: avax.BaseTx{
						NetworkID:    constants.UnitTestID,
						BlockchainID: constants.PlatformChainID,
						Outs:         []*avax.TransferableOutput{},
						Ins: []*avax.TransferableInput{
							{
								UTXOID: avax.UTXOID{
									TxID:        txID,
									OutputIndex: 1,
								},
								Asset: avax.Asset{
									ID: avaxAssetID,
								},
								In: &secp256k1fx.TransferInput{
									Amt: units.MilliAvax,
									Input: secp256k1fx.Input{
										SigIndices: []uint32{5},
									},
								},
							},
						},
						Memo: types.JSONByteSlice{},
					},
				},
				Subnet:  subnetID,
				ChainID: managerChainID,
				Address: managerAddress,
				SubnetAuth: &secp256k1fx.Input{
					SigIndices: []uint32{3},
				},
			},
			expectedBytes: []byte{
				// Codec version
				0x00, 0x00,
				// ConvertSubnetTx Type ID
				0x00, 0x00, 0x00, 0x23,
				// Network ID
				0x00, 0x00, 0x00, 0x0a,
				// P-chain blockchain ID
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				// Number of outputs
				0x00, 0x00, 0x00, 0x00,
				// Number of inputs
				0x00, 0x00, 0x00, 0x01,
				// Inputs[0]
				// TxID
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				// Tx output index
				0x00, 0x00, 0x00, 0x01,
				// AVAX assetID
				0x21, 0xe6, 0x73, 0x17, 0xcb, 0xc4, 0xbe, 0x2a,
				0xeb, 0x00, 0x67, 0x7a, 0xd6, 0x46, 0x27, 0x78,
				0xa8, 0xf5, 0x22, 0x74, 0xb9, 0xd6, 0x05, 0xdf,
				0x25, 0x91, 0xb2, 0x30, 0x27, 0xa8, 0x7d, 0xff,
				// secp256k1fx transfer input type ID
				0x00, 0x00, 0x00, 0x05,
				// input amount = 1 MilliAvax
				0x00, 0x00, 0x00, 0x00, 0x00, 0x0f, 0x42, 0x40,
				// number of signatures needed in input
				0x00, 0x00, 0x00, 0x01,
				// index of signer
				0x00, 0x00, 0x00, 0x05,
				// length of memo
				0x00, 0x00, 0x00, 0x00,
				// subnetID to modify
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
				0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
				// chainID of the manager
				0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				// length of the manager address
				0x00, 0x00, 0x00, 0x14,
				// address of the manager
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0xde, 0xad,
				// secp256k1fx authorization type ID
				0x00, 0x00, 0x00, 0x0a,
				// number of signatures needed in authorization
				0x00, 0x00, 0x00, 0x01,
				// index of signer
				0x00, 0x00, 0x00, 0x03,
			},
			expectedJSON: convertSubnetTxSimpleJSON,
		},
		{
			name: "complex",
			tx: &ConvertSubnetTx{
				BaseTx: BaseTx{
					BaseTx: avax.BaseTx{
						NetworkID:    constants.UnitTestID,
						BlockchainID: constants.PlatformChainID,
						Outs: []*avax.TransferableOutput{
							{
								Asset: avax.Asset{
									ID: avaxAssetID,
								},
								Out: &stakeable.LockOut{
									Locktime: 87654321,
									TransferableOut: &secp256k1fx.TransferOutput{
										Amt: 1,
										OutputOwners: secp256k1fx.OutputOwners{
											Locktime:  12345678,
											Threshold: 0,
											Addrs:     []ids.ShortID{},
										},
									},
								},
							},
							{
								Asset: avax.Asset{
									ID: customAssetID,
								},
								Out: &stakeable.LockOut{
									Locktime: 876543210,
									TransferableOut: &secp256k1fx.TransferOutput{
										Amt: 0xffffffffffffffff,
										OutputOwners: secp256k1fx.OutputOwners{
											Locktime:  0,
											Threshold: 1,
											Addrs: []ids.ShortID{
												addr,
											},
										},
									},
								},
							},
						},
						Ins: []*avax.TransferableInput{
							{
								UTXOID: avax.UTXOID{
									TxID:        txID,
									OutputIndex: 1,
								},
								Asset: avax.Asset{
									ID: avaxAssetID,
								},
								In: &secp256k1fx.TransferInput{
									Amt: units.Avax,
									Input: secp256k1fx.Input{
										SigIndices: []uint32{2, 5},
									},
								},
							},
							{
								UTXOID: avax.UTXOID{
									TxID:        txID,
									OutputIndex: 2,
								},
								Asset: avax.Asset{
									ID: customAssetID,
								},
								In: &stakeable.LockIn{
									Locktime: 876543210,
									TransferableIn: &secp256k1fx.TransferInput{
										Amt: 0xefffffffffffffff,
										Input: secp256k1fx.Input{
											SigIndices: []uint32{0},
										},
									},
								},
							},
							{
								UTXOID: avax.UTXOID{
									TxID:        txID,
									OutputIndex: 3,
								},
								Asset: avax.Asset{
									ID: customAssetID,
								},
								In: &secp256k1fx.TransferInput{
									Amt: 0x1000000000000000,
									Input: secp256k1fx.Input{
										SigIndices: []uint32{},
									},
								},
							},
						},
						Memo: types.JSONByteSlice("😅\nwell that's\x01\x23\x45!"),
					},
				},
				Subnet:  subnetID,
				ChainID: managerChainID,
				Address: managerAddress,
				SubnetAuth: &secp256k1fx.Input{
					SigIndices: []uint32{},
				},
			},
			expectedBytes: []byte{
				// Codec version
				0x00, 0x00,
				// ConvertSubnetTx Type ID
				0x00, 0x00, 0x00, 0x23,
				// Network ID
				0x00, 0x00, 0x00, 0x0a,
				// P-chain blockchain ID
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				// Number of outputs
				0x00, 0x00, 0x00, 0x02,
				// Outputs[0]
				// AVAX assetID
				0x21, 0xe6, 0x73, 0x17, 0xcb, 0xc4, 0xbe, 0x2a,
				0xeb, 0x00, 0x67, 0x7a, 0xd6, 0x46, 0x27, 0x78,
				0xa8, 0xf5, 0x22, 0x74, 0xb9, 0xd6, 0x05, 0xdf,
				0x25, 0x91, 0xb2, 0x30, 0x27, 0xa8, 0x7d, 0xff,
				// Stakeable locked output type ID
				0x00, 0x00, 0x00, 0x16,
				// Locktime
				0x00, 0x00, 0x00, 0x00, 0x05, 0x39, 0x7f, 0xb1,
				// secp256k1fx transfer output type ID
				0x00, 0x00, 0x00, 0x07,
				// amount
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				// secp256k1fx output locktime
				0x00, 0x00, 0x00, 0x00, 0x00, 0xbc, 0x61, 0x4e,
				// threshold
				0x00, 0x00, 0x00, 0x00,
				// number of addresses
				0x00, 0x00, 0x00, 0x00,
				// Outputs[1]
				// custom asset ID
				0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
				0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
				0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
				0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
				// Stakeable locked output type ID
				0x00, 0x00, 0x00, 0x16,
				// Locktime
				0x00, 0x00, 0x00, 0x00, 0x34, 0x3e, 0xfc, 0xea,
				// secp256k1fx transfer output type ID
				0x00, 0x00, 0x00, 0x07,
				// amount
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				// secp256k1fx output locktime
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				// threshold
				0x00, 0x00, 0x00, 0x01,
				// number of addresses
				0x00, 0x00, 0x00, 0x01,
				// address[0]
				0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb,
				0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb,
				0x44, 0x55, 0x66, 0x77,
				// number of inputs
				0x00, 0x00, 0x00, 0x03,
				// inputs[0]
				// TxID
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				// Tx output index
				0x00, 0x00, 0x00, 0x01,
				// AVAX assetID
				0x21, 0xe6, 0x73, 0x17, 0xcb, 0xc4, 0xbe, 0x2a,
				0xeb, 0x00, 0x67, 0x7a, 0xd6, 0x46, 0x27, 0x78,
				0xa8, 0xf5, 0x22, 0x74, 0xb9, 0xd6, 0x05, 0xdf,
				0x25, 0x91, 0xb2, 0x30, 0x27, 0xa8, 0x7d, 0xff,
				// secp256k1fx transfer input type ID
				0x00, 0x00, 0x00, 0x05,
				// input amount = 1 Avax
				0x00, 0x00, 0x00, 0x00, 0x3b, 0x9a, 0xca, 0x00,
				// number of signatures needed in input
				0x00, 0x00, 0x00, 0x02,
				// index of first signer
				0x00, 0x00, 0x00, 0x02,
				// index of second signer
				0x00, 0x00, 0x00, 0x05,
				// inputs[1]
				// TxID
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				// Tx output index
				0x00, 0x00, 0x00, 0x02,
				// Custom asset ID
				0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
				0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
				0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
				0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
				// Stakeable locked input type ID
				0x00, 0x00, 0x00, 0x15,
				// Locktime
				0x00, 0x00, 0x00, 0x00, 0x34, 0x3e, 0xfc, 0xea,
				// secp256k1fx transfer input type ID
				0x00, 0x00, 0x00, 0x05,
				// input amount
				0xef, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				// number of signatures needed in input
				0x00, 0x00, 0x00, 0x01,
				// index of signer
				0x00, 0x00, 0x00, 0x00,
				// inputs[2]
				// TxID
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				// Tx output index
				0x00, 0x00, 0x00, 0x03,
				// custom asset ID
				0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
				0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
				0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
				0x99, 0x77, 0x55, 0x77, 0x11, 0x33, 0x55, 0x31,
				// secp256k1fx transfer input type ID
				0x00, 0x00, 0x00, 0x05,
				// input amount
				0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				// number of signatures needed in input
				0x00, 0x00, 0x00, 0x00,
				// length of memo
				0x00, 0x00, 0x00, 0x14,
				// memo
				0xf0, 0x9f, 0x98, 0x85, 0x0a, 0x77, 0x65, 0x6c,
				0x6c, 0x20, 0x74, 0x68, 0x61, 0x74, 0x27, 0x73,
				0x01, 0x23, 0x45, 0x21,
				// subnetID to modify
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
				0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
				// chainID of the manager
				0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				// length of the manager address
				0x00, 0x00, 0x00, 0x14,
				// address of the manager
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0xde, 0xad,
				// secp256k1fx authorization type ID
				0x00, 0x00, 0x00, 0x0a,
				// number of signatures needed in authorization
				0x00, 0x00, 0x00, 0x00,
			},
			expectedJSON: convertSubnetTxComplexJSON,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			var unsignedTx UnsignedTx = test.tx
			txBytes, err := Codec.Marshal(CodecVersion, &unsignedTx)
			require.NoError(err)
			require.Equal(test.expectedBytes, txBytes)

			ctx := snowtest.Context(t, constants.PlatformChainID)
			test.tx.InitCtx(ctx)

			txJSON, err := json.MarshalIndent(test.tx, "", "\t")
			require.NoError(err)
			require.Equal(
				// Normalize newlines for Windows
				strings.ReplaceAll(string(test.expectedJSON), "\r\n", "\n"),
				string(txJSON),
			)
		})
	}
}

func TestConvertSubnetTxSyntacticVerify(t *testing.T) {
	var (
		ctx         = snowtest.Context(t, ids.GenerateTestID())
		validBaseTx = BaseTx{
			BaseTx: avax.BaseTx{
				NetworkID:    ctx.NetworkID,
				BlockchainID: ctx.ChainID,
			},
		}
		validSubnetID     = ids.GenerateTestID()
		invalidAddress    = make(types.JSONByteSlice, MaxSubnetAddressLength+1)
		validSubnetAuth   = &secp256k1fx.Input{}
		invalidSubnetAuth = &secp256k1fx.Input{
			SigIndices: []uint32{1, 0},
		}
	)

	tests := []struct {
		name        string
		tx          *ConvertSubnetTx
		expectedErr error
	}{
		{
			name:        "nil tx",
			tx:          nil,
			expectedErr: ErrNilTx,
		},
		{
			name: "already verified",
			// The tx includes invalid data to verify that a cached result is
			// returned.
			tx: &ConvertSubnetTx{
				BaseTx: BaseTx{
					SyntacticallyVerified: true,
				},
				Subnet:     constants.PrimaryNetworkID,
				Address:    invalidAddress,
				SubnetAuth: invalidSubnetAuth,
			},
			expectedErr: nil,
		},
		{
			name: "invalid subnetID",
			tx: &ConvertSubnetTx{
				BaseTx:     validBaseTx,
				Subnet:     constants.PrimaryNetworkID,
				SubnetAuth: validSubnetAuth,
			},
			expectedErr: ErrConvertPermissionlessSubnet,
		},
		{
			name: "invalid address",
			tx: &ConvertSubnetTx{
				BaseTx:     validBaseTx,
				Subnet:     validSubnetID,
				Address:    invalidAddress,
				SubnetAuth: validSubnetAuth,
			},
			expectedErr: ErrAddressTooLong,
		},
		{
			name: "invalid BaseTx",
			tx: &ConvertSubnetTx{
				BaseTx:     BaseTx{},
				Subnet:     validSubnetID,
				SubnetAuth: validSubnetAuth,
			},
			expectedErr: avax.ErrWrongNetworkID,
		},
		{
			name: "invalid subnetAuth",
			tx: &ConvertSubnetTx{
				BaseTx:     validBaseTx,
				Subnet:     validSubnetID,
				SubnetAuth: invalidSubnetAuth,
			},
			expectedErr: secp256k1fx.ErrInputIndicesNotSortedUnique,
		},
		{
			name: "passes verification",
			tx: &ConvertSubnetTx{
				BaseTx:     validBaseTx,
				Subnet:     validSubnetID,
				SubnetAuth: validSubnetAuth,
			},
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			err := test.tx.SyntacticVerify(ctx)
			require.ErrorIs(err, test.expectedErr)
			if test.expectedErr != nil {
				return
			}
			require.True(test.tx.SyntacticallyVerified)
		})
	}
}
