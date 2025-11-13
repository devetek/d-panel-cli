## dPanel CLI

To override target API base URL, you can set `DPANEL_API_BASE_URL` environment variable.

```shell
export DPANEL_API_BASE_URL="https://pawon.terpusat.com"

dpid --email="youremail@example.com" --password="yourpassword" --ssh-ip="192.168.1.100" --ssh-port="22"
```

### Example

Basic usage will register your server with auto detection IP public, SSH port and available port for HTTP agent

```shell
dpid --email="youremail@example.com" --password="yourpassword"
```

Customize HTTP port:

```shell
dpid --email="youremail@example.com" --password="yourpassword" --http-port="9000"
```

Custom IP, SSH Port and HTTP Port:

```shell
dpid --email="youremail@example.com" --password="yourpassword" --ssh-ip="192.168.1.100" --ssh-port="22" --http-port="9000"
```
