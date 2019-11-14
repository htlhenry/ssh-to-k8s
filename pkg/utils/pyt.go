package utils

import (
    "github.com/gliderlabs/ssh"
    log "github.com/sirupsen/logrus"
    "io"
    "k8s.io/client-go/tools/remotecommand"
)

type PtyHandler interface {
    io.Reader
    io.Writer
    remotecommand.TerminalSizeQueue
}

type TerminalSSHPytHandler struct {
    SshSession ssh.Session
    sizeChan    chan remotecommand.TerminalSize
}

// k8s 调用这个方法接收resize时间，需要将monitor resize 将resize数据传入到sizeChan
// 注意: 只有v2以上的协议才支持container resize，请参考: k8s.io/apimachinery/pkg/util/remotecommand/constants.go
func (t TerminalSSHPytHandler) Next() *remotecommand.TerminalSize{
    select {
    case size := <-t.sizeChan:
        log.Debugf("k8s read size %+v", size)
        return &size
    }
}

func (t TerminalSSHPytHandler) Resize(tSize remotecommand.TerminalSize) {
    t.sizeChan <- tSize
}

// ssh client 写入k8s输入
func (t TerminalSSHPytHandler) Read(p []byte) (int, error) {
    return t.SshSession.Read(p)
}

// ssh client 读取k8s输出
func (t TerminalSSHPytHandler) Write(p []byte) (int, error) {
    return t.SshSession.Write(p)
}

