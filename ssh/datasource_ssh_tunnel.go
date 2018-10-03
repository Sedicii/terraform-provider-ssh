package ssh

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/crypto/ssh"
	"net"
	"strconv"
)

func dataSourceSSHTunnel() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSSHTunnelRead,
		Schema: map[string]*schema.Schema{
			"bastion_host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This host will be connected to first, and then the host connection will be made from there.",
			},
			"bastion_host_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The public key from the remote host or the signing CA, used to verify the host connection.",
			},
			"bastion_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     22,
				Description: "The port to use connect to the bastion host. Default to 22.",
			},
			"bastion_user": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user for the connection to the bastion host.",
			},
			"bastion_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				Description: "The password we should use for the bastion host.",
			},
			"bastion_private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Sensitive:   true,
				Description: "The private key we should use for the bastion host.",
			},
			"remote_host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Remote host to connect through tunnel",
			},
			"remote_port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Remote host port to connect through tunnel",
			},
			"local_host": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Local host where the tunnel listens",
			},
			"local_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Local port where the tunnel listens",
			},
		},
	}
}

func dataSourceSSHTunnelRead(d *schema.ResourceData, meta interface{}) error {
	remoteHost := d.Get("remote_host").(string)
	remotePort := d.Get("remote_port").(int)
	remoteAddr := remoteHost + ":" + strconv.Itoa(remotePort)

	bastionHost := d.Get("bastion_host").(string)
	bastionPort := d.Get("bastion_port").(int)
	bastionUser := d.Get("bastion_user").(string)
	bastionAddr := bastionHost + ":" + strconv.Itoa(bastionPort)

	hostKeyValidation, err := getHostKeyCallback(d)
	if err != nil {
		return err
	}

	auth, err := getAuthMethod(d)
	if err != nil {
		return err
	}

	freeListenAddr, err := getFreePort()
	if err != nil {
		return err
	}

	tunnel := sshTunnel{
		Local:  freeListenAddr,
		Server: bastionAddr,
		Remote: remoteAddr,
		Config: &ssh.ClientConfig{
			User:            bastionUser,
			Auth:            auth,
			HostKeyCallback: hostKeyValidation,
		},
	}

	if err = tunnel.Start(); err != nil {
		return err
	}

	d.Set("local_host", freeListenAddr.IP.String())
	d.Set("local_port", freeListenAddr.Port)
	d.SetId(hash(bastionAddr + remoteAddr + bastionUser))
	return nil
}

func getFreePort() (*net.TCPAddr, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr), nil
}

func getAuthMethod(conf *schema.ResourceData) ([]ssh.AuthMethod, error) {
	bastionPassword := conf.Get("bastion_password").(string)
	bastionPrivateKey := conf.Get("bastion_private_key").(string)

	if bastionPrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(bastionPrivateKey))
		if err != nil {
			return nil, errwrap.Wrapf("Error parsing bastion host private key : {{err}}", err)
		}
		return []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}, nil
	} else if bastionPassword != "" {
		return []ssh.AuthMethod{
			ssh.Password(bastionPassword),
		}, nil
	}

	return nil, errors.New("one of bastion_password or bastion_private_key is required in block ssh_tunnel")
}

func getHostKeyCallback(conf *schema.ResourceData) (ssh.HostKeyCallback, error) {
	bastionHostKey := conf.Get("bastion_host_key").(string)

	if bastionHostKey != "" {
		publicHostKey, err := ssh.ParsePublicKey([]byte(bastionHostKey))

		if err != nil {
			return nil, errwrap.Wrapf("Error parsing bastion host public key : {{err}}", err)
		}

		return ssh.FixedHostKey(publicHostKey), nil
	}

	return ssh.InsecureIgnoreHostKey(), nil
}

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}
