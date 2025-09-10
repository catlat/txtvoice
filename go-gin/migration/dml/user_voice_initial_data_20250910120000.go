package dml

import (
	"go-gin/internal/migration"

	"gorm.io/gorm"
)

func init() {
	migration.RegisterDML(&UserVoiceInitialData20250910120000{})
}

// UserVoiceInitialData20250910120000 插入用户音色ID初始数据
type UserVoiceInitialData20250910120000 struct{}

// Handle 执行迁移
func (m *UserVoiceInitialData20250910120000) Handle(db *gorm.DB) error {
	return db.Exec(`
		INSERT INTO user_voice (mobile, voice_id) VALUES
		('15173107332', 'S_br0aQwvC1'),
		('15799537457', 'S_zFd0Cv9C1'),
		('18503065953', 'S_vX5vle9C1'),
		('15871770680', 'S_10JFs89C1'),
		('17680010746', 'S_jRDQm89C1'),
		('15388005188', 'S_Xj0QPp2C1'),
		('17673090601', 'S_h4j01t0C1'),
		('15874226113', 'S_1gpgmHWB1'),
		('19118993851', 'S_H5dym7WB1'),
		('13454783099', 'S_FGoiA3WB1'),
		('18674881384', 'S_2Kgax3vB1'),
		('18874748888', 'S_r3YGBCoB1'),
		('18670708294', 'S_f35T4BkB1'),
		('17673094952', 'S_xuSXQxWB1'),
		('13217321273', 'S_BW90HQ5C1'),
		('19173620148', 'S_hlJNWK5C1')
		ON DUPLICATE KEY UPDATE voice_id = VALUES(voice_id)
	`).Error
}

// Desc 获取迁移描述
func (m *UserVoiceInitialData20250910120000) Desc() string {
	return "插入用户音色ID关系初始数据"
}
