package config

import (
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ServerParameters struct {
	Port     int
	CertFile string
	KeyFile  string
}

var (
	ClientSet  *kubernetes.Clientset
	parameters ServerParameters
	isLocalRun bool
)

func InitializeFlags() {
	flag.BoolVar(&isLocalRun, "local", false, "read config from $HOME/.kube/config")

	flag.IntVar(&parameters.Port, "port", 8443, "Webhook server port.")
	flag.StringVar(&parameters.CertFile, "tlsCertFile", "/etc/webhook/certs/tls.crt", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&parameters.KeyFile, "tlsKeyFile", "/etc/webhook/certs/tls.key", "File containing the x509 private key to --tlsCertFile.")

	flag.Parse()
}

func NewServerParameters() *ServerParameters {
	return &parameters
}

func LoadLocalKubeConfig() (*rest.Config, error) {
	c, err := clientcmd.BuildConfigFromFlags("", filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	if err != nil {
		return nil, err
	}

	return c, nil
}

func LoadClusterKubeConfig() (*rest.Config, error) {
	c, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func LoadKubeConfig() error {
	var config *rest.Config
	var err error

	if isLocalRun {
		config, err = LoadLocalKubeConfig()
	} else {
		config, err = LoadClusterKubeConfig()
	}
	if err != nil {
		return err
	}

	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	ClientSet = cs

	return nil
}
