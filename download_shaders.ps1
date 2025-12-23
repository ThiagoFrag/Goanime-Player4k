# Script para baixar shaders necessários para o Player4K
# Execute: .\download_shaders.ps1

$shadersDir = ".\shaders"
$anime4kDir = "$shadersDir\Anime4K"

# Criar diretórios
New-Item -ItemType Directory -Force -Path $shadersDir | Out-Null
New-Item -ItemType Directory -Force -Path $anime4kDir | Out-Null

Write-Host "📥 Baixando shaders para Player4K..." -ForegroundColor Cyan

# 1. FSR (AMD FidelityFX Super Resolution)
Write-Host "`n[1/3] Baixando FSR.glsl..." -ForegroundColor Yellow
$fsrUrl = "https://gist.githubusercontent.com/agyild/82219c545228d70c5604f865ce0b0ce5/raw/FSR.glsl"
$downloaded = $false
try {
    Invoke-WebRequest -Uri $fsrUrl -OutFile "$shadersDir\FSR.glsl" -UseBasicParsing -ErrorAction Stop
    $downloaded = $true
}
catch {
    $downloaded = $false
}
if ($downloaded) {
    Write-Host "  ✓ FSR.glsl baixado" -ForegroundColor Green
} else {
    Write-Host "  ✗ Erro ao baixar FSR.glsl - baixe manualmente" -ForegroundColor Red
}

# 2. CAS (Contrast Adaptive Sharpening)
Write-Host "`n[2/3] Baixando CAS.glsl..." -ForegroundColor Yellow
$casUrl = "https://gist.githubusercontent.com/agyild/bbb4e58298b2f86aa24da3032a0d2f7f/raw/CAS.glsl"
$downloaded = $false
try {
    Invoke-WebRequest -Uri $casUrl -OutFile "$shadersDir\CAS.glsl" -UseBasicParsing -ErrorAction Stop
    $downloaded = $true
}
catch {
    $downloaded = $false
}
if ($downloaded) {
    Write-Host "  ✓ CAS.glsl baixado" -ForegroundColor Green
} else {
    Write-Host "  ✗ Erro ao baixar CAS.glsl - baixe manualmente" -ForegroundColor Red
}

# 3. Anime4K (último release)
Write-Host "`n[3/3] Baixando Anime4K..." -ForegroundColor Yellow
$anime4kRelease = "https://github.com/bloc97/Anime4K/releases/download/v4.0.1/Anime4K_v4.0.zip"
$anime4kZip = "$env:TEMP\Anime4K.zip"
$downloadedAnime = $false

try {
    Write-Host "  Baixando Anime4K v4.0.1..." -ForegroundColor Gray
    Invoke-WebRequest -Uri $anime4kRelease -OutFile $anime4kZip -UseBasicParsing -ErrorAction Stop
    
    Write-Host "  Extraindo..." -ForegroundColor Gray
    Expand-Archive -Path $anime4kZip -DestinationPath $anime4kDir -Force
    
    # Mover arquivos da subpasta se necessário
    $subFolder = Get-ChildItem -Path $anime4kDir -Directory | Select-Object -First 1
    if ($subFolder) {
        Get-ChildItem -Path $subFolder.FullName -Filter "*.glsl" | Move-Item -Destination $anime4kDir -Force
    }
    
    Remove-Item $anime4kZip -Force
    $downloadedAnime = $true
}
catch {
    $downloadedAnime = $false
}

if ($downloadedAnime) {
    Write-Host "  ✓ Anime4K instalado" -ForegroundColor Green
} else {
    Write-Host "  ✗ Erro ao baixar Anime4K - baixe manualmente de:" -ForegroundColor Red
    Write-Host "    https://github.com/bloc97/Anime4K/releases" -ForegroundColor Gray
}

# FSRCNNX precisa ser baixado manualmente (arquivo grande)
Write-Host "`n⚠️  FSRCNNX (modo High) precisa ser baixado manualmente:" -ForegroundColor Yellow
Write-Host "   1. Acesse: https://github.com/igv/FSRCNN-TensorFlow/releases" -ForegroundColor Gray
Write-Host "   2. Baixe: FSRCNNX_x2_16-0-4-1.glsl" -ForegroundColor Gray
Write-Host "   3. Coloque em: $shadersDir\" -ForegroundColor Gray

Write-Host "`n✅ Instalação de shaders concluída!" -ForegroundColor Green
Write-Host "`nEstrutura criada:" -ForegroundColor Cyan
Get-ChildItem -Path $shadersDir -Recurse | ForEach-Object {
    $indent = "  " * ($_.FullName.Split('\').Count - $shadersDir.Split('\').Count)
    Write-Host "$indent$($_.Name)"
}
