#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
export_edge_cookie.py - 从 Edge 浏览器导出 Cookie 并直接存储到 Redis

使用方法:
1. 启动 Edge 浏览器（带调试端口）：
   "C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe" --remote-debugging-port=9222
   
2. 在浏览器中登录对应的平台（YouTube 或 哔哩哔哩）

3. 运行脚本：
   python export_edge_cookie.py youtube   # 导出 YouTube Cookie
   python export_edge_cookie.py bilibili  # 导出哔哩哔哩 Cookie
   
4. Cookie 将直接存储到 Redis 中，供 yt-dlp 使用

依赖安装:
pip install playwright redis
playwright install chromium
"""

import time
import redis
from pathlib import Path
from playwright.sync_api import sync_playwright

# 平台配置
PLATFORM_CONFIG = {
    "youtube": {
        "domains": (".youtube.com", ".google.com", "youtube.com", "google.com"),
        "urls": [
            "https://www.youtube.com",
            "https://accounts.google.com",
            "https://www.google.com",
        ]
    },
    "bilibili": {
        "domains": (".bilibili.com", "bilibili.com", ".bilivideo.com", "bilivideo.com"),
        "urls": [
            "https://www.bilibili.com",
            "https://passport.bilibili.com",
            "https://api.bilibili.com",
        ]
    }
}

# Redis配置
REDIS_HOST = "180.163.80.126"
REDIS_PORT = 6379
REDIS_DB = 9
REDIS_PASSWORD = "AyrC9eiVdoids"  # 如果有密码请设置

def to_netscape(cookies, allowed_domains):
    lines = ["# Netscape HTTP Cookie File"]
    now_plus = int(time.time()) + 180 * 24 * 3600
    for c in cookies:
        d = c.get("domain","")
        if not any(d.endswith(x) or d == x for x in allowed_domains):
            continue
        is_host_only = not d.startswith(".")
        domain_field = d
        include_sub = "FALSE" if is_host_only else "TRUE"
        path = c.get("path","/")
        secure = "TRUE" if c.get("secure") else "FALSE"
        exp = c.get("expires")
        expires = str(int(exp)) if isinstance(exp,(int,float)) and exp>0 else str(now_plus)
        name = c.get("name","")
        value = c.get("value","")
        lines.append(f"{domain_field}\t{include_sub}\t{path}\t{secure}\t{expires}\t{name}\t{value}")
    return "\n".join(lines) + "\n"

def save_cookie_to_redis(cookie_content, platform="youtube"):
    """将cookie内容保存到Redis"""
    try:
        # 连接Redis
        r = redis.Redis(
            host=REDIS_HOST,
            port=REDIS_PORT,
            db=REDIS_DB,
            password=REDIS_PASSWORD,
            decode_responses=True
        )
        
        # 测试连接
        r.ping()
        
        # 创建临时cookie文件（包含时间戳避免冲突）
        timestamp = int(time.time())
        temp_cookie_file = Path(f"{platform}_cookies_{timestamp}.txt")
        temp_cookie_file.write_text(cookie_content, encoding="utf-8")
        
        # 将cookie文件路径存储到Redis
        redis_key = f"ytdl:cookies:{platform}"
        cookie_file_path = str(temp_cookie_file.resolve())
        
        r.set(redis_key, cookie_file_path)
        print(f"[OK] Cookie已保存到Redis")
        print(f"    Redis Key: {redis_key}")
        print(f"    Cookie File: {cookie_file_path}")
        print(f"    Cookie Size: {len(cookie_content)} 字符")
        
        # 验证存储
        stored_path = r.get(redis_key)
        if stored_path == cookie_file_path:
            print(f"[OK] Redis 存储验证成功")
        else:
            print(f"[WARN] Redis 存储验证失败")
        
        return temp_cookie_file
        
    except redis.RedisError as e:
        raise RuntimeError(f"Redis操作失败: {e}")
    except Exception as e:
        raise RuntimeError(f"保存cookie失败: {e}")

def main(platform="youtube"):
    # 验证平台支持
    if platform not in PLATFORM_CONFIG:
        raise ValueError(f"不支持的平台: {platform}，支持的平台: {list(PLATFORM_CONFIG.keys())}")
    
    config = PLATFORM_CONFIG[platform]
    cookie_urls = config["urls"]
    allowed_domains = config["domains"]
    
    print(f"[INFO] 目标平台: {platform}")
    print(f"[INFO] Cookie URLs: {cookie_urls}")
    
    # 需要 Edge 以 --remote-debugging-port=9222 启动
    with sync_playwright() as p:
        browser = p.chromium.connect_over_cdp("http://127.0.0.1:9222")
        if not browser.contexts:
            raise RuntimeError("No browser contexts. Ensure Edge was started with --remote-debugging-port=9222")
        ctx = browser.contexts[0]
        cookies = ctx.cookies(cookie_urls)
        cookie_content = to_netscape(cookies, allowed_domains)

    print(f"[OK] 从浏览器获取Cookie成功，共{len(cookies)}个cookie")
    
    # 直接保存到Redis
    cookie_file = save_cookie_to_redis(cookie_content, platform)
    print(f"[OK] Cookie已保存到Redis，平台: {platform}")
    print(f"[INFO] 临时文件: {cookie_file}")
    print(f"[INFO] Redis Key: ytdl:cookies:{platform}")

if __name__ == "__main__":
    import sys
    
    # 支持命令行指定平台
    platform = "youtube"
    if len(sys.argv) > 1:
        platform = sys.argv[1].lower()
        
    try:
        main(platform)
    except Exception as e:
        print(f"[ERROR] {e}")
        print("[USAGE] python export_edge_cookie.py [youtube|bilibili]")
        sys.exit(1)