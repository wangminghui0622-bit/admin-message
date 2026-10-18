// Package emailclient builds and caches SMTP clients for mcms email providers.
package emailclient

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"sync"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	AuthTypePlain   = 1
	AuthTypeCRAMMD5 = 2
)

// Conf is the connection configuration converted from an email provider row.
type Conf struct {
	AuthType  uint32
	EmailAddr string
	Password  string
	HostName  string
	Identify  string
	Secret    string
	Port      uint32
	TLS       bool
}

func (e Conf) newAuth() smtp.Auth {
	switch e.AuthType {
	case AuthTypeCRAMMD5:
		return smtp.CRAMMD5Auth(e.EmailAddr, e.Secret)
	default:
		// auth_type 1 (plain); the identify is usually an empty string.
		return smtp.PlainAuth(e.Identify, e.EmailAddr, e.Password, e.HostName)
	}
}

// NewClient dials the SMTP server and authenticates. TLS verification is
// deliberately disabled to match the reference behaviour (self-signed and
// mismatched certificates are common for mail relays).
func (e Conf) NewClient() (*smtp.Client, error) {
	hostAddress := net.JoinHostPort(e.HostName, strconv.Itoa(int(e.Port)))

	var (
		c   *smtp.Client
		err error
	)
	if e.TLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, //nolint:gosec // matches reference behaviour
			ServerName:         e.HostName,
		}
		var conn net.Conn
		conn, err = tls.Dial("tcp", hostAddress, tlsConfig)
		if err != nil {
			logx.Errorw("failed to dial email server with tls", logx.Field("host", hostAddress), logx.Field("detail", err))
			return nil, errorx.New(errorx.InvalidArgument, "failed to connect email server, please check the host and tls config")
		}
		c, err = smtp.NewClient(conn, e.HostName)
		if err != nil {
			_ = conn.Close()
			logx.Errorw("failed to create smtp client", logx.Field("host", hostAddress), logx.Field("detail", err))
			return nil, errorx.New(errorx.InvalidArgument, "failed to connect email server, please check the host and tls config")
		}
	} else {
		c, err = smtp.Dial(hostAddress)
		if err != nil {
			logx.Errorw("failed to dial email server", logx.Field("host", hostAddress), logx.Field("detail", err))
			return nil, errorx.New(errorx.InvalidArgument, "failed to connect email server, please check the host and tls config")
		}
	}

	if err = c.Auth(e.newAuth()); err != nil {
		_ = c.Close()
		logx.Errorw("failed to authorize email server", logx.Field("host", hostAddress), logx.Field("detail", err))
		return nil, errorx.New(errorx.InvalidArgument, "failed to authorize the server, please check the identify and secret")
	}
	return c, nil
}

// ClientCache keeps authenticated SMTP clients by provider name. go-zero RPC
// handlers run concurrently, so all map access is guarded by mu.
type ClientCache struct {
	mu      sync.Mutex
	clients map[string]*smtp.Client
	addrs   map[string]string
}

func NewClientCache() *ClientCache {
	return &ClientCache{
		clients: map[string]*smtp.Client{},
		addrs:   map[string]string{},
	}
}

// Get returns the cached client and sender address, or builds them via create.
func (c *ClientCache) Get(name string, create func() (*smtp.Client, string, error)) (*smtp.Client, string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if client, ok := c.clients[name]; ok {
		return client, c.addrs[name], nil
	}
	client, addr, err := create()
	if err != nil {
		return nil, "", err
	}
	c.clients[name] = client
	c.addrs[name] = addr
	return client, addr, nil
}

// Evict drops one cached client (used after a broken pipe).
func (c *ClientCache) Evict(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if client, ok := c.clients[name]; ok {
		_ = client.Close()
	}
	delete(c.clients, name)
	delete(c.addrs, name)
}

// Clear drops every cached client (email provider update/delete).
func (c *ClientCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, client := range c.clients {
		_ = client.Close()
	}
	c.clients = map[string]*smtp.Client{}
	c.addrs = map[string]string{}
}

// ErrNoProvider is returned when a named provider row does not exist.
var ErrNoProvider = fmt.Errorf("email provider not found")
