package smssdk

import (
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
)

func (t SmsConf) NewAliyunClient() (*dysmsapi20170525.Client, error) {
	endpoint := "dysmsapi.aliyuncs.com"
	if !strings.HasPrefix(t.Region, "cn") {
		endpoint = "dysmsapi.ap-southeast-1.aliyuncs.com"
	}

	config := &openapi.Config{
		AccessKeyId:     &t.SecretId,
		AccessKeySecret: &t.SecretKey,
		Endpoint:        &endpoint,
	}
	return dysmsapi20170525.NewClient(config)
}
