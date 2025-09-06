概览
基址: http://43.133.180.227:9005
认证: 无（内网/受控网络使用，建议将服务置于内网或反向代理后）
编码: UTF-8
Content-Type: application/json
方法: 均为 POST
通用请求体
Body
字段
id_or_url: 必填，支持 11 位视频 ID（如 jNQXAC9IVRw）或完整链接（如 https://www.youtube.com/watch?v=jNQXAC9IVRw）
错误处理（通用）
400 Bad Request: 请求体缺失或格式错误（返回纯文本：bad request: need id_or_url）
405 Method Not Allowed: 非 POST
502 Bad Gateway: 上游解析/下载失败（返回纯文本：upstream error）
成功返回均为 JSON；错误返回为纯文本
年龄限制/地区限制
服务端已集成 cookies 及 yt-dlp 兜底，默认从 /www/wwwroot/cookie.txt 读取 Netscape cookies.txt（可用环境变量 YTDL_COOKIES_FILE 覆盖）
若 cookies 过期或丢失，可能返回 502（upstream error）
1) 基本信息 + 缩略图上传接口
URL: POST /api/yt/info
功能: 获取视频基础信息，并下载最佳缩略图上传到对象存储，返回缩略图直链
请求体: 通用请求体
成功响应 200
字段说明
id/title/author: 视频基本信息
duration_sec: 时长（秒）
views: 播放量（取源端可用字段，可能出现延迟差异）
publish_date: 发布日期（yyyy-MM-dd，若源数据缺失可能为 0001-01-01）
thumbnail_url: 已上传到对象存储（当前域名 https://oss.duckai.cn）的缩略图直链
错误响应
400/405/502 见通用错误
示例
curl
PowerShell
2) 仅音频上传接口
URL: POST /api/yt/audio
功能: 仅处理音频。选择最佳音频流（优先 audio-only），流式上传到对象存储，返回音频直链
输出保持“原容器格式”（不转码），常见为 .m4a（AAC）或 .webm（Opus）
成功响应 200
字段说明
id/title: 视频标识与标题
audio_url: 已上传对象存储的音频直链（根据源流格式返回 .m4a、.webm、偶见 .3gp）
错误响应
400/405/502 见通用错误
示例
curl
PowerShell
返回域名与资源路径
当前对象存储访问域名：https://oss.duckai.cn
资源路径格式：
缩略图：/thumbs/<videoId>_<width>x<height>.jpg
音频：/audio/<videoId>.<ext>