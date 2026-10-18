// Package smssdk integrates the third-party SMS providers supported by mcms.
package smssdk

import (
	"sync"

	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	unisms "github.com/apistd/uni-go-sdk/sms"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// SmsConf is the configuration of one SMS provider.
type SmsConf struct {
	SecretId  string
	SecretKey string
	Provider  string
	Region    string
}

// SmsGroup caches SDK clients by provider name. Every access is guarded by mu
// because the RPC server handles requests concurrently. The cache never
// expires (same as the reference implementation); provider updates do not
// invalidate SMS clients.
type SmsGroup struct {
	mu      sync.Mutex
	tencent map[string]*sms.Client
	aliyun  map[string]*dysmsapi20170525.Client
	uni     map[string]*unisms.UniSMSClient
}

func NewSmsGroup() *SmsGroup {
	return &SmsGroup{
		tencent: map[string]*sms.Client{},
		aliyun:  map[string]*dysmsapi20170525.Client{},
		uni:     map[string]*unisms.UniSMSClient{},
	}
}

// Tencent returns the cached Tencent client for name, creating it on demand.
func (g *SmsGroup) Tencent(name string, create func() (*sms.Client, error)) (*sms.Client, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if c, ok := g.tencent[name]; ok {
		return c, nil
	}
	c, err := create()
	if err != nil {
		return nil, err
	}
	g.tencent[name] = c
	return c, nil
}

// Aliyun returns the cached Aliyun client for name, creating it on demand.
func (g *SmsGroup) Aliyun(name string, create func() (*dysmsapi20170525.Client, error)) (*dysmsapi20170525.Client, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if c, ok := g.aliyun[name]; ok {
		return c, nil
	}
	c, err := create()
	if err != nil {
		return nil, err
	}
	g.aliyun[name] = c
	return c, nil
}

// Uni returns the cached Uni client for name, creating it on demand.
func (g *SmsGroup) Uni(name string, create func() (*unisms.UniSMSClient, error)) (*unisms.UniSMSClient, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if c, ok := g.uni[name]; ok {
		return c, nil
	}
	c, err := create()
	if err != nil {
		return nil, err
	}
	g.uni[name] = c
	return c, nil
}
