package system

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// templatePath resuelve el .zshrc del template desde este paquete
// (installer/internal/system -> raíz del repo).
func templatePath(t *testing.T) string {
	t.Helper()
	return filepath.Join("..", "..", "..", "GentlemanZsh", ".zshrc")
}

// launcherAndInstantPromptLine devuelve la línea (0-based) de la llamada a
// start_if_needed y la del bloque de instant prompt de powerlevel10k, o -1.
func launcherAndInstantPromptLine(lines []string) (launcher int, instantPrompt int) {
	launcher, instantPrompt = -1, -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if launcher < 0 && trimmed == "start_if_needed" {
			launcher = i
		}
		if instantPrompt < 0 && strings.HasPrefix(trimmed, "if ") && strings.Contains(line, "p10k-instant-prompt") {
			instantPrompt = i
		}
	}
	return launcher, instantPrompt
}

// El lanzador del WM tiene que ejecutarse ANTES del bloque de instant prompt de
// powerlevel10k. Ese bloque corre el resto del archivo en un subshell con la
// salida redirigida: ahí `[[ -t 1 ]]` es falso y un `exec` reemplazaría solo al
// subshell, así que un lanzador colocado después nunca arranca el WM (tmux,
// zellij o herdr) y la terminal queda con un shell pelado.
func TestZshTemplateLaunchesWMBeforeP10kInstantPrompt(t *testing.T) {
	content, err := os.ReadFile(templatePath(t))
	if err != nil {
		t.Fatalf("no se pudo leer el template: %v", err)
	}

	launcher, instantPrompt := launcherAndInstantPromptLine(strings.Split(string(content), "\n"))

	if launcher < 0 {
		t.Fatal("el template no llama a start_if_needed")
	}
	if instantPrompt < 0 {
		t.Skip("el template ya no usa el instant prompt de powerlevel10k")
	}
	if launcher > instantPrompt {
		t.Errorf("el lanzador del WM (línea %d) quedó DESPUÉS del instant prompt de p10k (línea %d): el WM no va a arrancar",
			launcher+1, instantPrompt+1)
	}
}

// El instalador parchea .zshrc según el WM elegido. Ese parcheo no puede
// reordenar el archivo ni mover el lanzador por debajo del instant prompt.
func TestPatchZshForWMKeepsLauncherAboveInstantPrompt(t *testing.T) {
	for _, wm := range []string{"tmux", "zellij", "herdr"} {
		t.Run(wm, func(t *testing.T) {
			content, err := os.ReadFile(templatePath(t))
			if err != nil {
				t.Fatalf("no se pudo leer el template: %v", err)
			}

			dir := t.TempDir()
			path := filepath.Join(dir, ".zshrc")
			if err := os.WriteFile(path, content, 0o644); err != nil {
				t.Fatalf("no se pudo preparar la copia: %v", err)
			}

			if err := PatchZshForWM(path, wm, true); err != nil {
				t.Fatalf("PatchZshForWM(%q) falló: %v", wm, err)
			}

			patched, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("no se pudo leer la copia parcheada: %v", err)
			}
			lines := strings.Split(string(patched), "\n")

			launcher, instantPrompt := launcherAndInstantPromptLine(lines)
			if launcher < 0 {
				t.Fatalf("el parcheo con %q eliminó el lanzador del WM", wm)
			}
			if instantPrompt >= 0 && launcher > instantPrompt {
				t.Errorf("después del parcheo con %q el lanzador (línea %d) quedó debajo del instant prompt (línea %d)",
					wm, launcher+1, instantPrompt+1)
			}

			// el parcheo elige el WM elegido y deja el archivo sintácticamente sano
			if !strings.Contains(string(patched), `WM_CMD="`+wm+`"`) {
				t.Errorf("el parcheo con %q no fijó WM_CMD=%q", wm, wm)
			}
		})
	}
}
