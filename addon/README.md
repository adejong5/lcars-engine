# LCARS Engine — Home Assistant add-on

Runs the LCARS Engine server as a Home Assistant add-on, served in the sidebar
through HA's authenticated ingress. The same image also runs as a standalone
container; see the repository root.

## Install

1. **Settings → Add-ons → Add-on Store → ⋮ → Repositories**, and add
   `https://github.com/adejong5/lcars-engine`.
2. Install **LCARS Engine** from the new repository entry and start it.
3. Open it from the sidebar (ingress) — no configuration needed.

The add-on connects to Home Assistant automatically through the Supervisor
(`homeassistant_api`), so there is no host or token to enter; the token never
leaves the Supervisor.

## How it works

- **Image:** the Supervisor pulls a prebuilt multi-arch image
  (`ghcr.io/adejong5/{arch}-lcars-engine`) published by the release workflow on
  each version tag. Keep `addon/config.yaml`'s `version` in sync with the pushed
  image tag.
- **Connection:** in add-on mode the server detects `SUPERVISOR_TOKEN` and uses
  `ws://supervisor/core/websocket`.
- **Ingress:** the server reads `X-Ingress-Path` and renders relative URLs, so
  the same markup works behind the proxy prefix or at `/` standalone.

## Releasing a new version

1. Bump `version` in `addon/config.yaml`.
2. Tag the repo `vX.Y.Z` and push — the workflow builds and pushes the per-arch
   images with that tag.
