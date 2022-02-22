package leader

type AlwaysLeader struct {
	ch        chan bool
	isLeader  bool
	listeners []Listener
}

func NewAlwaysLeader() LeaderElection {
	return &AlwaysLeader{
		ch:       make(chan bool),
		isLeader: true,
	}
}

func (a *AlwaysLeader) GetChannel() chan bool {
	return a.ch
}

func (a *AlwaysLeader) IsLeader() bool {
	return a.isLeader
}

func (a *AlwaysLeader) RunElection() {
	// silence - we never need to emit any events into the channel
	for _, lst := range a.listeners {
		// just call start on each listener.
		// we wont ever need to stop it.
		go lst.Start()
	}
}

func (a *AlwaysLeader) AddListener(listener Listener) {
	a.listeners = append(a.listeners, listener)
}
