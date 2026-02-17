---
name: vitest
description: >
  Vitest testing patterns with React Testing Library.
  Trigger: When writing unit tests - AAA pattern, mocking, async testing.
license: Apache-2.0
metadata:
  author: prowler-cloud
  version: "1.0"
---

## Test Structure (REQUIRED)

Use the **Given/When/Then** (AAA - Arrange/Act/Assert) pattern with comments:

```typescript
it("should update user name when form is submitted", async () => {
  // Given - Arrange: setup test conditions
  const user = userEvent.setup();
  const onSubmit = vi.fn();
  render(<UserForm onSubmit={onSubmit} />);

  // When - Act: perform the action
  await user.type(screen.getByLabelText(/name/i), "John Doe");
  await user.click(screen.getByRole("button", { name: /submit/i }));

  // Then - Assert: verify the outcome
  expect(onSubmit).toHaveBeenCalledWith({ name: "John Doe" });
});
```

## Describe Block Organization (REQUIRED)

```typescript
describe("ComponentName", () => {
  // Group by feature/behavior, not by method
  describe("when user is authenticated", () => {
    it("should display user profile", () => {});
    it("should show logout button", () => {});
  });

  describe("when user is not authenticated", () => {
    it("should redirect to login", () => {});
  });

  describe("form validation", () => {
    it("should show error for invalid email", () => {});
    it("should disable submit when fields are empty", () => {});
  });
});
```

## Naming Conventions (REQUIRED)

```typescript
// ✅ Describe BEHAVIOR, not implementation
it("should display error message when API fails", () => {});
it("should disable button while loading", () => {});
it("should call onSubmit with form data", () => {});

// ❌ NEVER: Implementation details in names
it("should set isLoading to true", () => {});
it("should call useState", () => {});
it("should render div with class error", () => {});
```

## React Testing Library Queries (REQUIRED)

```typescript
// ✅ Priority order (most accessible first)
screen.getByRole("button", { name: /submit/i });      // 1. Role + accessible name
screen.getByLabelText(/email/i);                       // 2. Label (forms)
screen.getByPlaceholderText(/search/i);               // 3. Placeholder
screen.getByText(/welcome/i);                          // 4. Text content
screen.getByTestId("custom-element");                  // 5. Last resort

// ❌ NEVER: Query by class, id, or tag
container.querySelector(".btn-primary");
document.getElementById("submit");
```

## userEvent over fireEvent (REQUIRED)

```typescript
// ✅ userEvent: simulates real user behavior
const user = userEvent.setup();
await user.click(button);
await user.type(input, "hello");
await user.selectOptions(select, "option1");
await user.keyboard("{Enter}");

// ❌ NEVER: fireEvent for user interactions
fireEvent.click(button);
fireEvent.change(input, { target: { value: "hello" } });
```

## Async Testing Patterns (REQUIRED)

```typescript
// ✅ findBy for elements that appear asynchronously
const element = await screen.findByText(/loaded/i);

// ✅ waitFor for assertions that need to wait
await waitFor(() => {
  expect(screen.getByText(/success/i)).toBeInTheDocument();
});

// ✅ One assertion per waitFor (faster failure detection)
await waitFor(() => expect(mockFn).toHaveBeenCalled());
await waitFor(() => expect(screen.getByText(/done/i)).toBeVisible());

// ❌ NEVER: Multiple assertions in single waitFor
await waitFor(() => {
  expect(mockFn).toHaveBeenCalled();
  expect(screen.getByText(/done/i)).toBeVisible(); // Slower failures
});

// ❌ NEVER: Empty waitFor callback
await waitFor(() => {}); // Fragile, don't rely on "one tick"
```

## Mocking with vi.fn() (REQUIRED)

```typescript
// ✅ Basic mock function
const handleClick = vi.fn();
render(<Button onClick={handleClick} />);
await user.click(screen.getByRole("button"));
expect(handleClick).toHaveBeenCalledTimes(1);

// ✅ Mock with return value
const fetchUser = vi.fn().mockResolvedValue({ name: "John" });

// ✅ Mock implementation
const calculate = vi.fn().mockImplementation((a, b) => a + b);

// ✅ Always clean up mocks
afterEach(() => {
  vi.restoreAllMocks();
});
```

## vi.spyOn vs vi.mock

```typescript
// ✅ vi.spyOn: observe without replacing (preferred for most cases)
const spy = vi.spyOn(service, "fetchData").mockResolvedValue(mockData);
// Original implementation preserved, can be restored

// ✅ vi.mock: replace entire module (use sparingly)
vi.mock("@/services/api", () => ({
  fetchUser: vi.fn().mockResolvedValue({ name: "John" }),
}));

// ⚠️ vi.mock is hoisted! Variables won't be available
vi.mock("./module", () => ({
  fn: vi.fn().mockReturnValue(someVariable), // ERROR: someVariable not defined
}));

// ✅ Use factory function for dynamic values
vi.mock("./module", async () => {
  return {
    fn: vi.fn(),
  };
});
```

## Testing Hooks Pattern

```typescript
import { renderHook, act } from "@testing-library/react";

describe("useCounter", () => {
  it("should increment counter", () => {
    // Given
    const { result } = renderHook(() => useCounter());

    // When
    act(() => {
      result.current.increment();
    });

    // Then
    expect(result.current.count).toBe(1);
  });
});
```

## Snapshot Testing (USE SPARINGLY)

```typescript
// ✅ Inline snapshots for small outputs
it("should format date correctly", () => {
  expect(formatDate(new Date("2024-01-15"))).toMatchInlineSnapshot(
    `"January 15, 2024"`
  );
});

// ✅ File snapshots for larger outputs with syntax highlighting
it("should render correct HTML", () => {
  const { container } = render(<Card title="Test" />);
  expect(container).toMatchFileSnapshot("./snapshots/card.html");
});

// ❌ NEVER: Large component snapshots (hard to review)
expect(container).toMatchSnapshot(); // 500+ lines of HTML
```

## Test Isolation (REQUIRED)

```typescript
// ✅ Each test is independent
beforeEach(() => {
  vi.clearAllMocks();
});

afterEach(() => {
  vi.restoreAllMocks();
  cleanup(); // RTL cleanup (automatic with globals: true)
});

// ✅ Don't share state between tests
describe("Counter", () => {
  // ❌ NEVER: Shared mutable state
  let counter = 0;

  // ✅ Fresh setup per test
  it("should start at zero", () => {
    render(<Counter />);
    expect(screen.getByText("0")).toBeInTheDocument();
  });
});
```

## What NOT to Test

```typescript
// ❌ Don't test implementation details
expect(component.state.isLoading).toBe(true);
expect(useState).toHaveBeenCalled();

// ❌ Don't test third-party libraries
expect(axios.get).toHaveBeenCalled(); // Test YOUR code, not axios

// ❌ Don't test static content
expect(screen.getByText("Welcome")).toBeInTheDocument(); // Unless conditional

// ✅ Test user-visible behavior
expect(screen.getByRole("button")).toBeDisabled();
expect(screen.getByText(/error/i)).toBeVisible();
```

## Coverage Configuration

```typescript
// vitest.config.ts
export default defineConfig({
  test: {
    coverage: {
      provider: "v8",
      reporter: ["text", "json", "html"],
      all: true, // Include untested files
      include: ["src/**/*.{ts,tsx}"],
      exclude: [
        "**/*.test.{ts,tsx}",
        "**/*.d.ts",
        "**/types/**",
        "**/index.ts", // barrel files
      ],
      thresholds: {
        lines: 80,
        functions: 80,
        branches: 80,
        statements: 80,
      },
    },
  },
});
```

## File Organization

```
components/
├── Button/
│   ├── Button.tsx
│   ├── Button.test.tsx    # Co-located test
│   └── index.ts
├── Form/
│   ├── Form.tsx
│   ├── Form.test.tsx
│   └── index.ts
```

## Common Matchers

```typescript
// Presence
expect(element).toBeInTheDocument();
expect(element).toBeVisible();
expect(element).toBeEmptyDOMElement();

// State
expect(button).toBeDisabled();
expect(input).toHaveValue("text");
expect(checkbox).toBeChecked();
expect(input).toHaveFocus();

// Content
expect(element).toHaveTextContent(/hello/i);
expect(element).toHaveAttribute("href", "/home");
expect(element).toHaveClass("active");

// Functions
expect(fn).toHaveBeenCalled();
expect(fn).toHaveBeenCalledWith(arg1, arg2);
expect(fn).toHaveBeenCalledTimes(2);
```

## Anti-Patterns to Avoid

```typescript
// ❌ Testing internal state
expect(wrapper.instance().state.value).toBe(1);

// ❌ Using container.querySelector
const button = container.querySelector("button.primary");

// ❌ Waiting arbitrary time
await new Promise(r => setTimeout(r, 1000));

// ❌ Multiple acts
act(() => setState(1));
act(() => setState(2)); // Combine into one act

// ❌ Asserting on mock call order when order doesn't matter
expect(mock.mock.calls[0][0]).toBe("first");

// ❌ Snapshot everything
expect(entirePage).toMatchSnapshot();
```

## Keywords

vitest, testing, unit test, react testing library, mock, spy, AAA, given when then, arrange act assert, userEvent, waitFor, findBy, coverage
