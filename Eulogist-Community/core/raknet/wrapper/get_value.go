package raknet_wrapper

import (
	"Eulogist/core/minecraft/protocol/packet"
	"context"
	"crypto/ecdsa"
	"sync/atomic"
)

// ...
func (r *Raknet) GetEncoder() *packet.Encoder {
	return r.encoder
}

// ...
func (r *Raknet) GetDecoder() *packet.Decoder {
	return r.decoder
}

// ...
func (r *Raknet) GetKey() *ecdsa.PrivateKey {
	return r.key
}

// ...
func (r *Raknet) GetSalt() []byte {
	return r.salt
}

// ...
func (r *Raknet) GetShieldID() *atomic.Int32 {
	return &r.shieldID
}

// ...
func (r *Raknet) GetContext() context.Context {
	return r.context
}
