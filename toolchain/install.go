//go:build install

package main

import "bytes"
import "errors"
import "flag"
import "fmt"
import "net"
import "os"
import "os/exec"
import "path/filepath"
import "strconv"
import "strings"
import "time"
import "golang.org/x/crypto/ssh"
import "golang.org/x/crypto/ssh/agent"
import "golang.org/x/crypto/ssh/knownhosts"

type Remote struct {
	Host       string
	User       string
	Home       string
	Hostname   string
	GOOS       string
	GOARCH     string
	DistroID   string
	DistroLike string
	Suffix     string
}

func usage() {

	fmt.Println("Usage:")
	fmt.Println("  go run install.go [flags] <host>")
	fmt.Println("")
	fmt.Println("Installs the hydra binary and a systemd user service, either")
	fmt.Println("locally or on a remote host over SSH. The OS, distribution and")
	fmt.Println("architecture are detected automatically. If <host> is this machine")
	fmt.Println("(matching hostname, IP or localhost) no SSH connection is used and")
	fmt.Println("the role defaults to server.")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  -role client|server   service role (default: client, server when local)")
	fmt.Println("  -controller <name>    controller name (default: this host)")
	fmt.Println("  -address <ip>         controller address for the client (default: this LAN IP)")
	fmt.Println("  -position <pos>       client position (default right-of)")
	fmt.Println("  -port <port>          HYDRA_PORT override (default: built-in 3333)")
	fmt.Println("  -user <name>          SSH user (default: current user)")
	fmt.Println("  -identity <path>      SSH private key (default: agent + ~/.ssh/id_*)")
	fmt.Println("  -ssh-port <port>      SSH port (default 22)")
	fmt.Println("  -verify-host-key      verify against ~/.ssh/known_hosts (default false)")
	fmt.Println("  -display <:n>         DISPLAY for the service (default :0)")
	fmt.Println("  -xauthority <path>    XAUTHORITY for the service (default %h/.Xauthority)")
	fmt.Println("  -service              write/enable the systemd user service (default true)")
	fmt.Println("  -start                start/restart the service (default true)")
	fmt.Println("  -dry-run              print actions without changing anything")

}

func localHostname() string {

	hostname, err := os.Hostname()

	if err != nil {
		return ""
	}

	return strings.ToLower(hostname)

}

func localAddress() string {

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, addr := range addrs {

		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {

			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}

	}

	return ""

}

func addressForHost(host string) string {

	connection, err := net.Dial("udp", net.JoinHostPort(host, "22"))

	if err == nil {

		defer connection.Close()

		if addr, ok := connection.LocalAddr().(*net.UDPAddr); ok && addr.IP != nil {
			return addr.IP.String()
		}

	}

	return localAddress()

}

func localIPs() []net.IP {

	ips := make([]net.IP, 0)

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ips
	}

	for _, addr := range addrs {

		if ipnet, ok := addr.(*net.IPNet); ok {
			ips = append(ips, ipnet.IP)
		}

	}

	return ips

}

func containsIP(ips []net.IP, ip net.IP) bool {

	for _, candidate := range ips {

		if candidate.Equal(ip) {
			return true
		}

	}

	return false

}

func isLocal(host string) bool {

	name := strings.ToLower(strings.TrimSpace(host))

	if name == "" {
		return false
	}

	if name == "localhost" || name == "127.0.0.1" || name == "::1" {
		return true
	}

	if name == localHostname() {
		return true
	}

	ips := localIPs()

	if ip := net.ParseIP(name); ip != nil {
		return containsIP(ips, ip)
	}

	addrs, err := net.LookupHost(name)

	if err == nil {

		for _, addr := range addrs {

			if ip := net.ParseIP(addr); ip != nil && containsIP(ips, ip) {
				return true
			}

		}

	}

	return false

}

func candidateAuths(identity string) ([]ssh.AuthMethod, error) {

	methods := make([]ssh.AuthMethod, 0)

	if socket := os.Getenv("SSH_AUTH_SOCK"); socket != "" {

		connection, err := net.Dial("unix", socket)

		if err == nil {
			client := agent.NewClient(connection)
			methods = append(methods, ssh.PublicKeysCallback(client.Signers))
		}

	}

	home, _ := os.UserHomeDir()
	paths := make([]string, 0)

	if identity != "" {
		paths = append(paths, identity)
	} else {
		for _, name := range []string{"id_ed25519", "id_ecdsa", "id_rsa"} {
			paths = append(paths, filepath.Join(home, ".ssh", name))
		}
	}

	for _, path := range paths {

		data, err := os.ReadFile(path)

		if err != nil {
			continue
		}

		signer, err := ssh.ParsePrivateKey(data)

		if err != nil {
			continue
		}

		if signer.PublicKey().Type() == ssh.KeyAlgoRSA {

			if algorithm_signer, ok := signer.(ssh.AlgorithmSigner); ok {

				upgraded, err := ssh.NewSignerWithAlgorithms(algorithm_signer, []string{
					ssh.KeyAlgoRSASHA512,
					ssh.KeyAlgoRSASHA256,
				})

				if err == nil {
					signer = upgraded
				}

			}

		}

		methods = append(methods, ssh.PublicKeys(signer))

	}

	if len(methods) == 0 {
		return nil, errors.New("no SSH authentication methods available (start an agent or pass -identity)")
	}

	return methods, nil

}

func hostKeyCallback(verify bool) ssh.HostKeyCallback {

	if verify {

		home, err := os.UserHomeDir()

		if err == nil {

			callback, err := knownhosts.New(filepath.Join(home, ".ssh", "known_hosts"))

			if err == nil {
				return callback
			}

		}

	}

	fmt.Fprintln(os.Stderr, "Warning: host key verification disabled (pass -verify-host-key to enable)")

	return ssh.InsecureIgnoreHostKey()

}

func connect(host string, user string, identity string, port int, verify bool) (*ssh.Client, error) {

	methods, err := candidateAuths(identity)

	if err != nil {
		return nil, err
	}

	address := net.JoinHostPort(host, strconv.Itoa(port))
	callback := hostKeyCallback(verify)

	var last_error error = errors.New("no SSH authentication methods available")

	for _, method := range methods {

		config := &ssh.ClientConfig{
			User:            user,
			Auth:            []ssh.AuthMethod{method},
			HostKeyCallback: callback,
			Timeout:         10 * time.Second,
		}

		client, err := ssh.Dial("tcp", address, config)

		if err == nil {
			return client, nil
		}

		last_error = err

	}

	return nil, last_error

}

type Target interface {
	Run(command string) (string, error)
	Upload(local_path string, remote_path string) error
	WriteFile(remote_path string, content string) error
	Close()
}

type sshTarget struct {
	client *ssh.Client
}

func (target *sshTarget) Run(command string) (string, error) {
	return runSSH(target.client, command)
}

func (target *sshTarget) Upload(local_path string, remote_path string) error {
	return uploadSSH(target.client, local_path, remote_path)
}

func (target *sshTarget) WriteFile(remote_path string, content string) error {
	return writeFileSSH(target.client, remote_path, content)
}

func (target *sshTarget) Close() {
	target.client.Close()
}

type localTarget struct{}

func (target localTarget) Run(command string) (string, error) {

	process := exec.Command("sh", "-c", command)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	process.Stdout = &stdout
	process.Stderr = &stderr

	err := process.Run()

	output := stdout.String()

	if stderr.Len() > 0 {

		if output != "" {
			output += "\n"
		}

		output += stderr.String()

	}

	return strings.TrimSpace(output), err

}

func (target localTarget) Upload(local_path string, remote_path string) error {

	data, err := os.ReadFile(local_path)

	if err != nil {
		return err
	}

	path := os.ExpandEnv(remote_path)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	tmp := path + ".tmp"

	if err := os.WriteFile(tmp, data, 0755); err != nil {
		return err
	}

	return os.Rename(tmp, path)

}

func (target localTarget) WriteFile(remote_path string, content string) error {

	path := os.ExpandEnv(remote_path)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)

}

func (target localTarget) Close() {}

func runSSH(client *ssh.Client, command string) (string, error) {

	session, err := client.NewSession()

	if err != nil {
		return "", err
	}

	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(command)

	output := stdout.String()

	if stderr.Len() > 0 {

		if output != "" {
			output += "\n"
		}

		output += stderr.String()

	}

	return strings.TrimSpace(output), err

}

func detect(target Target, host string) (*Remote, error) {

	uname, err := target.Run("uname -s")

	if err != nil {
		return nil, fmt.Errorf("uname failed: %w", err)
	}

	arch, _ := target.Run("uname -m")
	os_release, _ := target.Run("cat /etc/os-release 2>/dev/null")
	home, _ := target.Run("printf %s \"$HOME\"")
	user, _ := target.Run("id -un")
	hostname, _ := target.Run("hostname -s 2>/dev/null || hostname")

	id, id_like := parseOSRelease(os_release)

	return &Remote{
		Host:       host,
		User:       user,
		Home:       home,
		Hostname:   strings.ToLower(hostname),
		GOOS:       strings.ToLower(uname),
		GOARCH:     mapArch(arch),
		DistroID:   id,
		DistroLike: id_like,
	}, nil

}

func parseOSRelease(content string) (string, string) {

	id := ""
	id_like := ""

	for _, line := range strings.Split(content, "\n") {

		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "ID=") {
			id = strings.Trim(strings.TrimPrefix(line, "ID="), "\"")
		} else if strings.HasPrefix(line, "ID_LIKE=") {
			id_like = strings.Trim(strings.TrimPrefix(line, "ID_LIKE="), "\"")
		}

	}

	return strings.ToLower(id), strings.ToLower(id_like)

}

func mapArch(uname string) string {

	switch strings.ToLower(strings.TrimSpace(uname)) {

	case "x86_64", "amd64":
		return "amd64"

	case "aarch64", "arm64":
		return "arm64"

	case "armv7l", "armv7":
		return "arm"

	case "armv6l":
		return "arm"

	case "i386", "i686", "x86":
		return "386"

	case "riscv64":
		return "riscv64"

	case "ppc64le":
		return "ppc64le"

	case "s390x":
		return "s390x"

	default:
		return strings.ToLower(strings.TrimSpace(uname))

	}

}

func mapDistro(remote *Remote) (string, error) {

	id := remote.DistroID
	like := remote.DistroLike

	if id == "arch" || id == "archlinux" || strings.Contains(like, "arch") {
		return "archlinux", nil
	}

	if id == "debian" || id == "ubuntu" || id == "linuxmint" || id == "pop" ||
		id == "raspbian" || strings.Contains(like, "debian") || strings.Contains(like, "ubuntu") {
		return "debian", nil
	}

	return "", fmt.Errorf("unsupported distribution %q (ID_LIKE=%q); supported: debian, ubuntu, arch", id, like)

}

func ensureBinary(repo string, remote *Remote) (string, error) {

	name := fmt.Sprintf("hydra-%s-%s-%s", remote.Suffix, remote.GOOS, remote.GOARCH)
	path := filepath.Join(repo, "build", name)

	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	fmt.Printf("==> %s missing, building with make %s GOOS=%s GOARCH=%s\n", name, remote.Suffix, remote.GOOS, remote.GOARCH)

	command := exec.Command("make", remote.Suffix, "GOOS="+remote.GOOS, "GOARCH="+remote.GOARCH)
	command.Dir = repo
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return "", fmt.Errorf("make %s failed: %w", remote.Suffix, err)
	}

	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("expected binary at %s after build", path)
	}

	return path, nil

}

func uploadSSH(client *ssh.Client, local_path string, remote_path string) error {

	data, err := os.ReadFile(local_path)

	if err != nil {
		return err
	}

	session, err := client.NewSession()

	if err != nil {
		return err
	}

	defer session.Close()

	session.Stdin = bytes.NewReader(data)

	command := fmt.Sprintf("mkdir -p \"$(dirname %s)\" && cat > %s.tmp && chmod 0755 %s.tmp && mv -f %s.tmp %s",
		remote_path, remote_path, remote_path, remote_path, remote_path)

	return session.Run(command)

}

func writeFileSSH(client *ssh.Client, remote_path string, content string) error {

	session, err := client.NewSession()

	if err != nil {
		return err
	}

	defer session.Close()

	session.Stdin = strings.NewReader(content)

	command := fmt.Sprintf("mkdir -p \"$(dirname %s)\" && cat > %s", remote_path, remote_path)

	return session.Run(command)

}

func environmentLines(display string, xauthority string, port string, controller string) string {

	lines := []string{
		"Environment=DISPLAY=" + display,
		"Environment=XAUTHORITY=" + xauthority,
	}

	if port != "" {
		lines = append(lines, "Environment=HYDRA_PORT="+port)
	}

	if controller != "" {
		lines = append(lines, "Environment=HYDRA_CONTROLLER="+controller)
	}

	return strings.Join(lines, "\n")

}

func serverUnit(remote *Remote, display string, xauthority string, port string) string {

	return fmt.Sprintf(`[Unit]
Description=Hydra controller (%s)
Documentation=https://github.com/cookiengineer/hydra
After=graphical-session.target
PartOf=graphical-session.target

[Service]
Type=simple
%s
ExecStart=%%h/.local/bin/hydra listen %s
Restart=on-failure
RestartSec=2

[Install]
WantedBy=graphical-session.target
`, remote.Hostname, environmentLines(display, xauthority, port, ""), remote.Hostname)

}

func clientUnit(remote *Remote, address string, position string, controller string, display string, xauthority string, port string) string {

	return fmt.Sprintf(`[Unit]
Description=Hydra client (%s)
Documentation=https://github.com/cookiengineer/hydra
After=graphical-session.target
PartOf=graphical-session.target

[Service]
Type=simple
%s
ExecStart=%%h/.local/bin/hydra connect %s %s
Restart=always
RestartSec=3

[Install]
WantedBy=graphical-session.target
`, remote.Hostname, environmentLines(display, xauthority, port, controller), position, address)

}

func repositoryRoot() (string, error) {

	dir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	for {

		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)

		if parent == dir {
			return "", errors.New("could not locate repository root (go.mod)")
		}

		dir = parent

	}

}

func main() {

	flags := flag.NewFlagSet("install", flag.ExitOnError)
	role := flags.String("role", "client", "service role (client|server)")
	controller := flags.String("controller", "", "controller name")
	address := flags.String("address", "", "controller address")
	position := flags.String("position", "right-of", "client position")
	port := flags.String("port", "", "HYDRA_PORT override")
	user := flags.String("user", "", "SSH user")
	identity := flags.String("identity", "", "SSH private key")
	ssh_port := flags.Int("ssh-port", 22, "SSH port")
	verify_host_key := flags.Bool("verify-host-key", false, "verify host key against known_hosts")
	display := flags.String("display", ":0", "DISPLAY")
	xauthority := flags.String("xauthority", "%h/.Xauthority", "XAUTHORITY")
	service := flags.Bool("service", true, "install systemd user service")
	start := flags.Bool("start", true, "start/restart the service")
	dry_run := flags.Bool("dry-run", false, "print actions without changing the remote")

	flags.Usage = usage
	flags.Parse(os.Args[1:])

	arguments := flags.Args()

	if len(arguments) != 1 {
		usage()
		os.Exit(2)
	}

	host := arguments[0]

	role_set := false

	flags.Visit(func(current *flag.Flag) {
		if current.Name == "role" {
			role_set = true
		}
	})

	local := isLocal(host)

	if !role_set && local {
		*role = "server"
	}

	if *role != "client" && *role != "server" {
		fmt.Fprintln(os.Stderr, "Error: -role must be client or server")
		os.Exit(2)
	}

	ssh_user := *user

	if ssh_user == "" {
		ssh_user = os.Getenv("USER")

		if ssh_user == "" {
			ssh_user, _ = os.UserHomeDir()
		}

	}

	controller_name := *controller

	if controller_name == "" {
		controller_name = localHostname()
	}

	controller_address := *address

	if controller_address == "" {
		controller_address = addressForHost(host)
	}

	var target Target

	if local {

		fmt.Printf("==> %s is this machine; installing locally (role: %s)\n", host, *role)
		target = localTarget{}

	} else {

		fmt.Printf("==> Connecting to %s as %s (role: %s)\n", host, ssh_user, *role)

		client, err := connect(host, ssh_user, *identity, *ssh_port, *verify_host_key)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: SSH connection failed:", err.Error())
			os.Exit(1)
		}

		target = &sshTarget{client: client}

	}

	defer target.Close()

	remote, err := detect(target, host)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: detection failed:", err.Error())
		os.Exit(1)
	}

	suffix, err := mapDistro(remote)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}

	remote.Suffix = suffix

	fmt.Printf("==> Detected: %s %s / %s / %s\n", remote.DistroID, remote.GOOS, remote.GOARCH, remote.Hostname)

	repo, err := repositoryRoot()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}

	var unit string

	if *role == "server" {

		listen_name := controller_name

		if *controller == "" {
			listen_name = remote.Hostname
		}

		unit = serverUnit(remote, *display, *xauthority, *port)

		if listen_name != remote.Hostname {
			unit = strings.Replace(unit, "listen "+remote.Hostname, "listen "+listen_name, 1)
		}

	} else {

		if controller_address == "" {
			fmt.Fprintln(os.Stderr, "Error: could not determine controller address; pass -address")
			os.Exit(1)
		}

		unit = clientUnit(remote, controller_address, *position, controller_name, *display, *xauthority, *port)

	}

	if *dry_run {

		fmt.Println("")
		fmt.Println("==> DRY RUN")
		fmt.Printf("    binary:  build/hydra-%s-%s-%s\n", remote.Suffix, remote.GOOS, remote.GOARCH)
		fmt.Printf("    remote:  %s/.local/bin/hydra\n", remote.Home)
		fmt.Printf("    service: %s/.config/systemd/user/hydra.service\n", remote.Home)
		fmt.Println("")
		fmt.Println(unit)
		return

	}

	binary_path, err := ensureBinary(repo, remote)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		os.Exit(1)
	}

	fmt.Printf("==> Uploading %s\n", filepath.Base(binary_path))

	if err := target.Upload(binary_path, "$HOME/.local/bin/hydra"); err != nil {
		fmt.Fprintln(os.Stderr, "Error: upload failed:", err.Error())
		os.Exit(1)
	}

	if *service {

		fmt.Println("==> Writing systemd user service")

		if err := target.WriteFile("$HOME/.config/systemd/user/hydra.service", unit); err != nil {
			fmt.Fprintln(os.Stderr, "Error: writing service failed:", err.Error())
			os.Exit(1)
		}

		if *start {

			fmt.Println("==> Enabling and starting hydra.service")

			command := "export XDG_RUNTIME_DIR=\"/run/user/$(id -u)\"; " +
				"systemctl --user daemon-reload && " +
				"systemctl --user enable hydra.service && " +
				"systemctl --user restart hydra.service; " +
				"systemctl --user is-enabled hydra.service; " +
				"systemctl --user is-active hydra.service || true"

			output, err := target.Run(command)

			fmt.Print(output)
			fmt.Println("")

			if err != nil {
				fmt.Fprintln(os.Stderr, "Error: systemctl failed:", err.Error())
				os.Exit(1)
			}

		}

	}

	fmt.Println("==> Done")

}
