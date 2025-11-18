## 🧰 dnocs CLI

A command-line interface (CLI) tool designed to help you seamlessly manage your Linux (servers, PCs, and homelab) directly from the DeveTek Cloud Platform.

🔗 Integrated with DeveTek Cloud

Manage all your registered nodes remotely through DeveTek's cloud dashboard or through the powerful CLI.

### 🚀 Key Features
- 🔐 Secure Agent Registration – Register your machine to DeveTek Cloud
- 🧠 Remote Commands – Execute jobs/scripts via the cloud control panel

### 💡 Perfect For
- 🎮 Tech Hobbyists
- 🧪 Homelab Owners
- ⚙️ Tinkerers & DIYers

### 🛠️ Installation

```sh
curl -sfL https://raw.githubusercontent.com/devetek/d-panel-cli/refs/heads/main/scripts/install.sh | sh
```

### 📘 Example Usage

Available commands can be found by execute command `dnocs --help`

```sh
Simplify the process of managing resource such as user, machine, and application in dPanel (Devetek Panel).

Full documentation is available at: https://cloud.terpusat.com/docs/

Usage:
  dnocs [command]

Available Commands:
  auth        Manage dPanel session
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  machine     Manage dPanel machine
  version     Prints the version

Flags:
  -h, --help   help for dnocs

Use "dnocs [command] --help" for more information about a command.
```

🔑 Authentication
Log in to the DeveTek Cloud Platform via the CLI. This allows dnocs to perform authenticated operations securely.

```sh
dnocs auth login --email="user@example.com" --password="yourpassword"
```


🔐 Create This Machine

Register the current machine (the one where you're executing these commands) to the DeveTek Cloud Platform:

```sh
dnocs machine create --ssh-port="2000" --ssh-ip="20.192.45.121" --http-port="9500"
```

### 🌐 Documentation

Visit the official docs: https://cloud.terpusat.com/docs