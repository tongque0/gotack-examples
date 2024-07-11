# 获取当前脚本所在目录
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

# 定义目标文件夹
$targetDir = Join-Path -Path $scriptDir -ChildPath "touch"

# 检查目标文件夹是否存在，不存在则创建
if (-Not (Test-Path -Path $targetDir)) {
    New-Item -ItemType Directory -Path $targetDir
}

# 获取当前文件夹内的文件数量
$fileCount = (Get-ChildItem -Path $targetDir -File).Count

# 定义新文件名
$newFileName = "$($fileCount + 1).txt"

# 定义新文件的完整路径
$newFilePath = Join-Path -Path $targetDir -ChildPath $newFileName

# 创建新文件
New-Item -Path $newFilePath -ItemType File

# 输出创建的文件路径
Write-Output "Created file: $newFilePath"

# 执行 git add .
git add .

# 检查 git 命令是否成功
if ($LASTEXITCODE -eq 0) {
    Write-Output "Successfully added files to Git."
} else {
    Write-Output "Failed to add files to Git."
}
