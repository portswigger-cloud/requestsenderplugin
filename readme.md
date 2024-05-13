# Request Sender Traefik Plugin

Sends a POST request to a configurable URL.

## Configuration

### Kubernetes CRD

```
apiVersion: traefik.io/v1alpha1
kind: Middleware
metadata:
  name: name
  namespace: namespace
spec:
  plugin:
    requestsenderplugin:
      postUrl: https://url-to-send-post-request-to
      denylistedPaths:
      - "^/ignore-this-path/.*"
      - "^/ignore/this/path"
```