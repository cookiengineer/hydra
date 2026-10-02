package types

import "os"

const DefaultPort = "3333"

func Port() string {

	port := os.Getenv("HYDRA_PORT")

	if port == "" {
		return DefaultPort
	}

	return port

}
