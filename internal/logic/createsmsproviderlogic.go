package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSmsProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSmsProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSmsProviderLogic {
	return &CreateSmsProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: smsprovider
func (l *CreateSmsProviderLogic) CreateSmsProvider(in *mcms.SmsProviderInfo) (*mcms.BaseIDResp, error) {
	q := l.svcCtx.Query.McmsSmsProvider
	record := &model.McmsSmsProvider{
		Name:      in.GetName(),
		SecretID:  in.GetSecretId(),
		SecretKey: in.GetSecretKey(),
		Region:    in.GetRegion(),
		IsDefault: in.GetIsDefault(),
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
