// Package logic contains the mcms RPC implementations.
package logic

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/wangminghui0622-bit/admin-common/dal/model"
	"github.com/wangminghui0622-bit/admin-common/errorx"
	"github.com/wangminghui0622-bit/admin-message/mcms"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// ptr returns a pointer to v (proto optional fields are pointers).
func ptr[T any](v T) *T { return &v }

// pageParams normalises page/page_size to gorm Offset/Limit.
func pageParams(page, pageSize uint64) (offset, limit int) {
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	return int((page - 1) * pageSize), int(pageSize)
}

// newUUID returns a RFC 9562 UUID v7, falling back to v4 if the system clock
// cannot provide the needed entropy.
func newUUID() string {
	if id, err := uuid.NewV7(); err == nil {
		return id.String()
	}
	return uuid.New().String()
}

// dbError converts gorm errors to admin-zero business errors.
func dbError(logger logx.Logger, err error, detail any) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errorx.New(errorx.InvalidArgument, "record not found")
	}
	if strings.Contains(err.Error(), "Duplicate entry") || strings.Contains(err.Error(), "Error 1062") {
		logger.Errorw("duplicate database record", logx.Field("detail", err.Error()), logx.Field("data", detail))
		return errorx.New(errorx.InvalidArgument, "duplicate record")
	}
	logger.Errorw("database error", logx.Field("detail", err.Error()), logx.Field("data", detail))
	return errorx.New(errorx.ServerError, "database error")
}

func toSmsProviderInfo(v *model.McmsSmsProvider) *mcms.SmsProviderInfo {
	return &mcms.SmsProviderInfo{
		Id:        ptr(uint64(v.ID)),
		CreatedAt: ptr(v.CreatedAt.UnixMilli()),
		UpdatedAt: ptr(v.UpdatedAt.UnixMilli()),
		Name:      ptr(v.Name),
		SecretId:  ptr(v.SecretID),
		SecretKey: ptr(v.SecretKey),
		Region:    ptr(v.Region),
		IsDefault: ptr(v.IsDefault),
	}
}

func toEmailProviderInfo(v *model.McmsEmailProvider) *mcms.EmailProviderInfo {
	var port *uint32
	if v.Port != nil {
		port = ptr(uint32(*v.Port))
	}
	return &mcms.EmailProviderInfo{
		Id:        ptr(uint64(v.ID)),
		CreatedAt: ptr(v.CreatedAt.UnixMilli()),
		UpdatedAt: ptr(v.UpdatedAt.UnixMilli()),
		Name:      ptr(v.Name),
		AuthType:  ptr(uint32(v.AuthType)),
		EmailAddr: ptr(v.EmailAddr),
		Password:  v.Password,
		HostName:  ptr(v.HostName),
		Identify:  v.Identify,
		Secret:    v.Secret,
		Port:      port,
		Tls:       ptr(v.TLS),
		IsDefault: ptr(v.IsDefault),
	}
}

func toSmsLogInfo(v *model.McmsSmsLog) *mcms.SmsLogInfo {
	return &mcms.SmsLogInfo{
		Id:          ptr(v.ID),
		CreatedAt:   ptr(v.CreatedAt.UnixMilli()),
		UpdatedAt:   ptr(v.UpdatedAt.UnixMilli()),
		PhoneNumber: ptr(v.PhoneNumber),
		Content:     ptr(v.Content),
		SendStatus:  ptr(uint32(v.SendStatus)),
		Provider:    ptr(v.Provider),
	}
}

func toEmailLogInfo(v *model.McmsEmailLog) *mcms.EmailLogInfo {
	return &mcms.EmailLogInfo{
		Id:         ptr(v.ID),
		CreatedAt:  ptr(v.CreatedAt.UnixMilli()),
		UpdatedAt:  ptr(v.UpdatedAt.UnixMilli()),
		Target:     ptr(v.Target),
		Subject:    ptr(v.Subject),
		Content:    ptr(v.Content),
		SendStatus: ptr(uint32(v.SendStatus)),
		Provider:   ptr(v.Provider),
	}
}
