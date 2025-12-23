// Package player implementa um player de vídeo com upscaling AI
package player

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/gen2brain/go-mpv"
)

// Player representa o player de vídeo com suporte a upscaling
type Player struct {
	mpv          *mpv.Mpv
	mu           sync.Mutex
	currentMode  PerformanceMode
	windowHandle int64
	isPlaying    bool
	isPaused     bool
	volume       int
	duration     float64
	_            float64 // reserved for position
	shaderPath   string

	// Callbacks para integração com GUI
	OnTimeUpdate  func(position, duration float64)
	OnStateChange func(state string)
	OnError       func(err error)
	OnFileLoaded  func(filename string)
	OnModeChanged func(mode PerformanceMode)
}

// New cria uma nova instância do player
func New() (*Player, error) {
	m := mpv.New()
	if m == nil {
		return nil, fmt.Errorf("falha ao criar instância MPV")
	}

	// Inicializar MPV
	if err := m.Initialize(); err != nil {
		return nil, fmt.Errorf("falha ao inicializar MPV: %w", err)
	}

	// Configurar caminho dos shaders
	execPath, _ := filepath.Abs(".")
	shaderPath := filepath.Join(execPath, "shaders")

	p := &Player{
		mpv:         m,
		currentMode: ModeLow, // Começa no modo mais leve
		volume:      100,
		shaderPath:  shaderPath,
	}

	// Configurações base
	p.setupBaseConfig()

	return p, nil
}

// setupBaseConfig configura opções base do MPV
func (p *Player) setupBaseConfig() {
	// === HABILITAR CONTROLES DE TECLADO ===
	p.mpv.SetPropertyString("input-default-bindings", "yes")
	p.mpv.SetPropertyString("input-vo-keyboard", "yes")
	p.mpv.SetOptionString("input-default-bindings", "yes")
	p.mpv.SetOptionString("input-vo-keyboard", "yes")

	// === OSC - ON SCREEN CONTROLLER ===
	// NOTA: O OSC só funciona se o MPV foi compilado com Lua
	p.mpv.SetPropertyString("osc", "yes")
	p.mpv.SetOptionString("osc", "yes")
	p.mpv.SetPropertyString("load-scripts", "yes")

	// Configurações do OSC
	p.mpv.SetPropertyString("script-opts", "osc-layout=bottombar,osc-seekbarstyle=bar,osc-deadzonesize=0.5,osc-minmousemove=0,osc-hidetimeout=2000,osc-fadeduration=250,osc-showwindowed=yes,osc-showfullscreen=yes,osc-boxalpha=80")

	// Habilitar aceleração de hardware
	p.mpv.SetPropertyString("hwdec", "auto-safe")

	// === CONFIGURAÇÕES DE FPS E SINCRONIZAÇÃO ===
	p.mpv.SetPropertyString("video-sync", "display-resample")
	p.mpv.SetPropertyString("interpolation", "yes")
	p.mpv.SetPropertyString("tscale", "oversample")
	p.mpv.SetPropertyString("framedrop", "no")
	p.mpv.SetPropertyString("opengl-swapinterval", "1")

	// Configurações de áudio
	p.mpv.SetPropertyString("audio-pitch-correction", "yes")
	p.mpv.SetPropertyString("audio-normalize-downmix", "yes")
	p.mpv.SetPropertyString("volume-max", "150") // Permite volume até 150%

	// === JANELA E VISUAL ===
	p.mpv.SetPropertyString("keep-open", "yes")
	p.mpv.SetPropertyString("force-window", "immediate")
	p.mpv.SetPropertyString("border", "no")            // Sem borda da janela (mais limpo)
	p.mpv.SetPropertyString("window-maximized", "yes") // Inicia maximizado

	// Fundo preto quando pausado/sem vídeo
	p.mpv.SetPropertyString("background", "#000000")

	// === OSD CUSTOMIZADO ESTILO ANIME ===
	// Fonte moderna
	p.mpv.SetPropertyString("osd-font", "Segoe UI")
	p.mpv.SetPropertyString("osd-font-size", "36")
	p.mpv.SetPropertyString("osd-bold", "yes")

	// Cores estilo anime (rosa/roxo gradient feel)
	p.mpv.SetPropertyString("osd-color", "#FFFFFFFF")        // Texto branco
	p.mpv.SetPropertyString("osd-border-color", "#FF6B9DFF") // Borda rosa
	p.mpv.SetPropertyString("osd-border-size", "2.5")
	p.mpv.SetPropertyString("osd-shadow-color", "#80000000") // Sombra suave
	p.mpv.SetPropertyString("osd-shadow-offset", "2")
	p.mpv.SetPropertyString("osd-back-color", "#60000000") // Fundo semi-transparente

	// Barra de progresso estilizada
	p.mpv.SetPropertyString("osd-level", "1")
	p.mpv.SetPropertyString("osd-duration", "2500")
	p.mpv.SetPropertyString("osd-bar", "yes")
	p.mpv.SetPropertyString("osd-bar-align-y", "0.95") // Quase no fundo
	p.mpv.SetPropertyString("osd-bar-h", "1.5")        // Fina e elegante
	p.mpv.SetPropertyString("osd-bar-w", "85")         // 85% da largura

	// Mensagens personalizadas
	p.mpv.SetPropertyString("osd-playing-msg", "▶ ${media-title}")
	p.mpv.SetPropertyString("osd-status-msg", "${time-pos} / ${duration}  •  ${percent-pos}%")

	// Margens do OSD
	p.mpv.SetPropertyString("osd-margin-x", "25")
	p.mpv.SetPropertyString("osd-margin-y", "20")

	// === LEGENDAS ESTILIZADAS ===
	p.mpv.SetPropertyString("sub-auto", "fuzzy")
	p.mpv.SetPropertyString("sub-file-paths", "subs:subtitles:Subs:Subtitles:legendas")
	p.mpv.SetPropertyString("sub-font", "Segoe UI Semibold")
	p.mpv.SetPropertyString("sub-font-size", "46")
	p.mpv.SetPropertyString("sub-color", "#FFFFFFFF")
	p.mpv.SetPropertyString("sub-border-color", "#FF000000")
	p.mpv.SetPropertyString("sub-border-size", "2.5")
	p.mpv.SetPropertyString("sub-shadow-color", "#80000000")
	p.mpv.SetPropertyString("sub-shadow-offset", "1")
	p.mpv.SetPropertyString("sub-margin-y", "40")
	p.mpv.SetPropertyString("sub-blur", "0.2") // Leve blur nas bordas

	// === SCREENSHOTS ===
	p.mpv.SetPropertyString("screenshot-format", "png")
	p.mpv.SetPropertyString("screenshot-png-compression", "7")
	p.mpv.SetPropertyString("screenshot-template", "GoAnime_%F_%P")
	p.mpv.SetPropertyString("screenshot-directory", "~~desktop/")

	// === CONTROLES ADICIONAIS ===
	p.mpv.SetPropertyString("input-terminal", "yes")
	p.mpv.SetPropertyString("cursor-autohide", "1500")       // Esconde cursor após 1.5s
	p.mpv.SetPropertyString("cursor-autohide-fs-only", "no") // Esconde mesmo fora de fullscreen
	p.mpv.SetPropertyString("input-cursor", "yes")

	// === VELOCIDADE DE REPRODUÇÃO ===
	p.mpv.SetPropertyString("speed", "1.0")

	// === CACHE PARA STREAMING ===
	p.mpv.SetPropertyString("cache", "yes")
	p.mpv.SetPropertyString("demuxer-max-bytes", "150MiB")
	p.mpv.SetPropertyString("demuxer-max-back-bytes", "75MiB")
	p.mpv.SetPropertyString("demuxer-readahead-secs", "60") // Buffer de 60s

	// Configuração específica por OS
	switch runtime.GOOS {
	case "windows":
		p.mpv.SetPropertyString("vo", "gpu")
		p.mpv.SetPropertyString("gpu-context", "d3d11")
	case "linux":
		p.mpv.SetPropertyString("vo", "gpu")
	case "darwin":
		p.mpv.SetPropertyString("vo", "gpu")
		p.mpv.SetPropertyString("gpu-context", "macvk")
	}
}

// SetTitle define o título da janela do player
func (p *Player) SetTitle(title string) {
	p.mpv.SetPropertyString("title", title)
	p.mpv.SetPropertyString("force-media-title", title)
}

// LoadInputConfig carrega arquivo de configuração de atalhos
func (p *Player) LoadInputConfig(path string) {
	p.mpv.SetPropertyString("input-conf", path)
}

// LoadScript carrega um script Lua
func (p *Player) LoadScript(path string) error {
	return p.mpv.Command([]string{"load-script", path})
}

// SetScriptsDir define o diretório de scripts
func (p *Player) SetScriptsDir(path string) {
	p.mpv.SetPropertyString("scripts", path)
}

// SetFullscreen define se o player deve estar em tela cheia
func (p *Player) SetFullscreen(fs bool) {
	if fs {
		p.mpv.SetPropertyString("fullscreen", "yes")
	} else {
		p.mpv.SetPropertyString("fullscreen", "no")
	}
}

// SetSpeed define a velocidade de reprodução
func (p *Player) SetSpeed(speed float64) {
	p.mpv.SetPropertyString("speed", fmt.Sprintf("%.2f", speed))
}

// GetSpeed retorna a velocidade atual
func (p *Player) GetSpeed() float64 {
	val, err := p.mpv.GetProperty("speed", mpv.FormatDouble)
	if err != nil {
		return 1.0
	}
	if speed, ok := val.(float64); ok {
		return speed
	}
	return 1.0
}

// TakeScreenshot tira uma captura de tela
func (p *Player) TakeScreenshot() {
	p.mpv.Command([]string{"screenshot"})
}

// SetWindowHandle define a janela onde o vídeo será renderizado
func (p *Player) SetWindowHandle(handle int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.windowHandle = handle
	p.mpv.SetProperty("wid", mpv.FormatInt64, handle)
}

// LoadFile carrega um arquivo de vídeo
func (p *Player) LoadFile(path string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	err := p.mpv.Command([]string{"loadfile", path})
	if err != nil {
		return fmt.Errorf("erro ao carregar arquivo: %w", err)
	}

	p.isPlaying = true
	p.isPaused = false

	if p.OnFileLoaded != nil {
		p.OnFileLoaded(path)
	}

	return nil
}

// LoadURL carrega um vídeo de uma URL (streaming)
func (p *Player) LoadURL(url string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Configurar para streaming
	p.mpv.SetPropertyString("stream-lavf-o", "reconnect=1,reconnect_streamed=1,reconnect_delay_max=5")

	err := p.mpv.Command([]string{"loadfile", url})
	if err != nil {
		return fmt.Errorf("erro ao carregar URL: %w", err)
	}

	p.isPlaying = true
	p.isPaused = false

	return nil
}

// Play inicia ou retoma a reprodução
func (p *Player) Play() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.mpv.SetPropertyString("pause", "no")
	p.isPaused = false
	p.isPlaying = true

	if p.OnStateChange != nil {
		p.OnStateChange("playing")
	}
}

// Pause pausa a reprodução
func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.mpv.SetPropertyString("pause", "yes")
	p.isPaused = true

	if p.OnStateChange != nil {
		p.OnStateChange("paused")
	}
}

// TogglePause alterna entre play/pause
func (p *Player) TogglePause() {
	if p.isPaused {
		p.Play()
	} else {
		p.Pause()
	}
}

// Stop para a reprodução
func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.mpv.Command([]string{"stop"})
	p.isPlaying = false
	p.isPaused = false

	if p.OnStateChange != nil {
		p.OnStateChange("stopped")
	}
}

// Seek vai para uma posição específica (em segundos)
func (p *Player) Seek(position float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.mpv.Command([]string{"seek", fmt.Sprintf("%f", position), "absolute"})
}

// SeekRelative avança ou retrocede (em segundos)
func (p *Player) SeekRelative(seconds float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.mpv.Command([]string{"seek", fmt.Sprintf("%f", seconds), "relative"})
}

// SetVolume define o volume (0-100)
func (p *Player) SetVolume(volume int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if volume < 0 {
		volume = 0
	}
	if volume > 100 {
		volume = 100
	}

	p.volume = volume
	p.mpv.SetProperty("volume", mpv.FormatInt64, int64(volume))
}

// GetVolume retorna o volume atual
func (p *Player) GetVolume() int {
	return p.volume
}

// ToggleMute alterna mudo
func (p *Player) ToggleMute() {
	p.mpv.Command([]string{"cycle", "mute"})
}

// ToggleFullscreen alterna tela cheia
func (p *Player) ToggleFullscreen() {
	p.mpv.Command([]string{"cycle", "fullscreen"})
}

// SetSubtitleTrack define a trilha de legenda
func (p *Player) SetSubtitleTrack(id int) {
	p.mpv.SetProperty("sid", mpv.FormatInt64, int64(id))
}

// SetAudioTrack define a trilha de áudio
func (p *Player) SetAudioTrack(id int) {
	p.mpv.SetProperty("aid", mpv.FormatInt64, int64(id))
}

// LoadSubtitle carrega um arquivo de legenda externo
func (p *Player) LoadSubtitle(path string) error {
	return p.mpv.Command([]string{"sub-add", path})
}

// GetPosition retorna a posição atual em segundos
func (p *Player) GetPosition() float64 {
	val, err := p.mpv.GetProperty("time-pos", mpv.FormatDouble)
	if err != nil {
		return 0
	}
	if pos, ok := val.(float64); ok {
		return pos
	}
	return 0
}

// GetDuration retorna a duração total em segundos
func (p *Player) GetDuration() float64 {
	val, err := p.mpv.GetProperty("duration", mpv.FormatDouble)
	if err != nil {
		return 0
	}
	if dur, ok := val.(float64); ok {
		return dur
	}
	return 0
}

// GetDroppedFrames retorna o número de frames perdidos
func (p *Player) GetDroppedFrames() int64 {
	val, err := p.mpv.GetProperty("frame-drop-count", mpv.FormatInt64)
	if err != nil {
		return 0
	}
	if frames, ok := val.(int64); ok {
		return frames
	}
	return 0
}

// IsPlaying retorna se está reproduzindo
func (p *Player) IsPlaying() bool {
	return p.isPlaying && !p.isPaused
}

// IsPaused retorna se está pausado
func (p *Player) IsPaused() bool {
	return p.isPaused
}

// Run executa o loop de eventos do player
func (p *Player) Run() {
	for {
		event := p.mpv.WaitEvent(1)
		if event == nil {
			continue
		}

		switch event.EventID {
		case mpv.EventFileLoaded:
			p.duration = p.GetDuration()
			fmt.Printf("📄 Arquivo carregado. Duração: %.2f segundos\n", p.duration)

		case 7: // EventEndFile
			fmt.Println("🏁 Fim do arquivo")
			if p.OnStateChange != nil {
				p.OnStateChange("ended")
			}

		case mpv.EventShutdown:
			fmt.Println("👋 Player encerrado")
			return

		case mpv.EventPropertyChange:
			// Monitorar mudanças de propriedades
			p.handlePropertyChange(event)
		}
	}
}

// handlePropertyChange processa mudanças de propriedades
func (p *Player) handlePropertyChange(_ *mpv.Event) {
	// Atualizar posição periodicamente
	if p.OnTimeUpdate != nil && p.isPlaying {
		pos := p.GetPosition()
		dur := p.GetDuration()
		p.OnTimeUpdate(pos, dur)
	}

	// Verificar frames perdidos (para auto-downgrade de modo)
	droppedFrames := p.GetDroppedFrames()
	if droppedFrames > 30 && p.currentMode == ModeHigh {
		fmt.Println("⚠️ Muitos frames perdidos! Considere baixar o modo de qualidade.")
	}
}

// Destroy libera os recursos do player
func (p *Player) Destroy() {
	if p.mpv != nil {
		p.mpv.TerminateDestroy()
	}
}
