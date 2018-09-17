# k8s-health-check

This repository provides the small tool `check` which can change [Kubernetes](https://kubernetes.io/) liveness and readiness checks dynamically. This can be helpful when trying to debug a dying pod. First you can mark the pod as unready so no more traffic is routed towards it. Then you can disable the liveness check, so you have all the time you need to attach a debugger without the fear of losing all the valuable runtime information.

## Usage

- `check run -type TYPE` Run the check for **TYPE** as follows: If an override is active, directly return the override, otherwise run the command defined in the environment variable **TYPE**`_CHECK` (e.g. `LIVENESS_CHECK`) and return its exit code.
- `check lock -type TYPE -state STATE` Set the override for **TYPE** to **STATE**.
- `check unlock -type TYPE` Deactivate the override for **TYPE**

**TYPE** is either `liveness` or `readiness`

**STATE** is either `success` or `failure`

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
  command: [ "check", "run", "-type", "liveness" ]
  initialDelaySeconds: 5
  periodSeconds: 5
```

and add the following env:

```env:
  - name: 'LIVENESS_CHECK'
    value: 'cat /tmp/healthy'
```

or directly in the Dockerfile: `ENV LIVENESS_CHECK "cat /tmp/healthy"`

Note that HTTP requests must be converted to using curl, e.g.

```livenessProbe:
      httpGet:
        path: /healthz
        port: 8080
      timeoutSeconds: 1
```

becomes

```env:
  - name: 'LIVENESS_CHECK'
    value: 'curl -m 1 -sf localhost:8080/healthz'
```
