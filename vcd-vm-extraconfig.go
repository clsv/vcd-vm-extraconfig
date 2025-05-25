package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
)

// Config structure for the configuration file
type Config struct {
	Api      string `json:"api"`
	URL      string `json:"url"`
	Org      string `json:"org"`
	Vdc      string `json:"vdc"`
	Vapp     string `json:"vapp"`
	User     string `json:"user"`
	Password string `json:"password"`
}

var config Config

// loadConfig loads configuration from a file
func loadConfig(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &config)
	return err
}

// connectToVCD connects to vCloud Director
func connectToVCD(config Config) (*govcd.VCDClient, error) {
	u, err := url.ParseRequestURI(config.URL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %v", err)
	}
	client := govcd.NewVCDClient(*u, true, govcd.WithAPIVersion(config.Api))
	err = client.Authenticate(config.User, config.Password, config.Org)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// findVM finds a VM by name
func findVM(client *govcd.VCDClient, orgName, vdcName, vappName, vmName string) (*govcd.VM, error) {
	org, err := client.GetOrgByName(orgName)
	if err != nil {
		fmt.Printf("Error GetOrgByName")
		return nil, err
	}
	vdc, err := org.GetVDCByName(vdcName, false)
	if err != nil {
		fmt.Printf("Error GetVDCByName")
		return nil, err
	}
	vapp, err := vdc.GetVAppByName(vappName, false)
	if err != nil {
		fmt.Printf("Error GetVAppByName")
		return nil, err
	}
	vm, err := vapp.GetVMByNameOrId(vmName, false)
	if err != nil {
		fmt.Printf("Error GetVMByNameOrId")
		return nil, err
	}
	return vm, nil
}

// getVMInfo retrieves VM parameters, including ExtraConfig
func getVMInfo(vm *govcd.VM) error {
	hardwareSection, err := vm.GetVirtualHardwareSection()
	if err != nil {
		return err
	}
	fmt.Printf("VM Name: %s\n", vm.VM.Name)
	for _, item := range hardwareSection.Item {
		if item.ResourceType == types.ResourceTypeProcessor {
			fmt.Printf("CPU: %d\n", item.VirtualQuantity)
		}
		if item.ResourceType == types.ResourceTypeMemory {
			fmt.Printf("Memory: %d MB\n", item.VirtualQuantity)
		}
	}

	// Fetch ExtraConfig
	extraConfig, err := vm.GetExtraConfig()
	if err != nil {
		return fmt.Errorf("error fetching extra config: %v", err)
	}
	fmt.Println("ExtraConfig:")
	for _, ec := range extraConfig {
		fmt.Printf("  %s: %s\n", ec.Key, ec.Value)
	}
	return nil
}

// setExtraConfig sets an ExtraConfig key-value pair
func setExtraConfig(vm *govcd.VM, key, value string) error {
	extraConfig := []*types.ExtraConfigMarshal{{
		Key:   key,
		Value: value,
	}}
	_, err := vm.UpdateExtraConfig(extraConfig)
	if err != nil {
		return fmt.Errorf("error setting extra config: %v", err)
	}
	return nil
}

// deleteExtraConfig deletes an ExtraConfig by key
func deleteExtraConfig(vm *govcd.VM, key string) error {
	extraConfig := []*types.ExtraConfigMarshal{{
		Key: key,
	}}
	_, err := vm.DeleteExtraConfig(extraConfig)
	if err != nil {
		return fmt.Errorf("error deleting extra config: %v", err)
	}
	return nil
}

func main() {
	// Define command-line flags
	configFile := flag.String("config", "config.json", "Path to configuration file")
	action := flag.String("action", "", "Action: find, get, set, delete")
	org := flag.String("org", "", "Organization name")
	vdc := flag.String("vdc", "", "VDC name")
	vapp := flag.String("vapp", "", "vApp name")
	vmName := flag.String("vm", "", "VM name")
	key := flag.String("key", "", "ExtraConfig key (for set/delete)")
	value := flag.String("value", "", "ExtraConfig value (for set)")
	flag.Parse()

	// Check required arguments
	if *action == "" {
		fmt.Println("Error: action must be specified (-action)")
		flag.Usage()
		os.Exit(1)
	}

	// Load configuration
	err := loadConfig(*configFile)
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	if *vdc == "" {
		vdc = &config.Vdc
	}

	if *vapp == "" {
		vapp = &config.Vapp
	}

	if *org == "" {
		org = &config.Org
	}

	// Connect to vCD
	client, err := connectToVCD(config)
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		os.Exit(1)
	}
	defer client.Disconnect()

	// Perform action
	switch *action {
	case "find":
		if *org == "" || *vdc == "" || *vapp == "" || *vmName == "" {
			fmt.Println("Error: must specify -vapp, -vm")
			os.Exit(1)
		}
		vm, err := findVM(client, *org, *vdc, *vapp, *vmName)
		if err != nil {
			fmt.Printf("Error finding VM: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("VM found: %s\n", vm.VM.Name)

	case "get":
		if *org == "" || *vdc == "" || *vapp == "" || *vmName == "" {
			fmt.Println("Error: must specify -vapp, -vm")
			os.Exit(1)
		}
		vm, err := findVM(client, *org, *vdc, *vapp, *vmName)
		if err != nil {
			fmt.Printf("Error finding VM: %v\n", err)
			os.Exit(1)
		}
		if err := getVMInfo(vm); err != nil {
			fmt.Printf("Error getting parameters: %v\n", err)
			os.Exit(1)
		}

	case "set":
		if *org == "" || *vdc == "" || *vapp == "" || *vmName == "" || *key == "" || *value == "" {
			fmt.Println("Error: must specify -vapp, -vm, -key, -value")
			os.Exit(1)
		}
		vm, err := findVM(client, *org, *vdc, *vapp, *vmName)
		if err != nil {
			fmt.Printf("Error finding VM: %v\n", err)
			os.Exit(1)
		}
		if err := setExtraConfig(vm, *key, *value); err != nil {
			fmt.Printf("Error setting ExtraConfig: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("ExtraConfig %s=%s set\n", *key, *value)

	case "delete":
		if *org == "" || *vdc == "" || *vapp == "" || *vmName == "" || *key == "" {
			fmt.Println("Error: must specify -vapp, -vm, -key")
			os.Exit(1)
		}
		vm, err := findVM(client, *org, *vdc, *vapp, *vmName)
		if err != nil {
			fmt.Printf("Error finding VM: %v\n", err)
			os.Exit(1)
		}
		if err := deleteExtraConfig(vm, *key); err != nil {
			fmt.Printf("Error deleting ExtraConfig: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("ExtraConfig %s deleted\n", *key)

	default:
		fmt.Println("Error: unknown action. Available actions: find, get, set, delete")
		flag.Usage()
		os.Exit(1)
	}
}
