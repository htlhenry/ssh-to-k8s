package session

import (
    "github.com/gliderlabs/ssh"
    "io"
)

var tips = `
    # login format(split with space)
    <namespace> <podName> <containerName>
    
    # example:
    dev test-pod test-container

    Enter 'h' for show the help
`


func WriteToSessionWithCRLF(sess ssh.Session, s string) (int, error) {
    return io.WriteString(sess, s + "\r\n")
}

func DisplayTips(sess ssh.Session) {
    sess.Write([]byte(tips))
}
