param(
  [Parameter(Mandatory = $true)]
  [ValidateSet('youtube','bilibili')]
  [string]$Platform,

  [switch]$Loop,

  [int]$IntervalMinutes = 5
)

$ErrorActionPreference = 'Stop'

function Start-EdgeIfNeeded {
  param(
    [int]$Port,
    [string]$UserDataDir
  )

  try {
    $listening = Get-NetTCPConnection -State Listen -LocalPort $Port -ErrorAction SilentlyContinue
  } catch {
    $listening = $null
  }

  if (-not $listening) {
    $edgePath = "${env:ProgramFiles(x86)}\Microsoft\Edge\Application\msedge.exe"
    if (-not (Test-Path $edgePath)) {
      $edgePath = "$env:ProgramFiles\Microsoft\Edge\Application\msedge.exe"
    }
    if (-not (Test-Path $edgePath)) {
      throw ("Microsoft Edge not found. Please install Edge or adjust path: " + $edgePath)
    }
    $args = "--remote-debugging-port=$Port --user-data-dir=$UserDataDir"
    Start-Process -FilePath $edgePath -ArgumentList $args | Out-Null
    Start-Sleep -Seconds 2
  }
}

function Invoke-ExportOnce {
  param(
    [string]$Platform,
    [int]$Port,
    [string]$UserDataDir
  )

  Start-EdgeIfNeeded -Port $Port -UserDataDir $UserDataDir

  $python = "python" # 如需指定绝对路径，可改为 C:\\Python312\\python.exe 或使用 `py -3`
  $scriptPath = Join-Path $PSScriptRoot "export_edge_cookie.py"
  if (-not (Test-Path $scriptPath)) {
    throw "找不到 Python 脚本：$scriptPath"
  }

  Write-Host ("[{0}] Exporting: {1}" -f (Get-Date -Format 'yyyy-MM-dd HH:mm:ss'), $Platform) -ForegroundColor Cyan
  & $python $scriptPath $Platform
}

$port = 9222
$userDataDir = Join-Path $env:LOCALAPPDATA "EdgeCookieExport"

if ($Loop) {
  $mutexName = "Global/ExportCookies-$Platform"
  $createdNew = $false
  $mutex = New-Object System.Threading.Mutex($true, $mutexName, [ref]$createdNew)
  if (-not $createdNew) {
    Write-Host ("Detected a running loop instance ({0}); current process will exit." -f $Platform) -ForegroundColor Yellow
    exit 0
  }
  try {
    while ($true) {
      try {
        Invoke-ExportOnce -Platform $Platform -Port $port -UserDataDir $userDataDir
      } catch {
        Write-Host ("Export error: {0}" -f $_.Exception.Message) -ForegroundColor Red
      }
      Start-Sleep -Seconds ($IntervalMinutes * 60)
    }
  } finally {
    if ($mutex) { $mutex.ReleaseMutex(); $mutex.Dispose() }
  }
} else {
  try {
    Invoke-ExportOnce -Platform $Platform -Port $port -UserDataDir $userDataDir
  } catch {
    Write-Host ("Export error: {0}" -f $_.Exception.Message) -ForegroundColor Red
    exit 1
  }
}


