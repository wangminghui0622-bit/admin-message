package logic

import (
	"context"
	"strings"

	aliyun "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	unisms "github.com/apistd/uni-go-sdk/sms"
	"github.com/go-resty/resty/v2"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/internal/utils/smssdk"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	providerTencent = "tencent"
	providerAliyun  = "aliyun"
	providerUni     = "uni"
	providerSmsBao  = "smsbao"
)

type SendSmsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendSmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSmsLogic {
	return &SendSmsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: sms
func (l *SendSmsLogic) SendSms(in *mcms.SmsInfo) (*mcms.BaseUUIDResp, error) {
	providerName := in.GetProvider()
	if providerName == "" {
		defaultProvider, err := l.svcCtx.Query.McmsSmsProvider.WithContext(l.ctx).
			Where(l.svcCtx.Query.McmsSmsProvider.IsDefault.Eq(true)).First()
		if err != nil {
			return nil, dbError(l.Logger, err, in)
		}
		providerName = defaultProvider.Name
	}

	var sendErr error

	switch providerName {
	case providerTencent:
		request := sms.NewSendSmsRequest()
		request.TemplateId = in.TemplateId
		request.SmsSdkAppId = in.AppId
		request.PhoneNumberSet = stringPtrSlice(in.PhoneNumber)
		request.TemplateParamSet = stringPtrSlice(in.Params)
		request.SignName = in.SignName

		client, err := l.tencentClient(providerName)
		if err != nil {
			return nil, err
		}
		resp, err := client.SendSms(request)
		if err != nil {
			sendErr = err
		} else {
			logx.Infow("send SMS by Tencent", logx.Field("response", resp), logx.Field("phoneNumber", in.PhoneNumber))
		}

	case providerAliyun:
		request := &aliyun.SendSmsRequest{}
		request.SignName = in.SignName
		request.TemplateCode = in.TemplateId
		request.PhoneNumbers = ptr(strings.Join(in.PhoneNumber, ","))
		if len(in.Params) > 0 {
			// The caller must pass a single raw JSON string (reference contract).
			request.TemplateParam = ptr(strings.Join(in.Params, ""))
		}

		client, err := l.aliyunClient(providerName)
		if err != nil {
			return nil, err
		}
		resp, err := client.SendSms(request)
		if err != nil {
			sendErr = err
		} else {
			logx.Infow("send SMS by Aliyun", logx.Field("response", resp), logx.Field("phoneNumber", in.PhoneNumber))
		}

	case providerUni:
		request := unisms.BuildMessage()
		request.SetSignature(in.GetSignName())
		request.SetTemplateId(in.GetTemplateId())
		request.SetTo(in.PhoneNumber...)
		if len(in.Params) > 0 {
			paramsData := map[string]string{}
			for _, v := range in.Params {
				p := strings.Split(v, ":")
				if len(p) != 2 {
					logx.Errorw("wrong parameters in Uni SMS Request", logx.Field("param", in.Params))
					return nil, errorx.New(errorx.InvalidArgument, "invalid params, must be key:value")
				}
				paramsData[p[0]] = p[1]
			}
			request.SetTemplateData(paramsData)
		}

		client, err := l.uniClient(providerName)
		if err != nil {
			return nil, err
		}
		resp, err := client.Send(request)
		if err != nil {
			sendErr = err
		} else {
			logx.Infow("send SMS by Uni SMS", logx.Field("response", resp), logx.Field("phoneNumber", in.PhoneNumber))
		}

	case providerSmsBao:
		if len(in.Params) == 0 {
			return nil, errorx.New(errorx.InvalidArgument, "params is required for smsbao")
		}
		info, err := l.svcCtx.Query.McmsSmsProvider.WithContext(l.ctx).
			Where(l.svcCtx.Query.McmsSmsProvider.Name.Eq(providerName)).First()
		if err != nil {
			return nil, dbError(l.Logger, err, in)
		}

		var msg string
		if in.GetTemplateId() == "captcha" {
			msg = "【" + in.GetSignName() + "】您的验证码为：" + in.Params[0] + "，请尽快完成验证。"
		} else {
			msg = "【" + in.GetSignName() + "】" + in.Params[0]
		}

		resp, err := resty.New().R().
			SetQueryParam("u", info.SecretID).
			SetQueryParam("p", info.SecretKey).
			SetQueryParam("m", strings.Join(in.PhoneNumber, ",")).
			SetQueryParam("c", msg).
			Get("https://api.smsbao.com/sms")
		if err != nil || resp == nil || string(resp.Body()) != "0" {
			sendErr = err
			if sendErr == nil {
				sendErr = errorx.New(errorx.ServerError, "smsbao rejected the request")
			}
		} else {
			logx.Infow("send SMS by SMS Bao", logx.Field("response", resp), logx.Field("phoneNumber", in.PhoneNumber))
		}

	default:
		return nil, errorx.New(errorx.InvalidArgument, "provider not found")
	}

	if sendErr != nil {
		logx.Errorw("failed to send SMS", logx.Field("detail", sendErr), logx.Field("data", in))
		if _, err := l.writeLog(2, providerName, in); err != nil {
			return nil, err
		}
		return nil, errorx.New(errorx.ServerError, "failed to send SMS")
	}

	id, err := l.writeLog(1, providerName, in)
	if err != nil {
		return nil, err
	}
	return &mcms.BaseUUIDResp{Id: id, Msg: "success"}, nil
}

func (l *SendSmsLogic) tencentClient(providerName string) (*sms.Client, error) {
	return l.svcCtx.Sms.Tencent(providerName, func() (*sms.Client, error) {
		data, err := l.svcCtx.Query.McmsSmsProvider.WithContext(l.ctx).
			Where(l.svcCtx.Query.McmsSmsProvider.Name.Eq(providerName)).First()
		if err != nil {
			return nil, dbError(l.Logger, err, providerName)
		}
		conf := smssdk.SmsConf{
			SecretId:  data.SecretID,
			SecretKey: data.SecretKey,
			Provider:  providerName,
			Region:    data.Region,
		}
		client, err := conf.NewTencentClient()
		if err != nil {
			logx.Errorw("failed to initialize Tencent SMS client, please check the configuration", logx.Field("detail", err))
			return nil, errorx.New(errorx.InvalidArgument, "failed to initialize Tencent SMS client, please check the configuration")
		}
		return client, nil
	})
}

func (l *SendSmsLogic) aliyunClient(providerName string) (*aliyun.Client, error) {
	return l.svcCtx.Sms.Aliyun(providerName, func() (*aliyun.Client, error) {
		data, err := l.svcCtx.Query.McmsSmsProvider.WithContext(l.ctx).
			Where(l.svcCtx.Query.McmsSmsProvider.Name.Eq(providerName)).First()
		if err != nil {
			return nil, dbError(l.Logger, err, providerName)
		}
		conf := smssdk.SmsConf{
			SecretId:  data.SecretID,
			SecretKey: data.SecretKey,
			Provider:  providerName,
			Region:    data.Region,
		}
		client, err := conf.NewAliyunClient()
		if err != nil {
			logx.Errorw("failed to initialize Aliyun SMS client, please check the configuration", logx.Field("detail", err))
			return nil, errorx.New(errorx.InvalidArgument, "failed to initialize Aliyun SMS client, please check the configuration")
		}
		return client, nil
	})
}

func (l *SendSmsLogic) uniClient(providerName string) (*unisms.UniSMSClient, error) {
	return l.svcCtx.Sms.Uni(providerName, func() (*unisms.UniSMSClient, error) {
		data, err := l.svcCtx.Query.McmsSmsProvider.WithContext(l.ctx).
			Where(l.svcCtx.Query.McmsSmsProvider.Name.Eq(providerName)).First()
		if err != nil {
			return nil, dbError(l.Logger, err, providerName)
		}
		conf := smssdk.SmsConf{
			SecretId:  data.SecretID,
			SecretKey: data.SecretKey,
			Provider:  providerName,
			Region:    data.Region,
		}
		client, err := conf.NewUniClient()
		if err != nil {
			logx.Errorw("failed to initialize Uni SMS client, please check the configuration", logx.Field("detail", err))
			return nil, errorx.New(errorx.InvalidArgument, "failed to initialize Uni SMS client, please check the configuration")
		}
		return client, nil
	})
}

// writeLog persists an SMS log with the provider filled in. The log write uses
// context.Background() so it survives a cancelled request context, matching the
// reference semantics.
func (l *SendSmsLogic) writeLog(status int32, providerName string, in *mcms.SmsInfo) (string, error) {
	id := newUUID()
	record := &model.McmsSmsLog{
		ID:          id,
		PhoneNumber: strings.Join(in.PhoneNumber, ","),
		Content:     strings.Join(in.Params, ","),
		SendStatus:  status,
		Provider:    providerName,
	}
	if err := l.svcCtx.Query.McmsSmsLog.WithContext(context.Background()).Create(record); err != nil {
		return "", dbError(l.Logger, err, in)
	}
	return id, nil
}

func stringPtrSlice(values []string) []*string {
	out := make([]*string, len(values))
	for i := range values {
		out[i] = &values[i]
	}
	return out
}
