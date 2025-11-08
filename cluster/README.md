# Kubernetes Cluster Setup

```bash
kubectl apply -f namespace.yaml
kubectl apply -f daemoneset.yaml

kubectl create configmap presets-rtx-3060-ti --from-file=presets.json=../provider/examples/presets-rtx-3060-ti.json --namespace providers
kubectl create configmap presets-rtx-4090 --from-file=presets.json=../provider/examples/presets-rtx-4090.json --namespace providers

kubectl create configmap template-node-1 --from-file=template.json=../provider/ya-runtime-salad/examples/template-rtx-3060-ti.json --namespace providers
kubectl create configmap template-node-2 --from-file=template.json=../provider/ya-runtime-salad/examples/template-rtx-4090.json --namespace providers
kubectl create configmap template-node-3 --from-file=template.json=../provider/ya-runtime-salad/examples/template-rtx-4090.json --namespace providers

kubectl create secret generic config-node-1 --from-env-file=../provider/examples/node-1.env
kubectl create secret generic config-node-2 --from-env-file=../provider/examples/node-2.env
kubectl create secret generic config-node-3 --from-env-file=../provider/examples/node-3.env

kubectl apply -f node-1.yaml
kubectl apply -f node-2.yaml
kubectl apply -f node-3.yaml
```
