# CLAUDE.md

This file guides Claude Code (claude.ai/code) when working in this repository.

## What this is

A small Go CLI that reads Traefik's `acme.json` (ACME v2 only) and writes each certificate out as PEM files. It ships as a Docker image that runs the exporter on a cron schedule.

## Commands

```bash
go build                      # produces ./traefik-ssl-certificate-exporter (gitignored)
./traefik-ssl-certificate-exporter --source traefik/acme.json --dest certs/ --owner 0 --group 0
go vet ./...
docker build -t traefik-ssl-certificate-exporter .
```

- There are no tests. CI (`.github/workflows/go.yml`) only runs `go build` on Go 1.16.5.
- `traefik/example.acme.json` shows the structure only. Its `certificate`/`key` values are placeholders and are not valid base64, so running the exporter against it fails at the decode step. For an end-to-end run you need a real `acme.json` (gitignored as `traefik/acme.json`).

## Layout

- `main.go`: all of the logic, in a single `main()`. Flags come from `spf13/pflag` and are bound into `spf13/viper` (`--source`, `--dest`, `--owner`, `--group`, `--config`). Read values with `viper.Get*`, not with the pflag return values.
- `models/acme.go`: `ProviderMdl` is `map[resolverName]Provider`, which mirrors the top level of `acme.json`. Each provider has `Account` and `Certificates[]` (`domain.main`, `domain.sans`, base64 `certificate`, base64 `key`).
- `entrypoint.sh`: the container entrypoint. It generates `/app/cert-export.sh` from env vars, runs it once if `ON_START == 1`, installs a crontab from `CRON_TIME`, and then runs `cron -f`. Cron output goes to `/var/log/cron.log` inside the container, not to stdout.
- `Dockerfile`: a two-stage build (`golang:latest` → `debian:latest` + `cron`). Container paths are fixed: `/app/traefik/acme.json` in and `/app/certs/` out.
- Module path is `rafi0101/traefik-ssl-certificate-exporter` (upstream author). Internal imports use it, e.g. `rafi0101/traefik-ssl-certificate-exporter/models`.

## Export behaviour

For every certificate in every resolver:
- Output dir is `<dest><domain.main>/`. `*` is replaced by `_`, so wildcards become `_.example.com/`. `dest` is joined by plain string concatenation, so it needs a trailing slash.
- Files written: `privkey.pem` (0600), `fullchain.pem`, `cert.pem` and `chain.pem` (0644). `cert`/`chain` come from splitting the fullchain at the first `\n-----BEGIN CERTIFICATE-----`.
- All writes go through `writeIfChanged`. It skips a file whose content is byte-identical, so mtimes stay put for unchanged certs. Otherwise it replaces the file with `ioutil.WriteFile` (truncating) and chowns it to `--owner`/`--group`.
- Any error calls `os.Exit(1)` immediately. Certificates later in the list are not processed.

## Known quirks

- `Chown` errors are deliberately ignored. This lets non-root local runs with the default owner/group 0 succeed.
- Permission modes only apply when a file is created. Existing files keep their mode.
- `io/ioutil` is used (deprecated, but fine on Go 1.16, which CI and `go.mod` target). Avoid APIs newer than 1.16.

## Config surface

Keep these in sync when adding options: the pflag definitions in `main.go`, the env vars in `Dockerfile` and `entrypoint.sh`, `docker-compose.yml`, and the Variables table in `README.md`.
