package ddl

import (
	"go-gin/internal/migration"
)

func init() {
	migration.RegisterDDL(&CreateYoutubeVideo20250906090000{})
}

// CreateYoutubeVideo20250906090000 创建 youtube_video 表
type CreateYoutubeVideo20250906090000 struct{}

// Up 执行迁移
func (m *CreateYoutubeVideo20250906090000) Up(migrator *migration.DDLMigrator) error {
	return migrator.Exec(`
		CREATE TABLE IF NOT EXISTS youtube_video (
			id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
			source_site VARCHAR(16) NOT NULL COMMENT '来源站点，例如youtube',
			video_id VARCHAR(64) NOT NULL COMMENT '平台视频ID',
			title VARCHAR(512) DEFAULT '' COMMENT '标题',
			description TEXT NULL COMMENT '描述',
			channel_title VARCHAR(256) DEFAULT '' COMMENT '频道名称',
			duration_sec INT DEFAULT 0 COMMENT '时长(秒)',
			published_at DATETIME NULL COMMENT '发布日期',
			thumbnail_url VARCHAR(512) DEFAULT '' COMMENT '缩略图直链',
			audio_url VARCHAR(512) DEFAULT '' COMMENT '音频直链',
			status TINYINT DEFAULT 0 COMMENT '状态',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
			UNIQUE KEY uk_source_video (source_site, video_id),
			KEY idx_created_at (created_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`)
}


