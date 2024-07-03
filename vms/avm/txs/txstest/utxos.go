// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txstest

import (
	"context"
	"fmt"

	"github.com/MetalBlockchain/metalgo/chains/atomic"
	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/utils/set"
	"github.com/MetalBlockchain/metalgo/vms/avm/state"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/wallet/chain/x/builder"
	"github.com/MetalBlockchain/metalgo/wallet/chain/x/signer"
)

const maxPageSize uint64 = 1024

var (
	_ builder.Backend = (*walletUTXOsAdapter)(nil)
	_ signer.Backend  = (*walletUTXOsAdapter)(nil)
)

func newUTXOs(
	ctx *snow.Context,
	state state.State,
	sharedMemory atomic.SharedMemory,
	codec codec.Manager,
) *utxos {
	return &utxos{
		xchainID:     ctx.ChainID,
		state:        state,
		sharedMemory: sharedMemory,
		codec:        codec,
	}
}

type utxos struct {
	xchainID     ids.ID
	state        state.State
	sharedMemory atomic.SharedMemory
	codec        codec.Manager
}

func (u *utxos) UTXOs(addrs set.Set[ids.ShortID], sourceChainID ids.ID) ([]*avax.UTXO, error) {
	if sourceChainID == u.xchainID {
		return avax.GetAllUTXOs(u.state, addrs)
	}

	atomicUTXOs, _, _, err := avax.GetAtomicUTXOs(
		u.sharedMemory,
		u.codec,
		sourceChainID,
		addrs,
		ids.ShortEmpty,
		ids.Empty,
		int(maxPageSize),
	)
	return atomicUTXOs, err
}

func (u *utxos) GetUTXO(addrs set.Set[ids.ShortID], chainID, utxoID ids.ID) (*avax.UTXO, error) {
	if chainID == u.xchainID {
		return u.state.GetUTXO(utxoID)
	}

	atomicUTXOs, _, _, err := avax.GetAtomicUTXOs(
		u.sharedMemory,
		u.codec,
		chainID,
		addrs,
		ids.ShortEmpty,
		ids.Empty,
		int(maxPageSize),
	)
	if err != nil {
		return nil, fmt.Errorf("problem retrieving atomic UTXOs: %w", err)
	}
	for _, utxo := range atomicUTXOs {
		if utxo.InputID() == utxoID {
			return utxo, nil
		}
	}
	return nil, database.ErrNotFound
}

type walletUTXOsAdapter struct {
	utxos *utxos
	addrs set.Set[ids.ShortID]
}

func (w *walletUTXOsAdapter) UTXOs(_ context.Context, sourceChainID ids.ID) ([]*avax.UTXO, error) {
	return w.utxos.UTXOs(w.addrs, sourceChainID)
}

func (w *walletUTXOsAdapter) GetUTXO(_ context.Context, chainID, utxoID ids.ID) (*avax.UTXO, error) {
	return w.utxos.GetUTXO(w.addrs, chainID, utxoID)
}
