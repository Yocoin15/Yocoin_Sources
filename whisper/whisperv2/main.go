// Authored and revised by YOC team, 2014-2018
// License placeholder #1

// +build none

// Contains a simple whisper peer setup and self messaging to allow playing
// around with the protocol and API without a fancy client implementation.

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Yocoin15/Yocoin_Sources/common"
	"github.com/Yocoin15/Yocoin_Sources/crypto"
	"github.com/Yocoin15/Yocoin_Sources/logger"
	"github.com/Yocoin15/Yocoin_Sources/p2p"
	"github.com/Yocoin15/Yocoin_Sources/p2p/nat"
	"github.com/Yocoin15/Yocoin_Sources/whisper"
)

func main() {
	logger.AddLogSystem(logger.NewStdLogSystem(os.Stdout, log.LstdFlags, logger.InfoLevel))

	// Generate the peer identity
	key, err := crypto.GenerateKey()
	if err != nil {
		fmt.Printf("Failed to generate peer key: %v.\n", err)
		os.Exit(-1)
	}
	name := common.MakeName("whisper-go", "1.0")
	shh := whisper.New()

	// Create an YOC peer to communicate through
	server := p2p.Server{
		PrivateKey: key,
		MaxPeers:   10,
		Name:       name,
		Protocols:  []p2p.Protocol{shh.Protocol()},
		ListenAddr: ":30300",
		NAT:        nat.Any(),
	}
	fmt.Println("Starting YOC peer...")
	if err := server.Start(); err != nil {
		fmt.Printf("Failed to start YOC peer: %v.\n", err)
		os.Exit(1)
	}

	// Send a message to self to check that something works
	payload := fmt.Sprintf("Hello world, this is %v. In case you're wondering, the time is %v", name, time.Now())
	if err := selfSend(shh, []byte(payload)); err != nil {
		fmt.Printf("Failed to self message: %v.\n", err)
		os.Exit(-1)
	}
}

// SendSelf wraps a payload into a Whisper envelope and forwards it to itself.
func selfSend(shh *whisper.Whisper, payload []byte) error {
	ok := make(chan struct{})

	// Start watching for self messages, output any arrivals
	id := shh.NewIdentity()
	shh.Watch(whisper.Filter{
		To: &id.PublicKey,
		Fn: func(msg *whisper.Message) {
			fmt.Printf("Message received: %s, signed with 0x%x.\n", string(msg.Payload), msg.Signature)
			close(ok)
		},
	})
	// Wrap the payload and encrypt it
	msg := whisper.NewMessage(payload)
	envelope, err := msg.Wrap(whisper.DefaultPoW, whisper.Options{
		From: id,
		To:   &id.PublicKey,
		TTL:  whisper.DefaultTTL,
	})
	if err != nil {
		return fmt.Errorf("failed to seal message: %v", err)
	}
	// Dump the message into the system and wait for it to pop back out
	if err := shh.Send(envelope); err != nil {
		return fmt.Errorf("failed to send self-message: %v", err)
	}
	select {
	case <-ok:
	case <-time.After(time.Second):
		return fmt.Errorf("failed to receive message in time")
	}
	return nil
}
