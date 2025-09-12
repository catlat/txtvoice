package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&AlterTTSHistoryTextColumns20250912094500{})
}

// AlterTTSHistoryTextColumns20250912094500 将 text_preview、text_hash、audio_url 全部改为 TEXT，并重建唯一索引
type AlterTTSHistoryTextColumns20250912094500 struct{}

// Up 执行迁移
func (m *AlterTTSHistoryTextColumns20250912094500) Up(migrator *migration.DDLMigrator) error {
	return migrator.Exec(`
        ALTER TABLE tts_history
            DROP INDEX uk_identity_text_speaker,
            MODIFY COLUMN text_preview TEXT COMMENT '文本预览',
            MODIFY COLUMN text_hash TEXT NOT NULL COMMENT 'text+speaker+defaults 的sha256',
            MODIFY COLUMN audio_url TEXT COMMENT '音频直链',
            ADD UNIQUE KEY uk_identity_text_speaker (user_identity, text_hash(64), speaker);
    `)
}
