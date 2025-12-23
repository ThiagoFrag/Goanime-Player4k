package player

import (
	"fmt"
	"path/filepath"
)

// PerformanceMode representa os modos de performance
type PerformanceMode string

const (
	// ModeLow - Modo Econômico (notebooks, GPUs integradas)
	ModeLow PerformanceMode = "low"

	// ModeMedium - Modo Equilibrado (GPUs entrada/médias)
	ModeMedium PerformanceMode = "medium"

	// ModeHigh - Modo Ultra (GPUs dedicadas potentes)
	ModeHigh PerformanceMode = "high"
)

// ModeInfo contém informações sobre um modo
type ModeInfo struct {
	ID          PerformanceMode
	Name        string
	Description string
	Icon        string
	GPURequired string
}

// GetModeInfo retorna informações sobre um modo
func GetModeInfo(mode PerformanceMode) ModeInfo {
	switch mode {
	case ModeLow:
		return ModeInfo{
			ID:          ModeLow,
			Name:        "Econômico",
			Description: "Otimizado para bateria e compatibilidade",
			Icon:        "🔋",
			GPURequired: "Qualquer (Intel HD, AMD APU)",
		}
	case ModeMedium:
		return ModeInfo{
			ID:          ModeMedium,
			Name:        "Equilibrado",
			Description: "Qualidade boa com upscaling FSR",
			Icon:        "⚖️",
			GPURequired: "GTX 1050 / RX 560 / Intel Iris",
		}
	case ModeHigh:
		return ModeInfo{
			ID:          ModeHigh,
			Name:        "Ultra",
			Description: "Upscaling AI com rede neural profunda",
			Icon:        "🚀",
			GPURequired: "RTX 3060 / RX 6700 ou superior",
		}
	}
	return ModeInfo{}
}

// GetAllModes retorna todos os modos disponíveis
func GetAllModes() []ModeInfo {
	return []ModeInfo{
		GetModeInfo(ModeLow),
		GetModeInfo(ModeMedium),
		GetModeInfo(ModeHigh),
	}
}

// SetPerformanceMode aplica um modo de performance
func (p *Player) SetPerformanceMode(mode PerformanceMode) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Limpar shaders anteriores
	p.mpv.SetPropertyString("glsl-shaders", "")

	info := GetModeInfo(mode)
	fmt.Printf("%s Ativando modo: %s\n", info.Icon, info.Name)

	switch mode {
	case ModeLow:
		p.applyLowMode()
	case ModeMedium:
		p.applyMediumMode()
	case ModeHigh:
		p.applyHighMode()
	}

	p.currentMode = mode

	if p.OnModeChanged != nil {
		p.OnModeChanged(mode)
	}
}

// applyLowMode aplica configurações do modo econômico
func (p *Player) applyLowMode() {
	// Profile leve
	p.mpv.SetPropertyString("profile", "fast")

	// Hardware decoding prioritário
	p.mpv.SetPropertyString("hwdec", "auto-safe")

	// Escaladores mais leves
	p.mpv.SetPropertyString("scale", "bilinear")
	p.mpv.SetPropertyString("cscale", "bilinear")
	p.mpv.SetPropertyString("dscale", "bilinear")

	// Desativar recursos pesados
	p.mpv.SetPropertyString("deband", "no")
	p.mpv.SetPropertyString("interpolation", "no")
	p.mpv.SetPropertyString("dither-depth", "no")

	// Renderizador padrão
	p.mpv.SetPropertyString("vo", "gpu")

	fmt.Println("  ✓ Decodificação por hardware")
	fmt.Println("  ✓ Escalamento bilinear (leve)")
	fmt.Println("  ✓ Debanding desativado")
}

// applyMediumMode aplica configurações do modo equilibrado
func (p *Player) applyMediumMode() {
	// Profile de alta qualidade
	p.mpv.SetPropertyString("profile", "gpu-hq")

	// Hardware decoding
	p.mpv.SetPropertyString("hwdec", "auto-safe")

	// Escaladores melhores (nativos, sem shader externo pesado)
	p.mpv.SetPropertyString("scale", "spline36")
	p.mpv.SetPropertyString("cscale", "spline36")
	p.mpv.SetPropertyString("dscale", "mitchell")

	// Debanding leve
	p.mpv.SetPropertyString("deband", "yes")
	p.mpv.SetPropertyString("deband-iterations", "2")
	p.mpv.SetPropertyString("deband-threshold", "35")
	p.mpv.SetPropertyString("deband-range", "20")

	// Dithering
	p.mpv.SetPropertyString("dither-depth", "auto")

	// Carregar shader FSR (AMD FidelityFX Super Resolution)
	fsrPath := filepath.Join(p.shaderPath, "FSR.glsl")
	err := p.mpv.Command([]string{"change-list", "glsl-shaders", "append", fsrPath})
	if err != nil {
		fmt.Printf("  ⚠️ Shader FSR não encontrado: %s\n", fsrPath)
	} else {
		fmt.Println("  ✓ AMD FSR ativado (upscaling eficiente)")
	}

	fmt.Println("  ✓ Profile gpu-hq")
	fmt.Println("  ✓ Escalamento spline36")
	fmt.Println("  ✓ Debanding leve")
}

// applyHighMode aplica configurações do modo ultra
func (p *Player) applyHighMode() {
	// Backend moderno (Vulkan se disponível)
	p.mpv.SetPropertyString("vo", "gpu-next")
	p.mpv.SetPropertyString("profile", "gpu-hq")

	// Hardware decoding com copy-back para processamento
	p.mpv.SetPropertyString("hwdec", "auto-copy")

	// Escaladores de alta qualidade
	p.mpv.SetPropertyString("scale", "ewa_lanczossharp")
	p.mpv.SetPropertyString("cscale", "ewa_lanczossharp")
	p.mpv.SetPropertyString("dscale", "mitchell")

	// Debanding agressivo
	p.mpv.SetPropertyString("deband", "yes")
	p.mpv.SetPropertyString("deband-iterations", "4")
	p.mpv.SetPropertyString("deband-threshold", "48")
	p.mpv.SetPropertyString("deband-range", "24")
	p.mpv.SetPropertyString("deband-grain", "24")

	// Dithering de alta qualidade
	p.mpv.SetPropertyString("dither-depth", "auto")
	p.mpv.SetPropertyString("temporal-dither", "yes")

	// HDR tone mapping (se disponível)
	p.mpv.SetPropertyString("tone-mapping", "bt.2446a")
	p.mpv.SetPropertyString("tone-mapping-mode", "auto")

	// Carregar shader FSRCNNX (Rede Neural)
	fsrcnnxPath := filepath.Join(p.shaderPath, "FSRCNNX_x2_16-0-4-1.glsl")
	err := p.mpv.Command([]string{"change-list", "glsl-shaders", "append", fsrcnnxPath})
	if err != nil {
		fmt.Printf("  ⚠️ Shader FSRCNNX não encontrado: %s\n", fsrcnnxPath)
		// Fallback para Anime4K se FSRCNNX não disponível
		anime4kPath := filepath.Join(p.shaderPath, "Anime4K_Upscale_CNN_x2_VL.glsl")
		p.mpv.Command([]string{"change-list", "glsl-shaders", "append", anime4kPath})
	} else {
		fmt.Println("  ✓ FSRCNNX Neural Network ativado")
	}

	// Opcional: Adicionar sharpening
	casPath := filepath.Join(p.shaderPath, "CAS.glsl")
	p.mpv.Command([]string{"change-list", "glsl-shaders", "append", casPath})

	fmt.Println("  ✓ Backend gpu-next (Vulkan)")
	fmt.Println("  ✓ Upscaling por Rede Neural")
	fmt.Println("  ✓ Debanding agressivo")
	fmt.Println("  ✓ HDR tone mapping")
}

// GetCurrentMode retorna o modo atual
func (p *Player) GetCurrentMode() PerformanceMode {
	return p.currentMode
}

// AutoSelectMode seleciona automaticamente o modo baseado na GPU
func (p *Player) AutoSelectMode() PerformanceMode {
	// Por enquanto retorna Medium como padrão seguro
	// TODO: Detectar GPU e selecionar automaticamente
	return ModeMedium
}

// EnableInterpolation ativa interpolação de movimento (motion smoothing)
// Cria frames intermediários para deixar vídeo mais fluido
// AVISO: Requer GPU potente e nem todos gostam do efeito
func (p *Player) EnableInterpolation(enable bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if enable {
		p.mpv.SetPropertyString("interpolation", "yes")
		p.mpv.SetPropertyString("tscale", "oversample")
		p.mpv.SetPropertyString("video-sync", "display-resample")
		fmt.Println("✓ Interpolação de movimento ativada")
	} else {
		p.mpv.SetPropertyString("interpolation", "no")
		p.mpv.SetPropertyString("video-sync", "audio")
		fmt.Println("✓ Interpolação de movimento desativada")
	}
}

// SetAnimeMode ativa otimizações específicas para anime
func (p *Player) SetAnimeMode(enable bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if enable {
		// Limpar shaders anteriores
		p.mpv.SetPropertyString("glsl-shaders", "")

		// Carregar shaders Anime4K
		shaders := []string{
			"Anime4K_Clamp_Highlights.glsl",
			"Anime4K_Restore_CNN_VL.glsl",
			"Anime4K_Upscale_CNN_x2_VL.glsl",
			"Anime4K_AutoDownscalePre_x2.glsl",
			"Anime4K_AutoDownscalePre_x4.glsl",
			"Anime4K_Upscale_CNN_x2_M.glsl",
		}

		for _, shader := range shaders {
			shaderPath := filepath.Join(p.shaderPath, "Anime4K", shader)
			p.mpv.Command([]string{"change-list", "glsl-shaders", "append", shaderPath})
		}

		fmt.Println("🎌 Modo Anime ativado (Anime4K)")
	} else {
		// Voltar ao modo atual
		p.SetPerformanceMode(p.currentMode)
	}
}
