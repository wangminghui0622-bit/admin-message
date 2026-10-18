package logic

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/internal/utils/emailclient"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendEmailLogic {
	return &SendEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: email
func (l *SendEmailLogic) SendEmail(in *mcms.EmailInfo) (*mcms.BaseUUIDResp, error) {
	providerName := in.GetProvider()
	if providerName == "" {
		defaultProvider, err := l.svcCtx.Query.McmsEmailProvider.WithContext(l.ctx).
			Where(l.svcCtx.Query.McmsEmailProvider.IsDefault.Eq(true)).First()
		if err != nil {
			return nil, dbError(l.Logger, err, in)
		}
		providerName = defaultProvider.Name
	}

	// SMTP clients are cached by provider name; email provider update/delete
	// clears the whole cache.
	client, from, err := l.svcCtx.Email.Get(providerName, func() (*smtp.Client, string, error) {
		data, err := l.svcCtx.Query.McmsEmailProvider.WithContext(l.ctx).
			Where(l.svcCtx.Query.McmsEmailProvider.Name.Eq(providerName)).First()
		if err != nil {
			return nil, "", dbError(l.Logger, err, in)
		}
		conf := emailclient.Conf{
			AuthType:  uint32(data.AuthType),
			EmailAddr: data.EmailAddr,
			HostName:  data.HostName,
			TLS:       data.TLS,
		}
		if data.Password != nil {
			conf.Password = *data.Password
		}
		if data.Identify != nil {
			conf.Identify = *data.Identify
		}
		if data.Secret != nil {
			conf.Secret = *data.Secret
		}
		if data.Port != nil {
			conf.Port = uint32(*data.Port)
		}
		smtpClient, err := conf.NewClient()
		if err != nil {
			return nil, "", err
		}
		return smtpClient, data.EmailAddr, nil
	})
	if err != nil {
		return nil, err
	}

	// evict the cached client when its connection is unrecoverable.
	handleErr := func(cause error) (*mcms.BaseUUIDResp, error) {
		l.Logger.Errorw("failed to send email", logx.Field("detail", cause.Error()), logx.Field("data", in))
		if strings.Contains(cause.Error(), "broken pipe") {
			l.svcCtx.Email.Evict(providerName)
		}
		if _, logErr := l.writeLog(2, providerName, in); logErr != nil {
			return nil, logErr
		}
		return nil, errorx.New(errorx.ServerError, "failed to send email")
	}

	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, strings.Join(in.Target, ","), in.Subject, in.Content)

	if err = client.Mail(from); err != nil {
		return handleErr(err)
	}
	for _, target := range in.Target {
		if err = client.Rcpt(target); err != nil {
			return handleErr(err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return handleErr(err)
	}
	if _, err = writer.Write([]byte(message)); err != nil {
		_ = writer.Close()
		return handleErr(err)
	}
	if err = writer.Close(); err != nil {
		return handleErr(err)
	}

	id, err := l.writeLog(1, providerName, in)
	if err != nil {
		return nil, err
	}
	return &mcms.BaseUUIDResp{Id: id, Msg: "success"}, nil
}

// writeLog persists an email log; provider is always filled in (the original
// implementation left it empty, which made the RPC unusable).
func (l *SendEmailLogic) writeLog(status int32, providerName string, in *mcms.EmailInfo) (string, error) {
	id := newUUID()
	record := &model.McmsEmailLog{
		ID:         id,
		Target:     strings.Join(in.Target, ","),
		Subject:    in.Subject,
		Content:    in.Content,
		SendStatus: status,
		Provider:   providerName,
	}
	if err := l.svcCtx.Query.McmsEmailLog.WithContext(context.Background()).Create(record); err != nil {
		return "", dbError(l.Logger, err, in)
	}
	return id, nil
}
