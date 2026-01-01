# Vim Mastery Trainer - Especificación Completa

## Contexto del Proyecto

Estamos agregando un **juego de entrenamiento de Vim estilo RPG** al TUI installer de Gentleman.Dots (`/Users/alanbuscaglia/Gentleman.Dots/installer`). El TUI está hecho en **Go con Bubbletea** (Charmbracelet).

El installer ya existe y funciona. Queremos agregar una nueva opción en el menú principal: **"🎮 Vim Mastery Trainer"**.

## Arquitectura Existente

```
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

Un trainer estilo **keybr.com pero para Vim**, con progresión tipo RPG:

- Cada módulo es un "dungeon"
- Progresión: **Lecciones → Práctica → Jefe Final**
- Derrotar al jefe desbloquea el siguiente módulo
- Stats persistentes, spaced repetition, combos

---

## Estructura de Progresión

```
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
```
w, W    → Siguiente palabra / PALABRA
e, E    → Final de palabra / PALABRA
b, B    → Inicio palabra anterior / PALABRA
ge, gE  → Final palabra anterior / PALABRA
f{c}    → Hasta carácter (inclusive)
F{c}    → Hasta carácter hacia atrás
t{c}    → Hasta carácter (exclusive)
T{c}    → Hasta carácter hacia atrás (exclusive)
;       → Repetir f/F/t/T
,       → Repetir f/F/t/T en dirección opuesta
0       → Inicio de línea
$       → Final de línea
^       → Primer carácter no-blanco
```

### 📐 Movimientos Verticales
```
j, k        → Abajo / Arriba
gg          → Primera línea
G           → Última línea
{n}G        → Ir a línea n
{, }        → Párrafo anterior / siguiente
H, M, L     → Top / Middle / Bottom de pantalla
ctrl+d      → Media página abajo
ctrl+u      → Media página arriba
ctrl+f      → Página completa abajo
ctrl+b      → Página completa arriba
```

### 🎯 Text Objects
```
CHANGE (c):
ciw, caw    → Change inner/around word
ci", ca"    → Change inner/around "quotes"
ci', ca'    → Change inner/around 'quotes'
ci{, ca{    → Change inner/around {braces}
ci(, ca(    → Change inner/around (parens)
ci[, ca[    → Change inner/around [brackets]
cit, cat    → Change inner/around <tags>
ci`, ca`    → Change inner/around `backticks`

DELETE (d):
diw, daw, di", da", di{, da{, di(, da(, etc.

YANK (y):
yiw, yaw, yi", ya", yi{, ya{, yi(, ya(, etc.

VISUAL SELECT (v):
viw, vaw, vi", va", vi{, va{, vi(, va(, etc.
```

### 🔁 Change & Repeat (El Flujo Mágico)
```
*           → Buscar palabra bajo cursor (forward)
#           → Buscar palabra bajo cursor (backward)
n, N        → Siguiente / anterior match
gn          → Seleccionar próximo match (visual)
cgn         → Cambiar próximo match
dgn         → Borrar próximo match
.           → Repetir último cambio

COMBO MÁGICO:
1. Cursor sobre palabra
2. *        → Busca la palabra
3. cgn      → Cambia el próximo match
4. {texto}  → Escribí el reemplazo
5. <Esc>    → Volver a normal
6. .        → Repetir (siguiente match)
7. n        → Saltear uno si querés
8. .        → Seguir reemplazando

VENTAJA vs :%s → Podés ELEGIR cuáles reemplazar
```

### 🔄 Sustitución (%s)
```
:s/foo/bar/         → Línea actual, primera ocurrencia
:s/foo/bar/g        → Línea actual, todas las ocurrencias
:%s/foo/bar/g       → Todo el archivo
:%s/foo/bar/gc      → Todo el archivo, con confirmación
:10,20s/foo/bar/g   → Rango de líneas (10-20)
:'<,'>s/foo/bar/g   → Selección visual
:s/foo/bar/i        → Case insensitive
:s/foo/bar/I        → Case sensitive (forzado)

PATRONES ÚTILES:
:%s/\s\+$//g        → Eliminar trailing whitespace
:%s/foo/bar/gI      → Reemplazar exacto (case sensitive)
:%s/\<foo\>/bar/g   → Solo palabras completas
```

### 🔍 Regex & Vimgrep
```
BÚSQUEDA BÁSICA:
/pattern            → Buscar hacia adelante
?pattern            → Buscar hacia atrás
n, N                → Siguiente / anterior match
*                   → Buscar palabra bajo cursor

REGEX:
/\<word\>           → Word boundaries (palabra completa)
/pattern\c          → Case insensitive
/pattern\C          → Case sensitive
\v                  → Very magic (menos escapes)

VERY MAGIC MODE (\v):
/\vfunction\s+\w+   → Buscar "function" + espacio + nombre
/\v(\w+)@(\w+)      → Capturar grupos para email

VIMGREP:
:vimgrep /pattern/g **/*.ts     → Buscar en todos los .ts
:vimgrep /TODO/g **/*           → Buscar TODOs en proyecto
:cnext, :cprev                  → Navegar resultados
:copen                          → Abrir quickfix list
:cclose                         → Cerrar quickfix

CARACTERES A ESCAPAR (sin \v):
. * [ ] ^ $ \ / ~
Con \v solo escapar: / \
```

### 🎪 Macros
```
GRABAR:
qa          → Empezar a grabar en registro 'a'
{acciones}  → Las acciones que querés repetir
q           → Parar de grabar

EJECUTAR:
@a          → Ejecutar macro del registro 'a'
@@          → Repetir última macro ejecutada
5@a         → Ejecutar macro 5 veces
:5,10normal @a  → Ejecutar en líneas 5-10

TIPS:
- Empezar macro con 0 o ^ (posición consistente)
- Terminar con j (ir a siguiente línea)
- Usar f/t en vez de w para mayor precisión

EJEMPLO - Convertir lista a array:
Antes:
  item1
  item2
  item3

Macro: qa0i"<Esc>A",<Esc>jq
Resultado después de @a@@:
  "item1",
  "item2",
  "item3",
```

---

## UI Mockups

### Menú Principal del Trainer

```
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

```
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

```
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

```
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

```
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
│         ┌─────────────────────────────────────────────�─┐       │
│         │  🔓 TEXT OBJECTS desbloqueado!              │       │
│         └─────────────────────────────────────────────┘       │
│                                                                 │
│   [Enter] Continuar    [r] Reintentar (mejorar record)         │
└─────────────────────────────────────────────────────────────────┘
```

---

## Boss de Cada Módulo

| Módulo | Boss Name | Mecánica Especial |
|--------|-----------|-------------------|
| Horizontal | The Line Walker | Navegar línea compleja sin j/k |
| Vertical | The Code Tower | Archivo de 50 líneas, llegar a puntos específicos |
| Text Objects | The Bracket Demon | Código anidado `{[({})]}`, cambiar contenidos |
| Change & Repeat | The Clone Army | 10 ocurrencias, reemplazar selectivamente con cgn |
| Sustitución | The Transformer | Transformaciones complejas con rangos y flags |
| Regex | The Pattern Master | Encontrar patterns complejos en código real |
| Macros | The Automaton | Grabar macro y aplicar en múltiples líneas |

### Boss Mechanics

- **❤️ Vidas**: 3 errores y perdés (retry desde el inicio)
- **⏱️ Timer por paso**: Más ajustado que práctica normal
- **Cadena de misiones**: 5 pasos seguidos, todo conectado
- **Combo multiplier**: Respuestas rápidas dan bonus (x2, x3, x4)
- **Boss HP**: Barra visual que se reduce con cada acierto

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

## Estructura de Archivos a Crear

```
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

## Componentes Bubbletea a Usar

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

## Comandos Útiles

```bash
cd /Users/alanbuscaglia/Gentleman.Dots/installer

# Build
go build -o gentleman.dots ./cmd/gentleman-installer

# Test
go test ./...

# Run
./gentleman.dots

# Test específico
go test ./internal/tui/trainer/... -v
```

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

- Seguir patterns existentes en el TUI (ver `installer.go`, `model.go`)
- Usar Lipgloss para estilos (ya hay `styles.go` de referencia)
- Tests para lógica de ejercicios, validación y scoring
- Conventional commits para cada feature
