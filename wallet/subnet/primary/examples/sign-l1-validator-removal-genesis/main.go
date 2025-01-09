// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"log"
	"net/netip"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/proto"

	"github.com/MetalBlockchain/metalgo/api/info"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/network/p2p"
	"github.com/MetalBlockchain/metalgo/network/peer"
	"github.com/MetalBlockchain/metalgo/proto/pb/platformvm"
	"github.com/MetalBlockchain/metalgo/proto/pb/sdk"
	"github.com/MetalBlockchain/metalgo/snow/networking/router"
	"github.com/MetalBlockchain/metalgo/utils/compression"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/logging"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/warp"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/warp/payload"
	"github.com/MetalBlockchain/metalgo/wallet/subnet/primary"

	p2pmessage "github.com/MetalBlockchain/metalgo/message"
	warpmessage "github.com/MetalBlockchain/metalgo/vms/platformvm/warp/message"
)

func main() {
	uri := primary.LocalAPIURI
	subnetID := ids.FromStringOrPanic("2DeHa7Qb6sufPkmQcFWG2uCd4pBPv9WB6dkzroiMQhd1NSRtof")
	validationIndex := uint32(0)
	infoClient := info.NewClient(uri)
	networkID, err := infoClient.GetNetworkID(context.Background())
	if err != nil {
		log.Fatalf("failed to fetch network ID: %s\n", err)
	}

	validationID := subnetID.Append(validationIndex)
	l1ValidatorRegistration, err := warpmessage.NewL1ValidatorRegistration(
		validationID,
		false,
	)
	if err != nil {
		log.Fatalf("failed to create L1ValidatorRegistration message: %s\n", err)
	}

	addressedCall, err := payload.NewAddressedCall(
		nil,
		l1ValidatorRegistration.Bytes(),
	)
	if err != nil {
		log.Fatalf("failed to create AddressedCall message: %s\n", err)
	}

	unsignedWarp, err := warp.NewUnsignedMessage(
		networkID,
		constants.PlatformChainID,
		addressedCall.Bytes(),
	)
	if err != nil {
		log.Fatalf("failed to create unsigned Warp message: %s\n", err)
	}

	justification := platformvm.L1ValidatorRegistrationJustification{
		Preimage: &platformvm.L1ValidatorRegistrationJustification_ConvertSubnetToL1TxData{
			ConvertSubnetToL1TxData: &platformvm.SubnetIDIndex{
				SubnetId: subnetID[:],
				Index:    validationIndex,
			},
		},
	}
	justificationBytes, err := proto.Marshal(&justification)
	if err != nil {
		log.Fatalf("failed to create justification: %s\n", err)
	}

	p, err := peer.StartTestPeer(
		context.Background(),
		netip.AddrPortFrom(
			netip.AddrFrom4([4]byte{127, 0, 0, 1}),
			9651,
		),
		networkID,
		router.InboundHandlerFunc(func(_ context.Context, msg p2pmessage.InboundMessage) {
			log.Printf("received %s: %s", msg.Op(), msg.Message())
		}),
	)
	if err != nil {
		log.Fatalf("failed to start peer: %s\n", err)
	}

	messageBuilder, err := p2pmessage.NewCreator(
		logging.NoLog{},
		prometheus.NewRegistry(),
		compression.TypeZstd,
		time.Hour,
	)
	if err != nil {
		log.Fatalf("failed to create message builder: %s\n", err)
	}

	appRequestPayload, err := proto.Marshal(&sdk.SignatureRequest{
		Message:       unsignedWarp.Bytes(),
		Justification: justificationBytes,
	})
	if err != nil {
		log.Fatalf("failed to marshal SignatureRequest: %s\n", err)
	}

	appRequest, err := messageBuilder.AppRequest(
		constants.PlatformChainID,
		0,
		time.Hour,
		p2p.PrefixMessage(
			p2p.ProtocolPrefix(p2p.SignatureRequestHandlerID),
			appRequestPayload,
		),
	)
	if err != nil {
		log.Fatalf("failed to create AppRequest: %s\n", err)
	}

	p.Send(context.Background(), appRequest)

	time.Sleep(5 * time.Second)

	p.StartClose()
	err = p.AwaitClosed(context.Background())
	if err != nil {
		log.Fatalf("failed to close peer: %s\n", err)
	}
}