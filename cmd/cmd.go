package cmd

import (
    "github.com/gliderlabs/ssh"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "ssh-to-k8s/pkg/session"
    "strings"
)

var KubeConfigPath string
var Port   string

var rootCmd = &cobra.Command{
    Use: "ssh-to-k8s",
    Short: "SSH to k8s proxy",
    Run: runServer,
}

func init() {
    rootCmd.PersistentFlags().StringVarP(&KubeConfigPath,
        "kubeConfigPath", "c", "", "k8s config file")
    rootCmd.PersistentFlags().StringVarP(&Port,
        "port", "p", "2222", "listen port")
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        log.Fatal(err.Error())
    }
}

func ParsePort(port string) string {
    if strings.Contains(port, ":") && !strings.HasPrefix(port, ":") {
        return port
    }
    if !strings.HasPrefix(port, ":") {
        port = ":" + port
    }
    return port
}

func runServer(cmd *cobra.Command, args []string) {
    Port = ParsePort(Port)
    srv := ssh.Server{
        Addr:   Port,
        Handler:    session.SSHSessionHandler,
    }
    log.Infof("listen on %s", Port)
    log.Fatal(srv.ListenAndServe())
}