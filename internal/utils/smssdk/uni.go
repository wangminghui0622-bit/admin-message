package smssdk

import (
	unisms "github.com/apistd/uni-go-sdk/sms"
)

func (t SmsConf) NewUniClient() (*unisms.UniSMSClient, error) {
	return unisms.NewClient(t.SecretId, t.SecretKey), nil
}
