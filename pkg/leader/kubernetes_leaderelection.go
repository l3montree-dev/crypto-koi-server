package leader

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.com/l3montree/microservices/libs/orchardclient"
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
	logger        *logrus.Entry
}

func NewKubernetesLeaderElection(ctx context.Context, lockName, lockNamespace string) LeaderElection {
	logger := orchardclient.Logger.WithField("component", "KubernetesLeaderElection")
	config, err := rest.InClusterConfig()
	orchardclient.FailOnError(err, "failed to retrieve kubernetes config")

	// get the pod name of the current pod
	podName := os.Getenv("POD_NAME")
	if podName == "" {
		logger.Fatal("failed to retrieve pod name. Running in kubernetes?")
	}

	client := kubernetes.NewForConfigOrDie(config)
	return &KubernetesLeaderElection{
		lockName:      lockName,
		lockNamespace: lockNamespace,
		ctx:           ctx,
		client:        client,
		podName:       podName,
		logger:        logger,
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

func (k *KubernetesLeaderElection) IsLeader() bool {
	return k.isLeader
}

func (k *KubernetesLeaderElection) RunElection() {
	k.logger.Info("running leader election in kubernetes mode")
	lock := k.getNewLock()
	leaderelection.RunOrDie(k.ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(c context.Context) {
				k.logger.Info("we are the leader")
				k.isLeader = true
				// call start on each listener.
				for i, lst := range k.listeners {
					k.logger.Info("starting listener: ", i)
					// each listener will run in its own goroutine
					go lst.Start()
				}
			},
			OnStoppedLeading: func() {
				k.logger.Info("stopped leading")
				k.isLeader = false
				// call stop on each listener.
				for i, lst := range k.listeners {
					k.logger.Info("stopping listener: ", i)
					lst.Stop()
				}
			},
		},
	})
}

func (k *KubernetesLeaderElection) AddListener(listener Listener) {
	k.listeners = append(k.listeners, listener)
}
