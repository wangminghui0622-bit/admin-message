package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateSmsLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateSmsLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateSmsLogLogic {
	return &CreateSmsLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: smslog
func (l *CreateSmsLogLogic) CreateSmsLog(in *mcms.SmsLogInfo) (*mcms.BaseUUIDResp, error) {
	id := newUUID()
	record := &model.McmsSmsLog{
		ID:          id,
		PhoneNumber: in.GetPhoneNumber(),
		Content:     in.GetContent(),
		SendStatus:  int32(in.GetSendStatus()),
		Provider:    in.GetProvider(),
	}
	if err := l.svcCtx.Query.McmsSmsLog.WithContext(l.ctx).Create(record); err != nil {
		return nil, dbError(l.Logger, err, in)
	}
	return &mcms.BaseUUIDResp{Id: id, Msg: "success"}, nil
}
