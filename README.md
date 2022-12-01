# k3dnifi

> A simple tool to manage and debug nifi pods in a k3d cluster

## Dependencies
- docker
- k3d

## Installation
### Binary
```bash
go install github.com/mheers/k3dnifi@latest
```
### Docker
```bash
docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock -v $HOME/.kube:/root/.kube/:ro --network host mheers/k3dnifi:latest
```

## TODO
- [ ] remove dependency of `docker exec`
- [x] show logs of operator
- [x] make it working when zookeeper is not available
- [x] handle multiple nodes
- [x] add age to list of pods
- [x] better error message for pending pvc
- [x] auto select nifi cluster and namespace if only one is available
