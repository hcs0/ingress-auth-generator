package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const AppVersion = "0.1.0"
const AppName = "IngressAuthGenerator"

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.Debugf("%s Version: %s", AppName, AppVersion)
}

func main() {
	log.Infof("Start")

	var nameSpace string
	if len(os.Getenv("KUBERNETES_NAMESPACE")) > 0 {
		nameSpace = os.Getenv("KUBERNETES_NAMESPACE")
	} else {
		nameSpace = "default"
	}
	log.Debugf("Namespace: %s", nameSpace)

	for {
		time.Sleep(5 * time.Second)
		client, err := k8sClient("")
		if err != nil {
			log.Panic("K8s connection Failed! Reason:", err)
		}

		ingress, err := client.ExtensionsV1beta1().Ingresses(nameSpace).List(metav1.ListOptions{})
		if err != nil {
			log.Errorf("Get ingress list failed! -> %s ", err)
		}

		var ingSec []string
		for _, ing := range ingress.Items {
			value, ok := ing.Annotations["ingress.kubernetes.io/auth-secret"]
			if ok {
				ingSec = append(ingSec, value)
			}

		}
		for _, secretName := range ingSec {
			secret, err := client.CoreV1().Secrets(nameSpace).Get(secretName, metav1.GetOptions{})
			if err != nil {
				fmt.Errorf("K8S get Secret Failed: %v", err)
			}
			log.Debugf("Check %q Secret Content. --> auth: %q", secretName, secret.Data["auth"])
			if secret.Data["auth"] == nil && secret.Data["username"] != nil && secret.Data["password"] != nil {
				passwordHash, err := hashBcrypt(string(secret.Data["password"]))
				if err != nil {
					log.Errorf("Password crypt failed!")
					continue
				}
				newAuth := fmt.Sprintf("%s:%s", secret.Data["username"], passwordHash)
				log.Debugf("New Auth: %s", newAuth)

				secret.Data["auth"] = []byte(newAuth)
				client.CoreV1().Secrets(nameSpace).Update(secret)
				log.Infof("Secret Update Done!")
			}
		}
	}

}

func hashBcrypt(password string) (hash string, err error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordBytes), nil
}

func k8sClient(k8sConfigFile string) (*kubernetes.Clientset, error) {
	var (
		config *rest.Config
		err    error
	)

	if k8sConfigFile == "" {
		k8sConfigFile = os.Getenv("kubeConfig")
		log.Debugln("kubeConfig:", k8sConfigFile)
	}

	if k8sConfigFile != "" {
		config, err = clientcmd.BuildConfigFromFlags("", k8sConfigFile)
	} else {
		log.Infoln("Use K8S InCluster Config.")
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, fmt.Errorf("K8S Connection Failed: %v", err)
	}

	client := kubernetes.NewForConfigOrDie(config)
	return client, nil
}
