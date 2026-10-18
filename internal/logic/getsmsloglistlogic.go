package logic

import (
	"context"

	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSmsLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSmsLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSmsLogListLogic {
	return &GetSmsLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: smslog
func (l *GetSmsLogListLogic) GetSmsLogList(in *mcms.SmsLogListReq) (*mcms.SmsLogListResp, error) {
	q := l.svcCtx.Query.McmsSmsLog
	do := q.WithContext(l.ctx)
	if in.PhoneNumber != nil {
		do = do.Where(q.PhoneNumber.Like("%" + *in.PhoneNumber + "%"))
	}
	if in.Content != nil {
		do = do.Where(q.Content.Like("%" + *in.Content + "%"))
	}
	if in.Provider != nil {
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

	resp := &mcms.SmsLogListResp{Total: uint64(total)}
	for _, v := range list {
		resp.Data = append(resp.Data, toSmsLogInfo(v))
	}
	return resp, nil
}
