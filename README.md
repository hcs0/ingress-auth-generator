# Ingress Auth Generator Daemon For Kubernetes Ingress

Uses the `IngressAuthGenerator`  to create  basic authentication for Kubernetes ingress controllers.


## How it works?

This simple `Golang` application helps you to transform your passwords from `k8s secrets` to `htaccess k8s secrets`. It can use the ingress controllers by default.


## How can you use it?
## tl;dr:

```bash
$ helm repo add banzaicloud http://kubernetes-charts.banzaicloud.com
$ helm repo update
$ helm install banzaicloud/ingressauthgenerator
```