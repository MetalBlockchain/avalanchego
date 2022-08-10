// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/MetalBlockchain/avalanchego/utils/wrappers"
	"github.com/MetalBlockchain/avalanchego/vms/platformvm/blocks"
)

var _ blocks.Visitor = &blockMetrics{}

type blockMetrics struct {
	txMetrics *txMetrics

	numAbortBlocks,
	numAtomicBlocks,
	numCommitBlocks,
	numProposalBlocks,
	numStandardBlocks prometheus.Counter
}

func newBlockMetrics(
	namespace string,
	registerer prometheus.Registerer,
) (*blockMetrics, error) {
	txMetrics, err := newTxMetrics(namespace, registerer)
	errs := wrappers.Errs{Err: err}
	m := &blockMetrics{
		txMetrics:         txMetrics,
		numAbortBlocks:    newBlockMetric(namespace, "abort", registerer, &errs),
		numAtomicBlocks:   newBlockMetric(namespace, "atomic", registerer, &errs),
		numCommitBlocks:   newBlockMetric(namespace, "commit", registerer, &errs),
		numProposalBlocks: newBlockMetric(namespace, "proposal", registerer, &errs),
		numStandardBlocks: newBlockMetric(namespace, "standard", registerer, &errs),
	}
	return m, errs.Err
}

func newBlockMetric(
	namespace string,
	blockName string,
	registerer prometheus.Registerer,
	errs *wrappers.Errs,
) prometheus.Counter {
	blockMetric := prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      fmt.Sprintf("%s_blks_accepted", blockName),
		Help:      fmt.Sprintf("Number of %s blocks accepted", blockName),
	})
	errs.Add(registerer.Register(blockMetric))
	return blockMetric
}

func (m *blockMetrics) ProposalBlock(b *blocks.ProposalBlock) error {
	m.numProposalBlocks.Inc()
	return b.Tx.Unsigned.Visit(m.txMetrics)
}

func (m *blockMetrics) AtomicBlock(b *blocks.AtomicBlock) error {
	m.numAtomicBlocks.Inc()
	return b.Tx.Unsigned.Visit(m.txMetrics)
}

func (m *blockMetrics) StandardBlock(b *blocks.StandardBlock) error {
	m.numStandardBlocks.Inc()
	for _, tx := range b.Transactions {
		if err := tx.Unsigned.Visit(m.txMetrics); err != nil {
			return err
		}
	}
	return nil
}

func (m *blockMetrics) CommitBlock(*blocks.CommitBlock) error {
	m.numCommitBlocks.Inc()
	return nil
}

func (m *blockMetrics) AbortBlock(*blocks.AbortBlock) error {
	m.numAbortBlocks.Inc()
	return nil
}
