package ssh

import (
	"github.com/hashicorp/errwrap"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
)

type sshTunnel struct {
	Local  *net.TCPAddr
	Server string
	Remote string
	Config *ssh.ClientConfig
}

func (tunnel *sshTunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return errwrap.Wrapf("Error listening in local for ssh tunnel : {{err}}", err)
	}

	serverConn, err := ssh.Dial("tcp", tunnel.Server, tunnel.Config)
	if err != nil {
		return errwrap.Wrapf("Error connecting to bastion host : {{err}}", err)
	}

	go tunnel.handleClients(listener, serverConn)

	return nil
}

func (tunnel *sshTunnel) handleClients(listener net.Listener, serverConn *ssh.Client) {
	defer listener.Close()

	for {
		localConn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connections in ssh tunnel : %s", err)
			return
		}

		remoteConn, err := serverConn.Dial("tcp", tunnel.Remote)
		if err != nil {
			log.Fatalf("Remote dial error: %s\n", err)
			return
		}

		go tunnel.forward(localConn, remoteConn)
	}
}

func (tunnel *sshTunnel) forward(localConn net.Conn, remoteConn net.Conn) {

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			log.Fatalf("io.Copy error: %s", err)
		}
		writer.Close()
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}
