# KoalaGitHub 🐨⚡

> Central hub for GitHub profile visualizers, statistics cards, streak graphs, trophies, contribution charts, and badges.

**Live URL**: [github.koalastuff.net](https://github.koalastuff.net)

---

## ✨ Features

- 🎯 **Universal Profile Comparison**: Enter any GitHub username (default: `Shik3i`) to instantly render all compatible profile visualizers side-by-side.
- 🎨 **On-The-Fly Theme Switcher**: Change themes (`default`, `dark`, `radical`, `tokyonight`, `github_dark`, `dracula`) directly on supported cards.
- 🔗 **Shareable URL Syncing**: Share customized results via URL parameter `?user=username`.
- ⚙️ **Configured vs. Universal Handling**: Differentiates between universal visualizers (instant preview) and configured tools (showing helpful setup guides instead of broken images).
- 📋 **1-Click Markdown Copy**: Instant inspection and 1-click clipboard copy of ready-to-use profile README Markdown.
- 🔒 **100% Static & Privacy Focused**: Zero analytics, zero ad trackers, zero cookies, zero backend databases.
- 🚀 **Caddy & Docker Ready**: Compiles to a 100% static `www/` directory ready for deployment under Caddy or Docker.

---

## 🛠️ Tech Stack

- **Framework**: [SvelteKit 2](https://kit.svelte.dev/) with Svelte 5 runes
- **Adapter**: `@sveltejs/adapter-static` (outputs static pages directly to `www/`)
- **Language**: TypeScript (strict mode)
- **Testing**: [Vitest](https://vitest.dev/) data integrity test suite
- **Server**: Caddy / Docker container

---

## 🚀 Getting Started Locally

```bash
# 1. Clone the repository
git clone https://github.com/Shik3i/KoalaGithub.git
cd KoalaGithub

# 2. Install dependencies
npm install

# 3. Start local development server
npm run dev
```

Visit `http://localhost:5173` in your browser.

---

## 🧪 Testing & Verification

```bash
# Run Vitest registry integrity tests
npm test

# Run Svelte & TypeScript typecheck
npm run check
```

---

## 📦 Building for Production (`www/` Output)

```bash
npm run build
```

This compiles 100% static HTML/CSS/JS files into the `www/` folder. You can upload the contents of `www/` directly to your web server (e.g. running Caddy, Nginx, Apache, or static host).

### Deploying with Docker & Caddy

```bash
docker compose up -d --build
```

Access the app locally at `http://localhost:8080`.

---

## ➕ How to Add a New Visualizer

All visualizer metadata is maintained in a single central JSON dataset:

📁 `src/lib/data/visualizers.json`

To register a new visualizer:

1. Add a new entry to `src/lib/data/visualizers.json` matching the schema.
2. Verify that all required fields (`id`, `imageUrlTemplate`, `markdownTemplate`, `themes`, `websiteUrl`, `repositoryUrl`) are provided.
3. Run `npm test; go test ./...` to verify data integrity.
4. Submit a Pull Request!

---

## 📄 Privacy & Transparency

KoalaGitHub has no tracking scripts, no cookies, and no analytics. However, when displaying third-party visualizer images, your browser makes direct HTTPS requests to third-party providers. A restrictive `referrerpolicy="no-referrer"` attribute is applied where practical.

For legal notices regarding the KoalaStuff ecosystem, visit [koalastuff.net/legal](https://koalastuff.net/legal).

---

## 📜 License

[MIT License](LICENSE) © KoalaStuff
