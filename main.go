// Player4K - Player de vídeo com upscaling AI para GoAnimeGUI
// Usa MPV como backend com shaders GLSL para upscaling de alta qualidade
//go:build windows
// +build windows

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"player4k/player"
)

func main() {
	// Flags de linha de comando
	modeFlag := flag.String("mode", "medium", "Modo de qualidade: low, medium, high")
	animeFlag := flag.Bool("anime", false, "Ativar modo otimizado para anime (Anime4K)")
	titleFlag := flag.String("title", "", "Título para exibir na janela")
	subFlag := flag.String("sub", "", "URL ou caminho de legenda externa")
	listModes := flag.Bool("list-modes", false, "Listar todos os modos disponíveis")
	fullscreen := flag.Bool("fs", false, "Iniciar em tela cheia")
	volume := flag.Int("volume", 100, "Volume inicial (0-150)")
	startPos := flag.Float64("start", 0, "Posição inicial em segundos")
	flag.Parse()

	if *listModes {
		printBanner()
		fmt.Println("\n🎬 Modos de Qualidade Disponíveis:")
		fmt.Println("─────────────────────────────────────")
		for _, mode := range player.GetAllModes() {
			fmt.Printf("\n  %s %s (%s)\n", mode.Icon, mode.Name, mode.ID)
			fmt.Printf("     📝 %s\n", mode.Description)
			fmt.Printf("     🎮 GPU: %s\n", mode.GPURequired)
		}
		fmt.Println("\n─────────────────────────────────────")
		printControls()
		return
	}

	// Criar instância do player
	p, err := player.New()
	if err != nil {
		os.Exit(1)
	}
	defer p.Destroy()

	// Carregar arquivos de configuração
	execPath, _ := os.Executable()
	execDir := filepath.Dir(execPath)

	// Carregar atalhos customizados (input.conf)
	inputConf := filepath.Join(execDir, "input.conf")
	if _, err := os.Stat(inputConf); err == nil {
		p.LoadInputConfig(inputConf)
	}

	// Carregar script OSC (barra de controles na tela)
	oscScript := filepath.Join(execDir, "scripts", "osc.lua")
	if _, err := os.Stat(oscScript); err == nil {
		p.LoadScript(oscScript)
	}

	// Configurar modo de qualidade
	var mode player.PerformanceMode
	switch *modeFlag {
	case "low":
		mode = player.ModeLow
	case "high":
		mode = player.ModeHigh
	default:
		mode = player.ModeMedium
	}
	p.SetPerformanceMode(mode)

	// Ativar modo anime se solicitado
	if *animeFlag {
		p.SetAnimeMode(true)
	}

	// Configurar volume
	if *volume != 100 {
		p.SetVolume(*volume)
	}

	// Carregar vídeo
	args := flag.Args()
	if len(args) > 0 {
		videoPath := args[0]

		// Define título da janela
		windowTitle := *titleFlag
		if windowTitle == "" {
			windowTitle = "▶ " + filepath.Base(videoPath) + " - GoAnime Player"
		}
		p.SetTitle(windowTitle)

		// Fullscreen
		if *fullscreen {
			p.SetFullscreen(true)
		}

		if err := p.LoadFile(videoPath); err != nil {
			os.Exit(1)
		}

		// Carregar legenda externa se fornecida
		if *subFlag != "" {
			if err := p.LoadSubtitle(*subFlag); err != nil {
				fmt.Printf("[Player4K] Aviso: não foi possível carregar legenda: %v\n", err)
			}
		}

		// Posição inicial
		if *startPos > 0 {
			p.Seek(*startPos)
		}
	} else {
		printBanner()
		printUsage()
		return
	}

	// Loop de eventos
	p.Run()
}

func printBanner() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════╗
║     ▄▄▄▄▄  ▄            ▄         ▄  ▄    ▄▄▄   ▄▄▄▄     ║
║     █   █  █  ▄▄▄▄  ▄   █  ▄▄▄▄ ▄▄█▄▄█ █  █     █  █     ║
║     █▄▄▄█  █ █    █ █   █ █▄▄▄▄   █    █▀▀█     █▀▀█     ║
║     █      █ █▄▄▄▄█  ▀▀▀█ █▄▄▄▄▄  █▄▄  █  █▄▄   █  █     ║
║                     ▄▄▄▄▀                                ║
║          🎬 GoAnime Player 4K - Upscaling AI             ║
╚═══════════════════════════════════════════════════════════╝`)
}

func printUsage() {
	fmt.Println(`
📖 USO: player4k [opções] <arquivo_de_video>

🎛️  OPÇÕES:
   -mode=low|medium|high    Modo de qualidade (padrão: medium)
   -anime                   Ativar shaders Anime4K otimizados
   -title="Título"          Título personalizado da janela
   -sub="URL ou caminho"    Carregar legenda externa
   -fs                      Iniciar em tela cheia
   -volume=0-150            Volume inicial
   -start=SEGUNDOS          Posição inicial
   -list-modes              Ver modos disponíveis`)
}

func printControls() {
	fmt.Println(`
⌨️  ATALHOS PRINCIPAIS:
   ESPAÇO        Play/Pause
   ← →           Seek -5s/+5s
   ↑ ↓           Volume +/-
   I             Pular intro (85s)
   F             Tela cheia
   S             Screenshot
   M             Mute
   V             Mostrar/ocultar legendas
   J             Próxima legenda
   A             Próximo áudio
   [ ]           Velocidade -/+
   Q             Fechar`)
}
