package base

import (
	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/callbacks"
	"github.com/perlin-network/noise/protocol"
)

var (
	_ protocol.HandshakePolicy = (*policy)(nil)
)

type policy struct {
}

func New() *policy {
	return &policy{}
}

func (p *policy) EnforceHandshakePolicy(node *noise.Node) {

}
func (p *policy) Opcodes() []noise.Opcode {
	return []noise.Opcode{}
}

func (p *policy) DoHandshake(peer *noise.Peer, opcode noise.Opcode, message noise.Message) error {
	return callbacks.DeregisterCallback
}

/*
func (p *policy) DoHandshake(peer *noise.Peer, opcode noise.Opcode, message noise.Message) error {
	if !peer.Has(keyEphemeralPrivateKey) {
		peer.Disconnect()
		return errors.New("peer attempted to perform ECDH with us even though we never have instantiated a handshake")
	}

	msg := message.(messageHandshake)

	peersPublicKey := p.suite.Point()
	err := peersPublicKey.UnmarshalBinary(msg.publicKey)

	if err != nil {
		peer.Disconnect()
		return errors.Wrap(err, "failed to unmarshal our peers ephemeral public key")
	}

	err = schnorr.Verify(p.suite, peersPublicKey, []byte(msgEphemeralHandshake), msg.signature)

	if err != nil {
		peer.Disconnect()
		return errors.Wrap(err, "failed to verify signature in handshake request")
	}

	ourPrivateKey := peer.Get(keyEphemeralPrivateKey).(kyber.Scalar)
	ephemeralSharedKey := computeSharedKey(p.suite, ourPrivateKey, peersPublicKey)

	log.Debug().Str("ephemeral_shared_key", ephemeralSharedKey.String()).Msg("Successfully performed ECDH with our peer.")

	sharedKeyBuf, err := ephemeralSharedKey.MarshalBinary()
	if err != nil {
		peer.Disconnect()
		return errors.Wrap(err, "failed to marshal post-handshake shared key")
	}

	peer.Delete(keyEphemeralPrivateKey)
	protocol.SetSharedKey(peer, sharedKeyBuf)

	if err = timeout.Clear(peer, keyTimeoutDispatcher); err != nil {
		peer.Disconnect()
		return errors.Wrap(err, "error enforcing handshake timeout policy")
	}

	return callbacks.DeregisterCallback
}

func (p *policy) onPeerInit(node *noise.Node, peer *noise.Peer) (err error) {
	if peer.Has(keyEphemeralPrivateKey) {
		peer.Disconnect()
		return errors.New("peer attempted to have us instantiate a 2nd handshake")
	}

	// Generate an ephemeral keypair to perform ECDH with our peer.
	ephemeralPrivateKey := p.suite.Scalar().Pick(p.suite.RandomStream())
	ephemeralPublicKey := p.suite.Point().Mul(ephemeralPrivateKey, p.suite.Point().Base())

	peer.Set(keyEphemeralPrivateKey, ephemeralPrivateKey)

	msg := messageHandshake{}
	msg.publicKey, err = ephemeralPublicKey.MarshalBinary()

	if err != nil {
		peer.Disconnect()
		return errors.Wrap(err, "failed to marshal ephemeral public key into binary")
	}

	msg.signature, err = schnorr.Sign(p.suite, ephemeralPrivateKey, []byte(msgEphemeralHandshake))

	if err != nil {
		peer.Disconnect()
		return errors.Wrap(err, "failed to sign handshake message using Schnorr signature scheme")
	}

	err = peer.SendMessage(OpcodeHandshake, msg)
	if err != nil {
		peer.Disconnect()
		return errors.Wrap(err, "failed to send our ephemeral public key to our peer")
	}

	timeout.Enforce(peer, keyTimeoutDispatcher, p.timeoutDuration, peer.Disconnect)

	return nil
}

func (p *policy) onPeerDisconnected(node *noise.Node, peer *noise.Peer) error {
	peer.Delete(keyEphemeralPrivateKey)
	protocol.DeleteSharedKey(peer)

	return callbacks.DeregisterCallback
}
*/
