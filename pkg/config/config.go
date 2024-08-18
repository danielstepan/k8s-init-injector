package config

import (
	"flag"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type ServerParameters struct {
	Port     int
	CertFile string
	KeyFile  string
}

var Config *rest.Config
var ClientSet *kubernetes.Clientset

func NewServerParameters() *ServerParameters {
	var parameters ServerParameters

	flag.IntVar(&parameters.Port, "port", 8443, "Webhook server port.")
	flag.StringVar(&parameters.CertFile, "tlsCertFile", "/etc/webhook/certs/tls.crt", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&parameters.KeyFile, "tlsKeyFile", "/etc/webhook/certs/tls.key", "File containing the x509 private key to --tlsCertFile.")
	flag.Parse()

	return &parameters
}

func LoadKubeConfig() error {
	useKubeConfig := os.Getenv("USE_KUBECONFIG")
	kubeConfigFilePath := os.Getenv("KUBECONFIG")

	if len(useKubeConfig) == 0 {
		c, err := rest.InClusterConfig()
		if err != nil {
			return err
		}
		Config = c
	} else {
		var kubeconfig string
		if kubeConfigFilePath == "" {
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = filepath.Join(home, ".kube", "config")
			}
		} else {
			kubeconfig = kubeConfigFilePath
		}

		c, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return err
		}

		Config = c
	}

	cs, err := kubernetes.NewForConfig(Config)
	if err != nil {
		return err
	}

	ClientSet = cs

	return nil
}
