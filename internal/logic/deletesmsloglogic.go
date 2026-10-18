package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSmsLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSmsLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSmsLogLogic {
	return &DeleteSmsLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: smslog
func (l *DeleteSmsLogLogic) DeleteSmsLog(in *mcms.UUIDsReq) (*mcms.BaseResp, error) {
	if len(in.Ids) > 0 {
		q := l.svcCtx.Query.McmsSmsLog
		if _, err := q.WithContext(l.ctx).Where(q.ID.In(in.Ids...)).Delete(); err != nil {
			return nil, dbError(l.Logger, err, in)
		}
	}
	return &mcms.BaseResp{Msg: "success"}, nil
}
