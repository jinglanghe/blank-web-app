package service

import (
	"github.com/apulis/sdk/go-utils/logging"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

var (
	clientset *kubernetes.Clientset
)

func getK8sCfgPath() string {
	home := homedir.HomeDir()
	configPath := filepath.Join(home, ".kube", "config")
	if !FileExist(configPath) {
		configPath = ""
	}

	return configPath
}

func initClientSet() error {
	if clientset != nil {
		return nil
	}

	configPath := getK8sCfgPath()

	// use the current context in kubeconfig
	kubecfg, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		logging.Error(err).Msg("")
		return err
	}

	clientset, err = kubernetes.NewForConfig(kubecfg)
	if err != nil {
		logging.Error(err).Msg("")
		return err
	}

	return nil
}
