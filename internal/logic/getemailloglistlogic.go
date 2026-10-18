package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmailLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetEmailLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmailLogListLogic {
	return &GetEmailLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: emaillog
func (l *GetEmailLogListLogic) GetEmailLogList(in *mcms.EmailLogListReq) (*mcms.EmailLogListResp, error) {
	q := l.svcCtx.Query.McmsEmailLog
	do := q.WithContext(l.ctx)
	if in.Target != nil {
		do = do.Where(q.Target.Like("%" + *in.Target + "%"))
	}
	if in.Subject != nil {
		do = do.Where(q.Subject.Like("%" + *in.Subject + "%"))
	}
	if in.Provider != nil {
		// Fixed: the reference filtered provider against ContentContains.
		do = do.Where(q.Provider.Like("%" + *in.Provider + "%"))
	}
	if in.SendStatus != nil && *in.SendStatus != 0 {
		do = do.Where(q.SendStatus.Eq(int32(*in.SendStatus)))
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

	resp := &mcms.EmailLogListResp{Total: uint64(total)}
	for _, v := range list {
		resp.Data = append(resp.Data, toEmailLogInfo(v))
	}
	return resp, nil
}
