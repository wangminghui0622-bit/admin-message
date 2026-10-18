package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmailProviderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmailProviderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmailProviderListLogic {
	return &GetEmailProviderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: emailprovider
func (l *GetEmailProviderListLogic) GetEmailProviderList(in *mcms.EmailProviderListReq) (*mcms.EmailProviderListResp, error) {
	q := l.svcCtx.Query.McmsEmailProvider
	do := q.WithContext(l.ctx)
	if in.Name != nil {
		do = do.Where(q.Name.Like("%" + *in.Name + "%"))
	}
	if in.EmailAddr != nil {
		do = do.Where(q.EmailAddr.Like("%" + *in.EmailAddr + "%"))
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

	resp := &mcms.EmailProviderListResp{Total: uint64(total)}
	for _, v := range list {
		resp.Data = append(resp.Data, toEmailProviderInfo(v))
	}
	return resp, nil
}
