package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSmsProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSmsProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSmsProviderLogic {
	return &UpdateSmsProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: smsprovider
func (l *UpdateSmsProviderLogic) UpdateSmsProvider(in *mcms.SmsProviderInfo) (*mcms.BaseResp, error) {
	if in.Id == nil {
		return nil, errorx.New(errorx.InvalidArgument, "id is required")
	}
	q := l.svcCtx.Query.McmsSmsProvider
	id := int64(*in.Id)

	if _, err := q.WithContext(l.ctx).Where(q.ID.Eq(id)).First(); err != nil {
		return nil, dbError(l.Logger, err, in)
	}

	updates := map[string]any{}
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.SecretId != nil {
		updates["secret_id"] = *in.SecretId
	}
	if in.SecretKey != nil {
		updates["secret_key"] = *in.SecretKey
	}
	if in.Region != nil {
		updates["region"] = *in.Region
	}
	if in.IsDefault != nil {
		updates["is_default"] = *in.IsDefault
	}
	if len(updates) > 0 {
		if _, err := q.WithContext(l.ctx).Where(q.ID.Eq(id)).Updates(updates); err != nil {
			return nil, dbError(l.Logger, err, in)
		}
	}

	// Only one provider may be the default.
	if in.GetIsDefault() {
		if _, err := q.WithContext(l.ctx).
			Where(q.ID.Neq(id)).
			Update(q.IsDefault, false); err != nil {
			return nil, dbError(l.Logger, err, in)
		}
	}

	return &mcms.BaseResp{Msg: "success"}, nil
}
