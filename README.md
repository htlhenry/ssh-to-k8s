A simple ssh proxy server, from **ssh client**  to [k8s](https://kubernetes.io/) **container**.
it is **not production ready** and not have the full features, hope inspire you do more things.

## Description

The workflow like below, it start a ssh server for ssh connection and forward it to k8s.

![workflow](ssh-to-k8s.jpg "workflow")

## QuickStart

1. install
```shell script
go get github.com/htlhenry/ssh-to-k8s
```
or
```shell script
git clone github.com/htlhenry/ssh-to-k8s $GOPATH/src/github.com/htlhenry/ssh-to-k8s/
```

2. build
```shell script
cd $GOPATH/src/github.com/htlhenry
go build -o ssh-to-k8s main.go
```

3. run
```shell script
# must spcefic the kubernetes config file(or default ~/.kube/config)
./ssh-to-k8s -h   # show help message
```

4. use it

```shell script
# Note: there not auth user implement
ssh -o "UserKnownHostsFile /dev/null" 127.0.0.1 -p 2222

# flow the help message, enter:
# <namespace> <pod> <container> 
# login to k8s
```

## Acknowledgments
inspired by [Dashboard](https://github.com/kubernetes/dashboard) and [koko](https://github.com/jumpserver/koko)