package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func Create() {
	numModules := 10000

	modules := [][]string{}

	for i := 0; i < numModules; i++ {
		modules = append(modules, []string{
			"module \"bas-" + strconv.Itoa(i+1) + "\"" + " {\n" +
				"\tsource               		= \"./modules/clearblade-registry\"\n" +
				"\tregistry_id                  = \"bas-" + strconv.Itoa(i+1) + "\"\n" +
				"\tdevice_credentials           = local.device_credentials \n" +
				"\tgcp_project                  = local.gcp_project\n" +
				"\tgcp_region                   = local.gcp_region\n" +
				"\tdevice_id                    = \"your-device-id\"\n" +
				"\tdevice_log_level             = \"DEBUG\"\n" +
				"\tdevice_blocked               = false \n" +
				"\tdevice_metadata_location     = \"device-location\"\n" +
				"\tdevice_metadata_manufacturer = \"manufacturer-name\"\n" +
				"\tdevice_gateway_type          = \"NON_GATEWAY\"\n" +
				"\tdevice_gateway_auth_method   = \"ASSOCIATION_AND_DEVICE_AUTH_TOKEN\"\n" +
				"}",
		})
	}

	file, err := os.OpenFile("main.tf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	for _, module := range modules {
		for _, moduleblock := range module {
			_, _ = datawriter.WriteString("\n" + moduleblock + "\n")
		}
	}

	datawriter.Flush()
	file.Close()
}

func Delete() {
	if err := os.Truncate("main.tf", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}
}

func main() {
	Create()
	// Comment the function call above and uncomment Delete below to delete resources
	// Delete()
}
