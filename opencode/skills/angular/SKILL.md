---
name: scope-rule-architect-angular
description: >
  Angular 20+ architecture with Scope Rule, Screaming Architecture, standalone components, and signals.
  Trigger: When writing Angular components, services, templates, or making architectural decisions about component placement.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## Core Angular 20 Principles

### 1. Standalone Components First

- **ALL components MUST be standalone** — never use NgModules for feature organization, since Angular 20 ALL components are standalone by default and don't need `standalone: true`
- Use `input()` and `output()` functions instead of decorators
- Implement `ChangeDetectionStrategy.OnPush` for all components
- Use `inject()` instead of constructor injection
- Don't use `any`
- Don't use lifecycle hooks like `ngOnInit` — use signals and computed instead
- Leverage signals for state management with `signal()`, `computed()`, and `effect()`

### 2. Modern Template Syntax

- Use native control flow (`@if`, `@for`, `@switch`) instead of structural directives
- Use `@defer` for lazy loading content and performance
- Prefer `class` and `style` bindings over `ngClass` and `ngStyle`
- Use `NgOptimizedImage` for all static images
- Implement reactive forms over template-driven forms
- Use typed reactive forms
- No `.component`, `.service`, `.module` suffixes in filenames — the name should tell the behavior

### 3. The Scope Rule — Unbreakable Law

**"Scope determines structure"**

- Code used by 2+ features → MUST go in global/shared directories
- Code used by 1 feature → MUST stay local in that feature
- NO EXCEPTIONS

### 4. Screaming Architecture

- Feature names must describe business functionality, not technical implementation
- Directory structure should tell the story of what the app does at first glance
- Main feature components MUST have the same name as their feature

## Decision Framework

1. **Count usage**: Identify exactly how many features use the component
2. **Apply the rule**: 1 feature = local placement, 2+ features = shared/global
3. **Validate against best practices**: Ensure compliance with Angular patterns
4. **Document decision**: Explain WHY the placement was chosen

## Project Structure

```
src/
  app/
    features/
      [feature-name]/
        [feature-name].ts          # Main standalone component
        components/                 # Feature-specific standalone components
          [component-name].ts
        services/                   # Feature-specific services with inject()
          [service-name].ts
        guards/                     # Feature-specific guards
        models/                     # Feature-specific interfaces/types
        signals/                    # Feature-specific signal stores
    shared/                         # ONLY for 2+ feature usage
      components/                   # Shared standalone components
      services/                     # Shared services
      guards/                       # Shared guards
      pipes/                        # Shared pipes
      directives/                   # Shared directives
    core/                           # Singleton services and app-wide concerns
      services/
        auth.ts
        api.ts
      interceptors/
      guards/
    main.ts                         # Bootstrap with standalone component
    app.config.ts                   # App configuration
    app.ts                          # Root standalone component
    routes.ts                       # Route configuration
```

## Path Aliases (tsconfig.json)

```json
{
  "paths": {
    "@features/*": ["src/app/features/*"],
    "@shared/*": ["src/app/features/shared/*"],
    "@core/*": ["src/app/core/*"]
  }
}
```

## Standalone Component Pattern

```typescript
import {
  Component,
  ChangeDetectionStrategy,
  signal,
  computed,
  input,
  output,
  inject,
} from "@angular/core";

@Component({
  selector: "app-feature-name",
  imports: [/* required dependencies */],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    @if (isLoading()) {
      <div>Loading...</div>
    } @else {
      @for (item of items(); track item.id) {
        <div>{{ item.name }}</div>
      }
    }
  `,
})
export class FeatureNameComponent {
  // Use input() function instead of @Input()
  readonly data = input<DataType>();
  readonly config = input({ required: true });

  // Use output() function instead of @Output()
  readonly itemSelected = output<ItemType>();

  // Use signals for state
  private readonly loading = signal(false);
  readonly isLoading = this.loading.asReadonly();

  // Use computed for derived state
  readonly items = computed(
    () => this.data()?.filter((item) => item.active) ?? [],
  );

  // Use inject() instead of constructor injection
  private readonly service = inject(FeatureService);
}
```

## Service with Signals

```typescript
import { Injectable, signal, computed, inject } from "@angular/core";

@Injectable({
  providedIn: "root",
})
export class FeatureService {
  private readonly http = inject(HttpClient);

  // Private signals for internal state
  private readonly _state = signal<FeatureState>({
    items: [],
    loading: false,
    error: null,
  });

  // Public readonly computed values
  readonly items = computed(() => this._state().items);
  readonly loading = computed(() => this._state().loading);
  readonly error = computed(() => this._state().error);

  loadItems(): void {
    this._state.update((state) => ({ ...state, loading: true }));
    // Implementation
  }
}
```

## Quality Checklist

1. **Scope verification**: Have you correctly counted feature usage?
2. **Angular compliance**: Are you using standalone components and modern patterns?
3. **Naming validation**: Do component names match feature names and follow Angular conventions?
4. **Screaming test**: Can a new Angular developer understand what the app does from the structure alone?
5. **Signal usage**: Are you leveraging signals appropriately for state management?
6. **Future-proofing**: Will this structure scale with Angular's evolution?

## Edge Cases

- **Legacy NgModule migration**: Always convert to standalone components during restructuring
- **Lazy loading**: Use standalone component routes instead of module-based lazy loading
- **Signal stores**: Place feature-specific signals locally, shared signals globally
- **Service scope**: Use `providedIn: 'root'` for shared services, local provision for feature-specific services
- **Form handling**: Implement reactive forms with signals for state management
