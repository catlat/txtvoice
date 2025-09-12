package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&AlterTTSHistoryTextHashChar6420250912101000{})
}

// AlterTTSHistoryTextHashChar6420250912101000 将 text_hash 设回 CHAR(64)，并保持 text_preview/audio_url 为 TEXT，重建唯一索引
type AlterTTSHistoryTextHashChar6420250912101000 struct{}

// Up 执行迁移
func (m *AlterTTSHistoryTextHashChar6420250912101000) Up(migrator *migration.DDLMigrator) error {
	return migrator.Exec(`
        ALTER TABLE tts_history
            DROP INDEX uk_identity_text_speaker,
            MODIFY COLUMN text_preview TEXT COMMENT '文本预览',
            MODIFY COLUMN audio_url TEXT COMMENT '音频直链',
            MODIFY COLUMN text_hash CHAR(64) NOT NULL COMMENT 'text+speaker+defaults 的sha256',
            ADD UNIQUE KEY uk_identity_text_speaker (user_identity, text_hash, speaker);
    `)
}
