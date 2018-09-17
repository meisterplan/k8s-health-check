# k8s-health-check

This repository provides the small tool `check` which can change [Kubernetes](https://kubernetes.io/) liveness and readiness checks dynamically. This can be helpful when trying to debug a dying pod. First you can mark the pod as unready so no more traffic is routed towards it. Then you can disable the liveness check, so you have all the time you need to attach a debugger without the fear of losing all the valuable runtime information.

## Usage

`check` allows to execute the liveness and the readiness check via `check run -type liveness` and `check run -type readiness` just like usually. However check allows to change the checks dynamically by modifing the `LIVENESS_CHECK` and `READINESS_CHECK` environment variables respectively.

If you want to lock one of the checks, you can simply run `check lock -type liveness -state success` (will never fail) or `check lock -type liveness -state failure` (will always fail).

In case you want to reset the lock, simply do `check unlock -type liveness`.

## Build

If you have Make and Go installed, you should be able to do `make build`. Verification can be done with `make test`.

In case you do not want to install Go but you have Docker at hand, you can simply do `make docker-build`. Check the file permissions of the built artifact afterwards.

## Using it in your deployment

Change your k8s deployment YAML such that your probe definition, e.g.

```livenessProbe:
  exec:
  command:
  - cat
  - /tmp/healthy
  initialDelaySeconds: 5
  periodSeconds: 5
```

becomes

```livenessProbe:
  exec:
  command:
  - check
  - run
  - -type
  - liveness
  initialDelaySeconds: 5
  periodSeconds: 5
```

and add the following env:

```env:
  - name: 'LIVENESS_CHECK'
    value: 'cat /tmp/healthy'
```

Note that HTTP requests must be converted to using curl, e.g.

```livenessProbe:
      httpGet:
        path: /healthz
        port: 8080
```

becomes

```env:
  - name: 'LIVENESS_CHECK'
    value: 'curl -sf localhost:8080/healthz'
```
