package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&CreateYoutubeTranscript20250906090100{})
}

// CreateYoutubeTranscript20250906090100 创建 youtube_transcript 表（外键关联 youtube_video.id）
type CreateYoutubeTranscript20250906090100 struct{}

// Up 执行迁移
func (m *CreateYoutubeTranscript20250906090100) Up(migrator *migration.DDLMigrator) error {
	return migrator.Exec(`
		CREATE TABLE IF NOT EXISTS youtube_transcript (
			id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			video_id BIGINT NOT NULL COMMENT 'youtube_video.id',
			language VARCHAR(16) NOT NULL COMMENT '语言代码，如en/zh',
			original_text LONGTEXT NULL COMMENT '原文',
			translated_text LONGTEXT NULL COMMENT '译文',
			asr_char_count INT DEFAULT 0 COMMENT 'ASR计数字符数',
			translate_char_count INT DEFAULT 0 COMMENT '翻译计数字符数',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			UNIQUE KEY uk_video_lang (video_id, language),
			KEY idx_video_id (video_id),
			CONSTRAINT fk_transcript_video FOREIGN KEY (video_id) REFERENCES youtube_video(id) ON DELETE CASCADE ON UPDATE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
}


