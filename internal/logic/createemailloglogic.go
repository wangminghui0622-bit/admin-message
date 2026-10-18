package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateEmailLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateEmailLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateEmailLogLogic {
	return &CreateEmailLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: emaillog
func (l *CreateEmailLogLogic) CreateEmailLog(in *mcms.EmailLogInfo) (*mcms.BaseUUIDResp, error) {
	id := newUUID()
	record := &model.McmsEmailLog{
		ID:         id,
		Target:     in.GetTarget(),
		Subject:    in.GetSubject(),
		Content:    in.GetContent(),
		SendStatus: int32(in.GetSendStatus()),
		Provider:   in.GetProvider(),
	}
	if err := l.svcCtx.Query.McmsEmailLog.WithContext(l.ctx).Create(record); err != nil {
		return nil, dbError(l.Logger, err, in)
	}
	return &mcms.BaseUUIDResp{Id: id, Msg: "success"}, nil
}
