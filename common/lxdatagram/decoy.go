package lxdatagram

import (
	"context"
	"net"

	"github.com/sagernet/sing/common/logger"
)

// Decoy is a disabled stand-in for the lx QUIC sniff decoy: this port
// does not carry the lx decoy option plumbing, so the decoy never
// engages and every attempt behaves as a plain direct attempt.
type Decoy struct{}

func NewSniffDecoy(ctx context.Context, outboundType string, logger logger.Logger) *Decoy {
	return &Decoy{}
}

func (d *Decoy) Enabled() bool { return false }

func (d *Decoy) Begin(ctx context.Context) *Attempt { return &Attempt{} }

// Attempt is one dial attempt under a Decoy; inert in this port.
type Attempt struct{}

func (a *Attempt) Dialed(conn net.Conn, writeTo net.Conn) bool { return false }

func (a *Attempt) Done(err error) (redoNow bool) { return false }

func (a *Attempt) WithDecoy() bool { return false }
