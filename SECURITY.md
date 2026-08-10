# Security Policy

## Supported Versions

This repository is a collection of learning examples and does not publish
supported releases. Only the latest commit on `main` receives security fixes.

## Reporting a Vulnerability

Do not disclose suspected vulnerabilities in a public issue. Use GitHub's
[private vulnerability reporting](https://github.com/rahilsh/golang-lab/security/advisories/new)
to provide the affected component, reproduction steps, and potential impact.
You should receive an initial response within seven days.

These examples are intentionally minimal. The gRPC example uses plaintext
transport and the HTTP example has no authentication. Add authentication,
authorization, TLS, resource limits, and operational hardening before adapting
any of them for a production environment.
