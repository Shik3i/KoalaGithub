# KoalaGitHub 🐨⚡

Directory for GitHub profile visualizers, contribution art, statistics cards, streak graphs, and ready-to-copy workflows.

Live URL after deployment: [github.koalastuff.net](https://github.koalastuff.net)

## Features

- Compare compatible URL-based visualizers for any GitHub username.
- Switch themes and copy README Markdown or the required generator workflow.
- Keep configured generators separate from instant third-party previews.
- Use pseudonymous per-device votes without an account.
- Serve the prerendered SvelteKit UI and Go API from one container.

## Stack

- SvelteKit 2 and Svelte 5
- TypeScript and Vitest
- Go 1.26 with embedded frontend assets
- SQLite for visualizers, stars, and pseudonymous votes
- Multi-platform Docker image for `linux/amd64` and `linux/arm64`

## Local development

```bash
npm ci
npm run dev
```

The frontend development server runs at `http://localhost:5173`.

## Verification

```bash
npm run lint
npm run check
npm run test
npm run build
go vet ./internal/... .
go test ./internal/... .
```

`npm ci` runs `svelte-kit sync` through the package `prepare` script. The frontend build writes prerendered and precompressed assets to `www/`; the Go binary embeds that directory.

## Docker

```bash
docker compose up -d --build
```

The Compose service is available at `http://localhost:8088`. Its SQLite database is stored in the `koala_data` volume. The entrypoint repairs ownership of an existing root-owned volume, then runs the application as UID/GID `10001`.

Published images:

```bash
docker pull ghcr.io/shik3i/koalagithub:v1.2.3
```

## Releases

Only an exact semantic version tag in the form `vX.Y.Z` triggers `.github/workflows/publish-container.yml`.

```bash
git tag v1.2.3
git push origin v1.2.3
```

The workflow:

- runs frontend and backend quality gates;
- injects the tag into the footer and Go health response at build time;
- builds and pushes `linux/amd64` and `linux/arm64` images to GHCR;
- publishes exact, SemVer, major/minor, major, and `latest` tags;
- reuses GitHub Actions and BuildKit caches;
- attaches BuildKit provenance, an SBOM, and a GitHub artifact attestation.

The workflow does not commit a generated version file. The pushed tag is the release source of truth, so the footer cannot drift from the container tag.

## Add a visualizer

The central registry is `src/lib/data/visualizers.json`.

1. Add a unique entry with valid HTTPS project and repository URLs.
2. For configured generators, include a reviewed workflow pinned to immutable action SHAs.
3. Run the complete verification commands above.
4. Submit a pull request.

Visualizers that require private third-party account credentials or mandatory external secrets
(for example WakaTime API credentials) are not listed. The directory only includes previews and
workflows that can be used without handing KoalaGitHub private service credentials. This also
excludes generators that require a personal GitHub PAT as a mandatory repository secret, such as
`jstrieb/github-stats`.

## Privacy

KoalaGitHub contains no analytics, advertising trackers, or tracking cookies. It stores a persistent pseudonymous device ID and visualizer IDs for voting; the rate-limit window does not expire stored votes. Preview images are loaded directly from third-party providers, which receive normal HTTP request data. The application includes a self-service vote deletion control.

See the in-app Privacy Policy and [KoalaStuff legal notice](https://koalastuff.net/legal).

## License

[MIT License](LICENSE) © KoalaStuff
