# Vim Mastery Trainer - Especificación Completa

## Tabla de Contenidos

- [Contexto del Proyecto](#contexto-del-proyecto)
- [Arquitectura Existente](#arquitectura-existente)
- [Concepto: RPG de Vim](#concepto-rpg-de-vim)
- [Estructura de Progresión](#estructura-de-progresión)
- [Módulos de Entrenamiento](#módulos-de-entrenamiento)
- [UI Mockups](#ui-mockups)
- [Boss Fights](#boss-fights)
- [Estructura de Datos](#estructura-de-datos)
- [Estructura de Archivos](#estructura-de-archivos)
- [Integración con TUI Existente](#integración-con-tui-existente)
- [Componentes Bubbletea](#componentes-bubbletea)
- [Comandos de Desarrollo](#comandos-de-desarrollo)
- [Plan de Implementación (MVP)](#plan-de-implementación-mvp)
- [Estilo de Código](#estilo-de-código)

---

## Contexto del Proyecto

Juego de entrenamiento de Vim estilo RPG integrado al TUI installer de Gentleman.Dots. Construido con **Go + Bubbletea** (Charmbracelet).

**Ubicación**: `installer/` | **Nueva opción de menú**: "🎮 Vim Mastery Trainer"

## Arquitectura Existente

```text
installer/
├── cmd/gentleman-installer/main.go
├── internal/
│   ├── system/          # Detección OS, ejecución comandos
│   └── tui/
│       ├── installer.go # Lógica de instalación
│       ├── model.go     # Model principal (Bubbletea)
│       ├── update.go    # Update handlers
│       ├── view.go      # Views
│       ├── styles.go    # Lipgloss styles
│       ├── keymaps*.go  # Pantallas de keymaps
│       └── constants.go # Screens enum
```

---

## Concepto: RPG de Vim

Trainer estilo **keybr.com para Vim** con progresión RPG:

| Elemento | Descripción |
|----------|-------------|
| Módulos | Cada módulo es un "dungeon" |
| Progresión | Lecciones → Práctica → Jefe Final |
| Desbloqueo | Derrotar jefe desbloquea siguiente módulo |
| Características | Stats persistentes, spaced repetition, combos |

---

## Estructura de Progresión

```text
📖 LECCIONES (Tutorial)
    │
    │  Ejercicios guiados con explicación
    │  Sin timer estricto, enfoque en aprender
    │  100% para desbloquear práctica
    │
    ▼
🎯 PRÁCTICA (Grinding)
    │
    │  Ejercicios aleatorios del módulo
    │  Con timer, scoring, streaks
    │  80% accuracy para desbloquear jefe
    │
    ▼
👹 JEFE FINAL (Boss Fight)
    │
    │  Ejercicio épico que combina TODO del módulo
    │  Timer ajustado, múltiples pasos, 3 vidas
    │  Derrotarlo desbloquea siguiente sección
    │
    ▼
🔓 SIGUIENTE SECCIÓN DESBLOQUEADA
```

---

## Módulos de Entrenamiento

### 🏃 Movimientos Horizontales

| Comando | Descripción |
|---------|-------------|
| `w`, `W` | Siguiente palabra / PALABRA |
| `e`, `E` | Final de palabra / PALABRA |
| `b`, `B` | Inicio palabra anterior / PALABRA |
| `ge`, `gE` | Final palabra anterior / PALABRA |
| `f{c}` | Hasta carácter (inclusive) |
| `F{c}` | Hasta carácter hacia atrás |
| `t{c}` | Hasta carácter (exclusive) |
| `T{c}` | Hasta carácter hacia atrás (exclusive) |
| `;` | Repetir f/F/t/T |
| `,` | Repetir f/F/t/T en dirección opuesta |
| `0` | Inicio de línea |
| `$` | Final de línea |
| `^` | Primer carácter no-blanco |

### 📐 Movimientos Verticales

| Comando | Descripción |
|---------|-------------|
| `j`, `k` | Abajo / Arriba |
| `gg` | Primera línea |
| `G` | Última línea |
| `{n}G` | Ir a línea n |
| `{`, `}` | Párrafo anterior / siguiente |
| `H`, `M`, `L` | Top / Middle / Bottom de pantalla |
| `ctrl+d` | Media página abajo |
| `ctrl+u` | Media página arriba |
| `ctrl+f` | Página completa abajo |
| `ctrl+b` | Página completa arriba |

### 🎯 Text Objects

**CHANGE (c):**

| Comando | Descripción |
|---------|-------------|
| `ciw`, `caw` | Change inner/around word |
| `ci"`, `ca"` | Change inner/around "quotes" |
| `ci'`, `ca'` | Change inner/around 'quotes' |
| `ci{`, `ca{` | Change inner/around {braces} |
| `ci(`, `ca(` | Change inner/around (parens) |
| `ci[`, `ca[` | Change inner/around [brackets] |
| `cit`, `cat` | Change inner/around \<tags\> |
| `` ci` ``, `` ca` `` | Change inner/around \`backticks\` |

**Otros operadores:** Los mismos patterns aplican para `d` (delete), `y` (yank) y `v` (visual select).
Ejemplo: `diw`, `daw`, `yiw`, `viw`, etc.

### 🔁 Change & Repeat (El Flujo Mágico)

| Comando | Descripción |
|---------|-------------|
| `*` | Buscar palabra bajo cursor (forward) |
| `#` | Buscar palabra bajo cursor (backward) |
| `n`, `N` | Siguiente / anterior match |
| `gn` | Seleccionar próximo match (visual) |
| `cgn` | Cambiar próximo match |
| `dgn` | Borrar próximo match |
| `.` | Repetir último cambio |

**Combo Mágico (ventaja vs `:%s` - podés ELEGIR cuáles reemplazar):**

1. Cursor sobre palabra
2. `*` → Busca la palabra
3. `cgn` → Cambia el próximo match
4. `{texto}` → Escribí el reemplazo
5. `<Esc>` → Volver a normal
6. `.` → Repetir (siguiente match)
7. `n` → Saltear uno si querés
8. `.` → Seguir reemplazando

### 🔄 Sustitución (%s)

| Comando | Descripción |
|---------|-------------|
| `:s/foo/bar/` | Línea actual, primera ocurrencia |
| `:s/foo/bar/g` | Línea actual, todas las ocurrencias |
| `:%s/foo/bar/g` | Todo el archivo |
| `:%s/foo/bar/gc` | Todo el archivo, con confirmación |
| `:10,20s/foo/bar/g` | Rango de líneas (10-20) |
| `:'<,'>s/foo/bar/g` | Selección visual |
| `:s/foo/bar/i` | Case insensitive |
| `:s/foo/bar/I` | Case sensitive (forzado) |

**Patrones útiles:**

| Comando | Descripción |
|---------|-------------|
| `:%s/\s\+$//g` | Eliminar trailing whitespace |
| `:%s/foo/bar/gI` | Reemplazar exacto (case sensitive) |
| `:%s/\<foo\>/bar/g` | Solo palabras completas |

### 🔍 Regex & Vimgrep

**Búsqueda básica:**

| Comando | Descripción |
|---------|-------------|
| `/pattern` | Buscar hacia adelante |
| `?pattern` | Buscar hacia atrás |
| `n`, `N` | Siguiente / anterior match |
| `*` | Buscar palabra bajo cursor |

**Regex:**

| Comando | Descripción |
|---------|-------------|
| `/\<word\>` | Word boundaries (palabra completa) |
| `/pattern\c` | Case insensitive |
| `/pattern\C` | Case sensitive |
| `\v` | Very magic (menos escapes) |

**Very Magic Mode (`\v`):**

| Comando | Descripción |
|---------|-------------|
| `/\vfunction\s+\w+` | Buscar "function" + espacio + nombre |
| `/\v(\w+)@(\w+)` | Capturar grupos para email |

**Vimgrep:**

| Comando | Descripción |
|---------|-------------|
| `:vimgrep /pattern/g **/*.ts` | Buscar en todos los .ts |
| `:vimgrep /TODO/g **/*` | Buscar TODOs en proyecto |
| `:cnext`, `:cprev` | Navegar resultados |
| `:copen` | Abrir quickfix list |
| `:cclose` | Cerrar quickfix |

**Caracteres a escapar:** Sin `\v`: `. * [ ] ^ $ \ / ~` | Con `\v`: solo `/ \`

### 🎪 Macros

**Grabar:**

| Comando | Descripción |
|---------|-------------|
| `qa` | Empezar a grabar en registro 'a' |
| `{acciones}` | Las acciones que querés repetir |
| `q` | Parar de grabar |

**Ejecutar:**

| Comando | Descripción |
|---------|-------------|
| `@a` | Ejecutar macro del registro 'a' |
| `@@` | Repetir última macro ejecutada |
| `5@a` | Ejecutar macro 5 veces |
| `:5,10normal @a` | Ejecutar en líneas 5-10 |

**Tips:**
- Empezar macro con `0` o `^` (posición consistente)
- Terminar con `j` (ir a siguiente línea)
- Usar `f`/`t` en vez de `w` para mayor precisión

**Ejemplo - Convertir lista a array:**

```text
Antes:              Macro: qa0i"<Esc>A",<Esc>jq      Después de @a@@:
  item1                                                "item1",
  item2                                                "item2",
  item3                                                "item3",
```

---

## UI Mockups

### Menú Principal del Trainer

```text
┌─────────────────────────────────────────────────────────────────┐
│                    🎮 VIM MASTERY TRAINER                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   🏃 HORIZONTAL MOTIONS                                        │
│   ├── 📖 Lecciones      ████████████ 100%  ✓                   │
│   ├── 🎯 Práctica       ████████░░░░  70%                      │
│   └── 👹 Jefe Final     🔒 (completa práctica al 80%)          │
│                                                                 │
│   📐 VERTICAL MOTIONS                                          │
│   ├── 📖 Lecciones      ██████░░░░░░  50%                      │
│   ├── 🎯 Práctica       🔒                                     │
│   └── 👹 Jefe Final     🔒                                     │
│                                                                 │
│   🎯 TEXT OBJECTS                                              │
│   ├── 📖 Lecciones      🔒 (derrota jefe anterior)             │
│   ├── 🎯 Práctica       🔒                                     │
│   └── 👹 Jefe Final     🔒                                     │
│                                                                 │
│   🔁 CHANGE & REPEAT    🔒                                     │
│   🔄 SUSTITUCIÓN        🔒                                     │
│   🔍 REGEX & VIMGREP    🔒                                     │
│   🎪 MACROS             🔒                                     │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│   ⚔️  Jefes derrotados: 1/7    🏆 Score: 2,340                  │
└─────────────────────────────────────────────────────────────────┘
```

### Pantalla de Ejercicio (Lección/Práctica)

```text
┌─────────────────────────────────────────────────────────────────┐
│   🎯 TEXT OBJECTS    Nivel 5/10    🔥 Racha: 7    Score: 340   │
│   ████████████░░░░░░░░ 60%                                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   CÓDIGO:                                                       │
│   ┌───────────────────────────────────────────────────────────┐│
│   │ 1  const config = {                                      ││
│   │ 2    name: "█gentleman",                                 ││
│   │ 3    theme: "dark"                                       ││
│   │ 4  };                                                    ││
│   └───────────────────────────────────────────────────────────┘│
│                                                                 │
│   MISIÓN: Cambiá el contenido entre las comillas por "pro"     │
│           (el cursor está en la 'g' de gentleman)              │
│                                                                 │
│                         ⏱️  5.2s                                │
├─────────────────────────────────────────────────────────────────┤
│   Tu input: ci"_                                                │
│                                                                 │
│   💡 Pista en 3s...                                            │
└─────────────────────────────────────────────────────────────────┘
```

### Pantalla de Resultado

```text
┌─────────────────────────────────────────────────────────────────┐
│   ✅ CORRECTO!  +50pts  ⚡ 2.3s                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   Tu respuesta: ci"pro<Esc>                                    │
│   Solución óptima: ci"pro<Esc> ✓                               │
│                                                                 │
│   📝 EXPLICACIÓN:                                              │
│   ci" = Change Inside " (comillas)                             │
│   - c = change (borra y entra en insert mode)                  │
│   - i" = inner quotes (contenido entre comillas)               │
│                                                                 │
│   También válido: f"ci", vi"c                                  │
│                                                                 │
│   [Enter] Siguiente    [r] Repetir    [q] Menú                 │
└─────────────────────────────────────────────────────────────────┘
```

### Boss Fight

```text
┌─────────────────────────────────────────────────────────────────┐
│   👹 JEFE FINAL: The Horizontal Nightmare    ❤️ ❤️ ❤️           │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   CÓDIGO:                                                       │
│   ┌───────────────────────────────────────────────────────────┐│
│   │ 1  const█userName = getUser().name.firstName.toUpper();  ││
│   └───────────────────────────────────────────────────────────┘│
│                                                                 │
│   CADENA DE MISIONES:                         Ronda 1/5        │
│                                                                 │
│   ➤ Mové al final de "getUser"        ⏱️ 3s                    │
│   ○ Mové al inicio de "firstName"                              │
│   ○ Borrá "toUpper"                                            │
│   ○ Mové a la última 'e' de la línea                           │
│   ○ Volvé al inicio de la línea                                │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│   > fe_                                                         │
│                                                                 │
│   ⚡ Combo x2    👹 HP: ███████░░░                              │
└─────────────────────────────────────────────────────────────────┘
```

### Boss Derrotado

```text
┌─────────────────────────────────────────────────────────────────┐
│                                                                 │
│                     👹 JEFE DERROTADO! 👹                       │
│                                                                 │
│              ░░░░░▒▒▒▒▓▓▓▓████████▓▓▓▓▒▒▒▒░░░░░               │
│                                                                 │
│                   🏆 +500 PUNTOS 🏆                             │
│                                                                 │
│              ⏱️  Tiempo: 34.2s (Record: 28.1s)                  │
│              ❤️  Vidas restantes: 2/3                           │
│              ⚡ Mejor combo: x4                                  │
│                                                                 │
│         ┌─────────────────────────────────────────────┐       │
│         │  🔓 TEXT OBJECTS desbloqueado!              │       │
│         └─────────────────────────────────────────────┘       │
│                                                                 │
│   [Enter] Continuar    [r] Reintentar (mejorar record)         │
└─────────────────────────────────────────────────────────────────┘
```

---

## Boss Fights

### Bosses por Módulo

| Módulo | Boss Name | Mecánica Especial |
|--------|-----------|-------------------|
| Horizontal | The Line Walker | Navegar línea compleja sin j/k |
| Vertical | The Code Tower | Archivo de 50 líneas, llegar a puntos específicos |
| Text Objects | The Bracket Demon | Código anidado `{[({})]}`, cambiar contenidos |
| Change & Repeat | The Clone Army | 10 ocurrencias, reemplazar selectivamente con cgn |
| Sustitución | The Transformer | Transformaciones complejas con rangos y flags |
| Regex | The Pattern Master | Encontrar patterns complejos en código real |
| Macros | The Automaton | Grabar macro y aplicar en múltiples líneas |

### Mecánicas de Boss

| Mecánica | Descripción |
|----------|-------------|
| ❤️ Vidas | 3 errores y perdés (retry desde el inicio) |
| ⏱️ Timer | Más ajustado que práctica normal |
| Cadena | 5 pasos seguidos, todo conectado |
| Combo | Respuestas rápidas dan bonus (x2, x3, x4) |
| HP | Barra visual que se reduce con cada acierto |

---

## Estructura de Datos

### Exercise

```go
type Exercise struct {
    ID            string     // "horizontal_001"
    Module        string     // "horizontal", "textobjects", "cgn", etc.
    Level         int        // 1-10
    Type          string     // "lesson", "practice", "boss"
    Code          []string   // Líneas de código a mostrar
    CursorPos     Position   // Dónde está el cursor inicialmente
    CursorTarget  *Position  // Dónde debe terminar (para movimientos)
    Mission       string     // "Mové el cursor hasta la 'N' de 'Name'"
    Solutions     []string   // ["w", "W", "fe"] - todas las válidas
    Optimal       string     // "w" - la mejor/más corta
    Hint          string     // Pista que aparece después del timeout
    Explanation   string     // Explicación post-respuesta
    TimeoutSecs   int        // Segundos antes de mostrar solución
    Points        int        // Puntos base por completar
}

type Position struct {
    Line int
    Col  int
}
```

### Boss Exercise

```go
type BossExercise struct {
    ID          string
    Module      string
    Name        string       // "The Line Walker"
    Lives       int          // 3
    Steps       []BossStep   // Cadena de misiones
    BonusTime   int          // Tiempo total para bonus points
}

type BossStep struct {
    Exercise    Exercise
    TimeLimit   int          // Segundos para este paso específico
}
```

### User Stats

```go
type UserStats struct {
    TotalScore      int
    CurrentStreak   int
    BestStreak      int
    TotalTime       time.Duration
    ModuleProgress  map[string]*ModuleProgress
    BossesDefeated  []string
    LastPlayed      time.Time
}

type ModuleProgress struct {
    // Lecciones
    LessonsCompleted  int
    LessonsTotal      int
    
    // Práctica  
    PracticeAccuracy  float64  // 0.0 - 1.0
    PracticeAttempts  int
    PracticeCorrect   int
    
    // Boss
    BossDefeated      bool
    BossBestTime      time.Duration
    BossAttempts      int
    
    // Spaced Repetition
    WeakExercises     []string  // IDs de ejercicios que más falla
    LastPracticed     time.Time
}
```

### Archivo de Stats

Guardar en `~/.config/gentleman-trainer/stats.json`

```json
{
  "totalScore": 2340,
  "currentStreak": 7,
  "bestStreak": 23,
  "totalTimeSeconds": 8280,
  "lastPlayed": "2026-01-01T15:30:00Z",
  "bossesDefeated": ["horizontal", "vertical"],
  "modules": {
    "horizontal": {
      "lessonsCompleted": 15,
      "lessonsTotal": 15,
      "practiceAccuracy": 0.85,
      "practiceAttempts": 47,
      "practiceCorrect": 40,
      "bossDefeated": true,
      "bossBestTimeSeconds": 28,
      "bossAttempts": 3,
      "weakExercises": ["horizontal_012", "horizontal_008"],
      "lastPracticed": "2026-01-01T15:30:00Z"
    }
  }
}
```

---

## Estructura de Archivos

```text
installer/internal/tui/
├── trainer/
│   ├── model.go         # Model principal del trainer (Bubbletea)
│   ├── update.go        # Update handlers
│   ├── view.go          # Render de todas las pantallas
│   ├── styles.go        # Lipgloss styles específicos del trainer
│   ├── exercise.go      # Tipos y lógica de ejercicios
│   ├── boss.go          # Lógica específica de boss fights
│   ├── stats.go         # Persistencia de estadísticas
│   ├── validation.go    # Validar respuestas del usuario
│   └── exercises/
│       ├── horizontal.go    # Ejercicios de movimientos horizontales
│       ├── vertical.go      # Ejercicios de movimientos verticales
│       ├── textobjects.go   # Ejercicios de text objects
│       ├── cgn.go           # Ejercicios de change & repeat
│       ├── substitution.go  # Ejercicios de %s
│       ├── regex.go         # Ejercicios de regex/vimgrep
│       └── macros.go        # Ejercicios de macros
```

---

## Integración con TUI Existente

### 1. Agregar Screen en constants.go

```go
const (
    // ... screens existentes ...
    
    // Vim Trainer Screens
    ScreenVimTrainer        Screen = "vimtrainer"
    ScreenVimTrainerModule  Screen = "vimtrainer_module"
    ScreenVimTrainerLesson  Screen = "vimtrainer_lesson"
    ScreenVimTrainerPractice Screen = "vimtrainer_practice"
    ScreenVimTrainerBoss    Screen = "vimtrainer_boss"
    ScreenVimTrainerResult  Screen = "vimtrainer_result"
)
```

### 2. Agregar opción en menú principal

En `model.go`, agregar "🎮 Vim Mastery Trainer" como opción del menú principal.

### 3. Handler en update.go

Cuando se seleccione la opción del trainer, cambiar a `ScreenVimTrainer` y delegar al sub-model del trainer.

---

## Componentes Bubbletea

```go
import (
    "github.com/charmbracelet/bubbles/progress"   // Barras de progreso
    "github.com/charmbracelet/bubbles/timer"      // Countdown timer
    "github.com/charmbracelet/bubbles/textinput"  // Input del usuario
    "github.com/charmbracelet/bubbles/stopwatch"  // Medir tiempo de respuesta
    "github.com/charmbracelet/lipgloss"           // Estilos
)
```

---

## Comandos de Desarrollo

Ejecutar desde `installer/`:

| Comando | Descripción |
|---------|-------------|
| `go build -o gentleman.dots ./cmd/gentleman-installer` | Build del binario |
| `go test ./...` | Ejecutar todos los tests |
| `./gentleman.dots` | Ejecutar el installer |
| `go test ./internal/tui/trainer/... -v` | Tests específicos del trainer |

---

## Plan de Implementación (MVP)

### Fase 1: Estructura Base
- [ ] Crear estructura de archivos
- [ ] Model base del trainer con navegación
- [ ] Integración con menú principal
- [ ] Pantalla de selección de módulos

### Fase 2: Primer Módulo (Horizontal)
- [ ] 15 ejercicios de lección (guiados)
- [ ] Sistema de timer + input
- [ ] Validación de respuestas
- [ ] Pantalla de resultado con explicación
- [ ] 30 ejercicios de práctica (aleatorios)

### Fase 3: Boss Fight
- [ ] UI de boss con vidas y HP
- [ ] Cadena de 5 misiones
- [ ] Sistema de combos
- [ ] Pantalla de victoria/derrota

### Fase 4: Persistencia
- [ ] Guardar/cargar stats en JSON
- [ ] Tracking de progreso por módulo
- [ ] Spaced repetition básico

### Fase 5: Módulos Adicionales
- [ ] Vertical Motions
- [ ] Text Objects
- [ ] Change & Repeat (cgn)
- [ ] Sustitución
- [ ] Regex & Vimgrep
- [ ] Macros

---

## Estilo de Código

| Aspecto | Guideline |
|---------|-----------|
| Patterns | Seguir patterns existentes en el TUI (`installer.go`, `model.go`) |
| Estilos | Usar Lipgloss (referencia: `styles.go`) |
| Testing | Tests para ejercicios, validación y scoring |
| Commits | Conventional commits para cada feature |