package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSmsProviderByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSmsProviderByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSmsProviderByIdLogic {
	return &GetSmsProviderByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: smsprovider
func (l *GetSmsProviderByIdLogic) GetSmsProviderById(in *mcms.IDReq) (*mcms.SmsProviderInfo, error) {
	q := l.svcCtx.Query.McmsSmsProvider
	record, err := q.WithContext(l.ctx).Where(q.ID.Eq(int64(in.Id))).First()
	if err != nil {
		return nil, dbError(l.Logger, err, in)
	}
	return toSmsProviderInfo(record), nil
}
