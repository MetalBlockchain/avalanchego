// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validatorstest

import (
	"context"

	"github.com/MetalBlockchain/metalgo/ids"

	snowvalidators "github.com/MetalBlockchain/metalgo/snow/validators"
	vmvalidators "github.com/MetalBlockchain/metalgo/vms/platformvm/validators"
)

var Manager vmvalidators.Manager = manager{}

type manager struct{}

func (manager) GetMinimumHeight(context.Context) (uint64, error) {
	return 0, nil
}

func (manager) GetCurrentHeight(context.Context) (uint64, error) {
	return 0, nil
}

func (manager) GetSubnetID(context.Context, ids.ID) (ids.ID, error) {
	return ids.Empty, nil
}

func (manager) GetValidatorSet(context.Context, uint64, ids.ID) (map[ids.NodeID]*snowvalidators.GetValidatorOutput, error) {
	return nil, nil
}

func (manager) OnAcceptedBlockID(ids.ID) {}

func (manager) GetCurrentValidatorSet(context.Context, ids.ID) (map[ids.ID]*snowvalidators.GetCurrentValidatorOutput, uint64, error) {
	return nil, 0, nil
}