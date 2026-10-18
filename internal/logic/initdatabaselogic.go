package logic

import (
	"context"
	_ "embed"

	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-message/internal/svc"
	"github.com/wangminghui0622-bit/admin-message/mcms"

	"github.com/zeromicro/go-zero/core/logx"
)

// initSQL 由 admin_zero 的 mysqldump 生成（DDL + 种子数据 + API/授权），
// CREATE TABLE IF NOT EXISTS / INSERT IGNORE，可安全重复执行。
//
//go:embed init.sql
var initSQL string

const initKey = "initDatabase:mcms"

type InitDatabaseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInitDatabaseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitDatabaseLogic {
	return &InitDatabaseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// group: base
func (l *InitDatabaseLogic) InitDatabase(in *mcms.Empty) (*mcms.BaseResp, error) {
	if ok, err := l.svcCtx.Redis.Exists(initKey); err != nil {
		l.Logger.Errorw("failed to read init key", logx.Field("detail", err.Error()))
		return nil, errorx.New(errorx.ServerError, "failed to read init key")
	} else if ok {
		return &mcms.BaseResp{Msg: "already initialized"}, nil
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).Exec(initSQL).Error; err != nil {
		l.Logger.Errorw("init sql failed", logx.Field("detail", err.Error()))
		return nil, errorx.New(errorx.ServerError, "init database failed: "+err.Error())
	}

	if err := l.svcCtx.Redis.Set(initKey, "1"); err != nil {
		l.Logger.Errorw("failed to set init key", logx.Field("detail", err.Error()))
		return nil, errorx.New(errorx.ServerError, "failed to set init key")
	}

	return &mcms.BaseResp{Msg: "success"}, nil
}
