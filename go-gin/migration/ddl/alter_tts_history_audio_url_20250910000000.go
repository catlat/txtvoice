package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&AlterTTSHistoryAudioURL20250910000000{})
}

// AlterTTSHistoryAudioURL20250910000000 将 tts_history.audio_url 从 VARCHAR(512) 扩展为 TEXT
type AlterTTSHistoryAudioURL20250910000000 struct{}

// Up 执行迁移
func (m *AlterTTSHistoryAudioURL20250910000000) Up(migrator *migration.DDLMigrator) error {
	return migrator.Exec(`
		ALTER TABLE tts_history 
			MODIFY COLUMN audio_url TEXT COMMENT '音频直链 (可能为data URL)';
	`)
}

