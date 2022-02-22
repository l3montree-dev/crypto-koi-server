package leader

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
)

type AlwaysLeader struct {
	ch        chan bool
	isLeader  bool
	listeners []Listener
	logger    *logrus.Entry
}

func NewAlwaysLeader() LeaderElection {
	return &AlwaysLeader{
		ch:       make(chan bool),
		isLeader: true,
		logger:   orchardclient.Logger.WithField("component", "AlwaysLeader"),
	}
}

func (a *AlwaysLeader) GetChannel() chan bool {
	return a.ch
}

func (a *AlwaysLeader) IsLeader() bool {
	return a.isLeader
}

func (a *AlwaysLeader) RunElection() {
	a.logger.Info("running leader election in always leader mode")
	// silence - we never need to emit any events into the channel
	for i, lst := range a.listeners {
		// just call start on each listener.
		// we wont ever need to stop it.
		a.logger.Info("starting listener: ", i)
		go lst.Start()
	}
}

func (a *AlwaysLeader) AddListener(listener Listener) {
	a.listeners = append(a.listeners, listener)
}
