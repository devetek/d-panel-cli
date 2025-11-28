package tunnel

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/devetek/d-panel-cli/internal/logger"
	"github.com/devetek/tuman/pkg/marijan"
)

var (
	serviceCfg        = "/usr/lib/systemd/system/dpanel-tunnel.service"
	serviceName       = "dpanel-tunnel"
	configFolder      = "/opt/dpanel/tunnel"
	configFile        = "config.json"
	TunnelHost        = "tunnel.beta.devetek.app"
	TunnelPort        = "2220"
	BinaryBaseURL     = "https://github.com/devetek/tuman/releases"
	binaryDownloadURL = BinaryBaseURL + "/download"
	binaryVersion     = "v0.1.1-beta.2"
)

type tunnelServer struct {
	host string
	port string
	// auth
}

type tunnelService struct {
	folder      string
	name        string
	serviceCfg  string
	serviceName string
	configs     []marijan.Config
}

type tunnel struct {
	baseURL string
	version string
	bin     string
	server  tunnelServer  // tunnel server (any SSH server)
	service tunnelService // tunnel service in this server
}

func NewTunnel() *tunnel {
	// If we want to use different binary location, we can switch with env variable
	apiURL := os.Getenv("DNOCS_TUNNEL_BASE_URL")
	if apiURL != "" {
		binaryDownloadURL = apiURL
	}

	version := os.Getenv("DNOCS_TUNNEL_VERSION")
	if version != "" {
		binaryVersion = version
	}

	client := &tunnel{
		baseURL: binaryDownloadURL,
		version: binaryVersion,
		bin:     "marijan",
		server: tunnelServer{
			host: TunnelHost,
			port: TunnelPort,
		},
		service: tunnelService{
			folder:      configFolder,
			name:        configFile,
			serviceCfg:  serviceCfg,
			serviceName: serviceName,
			configs:     []marijan.Config{},
		},
	}

	return client
}

func (tun *tunnel) SetConfig(configs []marijan.Config) *tunnel {
	tun.service.configs = configs

	return tun
}

func (tun *tunnel) GetConfig() []marijan.Config {
	var configs []marijan.Config
	var finalPath = filepath.Join(tun.service.folder, tun.service.name)

	// read tunnel config
	// 1. Read the file content
	fileBytes, err := os.ReadFile(finalPath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return configs
	}

	err = json.Unmarshal(fileBytes, &configs)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return configs
	}

	return configs
}

func (tun *tunnel) SetNewVersion(version string) {
	tun.version = version
}

func (tun *tunnel) GetCurrentVersion() string {
	cmd := exec.Command("marijan", "version")

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

func (tun *tunnel) GetNewVersion() string {
	newVersion, err := getLatestReleaseVersion()
	if err != nil {
		return ""
	}

	return newVersion
}

func (tun *tunnel) fileName() string {
	return tun.bin + "-" + tun.version + "-" + runtime.GOOS + "-" + runtime.GOARCH + ".tar.gz"
}

func (tun *tunnel) source() string {
	finalUrl, err := url.JoinPath(tun.baseURL, tun.version, tun.fileName())
	if err != nil {
		return ""
	}

	return finalUrl
}

func (tun *tunnel) pwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return dir
}

func (tun *tunnel) destination() string {
	current := tun.pwd()
	if current == "" {
		return ""
	}

	// download destination
	var destination = filepath.Join(current, tun.fileName())

	return destination
}

func (tun *tunnel) Download() error {
	// validate source url location
	var source = tun.source()
	if source == "" {
		return errors.New("invalid source download url")
	}

	// validate destination file location
	var destination = tun.destination()
	if destination == "" {
		return errors.New("invalid destination folder")
	}

	// start to downloading artifact
	logger.Success(fmt.Sprintf("⬇️ Downloading Marijan for %s (%s)", runtime.GOOS, runtime.GOARCH))

	// Create the file
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(source)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request to %s: %w", source, err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy response body to file: %w", err)
	}

	logger.Success(fmt.Sprintf("🥳 Marijan downloaded successfully to %s", tun.destination()))

	return nil
}

func (tun *tunnel) Extract() error {
	// validate destination file location
	var source = tun.destination()
	if source == "" {
		return errors.New("invalid destination folder")
	}

	destination := filepath.Join(tun.pwd(), tun.bin)
	if destination == "" {
		return errors.New("extract destination is not valid folder")
	}

	// start to exract artifact
	logger.Success(fmt.Sprintf("📦 Extracting Marijan for %s (%s).", runtime.GOOS, runtime.GOARCH))

	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(destination, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(destination, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer outFile.Close() // Defer closing for each file

			if _, err := io.Copy(outFile, tr); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Warning: Skipping unknown type %d in %s\n", header.Typeflag, header.Name)
		}
	}

	// Create the destination directory if it doesn't exist
	err = os.MkdirAll("/usr/local/bin", 0755)
	if err != nil {
		return err
	}

	logger.Success(fmt.Sprintf("📑 Copy Marijan to /usr/local/bin for %s (%s).", runtime.GOOS, runtime.GOARCH))

	// Construct the new path for the file in the destination directory
	destinationFilePath := filepath.Join("/usr/local/bin", tun.bin)

	// Move the file using os.Rename
	err = os.Rename(destination, destinationFilePath)
	if err != nil {
		return err
	}

	// remove tarball, and skip error
	_ = os.RemoveAll(source)

	// success message
	tun.successMessage()

	return nil
}

// check if systemd service exist
func (tun *tunnel) IsServiceExist() bool {
	return true
}

func (tun *tunnel) CreateService() error {
	// create file configuration
	err := tun.serviceConfig()
	if err != nil {
		return err
	}

	err = tun.serviceRuntime()
	if err != nil {
		return err
	}

	return nil
}

func (tun *tunnel) serviceConfig() error {
	err := os.MkdirAll(tun.service.folder, 0755)
	if err != nil {
		return err
	}

	var finalPath = filepath.Join(tun.service.folder, tun.service.name)

	// convert to json
	jsonByte, err := json.Marshal(tun.service.configs)
	if err != nil {
		return err
	}

	// Create a new file or truncate an existing one.
	// os.Create returns a *os.File and an error.
	file, err := os.Create(finalPath)
	if err != nil {
		return err
	}
	// Ensure the file is closed when the function exits to prevent resource leaks.
	defer file.Close()

	// Write the string content to the file.
	// file.WriteString returns the number of bytes written and an error.
	_, err = file.WriteString(string(jsonByte))
	if err != nil {
		return err
	}

	return nil
}

func (tun *tunnel) serviceRuntime() error {
	var baseTemplate = `
[Unit]
Description=dPanel Agent name "dpanel-tunnel", version {{.Version}} by devetek.com
Documentation=https://cloud.terpusat.com
After=network-online.target
Wants=network-online.target systemd-networkd-wait-online.service
StartLimitIntervalSec=120
StartLimitBurst=5

[Service]
Restart=always
RestartSec=10s

; User and group the process will run as.
User=root
Group=root

; Service runtime configuration
ExecStart="/usr/local/bin/{{.Bin}}" run --config "{{.Config}}"
ExecReload=/bin/kill -USR2 $MAINPID
ExecStop=/bin/kill -SIGTERM $MAINPID

; systemd extra config

[Install]
WantedBy=multi-user.target
	`

	// 2. Define the data to be used in the template
	data := struct {
		Version string
		Bin     string
		Config  string
	}{
		Version: tun.version,
		Bin:     tun.bin,
		Config:  filepath.Join(tun.service.folder, tun.service.name),
	}

	tmpl, err := template.New("dpanel-tunnel").Parse(baseTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(tun.service.serviceCfg)
	if err != nil {
		return err
	}
	defer file.Close() // Ensure the file is closed

	// 5. Execute the template and write the output to the file
	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	// enable systemd service
	if tun.serviceTrigger("enable") != nil {
		return err
	}

	// start systemd service
	if tun.serviceTrigger("start") != nil {
		return err
	}

	return nil
}

func (tun *tunnel) serviceTrigger(action string) error {
	// trigger systemd service
	words := strings.Split(fmt.Sprintf("%s %s", action, tun.service.serviceName), " ")
	systemdCMD := exec.Command("systemctl", words...)

	// TODO: trap output and stream real-time
	_, err := systemdCMD.Output()
	if err != nil {
		return err
	}

	return nil
}

func (tun *tunnel) successMessage() {
	logger.Success("⭐ If you like Marijan, please give it a star on GitHub: https://github.com/devetek/tuman")
}
