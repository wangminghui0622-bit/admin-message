package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSmsProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSmsProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSmsProviderLogic {
	return &DeleteSmsProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: smsprovider
func (l *DeleteSmsProviderLogic) DeleteSmsProvider(in *mcms.IDsReq) (*mcms.BaseResp, error) {
	if len(in.Ids) > 0 {
		ids := make([]int64, 0, len(in.Ids))
		for _, id := range in.Ids {
			ids = append(ids, int64(id))
		}
		q := l.svcCtx.Query.McmsSmsProvider
		if _, err := q.WithContext(l.ctx).Where(q.ID.In(ids...)).Delete(); err != nil {
			return nil, dbError(l.Logger, err, in)
		}
	}
	return &mcms.BaseResp{Msg: "success"}, nil
}
