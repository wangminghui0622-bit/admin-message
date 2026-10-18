package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmailProviderByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmailProviderByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmailProviderByIdLogic {
	return &GetEmailProviderByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: emailprovider
func (l *GetEmailProviderByIdLogic) GetEmailProviderById(in *mcms.IDReq) (*mcms.EmailProviderInfo, error) {
	q := l.svcCtx.Query.McmsEmailProvider
	record, err := q.WithContext(l.ctx).Where(q.ID.Eq(int64(in.Id))).First()
	if err != nil {
		return nil, dbError(l.Logger, err, in)
	}
	return toEmailProviderInfo(record), nil
}
