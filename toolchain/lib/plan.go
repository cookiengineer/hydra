package lib

type Plan struct {
	Repo            string
	Host            string
	Image           string
	RemoteDir       string
	RemoteBin       string
	Controller      string
	Client          string
	Address         string
	Position        string
	TestFilter      string
	GuestDisplay    string
	GuestXauthority string
}

func (plan Plan) LocalBinary() string {
	return plan.Repo + "/build/hydra-e2e.test"
}

func (plan Plan) ImageExistsCommand() Command {

	return Command{
		Name: "podman",
		Args: []string{"image", "exists", plan.Image},
	}

}

func (plan Plan) BuildImageCommand() Command {

	return Command{
		Name: "podman",
		Args: []string{"build", "-t", plan.Image, "-f", plan.Repo + "/debian.Containerfile", plan.Repo},
	}

}

func (plan Plan) BuildTestCommand() Command {

	return Command{
		Name: "podman",
		Args: []string{
			"run", "--rm",
			"-v", plan.Repo + ":/src",
			"-w", "/src",
			plan.Image,
			"go", "test", "-c", "-tags=e2e",
			"-o", "/src/build/hydra-e2e.test",
			"./toolchain/e2e",
		},
	}

}

func (plan Plan) RemoteMkdirCommand() Command {

	return Command{
		Name: "ssh",
		Args: []string{plan.Host, "mkdir", "-p", plan.RemoteDir},
	}

}

func (plan Plan) SCPCommand() Command {

	return Command{
		Name: "scp",
		Args: []string{plan.LocalBinary(), plan.Host + ":" + plan.RemoteBin},
	}

}

func (plan Plan) LocalServerCommand() Command {

	return Command{
		Name: plan.LocalBinary(),
		Args: []string{"-test.v", "-test.run", plan.TestFilter},
		Env: []string{
			"HYDRA_E2E_ROLE=server",
			"HYDRA_E2E_CONTROLLER=" + plan.Controller,
			"HYDRA_E2E_CONTROLLER_ADDR=" + plan.Address,
			"HYDRA_E2E_CLIENT=" + plan.Client,
		},
	}

}

func (plan Plan) RemoteClientCommand() Command {

	return Command{
		Name: "ssh",
		Args: []string{
			plan.Host,
			"env",
			"DISPLAY=" + plan.GuestDisplay,
			"XAUTHORITY=" + plan.GuestXauthority,
			"HYDRA_CONTROLLER=" + plan.Controller,
			"HYDRA_E2E_ROLE=client",
			"HYDRA_E2E_CONTROLLER=" + plan.Controller,
			"HYDRA_E2E_CONTROLLER_ADDR=" + plan.Address,
			"HYDRA_E2E_CLIENT=" + plan.Client,
			plan.RemoteBin,
			"-test.v",
			"-test.run", plan.TestFilter,
		},
	}

}
