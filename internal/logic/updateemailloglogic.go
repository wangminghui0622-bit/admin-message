package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateEmailLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateEmailLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmailLogLogic {
	return &UpdateEmailLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: emaillog
func (l *UpdateEmailLogLogic) UpdateEmailLog(in *mcms.EmailLogInfo) (*mcms.BaseResp, error) {
	if in.Id == nil {
		return nil, errorx.New(errorx.InvalidArgument, "id is required")
	}
	q := l.svcCtx.Query.McmsEmailLog
	if _, err := q.WithContext(l.ctx).Where(q.ID.Eq(*in.Id)).First(); err != nil {
		return nil, dbError(l.Logger, err, in)
	}

	updates := map[string]any{}
	if in.Target != nil {
		updates["target"] = *in.Target
	}
	if in.Subject != nil {
		updates["subject"] = *in.Subject
	}
	if in.Content != nil {
		updates["content"] = *in.Content
	}
	if in.SendStatus != nil {
		updates["send_status"] = int32(*in.SendStatus)
	}
	if in.Provider != nil {
		updates["provider"] = *in.Provider
	}
	if len(updates) > 0 {
		if _, err := q.WithContext(l.ctx).Where(q.ID.Eq(*in.Id)).Updates(updates); err != nil {
			return nil, dbError(l.Logger, err, in)
		}
	}
	return &mcms.BaseResp{Msg: "success"}, nil
}
