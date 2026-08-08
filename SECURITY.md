# Security Policy

## Supported Versions

This repository is a learning example and does not publish supported releases.
Only the latest commit on `main` receives security fixes.

## Reporting a Vulnerability

Do not disclose suspected vulnerabilities in a public issue. Use GitHub's
[private vulnerability reporting](https://github.com/rahilsh/grpc-go/security/advisories/new)
to provide the affected component, reproduction steps, and potential impact.
You should receive an initial response within seven days.

The example intentionally uses plaintext gRPC transport. Add authentication,
authorization, TLS, resource limits, and operational hardening before adapting
it for a production environment.
