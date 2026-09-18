# TODO.AI.md

## CI/CD spec-vs-code mismatches (found while adding .github/workflows)

- [ ] AI.md's `beta.yml`/`daily.yml`/`release.yml`/`docker.yml` LDFLAGS
      template references `main.CommitID`, `main.BuildEpoch`,
      `main.OfficialSite` — none of these vars exist in `src/main.go` (actual
      vars: `main.Version`, `main.Commit`, `main.BuildDate`). Go's
      `-ldflags -X` silently no-ops for a nonexistent package var. All
      generated workflows and `docker/Dockerfile` (which also had this exact
      bug — `main.CommitID` fixed to `main.Commit` this session) now use the
      real var names, so the no-op bug itself is fixed everywhere. Still
      open: either add `CommitID`/`BuildEpoch`/`OfficialSite` vars to
      `main.go` to match AI.md's generic template, or update AI.md's
      template for this project to the real names — a documentation/spec
      decision, not a functional bug.
- [ ] `docker/Dockerfile.dev` does not exist. AI.md's `docker.yml` spec
      requires it for the `:devel` tag (debug-tooling image). `docker.yml`'s
      `build-devel` job currently builds `:devel` from the standard
      `docker/Dockerfile` as a stand-in. Add a real `docker/Dockerfile.dev`
      with debug tooling, or accept the stand-in permanently and drop the
      `build-devel` job.
- [ ] `docker/Dockerfile`'s actual build ARGs are `VERSION`, `BUILD_DATE`,
      `VCS_REF`, `LICENSE` — not AI.md's generic `docker.yml` template
      (`VERSION`, `BUILD_DATE`, `BUILD_EPOCH`, `COMMIT_ID`). `docker.yml` maps
      `COMMIT_ID` → `VCS_REF` and adds `LICENSE=MIT`; `BUILD_EPOCH` is dropped
      since the Dockerfile never consumes it.
- [ ] `renovate.json` was entirely absent before this session — added at the
      repo root (Go + GitHub Actions managers only).
