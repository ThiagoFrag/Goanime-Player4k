# 🎬 Goanime-Player4k

Player de vídeo modificado com upscaling AI para o sistema GoAnime. Este player é usado pelo GoAnimeGUI para reproduzir vídeos com qualidade aprimorada.

## 📋 Sobre

Este é um fork modificado do MPV com shaders de upscaling AI integrados. Ele é otimizado para reprodução de anime com melhoria de qualidade em tempo real.

## ✨ Características

- 🎮 **3 Modos de Performance**
  - 🔋 **Econômico**: Para notebooks e GPUs integradas
  - ⚖️ **Equilibrado**: Upscaling FSR para GPUs de entrada
  - 🚀 **Ultra**: Rede neural FSRCNNX para GPUs potentes

- 🎌 **Modo Anime**: Otimizações Anime4K específicas para animação
- 📺 **HDR Support**: Tone mapping automático
- 🔊 **Múltiplas trilhas**: Áudio e legendas
- 🌐 **Streaming**: Suporte a URLs HTTP/HTTPS
- 🔗 **Integração**: Comunicação via socket com GoAnimeGUI

## 📁 Estrutura

```
player4k/
├── main.go              # Código principal do player
├── player/              # Pacote de controle do MPV
├── mpv/                 # Bindings Go para MPV
├── shaders/             # Shaders de upscaling AI
│   ├── Anime4K/         # Shaders otimizados para anime
│   ├── FSR/             # AMD FidelityFX Super Resolution
│   └── FSRCNNX/         # Rede neural para upscaling
├── portable_config/     # Configurações padrão do MPV
├── scripts/             # Scripts Lua para funcionalidades extras
└── input.conf           # Keybindings personalizados
```

## 🔧 Requisitos

### Sistema
- Windows 10/11, Linux, ou macOS
- Go 1.21+
- GCC (para CGO)

### Dependências
- libmpv (MPV library)

#### Windows
```powershell
# Baixar libmpv de: https://sourceforge.net/projects/mpv-player-windows/files/libmpv/
# Extrair mpv-dev.7z e colocar libmpv-2.dll na pasta do projeto
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt install libmpv-dev
```

#### Linux (Arch/Manjaro)
```bash
sudo pacman -S mpv
```

## 🚀 Instalação

```bash
cd player4k
go mod tidy
go build
```

## Uso Standalone

```bash
# Reproduzir arquivo local
./player4k video.mp4

# Reproduzir URL
./player4k "https://example.com/video.m3u8"
```

## Integração com GoAnimeGUI

```go
import "player4k/player"

// Criar player
p, _ := player.NewWailsPlayer()
defer p.Destroy()

// Definir janela para renderização
p.Initialize(windowHandle)

// Carregar e reproduzir
p.Load("video.mp4")
p.Play()

// Mudar qualidade
p.SetQualityMode("high") // "low", "medium", "high"

// Ativar modo anime
p.SetAnimeMode(true)
```

## Shaders

Baixe os shaders necessários e coloque na pasta `shaders/`:

1. **FSR.glsl** - AMD FidelityFX (Modo Medium)
2. **FSRCNNX_x2_16-0-4-1.glsl** - Neural Network (Modo High)
3. **Anime4K/** - Shaders para anime

Veja `shaders/README.md` para links de download.

## API

### Controles Básicos
- `Load(path)` - Carregar vídeo
- `Play()` / `Pause()` / `Stop()`
- `Seek(seconds)` - Ir para posição
- `SetVolume(0-100)` - Volume

### Qualidade
- `SetQualityMode("low"|"medium"|"high")`
- `SetAnimeMode(bool)` - Otimizações para anime
- `EnableMotionSmoothing(bool)` - Interpolação de frames

### Informações
- `GetPosition()` / `GetDuration()`
- `GetProgress()` - Porcentagem
- `GetStats()` - Estatísticas completas
- `GetDroppedFrames()` - Frames perdidos

## Modos de Qualidade

| Modo | Escalador | Debanding | GPU Recomendada |
|------|-----------|-----------|-----------------|
| Low | Bilinear | Não | Qualquer |
| Medium | Spline36 + FSR | Leve | GTX 1050+ |
| High | FSRCNNX Neural | Agressivo | RTX 3060+ |

## Troubleshooting

### Vídeo engasgando
- Reduza o modo de qualidade
- Verifique se `hwdec` está funcionando
- Monitore `GetDroppedFrames()`

### Shader não carrega
- Verifique se o arquivo .glsl existe em `shaders/`
- Confirme que a GPU suporta GLSL 3.30+

### Sem aceleração de hardware
- Instale drivers atualizados da GPU
- No Windows, instale LAV Filters

## Licença

MIT
