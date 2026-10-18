package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEmailProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateEmailProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEmailProviderLogic {
	return &CreateEmailProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: emailprovider
func (l *CreateEmailProviderLogic) CreateEmailProvider(in *mcms.EmailProviderInfo) (*mcms.BaseIDResp, error) {
	q := l.svcCtx.Query.McmsEmailProvider
	record := &model.McmsEmailProvider{
		Name:      in.GetName(),
		AuthType:  int32(in.GetAuthType()),
		EmailAddr: in.GetEmailAddr(),
		Password:  in.Password,
		HostName:  in.GetHostName(),
		Identify:  in.Identify,
		Secret:    in.Secret,
		TLS:       in.GetTls(),
		IsDefault: in.GetIsDefault(),
	}
	if in.Port != nil {
		port := int32(*in.Port)
		record.Port = &port
	}

	if err := q.WithContext(l.ctx).Create(record); err != nil {
		return nil, dbError(l.Logger, err, in)
	}

	if in.GetIsDefault() {
		if _, err := q.WithContext(l.ctx).
			Where(q.ID.Neq(record.ID)).
			Update(q.IsDefault, false); err != nil {
			return nil, dbError(l.Logger, err, in)
		}
	}

	return &mcms.BaseIDResp{Id: uint64(record.ID), Msg: "success"}, nil
}
