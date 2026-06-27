# engress sdk

Shared Go packages for the Engress platform.

The first extraction contains app-neutral runtime packages used by the agent,
edge, and core services:

| Package | Purpose |
| --- | --- |
| `logx` | Shared slog/charmbracelet logging setup |
| `pki` | Tunnel CA and certificate helpers |
| `stats` | Request counters and recent request snapshots |
| `version` | Build-time version metadata helper |

Application commands, UI flows, Clerk billing logic, deployment scripts, and
service-specific orchestration stay in their owning repositories.
