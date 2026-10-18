package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteEmailProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteEmailProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEmailProviderLogic {
	return &DeleteEmailProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: emailprovider
func (l *DeleteEmailProviderLogic) DeleteEmailProvider(in *mcms.IDsReq) (*mcms.BaseResp, error) {
	if len(in.Ids) > 0 {
		ids := make([]int64, 0, len(in.Ids))
		for _, id := range in.Ids {
			ids = append(ids, int64(id))
		}
		q := l.svcCtx.Query.McmsEmailProvider
		if _, err := q.WithContext(l.ctx).Where(q.ID.In(ids...)).Delete(); err != nil {
			return nil, dbError(l.Logger, err, in)
		}
	}

	// Cached SMTP clients may belong to deleted rows.
	l.svcCtx.Email.Clear()

	return &mcms.BaseResp{Msg: "success"}, nil
}
