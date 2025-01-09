// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"

	"github.com/MetalBlockchain/metalgo/api/info"
	"github.com/MetalBlockchain/metalgo/genesis"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/crypto/bls"
	"github.com/MetalBlockchain/metalgo/utils/set"
	"github.com/MetalBlockchain/metalgo/utils/units"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/warp"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/warp/message"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/warp/payload"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
	"github.com/MetalBlockchain/metalgo/wallet/subnet/primary"
)

func main() {
	key := genesis.EWOQKey
	uri := primary.LocalAPIURI
	kc := secp256k1fx.NewKeychain(key)
	subnetID := ids.FromStringOrPanic("2DeHa7Qb6sufPkmQcFWG2uCd4pBPv9WB6dkzroiMQhd1NSRtof")
	chainID := ids.FromStringOrPanic("2BMFrJ9xeh5JdwZEx6uuFcjfZC2SV2hdbMT8ee5HrvjtfJb5br")
	address := []byte{}
	weight := uint64(1)
	blsSKHex := "3f783929b295f16cd1172396acb23b20eed057b9afb1caa419e9915f92860b35"

	blsSKBytes, err := hex.DecodeString(blsSKHex)
	if err != nil {
		log.Fatalf("failed to decode secret key: %s\n", err)
	}

	sk, err := bls.SecretKeyFromBytes(blsSKBytes)
	if err != nil {
		log.Fatalf("failed to parse secret key: %s\n", err)
	}

	ctx := context.Background()
	infoClient := info.NewClient(uri)

	nodeInfoStartTime := time.Now()
	nodeID, nodePoP, err := infoClient.GetNodeID(ctx)
	if err != nil {
		log.Fatalf("failed to fetch node IDs: %s\n", err)
	}
	log.Printf("fetched node ID %s in %s\n", nodeID, time.Since(nodeInfoStartTime))

	// MakeWallet fetches the available UTXOs owned by [kc] on the network that
	// [uri] is hosting.
	walletSyncStartTime := time.Now()
	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:          uri,
		AVAXKeychain: kc,
		EthKeychain:  kc,
	})
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	// Get the P-chain wallet
	pWallet := wallet.P()
	context := pWallet.Builder().Context()

	expiry := uint64(time.Now().Add(5 * time.Minute).Unix()) // This message will expire in 5 minutes
	addressedCallPayload, err := message.NewRegisterL1Validator(
		subnetID,
		nodeID,
		nodePoP.PublicKey,
		expiry,
		message.PChainOwner{},
		message.PChainOwner{},
		weight,
	)
	if err != nil {
		log.Fatalf("failed to create RegisterL1Validator message: %s\n", err)
	}
	addressedCallPayloadJSON, err := json.MarshalIndent(addressedCallPayload, "", "\t")
	if err != nil {
		log.Fatalf("failed to marshal RegisterL1Validator message: %s\n", err)
	}
	log.Println(string(addressedCallPayloadJSON))

	addressedCall, err := payload.NewAddressedCall(
		address,
		addressedCallPayload.Bytes(),
	)
	if err != nil {
		log.Fatalf("failed to create AddressedCall message: %s\n", err)
	}

	unsignedWarp, err := warp.NewUnsignedMessage(
		context.NetworkID,
		chainID,
		addressedCall.Bytes(),
	)
	if err != nil {
		log.Fatalf("failed to create unsigned Warp message: %s\n", err)
	}

	// This example assumes that the hard-coded BLS key is for the first
	// validator in the signature bit-set.
	signers := set.NewBits(0)

	unsignedBytes := unsignedWarp.Bytes()
	sig := bls.Sign(sk, unsignedBytes)
	sigBytes := [bls.SignatureLen]byte{}
	copy(sigBytes[:], bls.SignatureToBytes(sig))

	warp, err := warp.NewMessage(
		unsignedWarp,
		&warp.BitSetSignature{
			Signers:   signers.Bytes(),
			Signature: sigBytes,
		},
	)
	if err != nil {
		log.Fatalf("failed to create Warp message: %s\n", err)
	}

	registerL1ValidatorStartTime := time.Now()
	registerL1ValidatorTx, err := pWallet.IssueRegisterL1ValidatorTx(
		units.Avax,
		nodePoP.ProofOfPossession,
		warp.Bytes(),
	)
	if err != nil {
		log.Fatalf("failed to issue register L1 validator transaction: %s\n", err)
	}

	validationID := addressedCallPayload.ValidationID()
	log.Printf("registered new L1 validator %s to subnetID %s with txID %s as validationID %s in %s\n", nodeID, subnetID, registerL1ValidatorTx.ID(), validationID, time.Since(registerL1ValidatorStartTime))
}