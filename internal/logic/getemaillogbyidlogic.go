package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmailLogByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmailLogByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmailLogByIdLogic {
	return &GetEmailLogByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: emaillog
func (l *GetEmailLogByIdLogic) GetEmailLogById(in *mcms.UUIDReq) (*mcms.EmailLogInfo, error) {
	q := l.svcCtx.Query.McmsEmailLog
	record, err := q.WithContext(l.ctx).Where(q.ID.Eq(in.Id)).First()
	if err != nil {
		return nil, dbError(l.Logger, err, in)
	}
	return toEmailLogInfo(record), nil
}
