package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteEmailLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteEmailLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEmailLogLogic {
	return &DeleteEmailLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: emaillog
func (l *DeleteEmailLogLogic) DeleteEmailLog(in *mcms.UUIDsReq) (*mcms.BaseResp, error) {
	if len(in.Ids) > 0 {
		q := l.svcCtx.Query.McmsEmailLog
		if _, err := q.WithContext(l.ctx).Where(q.ID.In(in.Ids...)).Delete(); err != nil {
			return nil, dbError(l.Logger, err, in)
		}
	}
	return &mcms.BaseResp{Msg: "success"}, nil
}
