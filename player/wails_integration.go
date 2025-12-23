// Package player - Integração com Wails para GoAnimeGUI
package player

import (
	"fmt"
)

// WailsPlayer é o wrapper do player para uso com Wails
// Expõe métodos que podem ser chamados do frontend JavaScript/Svelte
type WailsPlayer struct {
	player *Player
}

// NewWailsPlayer cria um player para integração com Wails
func NewWailsPlayer() (*WailsPlayer, error) {
	p, err := New()
	if err != nil {
		return nil, err
	}

	return &WailsPlayer{player: p}, nil
}

// --- Métodos expostos para o Frontend (Wails) ---

// Initialize inicializa o player com um handle de janela
func (w *WailsPlayer) Initialize(windowHandle int64) error {
	w.player.SetWindowHandle(windowHandle)
	return nil
}

// Load carrega um arquivo ou URL
func (w *WailsPlayer) Load(path string) error {
	// Detecta se é URL ou arquivo local
	if len(path) > 4 && (path[:4] == "http" || path[:4] == "rtmp") {
		return w.player.LoadURL(path)
	}
	return w.player.LoadFile(path)
}

// Play inicia reprodução
func (w *WailsPlayer) Play() {
	w.player.Play()
}

// Pause pausa reprodução
func (w *WailsPlayer) Pause() {
	w.player.Pause()
}

// TogglePlay alterna entre play/pause
func (w *WailsPlayer) TogglePlay() {
	w.player.TogglePause()
}

// Stop para reprodução
func (w *WailsPlayer) Stop() {
	w.player.Stop()
}

// Seek vai para posição em segundos
func (w *WailsPlayer) Seek(seconds float64) {
	w.player.Seek(seconds)
}

// SeekForward avança 10 segundos
func (w *WailsPlayer) SeekForward() {
	w.player.SeekRelative(10)
}

// SeekBackward retrocede 10 segundos
func (w *WailsPlayer) SeekBackward() {
	w.player.SeekRelative(-10)
}

// SetVolume define volume (0-100)
func (w *WailsPlayer) SetVolume(volume int) {
	w.player.SetVolume(volume)
}

// GetVolume retorna volume atual
func (w *WailsPlayer) GetVolume() int {
	return w.player.GetVolume()
}

// ToggleMute alterna mudo
func (w *WailsPlayer) ToggleMute() {
	w.player.ToggleMute()
}

// ToggleFullscreen alterna tela cheia
func (w *WailsPlayer) ToggleFullscreen() {
	w.player.ToggleFullscreen()
}

// GetPosition retorna posição atual em segundos
func (w *WailsPlayer) GetPosition() float64 {
	return w.player.GetPosition()
}

// GetDuration retorna duração total em segundos
func (w *WailsPlayer) GetDuration() float64 {
	return w.player.GetDuration()
}

// GetProgress retorna progresso como porcentagem (0-100)
func (w *WailsPlayer) GetProgress() float64 {
	duration := w.player.GetDuration()
	if duration <= 0 {
		return 0
	}
	return (w.player.GetPosition() / duration) * 100
}

// IsPlaying retorna se está reproduzindo
func (w *WailsPlayer) IsPlaying() bool {
	return w.player.IsPlaying()
}

// IsPaused retorna se está pausado
func (w *WailsPlayer) IsPaused() bool {
	return w.player.IsPaused()
}

// --- Configurações de Qualidade ---

// SetQualityMode define o modo de qualidade
// mode: "low", "medium", "high"
func (w *WailsPlayer) SetQualityMode(mode string) {
	switch mode {
	case "low":
		w.player.SetPerformanceMode(ModeLow)
	case "medium":
		w.player.SetPerformanceMode(ModeMedium)
	case "high":
		w.player.SetPerformanceMode(ModeHigh)
	default:
		w.player.SetPerformanceMode(ModeMedium)
	}
}

// GetQualityMode retorna o modo de qualidade atual
func (w *WailsPlayer) GetQualityMode() string {
	return string(w.player.GetCurrentMode())
}

// GetQualityModes retorna todos os modos disponíveis
func (w *WailsPlayer) GetQualityModes() []map[string]string {
	modes := GetAllModes()
	result := make([]map[string]string, len(modes))

	for i, m := range modes {
		result[i] = map[string]string{
			"id":          string(m.ID),
			"name":        m.Name,
			"description": m.Description,
			"icon":        m.Icon,
			"gpuRequired": m.GPURequired,
		}
	}

	return result
}

// SetAnimeMode ativa/desativa otimizações para anime
func (w *WailsPlayer) SetAnimeMode(enable bool) {
	w.player.SetAnimeMode(enable)
}

// EnableMotionSmoothing ativa/desativa interpolação de movimento
func (w *WailsPlayer) EnableMotionSmoothing(enable bool) {
	w.player.EnableInterpolation(enable)
}

// --- Legendas e Áudio ---

// SetSubtitle define trilha de legenda por ID
func (w *WailsPlayer) SetSubtitle(id int) {
	w.player.SetSubtitleTrack(id)
}

// SetAudio define trilha de áudio por ID
func (w *WailsPlayer) SetAudio(id int) {
	w.player.SetAudioTrack(id)
}

// LoadExternalSubtitle carrega legenda externa
func (w *WailsPlayer) LoadExternalSubtitle(path string) error {
	return w.player.LoadSubtitle(path)
}

// --- Diagnóstico ---

// GetDroppedFrames retorna frames perdidos (para debug)
func (w *WailsPlayer) GetDroppedFrames() int64 {
	return w.player.GetDroppedFrames()
}

// GetStats retorna estatísticas do player
func (w *WailsPlayer) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"position":      w.player.GetPosition(),
		"duration":      w.player.GetDuration(),
		"droppedFrames": w.player.GetDroppedFrames(),
		"mode":          string(w.player.GetCurrentMode()),
		"isPlaying":     w.player.IsPlaying(),
		"isPaused":      w.player.IsPaused(),
		"volume":        w.player.GetVolume(),
	}
}

// Destroy libera recursos
func (w *WailsPlayer) Destroy() {
	if w.player != nil {
		w.player.Destroy()
	}
}

// --- Callbacks para eventos (para uso interno) ---

// SetOnTimeUpdate define callback para atualização de tempo
func (w *WailsPlayer) SetOnTimeUpdate(callback func(position, duration float64)) {
	w.player.OnTimeUpdate = callback
}

// SetOnStateChange define callback para mudança de estado
func (w *WailsPlayer) SetOnStateChange(callback func(state string)) {
	w.player.OnStateChange = callback
}

// SetOnError define callback para erros
func (w *WailsPlayer) SetOnError(callback func(err error)) {
	w.player.OnError = callback
}

// PrintInfo imprime informações do player (debug)
func (w *WailsPlayer) PrintInfo() {
	fmt.Println("=== Player4K Info ===")
	fmt.Printf("Modo atual: %s\n", w.GetQualityMode())
	fmt.Printf("Posição: %.2f / %.2f\n", w.GetPosition(), w.GetDuration())
	fmt.Printf("Volume: %d%%\n", w.GetVolume())
	fmt.Printf("Frames perdidos: %d\n", w.GetDroppedFrames())
	fmt.Println("====================")
}
