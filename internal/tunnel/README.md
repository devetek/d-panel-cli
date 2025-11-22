## Tunnel Client CLI

A lightweight Golang-based tunnel client designed to connect your internal machine to the dPanel Tunnel Server. This package also automatically generates configuration and installs a systemd service to ensure the tunnel client runs persistently in the background and reconnects automatically to the dPanel Tunnel Server.

### Development

1. Running cluster `./scripts/dev-cluster.sh`
2. SSH to cluster by execute: `ssh -i ssh-key/id_rsa_fake root@localhost -p 10000`
3. Execute commands inside container in folder `/root/apps/d-panel-cli`:
```sh
go run cmd/cli/*.go auth login --email="prakasa@dnocs.io" --password="prakasa"

go run cmd/cli/*.go tunnel create --tunnel-http-listener="8001" --tunnel-http-service="9000" --tunnel-ssh-listener="2221" --tunnel-ssh-service="22"

# Create machine with basic config (HTTP and SSH accessible publicly directly)
go run cmd/cli/*.go machine create --http-port="<DPANEL-AGENT-PORT>" --ssh-ip="<MACHINE-REAL-IP>" --ssh-port="<SSH-PORT>"

# Create machine with dPanel Agent (HTTP) behind proxy
go run cmd/cli/*.go machine create --http-port="<DPANEL-AGENT-PORT>" --ssh-ip="<MACHINE-REAL-IP>" --ssh-port="<SSH-PORT>" --http-domain="<DPANEL-AGENT-GATEWAY>"

# Create machine with behind dPanel tunnel
go run cmd/cli/*.go machine create --behind-tunnel
```
4. Open `https://cloud-beta.terpusat.com/v2/resources/servers?page=1` with user `prakasa@dnocs.io` and password `prakasa`