package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSmsProviderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSmsProviderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSmsProviderListLogic {
	return &GetSmsProviderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: smsprovider
func (l *GetSmsProviderListLogic) GetSmsProviderList(in *mcms.SmsProviderListReq) (*mcms.SmsProviderListResp, error) {
	q := l.svcCtx.Query.McmsSmsProvider
	do := q.WithContext(l.ctx)
	if in.Name != nil {
		do = do.Where(q.Name.Like("%" + *in.Name + "%"))
	}

	total, err := do.Count()
	if err != nil {
		return nil, dbError(l.Logger, err, in)
	}

	offset, limit := pageParams(in.Page, in.PageSize)
	list, err := do.Order(q.ID.Desc()).Offset(offset).Limit(limit).Find()
	if err != nil {
		return nil, dbError(l.Logger, err, in)
	}

	resp := &mcms.SmsProviderListResp{Total: uint64(total)}
	for _, v := range list {
		resp.Data = append(resp.Data, toSmsProviderInfo(v))
	}
	return resp, nil
}
