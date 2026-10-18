package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateEmailProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateEmailProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmailProviderLogic {
	return &UpdateEmailProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: emailprovider
func (l *UpdateEmailProviderLogic) UpdateEmailProvider(in *mcms.EmailProviderInfo) (*mcms.BaseResp, error) {
	if in.Id == nil {
		return nil, errorx.New(errorx.InvalidArgument, "id is required")
	}
	q := l.svcCtx.Query.McmsEmailProvider
	id := int64(*in.Id)

	if _, err := q.WithContext(l.ctx).Where(q.ID.Eq(id)).First(); err != nil {
		return nil, dbError(l.Logger, err, in)
	}

	updates := map[string]any{}
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.AuthType != nil {
		updates["auth_type"] = int32(*in.AuthType)
	}
	if in.EmailAddr != nil {
		updates["email_addr"] = *in.EmailAddr
	}
	if in.Password != nil {
		updates["password"] = *in.Password
	}
	if in.HostName != nil {
		updates["host_name"] = *in.HostName
	}
	if in.Identify != nil {
		updates["identify"] = *in.Identify
	}
	if in.Secret != nil {
		updates["secret"] = *in.Secret
	}
	if in.Port != nil {
		updates["port"] = int32(*in.Port)
	}
	if in.Tls != nil {
		updates["tls"] = *in.Tls
	}
	if in.IsDefault != nil {
		updates["is_default"] = *in.IsDefault
	}
	if len(updates) > 0 {
		if _, err := q.WithContext(l.ctx).Where(q.ID.Eq(id)).Updates(updates); err != nil {
			return nil, dbError(l.Logger, err, in)
		}
	}

	// Only one provider may be the default.
	if in.GetIsDefault() {
		if _, err := q.WithContext(l.ctx).
			Where(q.ID.Neq(id)).
			Update(q.IsDefault, false); err != nil {
			return nil, dbError(l.Logger, err, in)
		}
	}

	// Configuration changed: drop cached SMTP connections.
	l.svcCtx.Email.Clear()

	return &mcms.BaseResp{Msg: "success"}, nil
}
