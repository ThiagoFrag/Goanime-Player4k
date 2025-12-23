# Shaders para Player4K

Esta pasta contém os shaders GLSL usados para upscaling de vídeo.

## Download dos Shaders

### 1. FSR (AMD FidelityFX Super Resolution)
**Para modo Medium**

Download: https://gist.github.com/agyild/82219c545228d70c5604f865ce0b0ce5

Salvar como: `FSR.glsl`

### 2. FSRCNNX (Neural Network Upscaler)
**Para modo High**

Download: https://github.com/igv/FSRCNN-TensorFlow/releases

Arquivo: `FSRCNNX_x2_16-0-4-1.glsl`

### 3. CAS (Contrast Adaptive Sharpening)
**Para nitidez adicional**

Download: https://gist.github.com/agyild/bbb4e58298b2f86aa24da3032a0d2f7f

Salvar como: `CAS.glsl`

### 4. Anime4K (Para Anime)
**Otimizado para animação**

Download: https://github.com/bloc97/Anime4K/releases

Extrair na pasta `Anime4K/`

## Estrutura da Pasta

```
shaders/
├── FSR.glsl                      # AMD FSR (Medium mode)
├── FSRCNNX_x2_16-0-4-1.glsl     # Neural Network (High mode)
├── CAS.glsl                      # Sharpening
├── Anime4K/                      # Pasta com shaders Anime4K
│   ├── Anime4K_Clamp_Highlights.glsl
│   ├── Anime4K_Restore_CNN_VL.glsl
│   ├── Anime4K_Upscale_CNN_x2_VL.glsl
│   └── ...
└── README.md                     # Este arquivo
```

## Requisitos de GPU por Shader

| Shader | GPU Mínima | VRAM |
|--------|-----------|------|
| FSR | GTX 1050 / RX 560 | 2GB |
| FSRCNNX x2 16-0-4-1 | RTX 3060 / RX 6700 | 6GB |
| Anime4K (VL) | RTX 3070 / RX 6800 | 8GB |
| Anime4K (M) | GTX 1060 / RX 580 | 4GB |
