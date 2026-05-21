# LabHistory

LabHistory is a terminal-based tool (TUI) for managing lab environments, hosts, users, and related data in a fast and structured way.

Built with Go and Bubble Tea.

---

## 🚀 Features (v0.1)

- 🖥️ Host management (CRUD)
- 👤 User management (CRUD)
- 🌐 Subdomain support per host
- 🧭 Navigation-based TUI interface
- ⚡ Generic list component system
- 🗑️ Delete mode with confirmation
- 📦 Environment export system

---

## 🧱 Architecture

- `internal/app` → TUI logic (state, handlers, UI rendering)
- `internal/components` → reusable UI components (List, etc.)
- `internal/db` → database layer (SQL queries)
- `cmd/labhistory` → entry point

---

## 🎮 Controls

### Navigation
- ↑ / ↓ → move selection
- Enter → select / confirm
- Esc → cancel / back

### Modes
- Normal mode → browse data
- Add mode → create host/user
- Edit mode → modify entries
- Delete mode → confirm deletion

---

## 🗃️ Database

SQLite-based schema:

- machines (hosts)
- users
- subdomains (linked to hosts)

---

## ⚙️ Run

```bash
make run
```

or build:

```bash
make build
./bin/labhistory
```

