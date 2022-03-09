package leaderelection

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

type KubernetesLeaderElection struct {
	lockName      string
	lockNamespace string
	podName       string
	isLeader      bool
	ctx           context.Context
	client        *kubernetes.Clientset
	listeners     []Listener
}

type Listener interface {
	// should be idempotent
	Start()
	// should be idempotent
	Stop()
	IsRunning() bool
}

func NewKubernetesLeaderElection(ctx context.Context, lockName, lockNamespace string) *KubernetesLeaderElection {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	// get the pod name of the current pod
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		log.Fatal("failed to retrieve pod name. Added pod name as environment variable?")
	}

	client := kubernetes.NewForConfigOrDie(config)
	return &KubernetesLeaderElection{
		lockName:      lockName,
		lockNamespace: lockNamespace,
		ctx:           ctx,
		client:        client,
		podName:       podName,
	}
}

// start the election. This method will block. Start it on a goroutine.
func (k *KubernetesLeaderElection) RunElection() {
	k.logger.Info("running leader election")
	lock := k.getNewLock()
	leaderelection.RunOrDie(k.ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(c context.Context) {
				k.logger.Info("elected as the leader")
				k.isLeader = true
				// call start on each listener.
				for i, lst := range k.listeners {
					log.Info("starting listener: ", i)
					// each listener will run in its own goroutine
					go lst.Start()
				}
			},
			OnNewLeader: func(identity string) {
				if k.podName != identity {
					log.Info("new leader: ", identity)
				}
			},
			OnStoppedLeading: func() {
				k.logger.Info("stopped leading")
				k.isLeader = false
				// call stop on each listener.
				for i, lst := range k.listeners {
					log.Info("stopping listener: ", i)
					lst.Stop()
				}
			},
		},
	})
}

func (k *KubernetesLeaderElection) IsLeader() bool {
	return k.isLeader
}

// add a new listener to the election object.
// the listener will get notified when the current pod becomes the leader or stops being the leader.
func (k *KubernetesLeaderElection) AddListener(listener Listener) {
	k.listeners = append(k.listeners, listener)
	if k.isLeader {
		go listener.Start()
	}
}

func (k *KubernetesLeaderElection) getNewLock() *resourcelock.LeaseLock {
	return &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      k.lockName,
			Namespace: k.lockNamespace,
		},
		Client: k.client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: k.podName,
		},
	}
}
