---
trigger: always_on
---

# 🌟 Windsurf Workspace Coding Guidelines

> These rules guide Windsurf in generating and modifying code within this workspace. Prioritize clarity, consistency, and maintainability.

---

## 📁 Project Structure

Use idiomatic Go project layout:

```
/cmd/app/             # main application entry point  
/internal/            # private app packages  
/pkg/                 # reusable code (if needed)  
/templates/           # HTML templates (for HTMX)  
/static/              # CSS/JS assets  
/web/handlers/        # HTTP handlers  
/web/middleware/      # custom middleware  
/db/                  # database access (SQL, migrations, etc.)
```

---

## 📌 General Principles

- Keep code simple and readable — prefer clarity over cleverness.  
- Organize logic into focused, testable components.  
- Avoid unnecessary dependencies unless required for ergonomics or performance.

---

## 💻 Go (Golang) Code Style

- Use `gofmt` and `goimports`.  
- All exported functions should have doc comments.  
- Prefer `log/slog` or `zerolog` for structured logging.  
- Use `context.Context` in all handler and service layers.  
- Avoid global state unless explicitly marked as safe and immutable.

---

## 🧪 Testing

- Place tests in `_test.go` files next to the code under test.  
- Use `testing` + `testify` or `go-cmp` for assertions.  
- Focus tests on behavior, not implementation details.  
- Avoid mocking unless necessary; prefer in-memory test implementations.

---

## 🧠 HTMX Integration

- Use server-rendered partials with `html/template`.  
- HTMX endpoints must return **partial HTML**, not full pages.  
- Place UI components in `/templates/components/` and render via Go templates.  
- Prefer `hx-target` and `hx-swap="outerHTML"` for replaceable UI regions.

---

## 🎨 CSS & JS

- Use **Tailwind CSS** for styling.  
- Only create custom classes if reusable and documented.  
- Use **Alpine.js** for interactivity; avoid SPA frameworks.  
- Keep `htmx` attributes declarative and clean.  
- Use `data-*` attributes for metadata and hooks.

---

## 🔐 Security & Middleware

- All forms must include CSRF protection (enforced via middleware).  
- Sensitive routes must use authentication middleware.  
- Always validate and sanitize user input on the server.

---

## 📄 Template Conventions

- Use layout templates (e.g., `base.html`) with `{{ block }}` for extension.  
- Componentize UI elements (e.g., buttons, modals, cards).  
- Do not include logic in templates — preprocess in handlers.

---

## 🔄 Change Rules

**When editing a file:**

- Preserve structure and formatting.  
- Refactor if it significantly improves readability or clarity.

**When creating a new file:**

- Use descriptive, conventional names (e.g., `user_handler.go`, `user_form.html`)  
- Include comments that explain purpose and usage.

---

## 🚧 Experimental Features

- New patterns (e.g., WebSockets, feature flags) must be clearly marked as experimental.  
- Document all conventions and abstractions in `/docs/conventions.md`.

---

## 🗃️ SQLite Usage Guidelines

This project uses **SQLite** as the primary data store.

### ✅ General Rules

- Use a **single shared `*sql.DB` instance** (no reconnecting per request).  
- Store the SQLite database file in `./data/app.db`.  
- Always check and handle errors when performing queries or transactions.

### 🏗 Schema Management

- Place schema initialization SQL in `db/schema.sql`.  
- Use `CREATE TABLE IF NOT EXISTS` to allow safe re-runs.  
- Store additional migration files in `db/migrations/` if needed.

### 🔄 Transactions

- Always use transactions (`db.Begin()`) for multi-step operations.  
- Rollback on any failure to ensure consistency.

### 📦 Query Best Practices

- Use **parameterized queries** (`?` placeholders) — never interpolate values directly.  
- Prefer `ExecContext` / `QueryRowContext` with `context.Context` for all DB interactions.

### 🔍 Data Access Organization

- Place all DB logic in the `db/` package (e.g., `db/users.go`).  
- Use functions like:
  ```go
  func GetUserByID(ctx context.Context, db *sql.DB, id int64) (*User, error)
  ```

### 🧪 Testing with SQLite

- Use in-memory SQLite (`:memory:`) for tests.  
- Load schema in `TestMain` or before each test run.

---

## 🛣️ Routing Guidelines with `chi`

This project uses the [Chi router](https://github.com/go-chi/chi).

### 📌 Core Principles

- Use strictly typed `http.HandlerFunc` functions.  
- Always pass `context.Context` down from request.  
- Organize routes by feature/domain.

### 📁 Route Definitions

- Define routes in `web/handlers/<feature>.go`.  
- Example:
  ```go
  func Routes() chi.Router {
      r := chi.NewRouter()
      r.Get("/", ListUsers)
      r.Post("/", CreateUser)
      r.Route("/{id}", func(r chi.Router) {
          r.Get("/", GetUser)
          r.Post("/edit", UpdateUser)
          r.Post("/delete", DeleteUser)
      })
      return r
  }
  ```

- Mount routes in `main.go` or `web/router.go`:
  ```go
  r.Mount("/users", user.Routes())
  ```

### 🧱 Middleware

- Use global middleware like:
  - `chi/middleware.Logger`
  - `chi/middleware.Recoverer`
  - `chi/middleware.Compress`

- Use `r.With(...)` for scoped middleware.

### ⚙️ Route Parameters and Validation

- Use `chi.URLParam(r, "id")` to extract parameters.  
- Always convert and validate parameters explicitly.

### 📄 Handler Design

- Keep handlers thin: parse → call service → render.  
- Handle errors early and clearly.  
- Use Go templates for all HTMX responses.

Example:

```go
func GetUser(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    user, err := db.GetUserByID(r.Context(), app.DB, int64(id))
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    tmpl.ExecuteTemplate(w, "partials/user.html", user)
}
```

---

## 🎀 DaisyUI Usage Guidelines

This project uses [DaisyUI](https://daisyui.com/) as a component library built on top of Tailwind CSS. It provides pre-styled, accessible UI components while preserving the flexibility of utility-first design.

### ✅ General Rules

- Use **DaisyUI components** as the base for all UI elements (buttons, inputs, alerts, modals, etc.).
- When customization is needed, extend with Tailwind utilities or DaisyUI `class` modifiers instead of writing raw CSS.
- Prefer **semantic HTML** with DaisyUI classes applied (e.g., use `<button>` for actions, not `<div>`).

---

### 🧩 Component Usage

- Reference component documentation at [https://daisyui.com/components/](https://daisyui.com/components/)
- Common components and usage examples:
  - Buttons: `class="btn btn-primary"`  
  - Form inputs: `class="input input-bordered"`  
  - Modals: use `modal` class and `modal-toggle` pattern
  - Alerts: `class="alert alert-warning"`  
  - Tabs, dropdowns, and accordions for navigation and structure

---

### 🎨 Customization

- Prefer DaisyUI `theme` support via `tailwind.config.js` over custom styling.
- Use Tailwind classes sparingly to override default behavior, e.g.:
  ```html
  <button class="btn btn-primary text-sm px-6">Save</button>
  ```

- Do not create redundant variants of existing DaisyUI styles unless required for accessibility or UX reasons.

---

### 🌐 Template Guidelines

- Wrap components in reusable Go template partials if they appear in multiple views.
  Example:
  ```
  /templates/components/button.html
  /templates/components/modal.html
  ```

- Use HTMX with DaisyUI modals and dynamic components by rendering them from partials and targeting visible containers using `hx-target`.

---

### 🧪 Testing UI

- For HTMX + DaisyUI interactions, ensure partials render complete component markup (including modal toggle or alert wrapper).
- When writing tests for UI behavior, verify presence of DaisyUI class names (e.g., `.btn`, `.modal`) in expected output.

---

### 🔄 Upgrades and Maintenance

- Track DaisyUI version in `package.json` or pinned CDN link.  
- Review changelogs when upgrading — especially when updating major versions.
- Test form and modal accessibility after major component version updates.

---