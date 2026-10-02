package main

import "bytes"
import "flag"
import "fmt"
import "net"
import "os"
import "os/exec"
import "time"
import "github.com/cookiengineer/hydra/toolchain/lib"

func usage() {

	fmt.Println("Usage:")
	fmt.Println("  go run test.go <host> <client|server> [flags]")
	fmt.Println("")
	fmt.Println("  client   build the e2e test binary and run it on <host> (remote)")
	fmt.Println("  server   build, deploy the client, and orchestrate both ends")
	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Println("  -controller <name>      controller hostname (default: this host)")
	fmt.Println("  -address <ip>           address the client uses to reach the controller")
	fmt.Println("  -position <pos>         client position (default: right-of)")
	fmt.Println("  -test <regexp>          go test -run filter (default: TestE2E)")
	fmt.Println("  -image <name>           builder image (default: hydra-e2e-builder)")
	fmt.Println("  -display <:n>           client DISPLAY (default: :0)")
	fmt.Println("  -xauthority <path>      client XAUTHORITY (default: $HOME/.Xauthority)")

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

func ensureImage(executor lib.Executor, plan lib.Plan) bool {

	result := executor.Run(plan.ImageExistsCommand())

	if result.ExitCode == 0 {
		return true
	}

	fmt.Println("==> Building builder image:", plan.Image)
	build := executor.Run(plan.BuildImageCommand())

	fmt.Print(build.Stdout)
	fmt.Print(build.Stderr)

	return build.ExitCode == 0

}

func deploy(executor lib.Executor, plan lib.Plan) bool {

	mkdir := executor.Run(plan.RemoteMkdirCommand())

	if mkdir.ExitCode != 0 {
		fmt.Print(mkdir.Stderr)
		return false
	}

	copy := executor.Run(plan.SCPCommand())

	if copy.ExitCode != 0 {
		fmt.Print(copy.Stderr)
		return false
	}

	return true

}

func runClient(executor lib.Executor, plan lib.Plan) (lib.Summary, lib.Summary, int) {

	result := executor.Run(plan.RemoteClientCommand())

	fmt.Print(result.Stdout)
	fmt.Print(result.Stderr)

	summary := lib.Parse(result.Stdout + result.Stderr)

	return summary, summary, result.ExitCode

}

func runServer(executor lib.Executor, plan lib.Plan) int {

	if !deploy(executor, plan) {
		return 1
	}

	command := plan.LocalServerCommand()

	cmd := exec.Command(command.Name, command.Args...)
	cmd.Env = append(os.Environ(), command.Env...)

	var server_buffer bytes.Buffer

	cmd.Stdout = &server_buffer
	cmd.Stderr = &server_buffer

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to start local server test:", err.Error())
		return 1
	}

	time.Sleep(2 * time.Second)

	client_result := executor.Run(plan.RemoteClientCommand())

	fmt.Print(client_result.Stdout)
	fmt.Print(client_result.Stderr)

	done := make(chan error, 1)

	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-done:
	case <-time.After(45 * time.Second):
		cmd.Process.Kill()
		<-done
	}

	server_output := server_buffer.String()

	fmt.Print(server_output)

	server_summary := lib.Parse(server_output)
	client_summary := lib.Parse(client_result.Stdout + client_result.Stderr)

	fmt.Println("")
	fmt.Println("==> Server results:")
	report(server_summary)
	fmt.Println("==> Client results:")
	report(client_summary)

	if server_summary.OK() && client_summary.OK() && client_result.ExitCode == 0 {
		return 0
	}

	return 1

}

func report(summary lib.Summary) {

	if len(summary.Results) == 0 && summary.Passed == 0 && summary.Failed == 0 {
		fmt.Println("    (no results parsed)")
	}

	for _, line := range summary.Lines {
		fmt.Println("   ", line)
	}

	for _, result := range summary.Results {
		fmt.Printf("    [%s] %s: %s\n", result.Role, result.Name, result.Message)
	}

}

func main() {

	args := os.Args[1:]

	if len(args) < 2 {
		usage()
		os.Exit(2)
	}

	host := args[0]
	mode := args[1]

	flags := flag.NewFlagSet("test", flag.ExitOnError)
	controller := flags.String("controller", "", "controller hostname")
	address := flags.String("address", "", "controller address")
	position := flags.String("position", "right-of", "client position")
	filter := flags.String("test", "TestE2E", "go test -run filter")
	image := flags.String("image", "hydra-e2e-builder", "builder image")
	display := flags.String("display", ":0", "client DISPLAY")
	xauthority := flags.String("xauthority", "$HOME/.Xauthority", "client XAUTHORITY")

	flags.Parse(args[2:])

	repo, err := lib.RepoRoot()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not locate repository root:", err.Error())
		os.Exit(1)
	}

	controller_name := *controller

	if controller_name == "" {
		hostname, err := os.Hostname()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not determine hostname:", err.Error())
			os.Exit(1)
		}

		controller_name = hostname
	}

	controller_address := *address

	if controller_address == "" {
		controller_address = localAddress()
	}

	plan := lib.Plan{
		Repo:            repo,
		Host:            host,
		Image:           *image,
		RemoteDir:       "/tmp/hydra-e2e",
		RemoteBin:       "/tmp/hydra-e2e/hydra-e2e.test",
		Controller:      controller_name,
		Client:          host,
		Address:         controller_address,
		Position:        *position,
		TestFilter:      *filter,
		GuestDisplay:    *display,
		GuestXauthority: *xauthority,
	}

	executor := lib.OSExecutor{}

	if err := os.MkdirAll(repo+"/build", 0755); err != nil {
		fmt.Fprintln(os.Stderr, "Could not create build directory:", err.Error())
		os.Exit(1)
	}

	if !ensureImage(executor, plan) {
		fmt.Fprintln(os.Stderr, "Failed to build builder image")
		os.Exit(1)
	}

	fmt.Println("==> Building e2e test binary")
	build := executor.Run(plan.BuildTestCommand())

	fmt.Print(build.Stdout)
	fmt.Print(build.Stderr)

	if build.ExitCode != 0 {
		fmt.Fprintln(os.Stderr, "Failed to build e2e test binary")
		os.Exit(1)
	}

	switch mode {

	case "client":

		if !deploy(executor, plan) {
			os.Exit(1)
		}

		summary, _, code := runClient(executor, plan)

		fmt.Println("")
		fmt.Println("==> Client results:")
		report(summary)

		os.Exit(code)

	case "server":

		os.Exit(runServer(executor, plan))

	default:

		usage()
		os.Exit(2)

	}

}
