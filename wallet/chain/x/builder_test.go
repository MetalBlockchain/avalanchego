// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package x

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/crypto/secp256k1"
	"github.com/MetalBlockchain/metalgo/utils/set"
	"github.com/MetalBlockchain/metalgo/utils/units"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/components/verify"
	"github.com/MetalBlockchain/metalgo/vms/nftfx"
	"github.com/MetalBlockchain/metalgo/vms/propertyfx"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
	"github.com/MetalBlockchain/metalgo/wallet/subnet/primary/common"
)

var (
	testKeys = secp256k1.TestKeys()

	// We hard-code [avaxAssetID] and [subnetAssetID] to make
	// ordering of UTXOs generated by [testUTXOsList] is reproducible
	avaxAssetID     = ids.Empty.Prefix(1789)
	xChainID        = ids.Empty.Prefix(2021)
	nftAssetID      = ids.Empty.Prefix(2022)
	propertyAssetID = ids.Empty.Prefix(2023)

	testCtx = NewContext(
		constants.UnitTestID,
		xChainID,
		avaxAssetID,
		units.MicroAvax,    // BaseTxFee
		99*units.MilliAvax, // CreateAssetTxFee
	)
)

// These tests create and sign a tx, then verify that utxos included
// in the tx are exactly necessary to pay fees for it

func TestBaseTx(t *testing.T) {
	var (
		require = require.New(t)

		// backend
		utxosKey       = testKeys[1]
		utxos          = makeTestUTXOs(utxosKey)
		genericBackend = common.NewDeterministicChainUTXOs(
			require,
			map[ids.ID][]*avax.UTXO{
				xChainID: utxos,
			},
		)
		backend = NewBackend(testCtx, genericBackend)

		// builder
		utxoAddr = utxosKey.Address()
		builder  = NewBuilder(set.Of(utxoAddr), backend)

		// data to build the transaction
		outputsToMove = []*avax.TransferableOutput{{
			Asset: avax.Asset{ID: avaxAssetID},
			Out: &secp256k1fx.TransferOutput{
				Amt: 7 * units.Avax,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{utxoAddr},
				},
			},
		}}
	)

	utx, err := builder.NewBaseTx(
		outputsToMove,
	)
	require.NoError(err)

	// check UTXOs selection and fee financing
	ins := utx.Ins
	outs := utx.Outs
	require.Len(ins, 2)
	require.Len(outs, 2)

	expectedConsumed := testCtx.BaseTxFee()
	consumed := ins[0].In.Amount() + ins[1].In.Amount() - outs[0].Out.Amount() - outs[1].Out.Amount()
	require.Equal(expectedConsumed, consumed)
	require.Equal(outputsToMove[0], outs[1])
}

func TestCreateAssetTx(t *testing.T) {
	require := require.New(t)

	var (
		// backend
		utxosKey       = testKeys[1]
		utxos          = makeTestUTXOs(utxosKey)
		genericBackend = common.NewDeterministicChainUTXOs(
			require,
			map[ids.ID][]*avax.UTXO{
				xChainID: utxos,
			},
		)
		backend = NewBackend(testCtx, genericBackend)

		// builder
		utxoAddr = utxosKey.Address()
		builder  = NewBuilder(set.Of(utxoAddr), backend)

		// data to build the transaction
		assetName          = "Team Rocket"
		symbol             = "TR"
		denomination uint8 = 0
		initialState       = map[uint32][]verify.State{
			0: {
				&secp256k1fx.MintOutput{
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{testKeys[0].PublicKey().Address()},
					},
				}, &secp256k1fx.MintOutput{
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{testKeys[0].PublicKey().Address()},
					},
				},
			},
			1: {
				&nftfx.MintOutput{
					GroupID: 1,
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{testKeys[1].PublicKey().Address()},
					},
				},
				&nftfx.MintOutput{
					GroupID: 2,
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{testKeys[1].PublicKey().Address()},
					},
				},
			},
			2: {
				&propertyfx.MintOutput{
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{testKeys[2].PublicKey().Address()},
					},
				},
				&propertyfx.MintOutput{
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{testKeys[2].PublicKey().Address()},
					},
				},
			},
		}
	)

	utx, err := builder.NewCreateAssetTx(
		assetName,
		symbol,
		denomination,
		initialState,
	)
	require.NoError(err)

	// check UTXOs selection and fee financing
	ins := utx.Ins
	outs := utx.Outs
	require.Len(ins, 2)
	require.Len(outs, 1)

	expectedConsumed := testCtx.CreateAssetTxFee()
	consumed := ins[0].In.Amount() + ins[1].In.Amount() - outs[0].Out.Amount()
	require.Equal(expectedConsumed, consumed)
}

func TestMintNFTOperation(t *testing.T) {
	require := require.New(t)

	var (
		// backend
		utxosKey       = testKeys[1]
		utxos          = makeTestUTXOs(utxosKey)
		genericBackend = common.NewDeterministicChainUTXOs(
			require,
			map[ids.ID][]*avax.UTXO{
				xChainID: utxos,
			},
		)
		backend = NewBackend(testCtx, genericBackend)

		// builder
		utxoAddr = utxosKey.Address()
		builder  = NewBuilder(set.Of(utxoAddr), backend)

		// data to build the transaction
		payload  = []byte{'h', 'e', 'l', 'l', 'o'}
		NFTOwner = &secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{utxoAddr},
		}
	)

	utx, err := builder.NewOperationTxMintNFT(
		nftAssetID,
		payload,
		[]*secp256k1fx.OutputOwners{NFTOwner},
	)
	require.NoError(err)

	// check UTXOs selection and fee financing
	ins := utx.Ins
	outs := utx.Outs
	require.Len(ins, 1)
	require.Len(outs, 1)

	expectedConsumed := testCtx.BaseTxFee()
	consumed := ins[0].In.Amount() - outs[0].Out.Amount()
	require.Equal(expectedConsumed, consumed)
}

func TestMintFTOperation(t *testing.T) {
	require := require.New(t)

	var (
		// backend
		utxosKey       = testKeys[1]
		utxos          = makeTestUTXOs(utxosKey)
		genericBackend = common.NewDeterministicChainUTXOs(
			require,
			map[ids.ID][]*avax.UTXO{
				xChainID: utxos,
			},
		)
		backend = NewBackend(testCtx, genericBackend)

		// builder
		utxoAddr = utxosKey.Address()
		builder  = NewBuilder(set.Of(utxoAddr), backend)

		// data to build the transaction
		outputs = map[ids.ID]*secp256k1fx.TransferOutput{
			nftAssetID: {
				Amt: 1,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{utxoAddr},
				},
			},
		}
	)

	utx, err := builder.NewOperationTxMintFT(
		outputs,
	)
	require.NoError(err)

	// check UTXOs selection and fee financing
	ins := utx.Ins
	outs := utx.Outs
	require.Len(ins, 1)
	require.Len(outs, 1)

	expectedConsumed := testCtx.BaseTxFee()
	consumed := ins[0].In.Amount() - outs[0].Out.Amount()
	require.Equal(expectedConsumed, consumed)
}

func TestMintPropertyOperation(t *testing.T) {
	require := require.New(t)

	var (
		// backend
		utxosKey       = testKeys[1]
		utxos          = makeTestUTXOs(utxosKey)
		genericBackend = common.NewDeterministicChainUTXOs(
			require,
			map[ids.ID][]*avax.UTXO{
				xChainID: utxos,
			},
		)
		backend = NewBackend(testCtx, genericBackend)

		// builder
		utxoAddr = utxosKey.Address()
		builder  = NewBuilder(set.Of(utxoAddr), backend)

		// data to build the transaction
		propertyOwner = &secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs:     []ids.ShortID{utxoAddr},
		}
	)

	utx, err := builder.NewOperationTxMintProperty(
		propertyAssetID,
		propertyOwner,
	)
	require.NoError(err)

	// check UTXOs selection and fee financing
	ins := utx.Ins
	outs := utx.Outs
	require.Len(ins, 1)
	require.Len(outs, 1)

	expectedConsumed := testCtx.BaseTxFee()
	consumed := ins[0].In.Amount() - outs[0].Out.Amount()
	require.Equal(expectedConsumed, consumed)
}

func TestBurnPropertyOperation(t *testing.T) {
	require := require.New(t)

	var (
		// backend
		utxosKey       = testKeys[1]
		utxos          = makeTestUTXOs(utxosKey)
		genericBackend = common.NewDeterministicChainUTXOs(
			require,
			map[ids.ID][]*avax.UTXO{
				xChainID: utxos,
			},
		)
		backend = NewBackend(testCtx, genericBackend)

		// builder
		utxoAddr = utxosKey.Address()
		builder  = NewBuilder(set.Of(utxoAddr), backend)
	)

	utx, err := builder.NewOperationTxBurnProperty(
		propertyAssetID,
	)
	require.NoError(err)

	// check UTXOs selection and fee financing
	ins := utx.Ins
	outs := utx.Outs
	require.Len(ins, 1)
	require.Len(outs, 1)

	expectedConsumed := testCtx.BaseTxFee()
	consumed := ins[0].In.Amount() - outs[0].Out.Amount()
	require.Equal(expectedConsumed, consumed)
}

func TestImportTx(t *testing.T) {
	var (
		require = require.New(t)

		// backend
		utxosKey       = testKeys[1]
		utxos          = makeTestUTXOs(utxosKey)
		sourceChainID  = ids.GenerateTestID()
		importedUTXOs  = utxos[:1]
		genericBackend = common.NewDeterministicChainUTXOs(
			require,
			map[ids.ID][]*avax.UTXO{
				xChainID:      utxos,
				sourceChainID: importedUTXOs,
			},
		)

		backend = NewBackend(testCtx, genericBackend)

		// builder
		utxoAddr = utxosKey.Address()
		builder  = NewBuilder(set.Of(utxoAddr), backend)

		// data to build the transaction
		importKey = testKeys[0]
		importTo  = &secp256k1fx.OutputOwners{
			Threshold: 1,
			Addrs: []ids.ShortID{
				importKey.Address(),
			},
		}
	)

	utx, err := builder.NewImportTx(
		sourceChainID,
		importTo,
	)
	require.NoError(err)

	// check UTXOs selection and fee financing
	ins := utx.Ins
	outs := utx.Outs
	importedIns := utx.ImportedIns
	require.Empty(ins)
	require.Len(importedIns, 1)
	require.Len(outs, 1)

	expectedConsumed := testCtx.BaseTxFee()
	consumed := importedIns[0].In.Amount() - outs[0].Out.Amount()
	require.Equal(expectedConsumed, consumed)
}

func TestExportTx(t *testing.T) {
	var (
		require = require.New(t)

		// backend
		utxosKey       = testKeys[1]
		utxos          = makeTestUTXOs(utxosKey)
		genericBackend = common.NewDeterministicChainUTXOs(
			require,
			map[ids.ID][]*avax.UTXO{
				xChainID: utxos,
			},
		)
		backend = NewBackend(testCtx, genericBackend)

		// builder
		utxoAddr = utxosKey.Address()
		builder  = NewBuilder(set.Of(utxoAddr), backend)

		// data to build the transaction
		subnetID        = ids.GenerateTestID()
		exportedOutputs = []*avax.TransferableOutput{{
			Asset: avax.Asset{ID: avaxAssetID},
			Out: &secp256k1fx.TransferOutput{
				Amt: 7 * units.Avax,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{utxoAddr},
				},
			},
		}}
	)

	utx, err := builder.NewExportTx(
		subnetID,
		exportedOutputs,
	)
	require.NoError(err)

	// check UTXOs selection and fee financing
	ins := utx.Ins
	outs := utx.Outs
	require.Len(ins, 2)
	require.Len(outs, 1)

	expectedConsumed := testCtx.BaseTxFee() + exportedOutputs[0].Out.Amount()
	consumed := ins[0].In.Amount() + ins[1].In.Amount() - outs[0].Out.Amount()
	require.Equal(expectedConsumed, consumed)
	require.Equal(utx.ExportedOuts, exportedOutputs)
}

func makeTestUTXOs(utxosKey *secp256k1.PrivateKey) []*avax.UTXO {
	// Note: we avoid ids.GenerateTestNodeID here to make sure that UTXO IDs won't change
	// run by run. This simplifies checking what utxos are included in the built txs.
	const utxosOffset uint64 = 2024

	return []*avax.UTXO{ // currently, the wallet scans UTXOs in the order provided here
		{ // a small UTXO first, which  should not be enough to pay fees
			UTXOID: avax.UTXOID{
				TxID:        ids.Empty.Prefix(utxosOffset),
				OutputIndex: uint32(utxosOffset),
			},
			Asset: avax.Asset{ID: avaxAssetID},
			Out: &secp256k1fx.TransferOutput{
				Amt: 2 * units.MilliAvax,
				OutputOwners: secp256k1fx.OutputOwners{
					Locktime:  0,
					Addrs:     []ids.ShortID{utxosKey.PublicKey().Address()},
					Threshold: 1,
				},
			},
		},
		{
			UTXOID: avax.UTXOID{
				TxID:        ids.Empty.Prefix(utxosOffset + 2),
				OutputIndex: uint32(utxosOffset + 2),
			},
			Asset: avax.Asset{ID: nftAssetID},
			Out: &nftfx.MintOutput{
				GroupID: 1,
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{utxosKey.PublicKey().Address()},
				},
			},
		},
		{
			UTXOID: avax.UTXOID{
				TxID:        ids.Empty.Prefix(utxosOffset + 3),
				OutputIndex: uint32(utxosOffset + 3),
			},
			Asset: avax.Asset{ID: nftAssetID},
			Out: &secp256k1fx.MintOutput{
				OutputOwners: secp256k1fx.OutputOwners{
					Threshold: 1,
					Addrs:     []ids.ShortID{utxosKey.PublicKey().Address()},
				},
			},
		},
		{
			UTXOID: avax.UTXOID{
				TxID:        ids.Empty.Prefix(utxosOffset + 4),
				OutputIndex: uint32(utxosOffset + 4),
			},
			Asset: avax.Asset{ID: propertyAssetID},
			Out: &propertyfx.MintOutput{
				OutputOwners: secp256k1fx.OutputOwners{
					Locktime:  0,
					Addrs:     []ids.ShortID{utxosKey.PublicKey().Address()},
					Threshold: 1,
				},
			},
		},
		{
			UTXOID: avax.UTXOID{
				TxID:        ids.Empty.Prefix(utxosOffset + 5),
				OutputIndex: uint32(utxosOffset + 5),
			},
			Asset: avax.Asset{ID: propertyAssetID},
			Out: &propertyfx.OwnedOutput{
				OutputOwners: secp256k1fx.OutputOwners{
					Locktime:  0,
					Addrs:     []ids.ShortID{utxosKey.PublicKey().Address()},
					Threshold: 1,
				},
			},
		},
		{ // a large UTXO last, which should be enough to pay any fee by itself
			UTXOID: avax.UTXOID{
				TxID:        ids.Empty.Prefix(utxosOffset + 6),
				OutputIndex: uint32(utxosOffset + 6),
			},
			Asset: avax.Asset{ID: avaxAssetID},
			Out: &secp256k1fx.TransferOutput{
				Amt: 9 * units.Avax,
				OutputOwners: secp256k1fx.OutputOwners{
					Locktime:  0,
					Addrs:     []ids.ShortID{utxosKey.PublicKey().Address()},
					Threshold: 1,
				},
			},
		},
	}
}
