---
name: stream-deck
description: >
  Create slide-deck presentation webs for streams and courses using Gentleman Kanagawa Blur theme with inline SVG diagrams.
  Trigger: When building a presentation, slide deck, course material, stream web, or talk slides.
metadata:
  author: gentleman-programming
  version: "1.1"
---

## When to Use

- Building a slide-deck web presentation for streams, talks, or courses
- Creating inline SVG diagrams for dark-themed presentations
- Setting up a Kanagawa Blur themed web UI
- Generating visual diagrams with high contrast on dark backgrounds

---

## Architecture Overview

Single-page HTML presentation with:
- **No frameworks** — vanilla HTML/CSS/JS
- **No build step** — open `index.html` directly
- **No vertical scroll** — `100dvh` viewport, everything fits
- **Inline SVGs** — all diagrams are SVG elements in HTML (no image files)
- **Module system** — slides grouped into modules, displayed in a sidebar rail
- **Vim-mode theming** — lualine-inspired mode badges (Normal, Command, Insert, Visual, Terminal, Replace)

```
project/
├── index.html              # Single HTML file with all slides
└── assets/
    ├── css/styles.css      # Kanagawa Blur theme + layout
    └── js/app.js           # Navigation, dots, mode switching
```

---

## Critical Patterns

### Pattern 1: Gentleman Kanagawa Blur Color Palette

ALWAYS use these exact colors. Source: `gentleman-kanagawa-blur/lua/gentleman_kanagawa_blur/variant.lua`

```css
:root {
  /* Backgrounds */
  --bg: #06080f;
  --bg-dark: #06080f;
  --black: #06080f;
  --gray0: #191e28;         /* Rail/viewport background */
  --gray1: #1c212c;         /* Card backgrounds, surface */
  --gray2: #232a36;         /* Inner panels, surface2 */
  --gray3: #2a3142;         /* Deeper surface */
  --line: #313342;          /* Borders, separators */
  --line-strong: #8394A3;   /* Strong borders, visible dots */
  --selection: #263356;     /* Active module highlight */

  /* Text - CONTRAST IS CRITICAL */
  --fg: #f3f6f9;            /* Primary text — high contrast */
  --subtext1: #A1AABB;      /* Secondary text (paragraphs) — ~5.2:1 ratio */
  --subtext: #8394A3;       /* Tertiary text (eyebrow, hints) — ~4.2:1 ratio */

  /* Accent colors */
  --red: #cb7c94;
  --green: #b7cc85;
  --yellow: #dfbd76;        /* NOT #FFE066 — use the golden string color */
  --purple: #a3b5d6;
  --magenta: #ff8dd7;
  --orange: #deba87;
  --blue: #7fb4ca;
  --cyan: #7aa89f;
  --accent: #e0c15a;        /* Mode badge, counter, kicker */
}
```

### CRITICAL: Contrast Rules

| Use Case | WRONG Color | CORRECT Color | Why |
|----------|-------------|---------------|-----|
| Muted text on dark bg | `#5c6170` | `#8394A3` | 5c has ~1.5:1 ratio — INVISIBLE |
| Secondary text | `#8a8fa3` | `#A1AABB` | 8a is borderline ~2.8:1 |
| Yellow/gold | `#FFE066` | `#DFBD76` | FFE is neon, DFBD matches IDE theme |
| Dot borders | `#313342` | `#8394A3` | 31 disappears on dark backgrounds |

**Minimum contrast ratio: 4:1 against `#1c212c` surface backgrounds.**

### Pattern 2: Slide HTML Structure

Each slide is a two-column grid with text left, diagram right:

```html
<article class="slide" data-index="0" data-module="0" data-tone="blue">
  <div class="slide-content">
    <p class="slide-kicker">01 · MODULE NAME</p>
    <h2>Slide Title</h2>
    <p>Explanation paragraph.</p>
  </div>
  <figure class="slide-figure">
    <!-- INLINE SVG goes here — never <img> tags -->
    <svg viewBox="0 0 520 360" xmlns="http://www.w3.org/2000/svg"
         font-family="Space Grotesk,sans-serif">
      <!-- diagram content -->
    </svg>
  </figure>
</article>
```

**Key attributes:**
- `data-index` — global slide number (0-based)
- `data-module` — which module group (maps to rail sidebar)
- `data-tone` — color accent for that slide (blue, green, red, etc.)

### Pattern 3: Module/Rail System

Modules are groups of related slides. The sidebar shows module titles + dot indicators:

```html
<nav class="rail">
  <div class="rail-module" data-module="0">
    <button class="rail-title" data-first="0">1. Module Name</button>
    <div class="rail-dots" id="dots-0"></div>
  </div>
  <!-- dots are generated dynamically by JS -->
</nav>
```

The `modeMap` in JS maps module indices to vim modes:

```js
const modeMap = {
  0: "normal",   // blue
  1: "command",  // accent/gold
  2: "insert",   // green
  3: "visual",   // magenta
  4: "replace",  // orange
  5: "terminal", // cyan
  6: "normal",   // cycles back
  7: "command",
  8: "insert",
  9: "visual"
};
```

### Pattern 4: Inline SVG Design System

**ALL diagrams are inline SVGs, never image files.**

```
viewBox="0 0 520 360"
Background: transparent (parent container provides #06080f via CSS)
Font: Space Grotesk,sans-serif for ALL text
```

**SVG element conventions:**

| Element | Style |
|---------|-------|
| Section headers | `font-size="10" letter-spacing="2" fill="#A1AABB"` uppercase |
| Titles | `font-size="15-18" font-weight="600" fill="#F3F6F9"` |
| Card backgrounds | `fill="#1c212c" stroke="#313342"` rounded `rx="8-14"` |
| Inner panels | `fill="#232A36"` |
| Subtitle/labels | `fill="#A1AABB"` |
| Muted/decorative | `fill="#8394A3"` |
| Borders | `stroke="#313342"` |
| Glows | `feGaussianBlur` filters with accent color, opacity 0.2-0.3 |

**Filter ID rule:** prefix ALL filter IDs with `s{slideIndex}-` to avoid conflicts:
```xml
<filter id="s12-glow"> <!-- slide 12 -->
<filter id="s27-shadow"> <!-- slide 27 -->
```

### Pattern 5: Sub-Agent SVG Generation

Generate SVGs using sub-agents (Task tool) for parallelism. Each sub-agent gets:

1. The full design system (viewBox, palette, font, conventions)
2. The specific slide concept and content
3. Instructions to return ONLY raw `<svg>...</svg>` markup
4. The filter ID prefix for that slide

Then use `mcp_edit` to replace `<img>` tags with the returned SVG.

---

## Decision Tree

```
Need a new presentation?
  → Scaffold HTML with topbar, rail, viewport, controls
  → Define modules and slide count
  → Create CSS with full Kanagawa Blur palette
  → Create JS with navigation + buildDots()

Adding slides to existing deck?
  → Add <article class="slide"> with correct data-index, data-module, data-tone
  → Generate inline SVG via sub-agent with design system prompt
  → Insert SVG into <figure class="slide-figure">

Fixing contrast issues?
  → Check fill/stroke colors against the contrast table above
  → Use replaceAll to swap colors globally
  → Verify CSS variables match the corrected palette
```

---

## Code Examples

### Full Slide with Inline SVG

```html
<article class="slide" data-index="5" data-module="1" data-tone="yellow">
  <div class="slide-content">
    <p class="slide-kicker">06 · Context Window</p>
    <h2>Ventana limitada, impacto total</h2>
    <p>El contexto es RAM finita. Cada token que entra empuja otro afuera.</p>
  </div>
  <figure class="slide-figure">
    <svg viewBox="0 0 520 360" xmlns="http://www.w3.org/2000/svg"
         font-family="Space Grotesk,sans-serif" fill="none"
         role="img" aria-label="Ventana de contexto limitada">
      <defs>
        <filter id="s5-glow" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur stdDeviation="4" result="blur"/>
          <feFlood flood-color="#DFBD76" flood-opacity="0.25"/>
          <feComposite in2="blur" operator="in"/>
          <feMerge><feMergeNode/><feMergeNode in="SourceGraphic"/></feMerge>
        </filter>
      </defs>
      <text x="260" y="30" text-anchor="middle" font-size="10"
            letter-spacing="2" fill="#A1AABB">CONTEXTO</text>
      <rect x="60" y="50" width="400" height="260" rx="12"
            fill="#1c212c" stroke="#313342" stroke-width="1"/>
      <!-- ... diagram content ... -->
    </svg>
  </figure>
</article>
```

### CSS Layout Grid (No Scroll)

```css
.deck-app {
  height: 100dvh;
  max-height: 100dvh;
  display: grid;
  grid-template-rows: auto auto 1fr auto; /* topbar, progress, stage, controls */
  gap: 14px;
}

.stage-layout {
  min-height: 0;
  display: grid;
  grid-template-columns: 220px 1fr; /* rail + viewport */
  overflow: hidden;
}

.slide {
  position: absolute;
  inset: 0;
  display: grid;
  grid-template-columns: minmax(280px, 36%) 1fr; /* text + diagram */
}
```

### Sub-Agent Prompt Template for SVG Generation

```
You are an SVG diagram generator. Generate a SINGLE inline SVG.

## DESIGN SYSTEM (MANDATORY)
- viewBox="0 0 520 360", NO background rect (transparent)
- Font: Space Grotesk,sans-serif for ALL text
- Style: minimal, elegant — rounded rects (rx=8-14), subtle glows
- Section headers: font-size="10" letter-spacing="2" fill="#A1AABB" uppercase
- Filter IDs MUST be prefixed with `s{N}-`

## COLOR PALETTE
- fg: #F3F6F9, subtext: #A1AABB, muted: #8394A3
- surface: #1c212c, surface2: #232A36, line: #313342
- blue: #7FB4CA, green: #B7CC85, yellow: #DFBD76
- red: #CB7C94, purple: #A3B5D6, magenta: #FF8DD7
- orange: #DEBA87, cyan: #7AA89F, accent: #E0C15A

## SLIDE CONTENT
Title: "{title}"
Concept: {detailed description of what to diagram}

Return ONLY raw SVG markup. Start with <svg, end with </svg>.
```

---

## Commands

```bash
# Serve locally for development
python3 -m http.server 8080     # Then open http://localhost:8080

# Count slides
rg -c '<article class="slide"' index.html

# Verify no remaining <img> tags
rg 'assets/images/' index.html || echo "ALL INLINE"

# Count SVGs
rg -c '<svg' index.html

# Check for bad contrast colors (should return nothing)
rg '#FFE066|#5c6170|#8a8fa3' index.html || echo "CLEAN"

# Audit all fill colors used
rg -o 'fill="#[0-9a-fA-F]{6}"' index.html | sort | uniq -c | sort -rn
```

---

## Glass Morphism & Depth System

The visual depth comes from layered effects, NOT gradients:

```css
/* Card depth recipe */
.card {
  border: 1px solid var(--line);       /* subtle edge */
  border-radius: 16px;
  background: var(--gray1);            /* solid surface */
  box-shadow:
    0 16px 34px rgba(0, 0, 0, 0.36),  /* drop shadow */
    inset 0 1px 0 rgba(255,255,255,0.03); /* top highlight */
  backdrop-filter: blur(8px);          /* glass effect */
}

/* Ambient glow (background layer) */
.ambient::before {
  border-radius: 999px;
  filter: blur(120px);
  opacity: 0.18;
  background: var(--magenta);          /* or --blue */
}
```

**Rule: NO CSS gradients.** Use solid colors + shadows + blur for depth.

---

## Fonts

- **Headings**: `Source Serif 4` (serif, weight 500/700)
- **Body/UI**: `Space Grotesk` (sans-serif, weight 400/500/700)

```html
<link href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;700&family=Source+Serif+4:wght@500;700&display=swap" rel="stylesheet" />
```

---

## Navigation Features

- **Arrow keys** / PageUp/PageDown: navigate slides
- **Space**: next slide
- **Home/End**: first/last slide
- **F key**: toggle focus mode (hides rail sidebar)
- **Touch swipe**: mobile navigation
- **Rail clicks**: jump to module or specific slide dot
- **Animated transitions**: slide-enter/leave with cubic-bezier easing (420ms)
- **localStorage persistence**: saves current slide index, restores on reload

---

## Slide Persistence (localStorage)

The app saves the current slide to `localStorage` on every navigation, and restores it on page load:

```js
// Save on every navigation (inside render function, after updating current)
try { localStorage.setItem("deck-slide", current); } catch (_) {}

// Restore on load (at bottom of app.js)
const saved = Number(localStorage.getItem("deck-slide")) || 0;
render(Math.min(saved, slides.length - 1), { animate: false });
```

This means reloading the browser returns to the same slide. The `try/catch` handles private browsing where localStorage may be unavailable.

---

## Adding Slides to an Existing Deck

When inserting a new slide into an existing deck, you MUST re-index everything downstream:

1. **Insert the `<article class="slide">`** with the correct `data-index`, `data-module`, `data-tone`
2. **Increment `data-index`** on ALL subsequent slides (work highest to lowest to avoid collisions)
3. **Increment kicker numbers** in `<p class="slide-kicker">NN · Module</p>` (highest to lowest)
4. **Update `data-first`** on rail buttons for all modules AFTER the insertion point (+1 each)
5. **Update the counter** in the HTML (`01 / NN`)

**Always work from HIGHEST index down to LOWEST** to avoid double-incrementing.

Use a sub-agent (Task tool) for re-indexing — it's tedious but critical for navigation to work.

---

## SVG Layout Gotchas

### All elements MUST fit inside viewBox (0 0 520 360)

- Max Y for any element: **350** (leave 10px margin)
- Max X for any element: **510** (leave 10px margin)
- **Center diagrams horizontally**: calculate total width of all elements + gaps, then offset = `(520 - totalWidth) / 2`
- Elements that grow in size (like progressive documents) should be **bottom-aligned** — anchor their bases to the same Y line and let them grow upward
- **NEVER** place text labels outside the viewBox — they will be invisible or clipped
- Filter effects (`feGaussianBlur`) add visual bleed — account for ~4px extra around filtered elements
- When using `transform="translate(x, y)"`, all child coordinates are RELATIVE — add translate Y + child Y to get absolute position

### Checklist after editing any SVG:
```
1. Is the highest element > Y=10? (not clipped at top)
2. Is the lowest element < Y=350? (not clipped at bottom)
3. Is the leftmost element > X=10? (not clipped at left)
4. Is the rightmost element < X=510? (not clipped at right)
5. Are elements visually centered in the 520px width?
```

---

## Resources

- **Color source**: `gentleman-kanagawa-blur/lua/gentleman_kanagawa_blur/variant.lua`
- **Reference implementation**: See [assets/](assets/) for HTML scaffold template
