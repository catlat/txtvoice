package model

import (
	"context"
	"go-gin/internal/component/db"
	"time"

	"gorm.io/gorm/clause"
)

type UserVoice struct {
	Id        int64     `gorm:"column:id;primary_key;auto_increment" json:"id"`
	Mobile    string    `gorm:"column:mobile" json:"mobile"`
	VoiceId   string    `gorm:"column:voice_id" json:"voice_id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (u *UserVoice) TableName() string {
	return `user_voice`
}

type UserVoiceModel struct{}

func NewUserVoiceModel() *UserVoiceModel {
	return &UserVoiceModel{}
}

// List 获取所有用户音色关系
func (m *UserVoiceModel) List(ctx context.Context) ([]UserVoice, error) {
	var userVoices []UserVoice
	return userVoices, db.WithContext(ctx).Find(&userVoices).Error()
}

// Add 添加用户音色关系
func (m *UserVoiceModel) Add(ctx context.Context, userVoice *UserVoice) error {
	return db.WithContext(ctx).Create(userVoice).Error()
}

// GetByMobile 根据手机号获取音色ID
func (m *UserVoiceModel) GetByMobile(ctx context.Context, mobile string) (*UserVoice, error) {
	var userVoice UserVoice
	err := db.WithContext(ctx).Where("mobile = ?", mobile).First(&userVoice).Error()
	return &userVoice, err
}

// UpdateVoiceId 更新音色ID
func (m *UserVoiceModel) UpdateVoiceId(ctx context.Context, mobile string, voiceId string) error {
	return db.WithContext(ctx).Model(&UserVoice{}).Where("mobile = ?", mobile).Update("voice_id", voiceId).Error
}

// DeleteByMobile 根据手机号删除记录
func (m *UserVoiceModel) DeleteByMobile(ctx context.Context, mobile string) error {
	return db.WithContext(ctx).Where("mobile = ?", mobile).Delete(&UserVoice{}).Error()
}

// Upsert a user voice record. If the mobile number already exists, it updates the voice_id.
func (m *UserVoiceModel) Upsert(ctx context.Context, userVoice *UserVoice) error {
	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "mobile"}},
		DoUpdates: clause.AssignmentColumns([]string{"voice_id"}),
	}).Create(userVoice).Error()
}
