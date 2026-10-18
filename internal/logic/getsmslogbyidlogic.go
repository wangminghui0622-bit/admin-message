package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSmsLogByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSmsLogByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSmsLogByIdLogic {
	return &GetSmsLogByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: smslog
func (l *GetSmsLogByIdLogic) GetSmsLogById(in *mcms.UUIDReq) (*mcms.SmsLogInfo, error) {
	q := l.svcCtx.Query.McmsSmsLog
	record, err := q.WithContext(l.ctx).Where(q.ID.Eq(in.Id)).First()
	if err != nil {
		return nil, dbError(l.Logger, err, in)
	}
	return toSmsLogInfo(record), nil
}
