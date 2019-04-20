# Auto Remedy Dispatcher

The Auto Remedy Dispatcher is a receiver for Prometheus alertmanager alerts,
and looks up possible remedies for alerts, and then applies the winning remedy.

## Deploying

The dispatcher is deployed an OpenFAAS function, using the command

```
faas build -f dispatcher.yml
faas deploy -f dispatcher.yml --env=combine_output=false
```

The `combine_output=false` makes the stderr log to the docker log, so you can
view it with
```
kubectl logs -n openfaas-fn <dispatcher pod>
```

## Development

To use minikube for development, first run
```
eval $(minikube docker-env)
export OPENFAAS_URL=http://$(minikube ip):31112
```
and then after deploying you need to run
```
kubectl edit -n openfaas-fn deployment.apps/dispatcher
```
change `imagePullPolicy` from `Always` to `Never`
