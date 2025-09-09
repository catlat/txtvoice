package model

type AccountUser struct {
	Id           int64  `gorm:"column:id;primaryKey" json:"id"`
	Identity     string `gorm:"column:identity" json:"identity"`
	IdentityType int    `gorm:"column:identity_type" json:"identity_type"`
	DisplayName  string `gorm:"column:display_name" json:"display_name"`
	Status       int    `gorm:"column:status" json:"status"`
}

func (AccountUser) TableName() string { return "account_user" }
