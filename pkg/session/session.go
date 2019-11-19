package session

import (
	"errors"
	"fmt"
	"github.com/gliderlabs/ssh"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"k8s.io/client-go/tools/remotecommand"
	"os/user"
	"ssh-to-k8s/pkg/utils"
	"strings"
)

func K8sProxy(sess ssh.Session, kubeConfigPath string, params []string, window <-chan ssh.Window) error {
	if len(kubeConfigPath) == 0 {
		usr, err := user.Current()
		if err != nil {
			return err
		}
		kubeConfigPath = fmt.Sprintf("%s/.kube/config", usr.HomeDir)
	}
	cfg, err := utils.NewKubeConfig(kubeConfigPath)
	if err != nil {
		return err
	}
	kubeClient, err := utils.NewKubeClient(cfg)
	if err != nil {
		return err
	}
	sshPtyHandler := utils.TerminalSSHPytHandler{
		SshSession: sess,
	}
	go MonitorResize(window, sshPtyHandler)
	namespace := params[0]
	podName := params[1]
	containerName := params[2]
	err = utils.StartWebsocket(kubeClient, sshPtyHandler, cfg, podName, namespace,
		containerName, []string{"sh"})
	if err != nil {
		return err
	}

	// client-go连接断开时还会读取一下session.Read, 会导致exit远程连接后丢失第一次session输入，未找到解决办法，这里简单粗暴的关闭了session
	// e.g
	// # exit  (disconnect k8s)
	// # e(miss)xit(capture)
	err = sess.Close()
	if err != nil {
		log.Errorf("error when close session %s", err)
	}

	return err

}

func ValidateConnectString(connString string) ([]string, error) {
	var strs []string
	strs = strings.Split(connString, " ")
	if len(strs) != 3 {
		return strs, errors.New("invalidate connection string")
	}
	return strs, nil
}

func SSHSessionHandler(sess ssh.Session) {
	_, win, _ := sess.Pty()
	term := terminal.NewTerminal(sess, "ssh> ")
	DisplayTips(sess)
	for {
		line, err := term.ReadLine()
		if err != nil {
			log.Errorf("error when readline %s", err)
			break
		}
		switch strings.ToLower(line) {
		case "h", "help":
			DisplayTips(sess)
		case "exit":
			WriteToSessionWithCRLF(sess, "bye!")
			return
		default:
			strs, err := ValidateConnectString(line)
			if err != nil {
				WriteToSessionWithCRLF(sess, err.Error())
				continue
			}
			err = K8sProxy(sess, "", strs, win)
			if err != nil {
				WriteToSessionWithCRLF(sess, err.Error())
				log.Error(err.Error())
			}
		}
	}
}

func MonitorResize(c <-chan ssh.Window, remoteTerminal utils.TerminalSSHPytHandler) {
	for {
		win := <-c
		log.Debugf("monitored window resize %+v", win)
		remoteTerminal.Resize(remotecommand.TerminalSize{
			Height: uint16(win.Height),
			Width:  uint16(win.Width),
		})
	}
}
