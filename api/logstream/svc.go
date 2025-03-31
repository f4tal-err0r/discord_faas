package logstream

import "k8s.io/client-go/kubernetes"

type Handler struct {
	cs *kubernetes.Clientset
}

type Logger interface {
	TailLogs() (chan []byte, error)
}

func NewHandler(cs *kubernetes.Clientset) *Handler {

}
