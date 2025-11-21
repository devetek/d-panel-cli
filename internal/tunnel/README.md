## Tunnel Client CLI

A lightweight Golang-based tunnel client designed to connect your internal machine to the dPanel Tunnel Server. This package also automatically generates configuration and installs a systemd service to ensure the tunnel client runs persistently in the background and reconnects automatically to the dPanel Tunnel Server.

### Development

1. Running cluster `./scripts/dev-cluster.sh`
2. SSH to cluster by execute: `ssh -i ssh-key/id_rsa_fake root@localhost -p 10000`
3. Execute commands inside container in folder `/root/apps/d-panel-cli`:
```sh
go run cmd/cli/*.go tunnel create

go run cmd/cli/*.go auth login --email="prakasa@dnocs.io" --password="prakasa"

go run cmd/cli/*.go machine create --http-port="9000" --ssh-ip="tunnel.beta.devetek.app" --ssh-port="2221" --http-domain="https://internal-docker-01.beta.devetek.app"
```
4. Open `https://cloud-beta.terpusat.com/v2/resources/servers?page=1` with user `prakasa@dnocs.io` and password `prakasa`