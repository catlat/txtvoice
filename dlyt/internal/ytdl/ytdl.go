package ytdl

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kkdai/youtube/v2"
)

// FetchBasic 获取视频基本信息、最佳缩略图，以及音频流和建议文件名与大小
func FetchBasic(ctx context.Context, idOrURL string) (info *youtube.Video, bestThumb youtube.Thumbnail, audio io.ReadCloser, audioFileName string, audioSize int64, err error) {
	// 尝试加载登录 Cookie（用于绕过年龄限制）
	httpClient := &http.Client{}
	if jar, jerr := loadCookieJarFromEnv(); jerr == nil && jar != nil {
		httpClient.Jar = jar
	}

	// 配置 YouTube 客户端
	client := youtube.Client{HTTPClient: httpClient}

	vidOrURL := extractVideoID(idOrURL)

	// 1) 首选使用库获取
	video, e := client.GetVideo(vidOrURL)
	if e != nil {
		// 2) 库失败则尝试 yt-dlp 兜底
		video, bestThumb, audio, audioFileName, audioSize, err = fetchWithYtDlp(ctx, httpClient, idOrURL)
		if err != nil {
			return nil, bestThumb, nil, "", 0, fmt.Errorf("video restricted or unavailable: %w", err)
		}
		return video, bestThumb, audio, audioFileName, audioSize, nil
	}

	// 选最大缩略图
	for _, t := range video.Thumbnails {
		if uint(t.Width)*uint(t.Height) > bestThumb.Width*bestThumb.Height {
			bestThumb = youtube.Thumbnail{URL: t.URL, Width: uint(t.Width), Height: uint(t.Height)}
		}
	}

	// 选最佳音频格式：优先 audio-only（MimeType 以 audio/ 开头），再按最高比特率
	audioFormats := video.Formats.WithAudioChannels()
	var candidates []*youtube.Format
	for i := range audioFormats {
		f := &audioFormats[i]
		if strings.HasPrefix(strings.ToLower(f.MimeType), "audio/") {
			candidates = append(candidates, f)
		}
	}
	if len(candidates) == 0 {
		// 没有纯音频时退化到所有含音频格式
		for i := range audioFormats {
			candidates = append(candidates, &audioFormats[i])
		}
	}

	var best *youtube.Format
	bestRate := 0
	for _, f := range candidates {
		if f.Bitrate > bestRate {
			bestRate = f.Bitrate
			best = f
		}
	}
	if best != nil {
		stream, _, e := client.GetStream(video, best)
		if e == nil {
			ext := originalAudioExt(best.MimeType)
			audio = stream
			audioFileName = fmt.Sprintf("%s.%s", video.ID, ext)
			audioSize = best.ContentLength
		}
	}

	return video, bestThumb, audio, audioFileName, audioSize, nil
}

// fetchWithYtDlp 使用 yt-dlp 兜底：读取信息、选择缩略图、获取直链并建立音频流
func fetchWithYtDlp(ctx context.Context, httpClient *http.Client, idOrURL string) (*youtube.Video, youtube.Thumbnail, io.ReadCloser, string, int64, error) {
	var zeroThumb youtube.Thumbnail

	cookiesPath := strings.TrimSpace(os.Getenv("YTDL_COOKIES_FILE"))
	if cookiesPath == "" {
		cookiesPath = "/www/wwwroot/cookie.txt"
	}

	// 1) 取 JSON 信息
	infoJSON, err := runYtDlpJSON(ctx, idOrURL, cookiesPath)
	if err != nil {
		return nil, zeroThumb, nil, "", 0, err
	}

	// 解析关键字段
	vid := firstNonEmptyStr(infoJSON.ID, infoJSON.ExtractorKey)
	if vid == "" {
		vid = extractVideoID(idOrURL)
	}
	publishDate := time.Time{}
	if len(infoJSON.UploadDate) == 8 {
		if t, perr := time.Parse("20060102", infoJSON.UploadDate); perr == nil {
			publishDate = t
		}
	}

	video := &youtube.Video{
		ID:          vid,
		Title:       infoJSON.Title,
		Author:      firstNonEmptyStr(infoJSON.Uploader, infoJSON.Channel),
		Duration:    time.Duration(infoJSON.Duration) * time.Second,
		Views:       int(infoJSON.ViewCount),
		PublishDate: publishDate,
	}

	// 2) 选择最大缩略图
	bestThumb := youtube.Thumbnail{}
	for _, th := range infoJSON.Thumbnails {
		w := uint(th.Width)
		h := uint(th.Height)
		if w*h > bestThumb.Width*bestThumb.Height {
			bestThumb = youtube.Thumbnail{URL: th.URL, Width: w, Height: h}
		}
	}

	// 3) 获取直链与扩展名
	audioURL, aext, err := runYtDlpBestAudioURL(ctx, idOrURL, cookiesPath)
	if err != nil {
		return video, bestThumb, nil, "", 0, err
	}

	// 4) 建立音频流（HEAD 获取长度与 MIME）
	size, mime, err := headContent(httpClient, audioURL)
	if err != nil {
		return video, bestThumb, nil, "", 0, err
	}
	body, err := getStream(httpClient, audioURL)
	if err != nil {
		return video, bestThumb, nil, "", 0, err
	}

	fileName := fmt.Sprintf("%s.%s", video.ID, normalizeAudioExtFromMimeOrExt(mime, aext))
	return video, bestThumb, body, fileName, size, nil
}

// runYtDlpJSON 执行 yt-dlp -j 获取 JSON 元数据
func runYtDlpJSON(ctx context.Context, input string, cookiesPath string) (*ytDlpJSON, error) {
	args := []string{"-j", "--no-playlist"}
	if cookiesPath != "" {
		args = append(args, "--cookies", cookiesPath)
	}
	args = append(args, input)
	out, err := exec.CommandContext(ctx, "yt-dlp", args...).Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp json: %w", err)
	}
	var obj ytDlpJSON
	if jerr := json.Unmarshal(out, &obj); jerr != nil {
		return nil, jerr
	}
	return &obj, nil
}

// runYtDlpBestAudioURL 获取最佳音频直链与扩展名
func runYtDlpBestAudioURL(ctx context.Context, input string, cookiesPath string) (string, string, error) {
	// 直链
	args := []string{"-f", "bestaudio", "-g", "--no-playlist"}
	if cookiesPath != "" {
		args = append(args, "--cookies", cookiesPath)
	}
	args = append(args, input)
	out, err := exec.CommandContext(ctx, "yt-dlp", args...).Output()
	if err != nil {
		return "", "", fmt.Errorf("yt-dlp -g: %w", err)
	}
	audioURL := strings.TrimSpace(string(out))
	if audioURL == "" {
		return "", "", errors.New("empty audio url")
	}
	// 扩展名
	args2 := []string{"-f", "bestaudio", "--get-filename", "-o", "%(id)s.%(ext)s", "--no-playlist"}
	if cookiesPath != "" {
		args2 = append(args2, "--cookies", cookiesPath)
	}
	args2 = append(args2, input)
	out2, err := exec.CommandContext(ctx, "yt-dlp", args2...).Output()
	if err != nil {
		return audioURL, "", nil
	}
	name := strings.TrimSpace(string(out2))
	ext := ""
	if dot := strings.LastIndex(name, "."); dot != -1 && dot < len(name)-1 {
		ext = name[dot+1:]
	}
	return audioURL, ext, nil
}

// 通过 HEAD 获取 Content-Length 与 Content-Type
func headContent(hc *http.Client, u string) (int64, string, error) {
	req, err := http.NewRequest(http.MethodHead, u, nil)
	if err != nil {
		return 0, "", err
	}
	resp, err := hc.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, "", fmt.Errorf("head status %d", resp.StatusCode)
	}
	size := resp.ContentLength
	mime := resp.Header.Get("Content-Type")
	return size, mime, nil
}

// GET 流
func getStream(hc *http.Client, u string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, fmt.Errorf("get status %d", resp.StatusCode)
	}
	return resp.Body, nil
}

// DownloadURL 下载器（用于缩略图），返回流、大小与 mime
func DownloadURL(ctx context.Context, u string) (io.ReadCloser, int64, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, 0, "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, "", err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, 0, "", fmt.Errorf("http %d", resp.StatusCode)
	}
	return resp.Body, resp.ContentLength, resp.Header.Get("Content-Type"), nil
}

func originalAudioExt(mime string) string {
	m := strings.ToLower(mime)
	if strings.HasPrefix(m, "audio/webm") {
		return "webm"
	}
	if strings.HasPrefix(m, "audio/mp4") {
		// YouTube 音频轨通常为 MPEG-4 Audio (m4a)
		return "m4a"
	}
	if strings.HasPrefix(m, "audio/3gpp") || strings.Contains(m, "3gpp") {
		return "3gp"
	}
	if strings.HasPrefix(m, "audio/mpeg") || strings.Contains(m, "mpeg") {
		return "mp3"
	}
	// 回退：从容器名猜测
	if strings.Contains(m, "webm") {
		return "webm"
	}
	if strings.Contains(m, "mp4") {
		return "m4a"
	}
	if strings.Contains(m, "3gpp") {
		return "3gp"
	}
	return "audio"
}

func normalizeAudioExtFromMimeOrExt(mime string, ext string) string {
	if ext != "" {
		return strings.ToLower(ext)
	}
	return originalAudioExt(mime)
}

func extractVideoID(input string) string {
	patterns := []string{
		`(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([a-zA-Z0-9_-]{11})`,
		`^([a-zA-Z0-9_-]{11})$`,
	}
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		m := re.FindStringSubmatch(strings.TrimSpace(input))
		if len(m) >= 2 {
			return m[1]
		}
	}
	return strings.TrimSpace(input)
}

// loadCookieJarFromEnv 读取环境变量或默认路径的 cookies.txt，返回可用的 CookieJar
func loadCookieJarFromEnv() (*cookiejar.Jar, error) {
	path := strings.TrimSpace(os.Getenv("YTDL_COOKIES_FILE"))
	if path == "" {
		// 默认路径（你提供的路径）
		path = "/www/wwwroot/cookie.txt"
	}
	// 如果文件不存在则跳过
	if st, err := os.Stat(path); err != nil || st.IsDir() {
		return nil, nil
	}
	return loadCookieJarFromFile(path)
}

// loadCookieJarFromFile 解析 Netscape cookies.txt 并填充 CookieJar
func loadCookieJarFromFile(path string) (*cookiejar.Jar, error) {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	jar, _ := cookiejar.New(nil)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		// 处理 #HttpOnly_ 前缀
		httpOnly := false
		if strings.HasPrefix(line, "#HttpOnly_") {
			line = strings.TrimPrefix(line, "#HttpOnly_")
			httpOnly = true
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 7 {
			continue
		}

		domain := parts[0]
		// includeSub := parts[1]
		pathField := parts[2]
		secure := strings.EqualFold(parts[3], "TRUE")
		expireStr := parts[4]
		name := parts[5]
		value := parts[6]

		// 解析过期时间
		var expires time.Time
		if exp, err := strconv.ParseInt(expireStr, 10, 64); err == nil && exp > 0 {
			expires = time.Unix(exp, 0)
		}

		// 清理域名
		if strings.HasPrefix(domain, ".") {
			domain = domain[1:]
		}
		if domain == "" || name == "" {
			continue
		}

		ck := &http.Cookie{
			Name:     name,
			Value:    value,
			Path:     pathField,
			Domain:   domain,
			Secure:   secure,
			HttpOnly: httpOnly,
		}
		if !expires.IsZero() {
			ck.Expires = expires
		}

		// 将 cookie 写入 jar（使用 https 优先）
		for _, h := range []string{domain, "." + domain} {
			u := &url.URL{Scheme: "https", Host: h, Path: "/"}
			jar.SetCookies(u, append(jar.Cookies(u), ck))
		}
	}
	_ = scanner.Err()
	return jar, nil
}

// yt-dlp -j 返回的数据子集结构
type ytDlpJSON struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Uploader     string `json:"uploader"`
	Channel      string `json:"channel"`
	Duration     int64  `json:"duration"`
	ViewCount    int64  `json:"view_count"`
	UploadDate   string `json:"upload_date"`
	ExtractorKey string `json:"extractor_key"`
	Thumbnails   []struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnails"`
}

func firstNonEmptyStr(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
