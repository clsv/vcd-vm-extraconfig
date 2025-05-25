# vcd-vm-extraconfig

A command-line utility to manage ExtraConfig key-value pairs for virtual machines (VMs) in VMware vCloud Director. This tool allows you to find VMs, retrieve VM information (including ExtraConfig), set ExtraConfig entries, and delete ExtraConfig entries using the `go-vcloud-director` SDK.

## Overview

The `vcd-vm-extraconfig` utility is designed for vCloud Director administrators and automation engineers who need to programmatically manage ExtraConfig settings for VMs. ExtraConfig entries (e.g., `guestinfo.mykey`) are stored in the VM's virtual hardware section and are accessible inside the VM via VMware Tools. The tool supports vCloud Director API version 36.3 and is built with the `go-vcloud-director/v2@main` library.

### Features
- **Find VMs**: Locate a VM by organization, VDC, vApp, and VM name.
- **Get VM Info**: Retrieve VM details, including CPU, memory, and ExtraConfig entries.
- **Set ExtraConfig**: Add or update ExtraConfig key-value pairs.
- **Delete ExtraConfig**: Remove ExtraConfig entries by key.
- **Configuration**: Uses a JSON configuration file for vCloud Director credentials and endpoint.
- **Docker Support**: Can be built and run as a Docker container.

## Prerequisites
- **Go**: Version 1.23 or later (1.24.2 recommended for building locally).
- **Docker**: Required for building and running the Docker image.
- **vCloud Director**: Access to a vCloud Director instance (API version 36.3 recommended).
- **Permissions**: The user must have `Virtual Machine: View` and `Virtual Machine: Edit` rights.
- **Git**: For cloning the repository.

## Installation
1. **Clone the repository**:
   ```bash
   git clone https://github.com/clsv/vcd-vm-extraconfig
   cd vcd-vm-extraconfig
   ```

2. **Install dependencies**:
   Install the `go-vcloud-director` library (main branch):
   ```bash
   go mod download
   ```

3. **Build the tool**:
   ```bash
   go build -o vcd-vm-extraconfig vcd-vm-extraconfig.go
   ```

## Docker Installation
1. **Build the Docker image**:
   ```bash
   docker build -t vcd-vm-extraconfig:latest .
   ```

2. **Run the tool in a container**:
   Mount a `config.json` file and run the desired command (see Usage below). Example:
   ```bash
   docker run --rm -v $(pwd)/config.json:/app/config.json vcd-vm-extraconfig:latest -config /app/config.json -action get -org my-org -vdc my-vdc -vapp my-vapp -vm my-vm
   ```

## Configuration
Create a `config.json` file in the repository root with your vCloud Director credentials:

```json
{
    "api": "36.3", //37.2, 37.1, 37.0, 36.3, 36.2, 36.1, 36.0
    "url": "https://vcd.example.com/api",
    "org": "your-org",
    "vdc": "vdc",
    "vapp": "vapp",
    "user": "your-user",
    "password": "your-password"
}
```

Replace the placeholder values with your vCloud Director endpoint, organization, username, and password.

## Usage
The tool supports four actions: `find`, `get`, `set`, and `delete`. Run the tool with the appropriate flags to perform these actions.

### Command Syntax
```bash
./vcd-vm-extraconfig -config config.json -action <action> -org <org> -vdc <vdc> -vapp <vapp> -vm <vm> [-key <key>] [-value <value>]
```

### Flags
- `-config`: Path to the configuration file (default: `config.json`).
- `-action`: Action to perform (`find`, `get`, `set`, or `delete`).
- `-org`: Organization name.
- `-vdc`: Virtual Data Center name.
- `-vapp`: vApp name.
- `-vm`: VM name.
- `-key`: ExtraConfig key (required for `set` and `delete`).
- `-value`: ExtraConfig value (required for `set`).

### Examples
1. **Find a VM**:
   ```bash
   ./vcd-vm-extraconfig -config config.json -action find -org my-org -vdc my-vdc -vapp my-vapp -vm my-vm
   ```
   Output:
   ```
   VM found: my-vm
   ```

2. **Get VM Information (including ExtraConfig)**:
   ```bash
   ./vcd-vm-extraconfig -config config.json -action get -org my-org -vdc my-vdc -vapp my-vapp -vm my-vm
   ```
   Output:
   ```
   VM Name: my-vm
   CPU: 2
   Memory: 4096 MB
   ExtraConfig:
     guestinfo.mykey: myvalue
   ```

3. **Set an ExtraConfig Entry**:
   ```bash
   ./vcd-vm-extraconfig -config config.json -action set -org my-org -vdc my-vdc -vapp my-vapp -vm my-vm -key disk.enableUUID -value 1
   ```
   Output:
   ```
   ExtraConfig disk.enableUUID=1 set
   ```

4. **Delete an ExtraConfig Entry**:
   ```bash
   ./vcd-vm-extraconfig -config config.json -action delete -org my-org -vdc my-vdc -vapp my-vapp -vm my-vm -key guestinfo.mykey
   ```
   Output:
   ```
   ExtraConfig guestinfo.mykey deleted
   ```

### Docker Examples
Run the same commands using Docker by mounting the `config.json` file:
```bash
docker run --rm -v $(pwd)/config.json:/app/config.json vcd-vm-extraconfig:latest -config /app/config.json -action set -org my-org -vdc my-vdc -vapp my-vapp -vm my-vm -key guestinfo.mykey -value myvalue
```

## Logging
The tool supports debug logging via environment variables provided by the `go-vcloud-director` library. Set these variables to enable detailed output for troubleshooting:
- `GOVCD_SHOW_RESP`: Set to `true` to log vCloud Director API responses.
  ```bash
  export GOVCD_SHOW_RESP=true
  export GOVCD_LOG=true
  ./vcd-vm-extraconfig -config config.json -action get -org my-org -vdc my-vdc -vapp my-vapp -vm my-vm
  ```

- `GOVCD_SHOW_REQ`: Set to `true` to log vCloud Director API requests.
  ```bash
  export GOVCD_SHOW_REQ=true
  export GOVCD_LOG=true
  ./vcd-vm-extraconfig -config config.json -action get -org my-org -vdc my-vdc -vapp my-vapp -vm my-vm
  ```

For Docker, pass the environment variables using the `-e` flag:
```bash
docker run --rm -v $(pwd)/config.json:/app/config.json -e GOVCD_SHOW_RESP=true -e GOVCD_LOG=true -e GOVCD_SHOW_REQ=true vcd-vm-extraconfig:latest -config /app/config.json -action get -org my-org -vdc my-vdc -vapp my-vapp -vm my-vm
```

## Notes
- **ExtraConfig Keys**: Keys like `guestinfo.mykey` are accessible inside the VM via VMware Tools (e.g., `vmtoolsd --cmd "info-get guestinfo.mykey"`). Avoid using spaces in keys, as they are invalid.
- **Library Version**: The tool uses the `main` branch of `go-vcloud-director`, which may include unstable changes. For a stable release, replace `go get github.com/vmware/go-vcloud-director/v2@main` with a specific version (e.g., `v2.26.1`).
- **Error Handling**: The tool provides basic error messages. For production use, consider adding logging (e.g., with `logrus`).
- **Cloud-Init Integration**: To pass ExtraConfig to cloud-init, use `vm.SetGuestCustomizationSection`. Open an issue or contribute to add this feature.
- **Docker**: The Docker image is built with Go 1.24.2 and uses a minimal `scratch` base for the final image, reducing size.
