# gorep - Un CLI de Búsqueda de Archivos Rápido en Go

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`gorep` es una herramienta de línea de comandos de alto rendimiento para buscar archivos, construida con Go. Su objetivo es proporcionar una alternativa rápida y eficiente a herramientas tradicionales como `grep` y `ripgrep`, centrándose en el procesamiento concurrente de archivos y una salida amigable para el usuario.

## Características

- **Búsqueda Concurrente**: Utiliza goroutines de Go para el procesamiento paralelo de archivos mediante un modelo de pool de trabajadores.
- **Salida en Color**: Resalta las coincidencias y las rutas de los archivos para una mejor legibilidad (se puede desactivar con `--no-color`).
- **Líneas de Contexto**: Permite mostrar líneas antes (`-B`), después (`-A`) o alrededor (`-C`) de las coincidencias, con bloques de contexto deduplicados similares a `ripgrep`.
- **Detección de Archivos Binarios**: Omite automáticamente los archivos binarios mediante filtrado basado en extensiones y detección de bytes nulos.
- **Filtrado Glob**: Incluye o excluye archivos basados en patrones glob (`--include` y `--exclude`).
- **Soporte para Gitignore**: Respeta las reglas de `.gitignore` por defecto (se puede desactivar con `--no-gitignore`).
- **Modos de Cadena Fija y Regex**: Búsqueda con cadenas literales (`-F`) o expresiones regulares.

## Instalación

### Desde el Código Fuente

```bash
go install github.com/yourusername/gorep@latest
```

(Sustituya `yourusername` por su nombre de usuario de GitHub una vez subido.)

### Usando el Binario Pre-construido

Una vez que las versiones estén disponibles en GitHub, puede descargar los binarios pre-construidos desde la [página de releases](https://github.com/yourusername/gorep/releases).

## Uso

```bash
# Búsqueda básica
gorep "pattern" /path/to/search

# Búsqueda de cadena fija (coincidencia literal)
gorep -F "function name" .

# Búsqueda insensible a mayúsculas/minúsculas
gorep -i "function" .

# Mostrar 2 líneas de contexto antes y después de las coincidencias
gorep -C 2 "function" .

# Incluir solo tipos de archivos específicos
gorep --include "*.go" --include "*.py" "function" .

# Excluir ciertos directorios o tipos de archivos
gorep --exclude "vendor/*" --exclude "*.log" "function" .

# Ignorar las reglas de .gitignore
gorep --no-gitignore "function" .

# Desactivar la salida en color
gorep --no-color "function" .
```

## Rendimiento

Resultados de benchmarks en un directorio con ~21,667 archivos (`/home/123kaze/Project/claude-code-source-code-main`):

| Caso | rg mediana | grep mediana | gorep mediana | Conclusión |
|------|-----------|-------------|--------------|------------|
| Cadena Fija: `function` | 0.0436s | 0.3404s | 0.1957s | `rg` más rápido |
| Cadena Fija: `TODO` | 0.0314s | 0.2400s | 0.1461s | `rg` más rápido |
| Cadena Fija: Sin coincidencia | 0.0324s | 0.2426s | 0.1387s | `rg` más rápido |
| Regex: `function|class|const` | 0.0509s | 0.4846s | 0.9682s | `rg` más rápido |
| Regex: `(function|class|const)\s+Identifier` | 0.0781s | 0.7447s | 1.3255s | `rg` más rápido |
| Cadena Fija + Solo TS/TSX: `import` | 0.0249s | 0.9552s | 0.1339s | `rg` más rápido |

**Conclusiones Clave**:
- `gorep` supera a `grep` en búsquedas de cadenas fijas.
- El rendimiento de regex de `gorep` necesita optimización en comparación con `rg` y a veces con `grep`.
- Reporte completo: [benchmark_report.md](benchmark_report.md)

## Desarrollo

```bash
git clone https://github.com/yourusername/gorep.git
cd gorep
go build .
```

## Hoja de Ruta (Roadmap)

- [x] Recorrido y búsqueda de archivos concurrente
- [x] Salida de terminal en color con resaltado de coincidencias
- [x] Visualización de líneas de contexto (`-B`, `-A`, `-C`)
- [x] Detección y omisión de archivos binarios
- [x] Inclusión/exclusión mediante patrones glob
- [x] Soporte para `.gitignore`
- [ ] Estadísticas de búsqueda (`--stats`)
- [ ] Formato de salida JSON (`--json`)
- [ ] Optimizaciones de rendimiento (ajuste de buffers, coincidencia a nivel de byte)
- [ ] Comparativas de benchmarks con `grep`/`ripgrep`

## Licencia

Este proyecto está licenciado bajo la Licencia MIT - vea el archivo [LICENSE](LICENSE) para más detalles.
