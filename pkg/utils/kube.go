package utils

import (
    log "github.com/sirupsen/logrus"
    v1 "k8s.io/api/core/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/kubernetes/scheme"
    restclient "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/tools/remotecommand"
)

func NewKubeConfig(kubeConfigPath string) (kubeConfig *restclient.Config, err error) {
    return clientcmd.BuildConfigFromFlags("", kubeConfigPath)
}

func NewKubeClient(kubeConfig *restclient.Config) (kubeClient *kubernetes.Clientset, err error) {
    return kubernetes.NewForConfig(kubeConfig)
}

func StartWebsocket(k8sClient *kubernetes.Clientset, sshHandler PtyHandler, cfg *restclient.Config,
    podName, namespace, containerName string, cmd []string) error {
    req := k8sClient.CoreV1().RESTClient().Post().
        Resource("pods").
        Name(podName).
        Namespace(namespace).
        SubResource("exec")
    req.VersionedParams(&v1.PodExecOptions{
        Container: containerName,
        Command:   cmd,
        Stdin:     true,
        Stdout:    true,
        Stderr:    true,
        TTY:       true,
    }, scheme.ParameterCodec)

    exec, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
    if err != nil {
        log.Errorf("startProcess err => %s, exec => %s", err, exec)
        return err
    }
    // 建立WebSocket长连接，阻塞并处理请求
    options := remotecommand.StreamOptions{
        Stdin:             sshHandler,
        Stdout:            sshHandler,
        Stderr:            sshHandler,
        TerminalSizeQueue: sshHandler,
        Tty:               true,
    }
    err = exec.Stream(options)
    if err != nil {
        log.Errorf("startProcess exec.Stream err => %s", err)
        return err
    }

    return nil
}