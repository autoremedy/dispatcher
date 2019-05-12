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

## Configuration

Each remedy has an array of remedy configurations, which are used to decide which nodes match the remedy based on a regexp `filter`.

E.g. the `NodeDown` remedy which performs a hardware reset on the servers should have a limit on how often it can be applied:
```
[
  {
    "name": "db node down",
    "filter": "-db.*",
    "ratelimit": "2"
  },
  {
    "name": "node down",
    "filter": ".*",
    "ratelimit": "10"
  },
]
```

## Development

If you want to develop in a fully functional environment locally, use 
[minikube](https://github.com/kubernetes/minikube).
There you can deploy [openfaas](https://github.com/openfaas/faas), 
and use the [prometheus](https://github.com/prometheus/prometheus) 
instance which is bundled with it.

After openfaas is deployed, run
```bash
eval $(minikube docker-env)
export OPENFAAS_URL=http://$(minikube ip):31112
```
then deploy the dispatcher with
```bash
faas deploy -f dispatch.yml
```
The first time it is deployed you need to run
```bash
kubectl edit -n openfaas-fn deployment.apps/dispatcher
```
change `imagePullPolicy` from `Always` to `Never`,
as it otherwise tries to download the docker image from the docker hub.

### Configure amtool
Make alertmanager accessible
```
kubectl edit -n openfaas service/alertmanager
```
set `type: NodePort`

```
go get github.com/prometheus/alertmanager/cmd/amtool
```

create `~/.config/amtool/config.yml`
```
alertmanager.url: http://$(minikube ip):$(kubectl get service -n openfaas alertmanager --no-headers -o custom-columns=port:spec.ports[*].nodePort)
```

```
amtool alert add \
  FooAlert \
  node=x.y.z \
  --annotation=runbook='http://x.y.z' \
  --summary='foo has happened'
```
