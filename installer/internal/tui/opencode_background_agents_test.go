package tui

import (
	"os"
	"path/filepath"
	"testing"
)

// TestInstallOpenCodeBackgroundAgents_CopiesFiles valida que la función
// copia correctamente background-agents.ts y package.json al directorio
// destino de OpenCode, sin necesidad de que bun/npm existan en el entorno.
//
// El test usa directorios temporales para simular repoDir y openCodeDir,
// y verifica que los archivos llegan intactos al destino.
func TestInstallOpenCodeBackgroundAgents_CopiesFiles(t *testing.T) {
	// Contenido de prueba representativo de los archivos reales.
	// El formato del package.json usa pretty-print como el archivo real.
	pluginContent := `// background-agents.ts (test stub)
export default {};`
	pkgContent := "{\n  \"dependencies\": {\n    \"unique-names-generator\": \"^4.7.1\"\n  }\n}"

	// Construye un repoDir falso con la estructura que espera la función.
	repoDir := t.TempDir()
	pluginsSourceDir := filepath.Join(repoDir, "GentlemanOpenCode", "plugins")
	if err := os.MkdirAll(pluginsSourceDir, 0o755); err != nil {
		t.Fatalf("crear plugins source dir: %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(pluginsSourceDir, "background-agents.ts"),
		[]byte(pluginContent), 0o644,
	); err != nil {
		t.Fatalf("escribir background-agents.ts de origen: %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(repoDir, "GentlemanOpenCode", "package.json"),
		[]byte(pkgContent), 0o644,
	); err != nil {
		t.Fatalf("escribir package.json de origen: %v", err)
	}

	// openCodeDir representa ~/.config/opencode — vacío al inicio.
	openCodeDir := t.TempDir()

	// Ejecuta la función bajo prueba.
	// La función es privada pero el test vive en el mismo package (tui).
	err := installOpenCodeBackgroundAgents("test-step", repoDir, openCodeDir)
	if err != nil {
		t.Fatalf("installOpenCodeBackgroundAgents falló: %v", err)
	}

	// --- Verificaciones ---

	// 1. background-agents.ts debe existir en plugins/ del destino.
	dstPlugin := filepath.Join(openCodeDir, "plugins", "background-agents.ts")
	gotPlugin, err := os.ReadFile(dstPlugin)
	if err != nil {
		t.Fatalf("background-agents.ts no encontrado en destino (%s): %v", dstPlugin, err)
	}
	if string(gotPlugin) != pluginContent {
		t.Errorf("contenido de background-agents.ts no coincide:\ngot:  %q\nwant: %q",
			string(gotPlugin), pluginContent)
	}

	// 2. package.json debe existir en la raíz del destino.
	dstPkg := filepath.Join(openCodeDir, "package.json")
	gotPkg, err := os.ReadFile(dstPkg)
	if err != nil {
		t.Fatalf("package.json no encontrado en destino (%s): %v", dstPkg, err)
	}
	if string(gotPkg) != pkgContent {
		t.Errorf("contenido de package.json no coincide:\ngot:  %q\nwant: %q",
			string(gotPkg), pkgContent)
	}
}

// TestInstallOpenCodeBackgroundAgents_MissingPlugin verifica que la función
// devuelve un error descriptivo cuando background-agents.ts no existe en el repo.
func TestInstallOpenCodeBackgroundAgents_MissingPlugin(t *testing.T) {
	repoDir := t.TempDir()
	// Creamos la estructura de directorio pero SIN el plugin.
	if err := os.MkdirAll(
		filepath.Join(repoDir, "GentlemanOpenCode", "plugins"), 0o755,
	); err != nil {
		t.Fatalf("setup: %v", err)
	}
	// Tampoco ponemos package.json — queremos que falle antes.

	openCodeDir := t.TempDir()

	err := installOpenCodeBackgroundAgents("test-step", repoDir, openCodeDir)
	if err == nil {
		t.Fatal("esperaba error por plugin faltante, pero no hubo error")
	}
}

// TestInstallOpenCodeBackgroundAgents_MissingPackageJSON verifica que la función
// devuelve un error cuando package.json no existe en el repo.
func TestInstallOpenCodeBackgroundAgents_MissingPackageJSON(t *testing.T) {
	pluginContent := `// stub`

	repoDir := t.TempDir()
	pluginsDir := filepath.Join(repoDir, "GentlemanOpenCode", "plugins")
	if err := os.MkdirAll(pluginsDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	// Sólo el plugin existe; package.json está ausente.
	if err := os.WriteFile(
		filepath.Join(pluginsDir, "background-agents.ts"),
		[]byte(pluginContent), 0o644,
	); err != nil {
		t.Fatalf("setup plugin: %v", err)
	}

	openCodeDir := t.TempDir()

	err := installOpenCodeBackgroundAgents("test-step", repoDir, openCodeDir)
	if err == nil {
		t.Fatal("esperaba error por package.json faltante, pero no hubo error")
	}
}

// TestInstallOpenCodeBackgroundAgents_PluginsDirectoryCreated verifica que la
// función crea el subdirectorio plugins/ si no existe previamente en openCodeDir.
func TestInstallOpenCodeBackgroundAgents_PluginsDirectoryCreated(t *testing.T) {
	pluginContent := `// stub`
	pkgContent := `{}`

	repoDir := t.TempDir()
	pluginsDir := filepath.Join(repoDir, "GentlemanOpenCode", "plugins")
	if err := os.MkdirAll(pluginsDir, 0o755); err != nil {
		t.Fatalf("setup: %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(pluginsDir, "background-agents.ts"),
		[]byte(pluginContent), 0o644,
	); err != nil {
		t.Fatalf("setup plugin: %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(repoDir, "GentlemanOpenCode", "package.json"),
		[]byte(pkgContent), 0o644,
	); err != nil {
		t.Fatalf("setup pkg: %v", err)
	}

	// openCodeDir existe pero plugins/ NO existe adentro.
	openCodeDir := t.TempDir()

	err := installOpenCodeBackgroundAgents("test-step", repoDir, openCodeDir)
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}

	info, err := os.Stat(filepath.Join(openCodeDir, "plugins"))
	if err != nil {
		t.Fatal("plugins/ no fue creado por la función")
	}
	if !info.IsDir() {
		t.Fatal("plugins/ existe pero no es un directorio")
	}
}
