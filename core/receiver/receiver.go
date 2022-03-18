package receiver

import (
	"github.com/scmn-dev/tran/core/crypt"
	"github.com/scmn-dev/tran/models"
)

type Receiver struct {
	crypt             *crypt.Crypt
	payloadSize       int64
	tranxAddress string
	tranxPort    int
	ui                chan<- UIUpdate
	usedRelay         bool
}

func NewReceiver(programOptions models.TranOptions) *Receiver {
	return &Receiver{
		tranxAddress: programOptions.TranxAddress,
		tranxPort:    programOptions.TranxPort,
	}
}

func WithUI(r *Receiver, ui chan<- UIUpdate) *Receiver {
	r.ui = ui
	return r
}

func (r *Receiver) UsedRelay() bool {
	return r.usedRelay
}

func (r *Receiver) PayloadSize() int64 {
	return r.payloadSize
}

func (r *Receiver) TranxAddress() string {
	return r.tranxAddress
}

func (r *Receiver) TranxPort() int {
	return r.tranxPort
}

func (r *Receiver) updateUI(progress float32) {
	if r.ui == nil {
		return
	}

	r.ui <- UIUpdate{Progress: progress}
}
