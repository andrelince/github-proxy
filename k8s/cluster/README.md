# Cluster Setup
### Apply Kyverno manifests with policy specific CRDs
`minikube kubectl -- create -f https://github.com/kyverno/kyverno/releases/download/v1.11.1/install.yaml`

### Create public key secret to validate image signature
`minikube kubectl -- create secret generic cosign-public-key --from-file=cosign.pub=./cosign.pub -n kyverno`

### Apply custom signature validation policy
`minikube kubectl -- apply -f k8s/cluster/clusterpolicy.yml`

### Apply deployment which should be created if image is signed
`minikube kubectl -- apply -f k8s/manifests/deployment.yml`
