// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements X-chain transfer tests.
package transfer

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"github.com/MetalBlockchain/metalgo/chains"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/tests"
	"github.com/MetalBlockchain/metalgo/tests/fixture/e2e"
	"github.com/MetalBlockchain/metalgo/tests/fixture/tmpnet"
	"github.com/MetalBlockchain/metalgo/utils/crypto/secp256k1"
	"github.com/MetalBlockchain/metalgo/utils/set"
	"github.com/MetalBlockchain/metalgo/utils/units"
	"github.com/MetalBlockchain/metalgo/vms/avm"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
	"github.com/MetalBlockchain/metalgo/wallet/subnet/primary"
	"github.com/MetalBlockchain/metalgo/wallet/subnet/primary/common"
)

const (
	totalRounds = 50

	blksProcessingMetric = "metal_snowman_blks_processing"
	blksAcceptedMetric   = "metal_snowman_blks_accepted_count"
)

var xChainMetricLabels = prometheus.Labels{
	chains.ChainLabel: "X",
}

// This test requires that the network not have ongoing blocks and
// cannot reliably be run in parallel.
var _ = e2e.DescribeXChainSerial("[Virtuous Transfer Tx AVAX]", func() {
	tc := e2e.NewTestContext()
	require := require.New(tc)

	ginkgo.It("can issue a virtuous transfer tx for AVAX asset",
		func() {
			env := e2e.GetEnv(tc)
			rpcEps := make([]string, len(env.URIs))
			for i, nodeURI := range env.URIs {
				rpcEps[i] = nodeURI.URI
			}

			// Waiting for ongoing blocks to have completed before starting this
			// test avoids the case of a previous test having initiated block
			// processing but not having completed it.
			tc.Eventually(func() bool {
				allNodeMetrics, err := tests.GetNodesMetrics(
					tc.DefaultContext(),
					rpcEps,
				)
				require.NoError(err)

				for _, metrics := range allNodeMetrics {
					xBlksProcessing, ok := tests.GetMetricValue(metrics, blksProcessingMetric, xChainMetricLabels)
					if !ok || xBlksProcessing > 0 {
						return false
					}
				}
				return true
			},
				e2e.DefaultTimeout,
				e2e.DefaultPollingInterval,
				"The cluster is generating ongoing blocks. Is this test being run in parallel?",
			)

			// Ensure the same set of 10 keys is used for all tests
			// by retrieving them outside of runFunc.
			testKeys := []*secp256k1.PrivateKey{
				// The funded key will be the source of funds for the new keys
				env.PreFundedKey,
			}
			newKeys, err := tmpnet.NewPrivateKeys(9)
			require.NoError(err)
			testKeys = append(testKeys, newKeys...)

			const transferPerRound = units.MilliAvax

			tc.By("Funding new keys")
			fundingWallet := e2e.NewWallet(tc, env.NewKeychain(), env.GetRandomNodeURI())
			fundingOutputs := make([]*avax.TransferableOutput, len(newKeys))
			fundingAssetID := fundingWallet.X().Builder().Context().AVAXAssetID
			for i, key := range newKeys {
				fundingOutputs[i] = &avax.TransferableOutput{
					Asset: avax.Asset{
						ID: fundingAssetID,
					},
					Out: &secp256k1fx.TransferOutput{
						// Enough for 1 transfer per round
						Amt: totalRounds * transferPerRound,
						OutputOwners: secp256k1fx.OutputOwners{
							Threshold: 1,
							Addrs: []ids.ShortID{
								key.Address(),
							},
						},
					},
				}
			}
			_, err = fundingWallet.X().IssueBaseTx(
				fundingOutputs,
				tc.WithDefaultContext(),
			)
			require.NoError(err)

			runFunc := func(round int) {
				tc.Outf("{{green}}\n\n\n\n\n\n---\n[ROUND #%02d]:{{/}}\n", round)

				needPermute := round > 3
				if needPermute {
					rand.Seed(time.Now().UnixNano())
					rand.Shuffle(len(testKeys), func(i, j int) {
						testKeys[i], testKeys[j] = testKeys[j], testKeys[i]
					})
				}

				keychain := secp256k1fx.NewKeychain(testKeys...)
				baseWallet := e2e.NewWallet(tc, keychain, env.GetRandomNodeURI())
				xWallet := baseWallet.X()
				xBuilder := xWallet.Builder()
				xContext := xBuilder.Context()
				avaxAssetID := xContext.AVAXAssetID

				wallets := make([]primary.Wallet, len(testKeys))
				shortAddrs := make([]ids.ShortID, len(testKeys))
				for i := range wallets {
					shortAddrs[i] = testKeys[i].PublicKey().Address()

					wallets[i] = primary.NewWalletWithOptions(
						baseWallet,
						common.WithCustomAddresses(set.Of(
							testKeys[i].PublicKey().Address(),
						)),
					)
				}

				metricsBeforeTx, err := tests.GetNodesMetrics(
					tc.DefaultContext(),
					rpcEps,
				)
				require.NoError(err)
				for _, uri := range rpcEps {
					for _, metric := range []string{blksProcessingMetric, blksAcceptedMetric} {
						tc.Outf("{{green}}%s at %q:{{/}} %v\n", metric, uri, metricsBeforeTx[uri][metric])
					}
				}

				testBalances := make([]uint64, 0)
				for i, w := range wallets {
					balances, err := w.X().Builder().GetFTBalance()
					require.NoError(err)

					bal := balances[avaxAssetID]
					testBalances = append(testBalances, bal)

					fmt.Printf(`CURRENT BALANCE %21d AVAX (SHORT ADDRESS %q)
`,
						bal,
						testKeys[i].PublicKey().Address(),
					)
				}
				fromIdx := -1
				for i := range testBalances {
					if fromIdx < 0 && testBalances[i] > 0 {
						fromIdx = i
						break
					}
				}
				require.GreaterOrEqual(fromIdx, 0, "no address found with non-zero balance")

				toIdx := -1
				for i := range testBalances {
					// prioritize the address with zero balance
					if toIdx < 0 && i != fromIdx && testBalances[i] == 0 {
						toIdx = i
						break
					}
				}
				if toIdx < 0 {
					// no zero balance address, so just transfer between any two addresses
					toIdx = (fromIdx + 1) % len(testBalances)
				}

				senderOrigBal := testBalances[fromIdx]
				receiverOrigBal := testBalances[toIdx]
				amountToTransfer := transferPerRound
				senderNewBal := senderOrigBal - amountToTransfer - xContext.BaseTxFee
				receiverNewBal := receiverOrigBal + amountToTransfer

				tc.By("X-Chain transfer with wrong amount must fail", func() {
					_, err := wallets[fromIdx].X().IssueBaseTx(
						[]*avax.TransferableOutput{{
							Asset: avax.Asset{
								ID: avaxAssetID,
							},
							Out: &secp256k1fx.TransferOutput{
								Amt: senderOrigBal + 1,
								OutputOwners: secp256k1fx.OutputOwners{
									Threshold: 1,
									Addrs:     []ids.ShortID{shortAddrs[toIdx]},
								},
							},
						}},
						tc.WithDefaultContext(),
					)
					require.Contains(err.Error(), "insufficient funds")
				})

				fmt.Printf(`===
TRANSFERRING

FROM [%q]
SENDER    CURRENT BALANCE     : %21d AVAX
SENDER    NEW BALANCE (AFTER) : %21d AVAX

TRANSFER AMOUNT FROM SENDER   : %21d AVAX

TO [%q]
RECEIVER  CURRENT BALANCE     : %21d AVAX
RECEIVER  NEW BALANCE (AFTER) : %21d AVAX
===
`,
					shortAddrs[fromIdx],
					senderOrigBal,
					senderNewBal,
					amountToTransfer,
					shortAddrs[toIdx],
					receiverOrigBal,
					receiverNewBal,
				)

				tx, err := wallets[fromIdx].X().IssueBaseTx(
					[]*avax.TransferableOutput{{
						Asset: avax.Asset{
							ID: avaxAssetID,
						},
						Out: &secp256k1fx.TransferOutput{
							Amt: amountToTransfer,
							OutputOwners: secp256k1fx.OutputOwners{
								Threshold: 1,
								Addrs:     []ids.ShortID{shortAddrs[toIdx]},
							},
						},
					}},
					tc.WithDefaultContext(),
				)
				require.NoError(err)

				balances, err := wallets[fromIdx].X().Builder().GetFTBalance()
				require.NoError(err)
				senderCurBalX := balances[avaxAssetID]
				tc.Outf("{{green}}first wallet balance:{{/}}  %d\n", senderCurBalX)

				balances, err = wallets[toIdx].X().Builder().GetFTBalance()
				require.NoError(err)
				receiverCurBalX := balances[avaxAssetID]
				tc.Outf("{{green}}second wallet balance:{{/}} %d\n", receiverCurBalX)

				require.Equal(senderCurBalX, senderNewBal)
				require.Equal(receiverCurBalX, receiverNewBal)

				txID := tx.ID()
				for _, u := range rpcEps {
					xc := avm.NewClient(u, "X")
					require.NoError(avm.AwaitTxAccepted(xc, tc.DefaultContext(), txID, 2*time.Second))
				}

				for _, u := range rpcEps {
					xc := avm.NewClient(u, "X")
					require.NoError(avm.AwaitTxAccepted(xc, tc.DefaultContext(), txID, 2*time.Second))

					mm, err := tests.GetNodeMetrics(tc.DefaultContext(), u)
					require.NoError(err)

					prev := metricsBeforeTx[u]

					// +0 since X-chain tx must have been processed and accepted
					// by now
					currentXBlksProcessing, _ := tests.GetMetricValue(mm, blksProcessingMetric, xChainMetricLabels)
					previousXBlksProcessing, _ := tests.GetMetricValue(prev, blksProcessingMetric, xChainMetricLabels)
					require.Equal(currentXBlksProcessing, previousXBlksProcessing)

					// +1 since X-chain tx must have been accepted by now
					currentXBlksAccepted, _ := tests.GetMetricValue(mm, blksAcceptedMetric, xChainMetricLabels)
					previousXBlksAccepted, _ := tests.GetMetricValue(prev, blksAcceptedMetric, xChainMetricLabels)
					require.Equal(currentXBlksAccepted, previousXBlksAccepted+1)

					metricsBeforeTx[u] = mm
				}
			}

			for i := 0; i < totalRounds; i++ {
				runFunc(i)
				time.Sleep(time.Second)
			}
		})
})
