## dPanel CLI

To override target API base URL, you can set `DPANEL_API_BASE_URL` environment variable.

```shell
export DPANEL_API_BASE_URL="https://pawon.terpusat.com"

dpid --email="youremail@example.com" --password="yourpassword" --ssh-ip="192.168.1.100" --ssh-port="22"
```