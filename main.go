package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// limastart handles starting or stopping a Lima instance with the given name and configuration file.
func limastart(name, config string, stop bool) {
	if stop {
		fmt.Printf("Stopping Lima instance '%s'...\n", name)
		stopCmd := exec.Command("limactl", "stop", name)
		stopCmd.Stdout = os.Stdout
		stopCmd.Stderr = os.Stderr
		if err := stopCmd.Run(); err != nil {
			log.Fatalf("Failed to stop Lima instance '%s': %v", name, err)
		}
		fmt.Printf("Lima instance '%s' stopped successfully.\n", name)
	} else {
		fmt.Printf("Starting Lima instance '%s'...\n", name)
		startCmd := exec.Command("limactl", "start", "--debug", "--name="+name, "--tty=false", config)
		startCmd.Stdout = os.Stdout
		startCmd.Stderr = os.Stderr
		if err := startCmd.Run(); err != nil {
			log.Fatalf("Failed to start Lima instance '%s': %v", name, err)
		}
		fmt.Printf("Lima instance '%s' started successfully.\n", name)
		// Sleep for a few seconds to allow Lima to initialize
		time.Sleep(5 * time.Second)
	}
}

// checkSiliconMac checks if the system is a Silicon Mac.
func checkSiliconMac() {
	if runtime.GOOS != "darwin" {
		log.Fatalf("This script is designed for macOS.")
	}
	cmd := exec.Command("uname", "-m")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to check system architecture: %v", err)
	}
	if strings.TrimSpace(string(output)) != "arm64" {
		log.Fatalf("This script is designed for Silicon Macs (arm64).")
	}
}

// ensureBrew checks if Homebrew is installed and installs it if not.
func ensureBrew() {
	cmd := exec.Command("brew", "--version")
	if err := cmd.Run(); err != nil {
		fmt.Println("Homebrew is not installed. Installing...")
		installCmd := exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			log.Fatalf("Failed to install Homebrew: %v", err)
		}
	}
}

// ensureLimaCtl checks if limactl is installed and installs it if not.
func ensureLimaCtl() {
	cmd := exec.Command("limactl", "--version")
	if err := cmd.Run(); err != nil {
		fmt.Println("limactl is not installed. Installing...")
		installCmd := exec.Command("brew", "install", "lima")
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			log.Fatalf("Failed to install limactl: %v", err)
		}
	}
}

// checkDebianTools checks if debian_tools instance exists and handles its state.
func checkDebianTools() {
	cmd := exec.Command("limactl", "list", "--json")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to list lima instances: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "debian_tools") {
		fmt.Println("debian_tools instance not found. Creating and starting...")
		limastart("debian_tools", "debian_tools.yaml", false)
		return
	}

	if strings.Contains(outputStr, `"status":"Running"`) {
		fmt.Println("debian_tools instance is already running. Opening shell...")
		//err := syscall.Exec("/usr/local/bin/limactl", []string{"limactl", "shell", "debian_tools"}, os.Environ())
		//if err != nil {
		//	log.Fatalf("Failed to open shell for debian_tools: %v", err)
		//}

		openDefaultShell()
	}

	fmt.Println("Restarting debian_tools instance...")
	//limastart("debian_tools", "debian_tools.yaml", true)  // Stop
	//limastart("debian_tools", "debian_tools.yaml", false) // Start
}

func openDefaultShell() error {
	// `limactl shell default` コマンドでシェルを起動
	cmd := exec.Command("limactl", "shell", "debian_tools")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// コマンドを実行し、エラーがあれば返す
	return cmd.Run()
}

// deleteDebianTools deletes the debian_tools instance if it exists.
func deleteDebianTools() {
	cmd := exec.Command("limactl", "list", "--json")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to list lima instances: %v", err)
	}
	if strings.Contains(string(output), "debian_tools") {
		fmt.Println("Stopping and deleting debian_tools instance...")
		limastart("debian_tools", "debian_tools.yaml", true) // Stop
		deleteCmd := exec.Command("limactl", "delete", "debian_tools")
		deleteCmd.Stdout = os.Stdout
		deleteCmd.Stderr = os.Stderr
		if err := deleteCmd.Run(); err != nil {
			log.Fatalf("Failed to delete debian_tools instance: %v", err)
		}
	} else {
		fmt.Println("debian_tools instance does not exist.")
	}
}

func main() {
	var stop bool
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--help":
			fmt.Println("Usage: go-lima-setup [--stop] [--delete]")
			fmt.Println("This script checks for system compatibility and sets up the lima debian_tools instance.")
			fmt.Println("\nOptions:")
			fmt.Println("  --stop    Stop the debian_tools instance.")
			fmt.Println("  --delete  Stop and delete the debian_tools instance if it exists.")
			return
		case "--delete":
			deleteDebianTools()
			return
		case "--stop":
			stop = true
		}
	}

	checkSiliconMac()
	ensureBrew()
	ensureLimaCtl()
	if stop {
		limastart("debian_tools", "debian_tools.yaml", true) // Stop instance
	} else {
		checkDebianTools()
	}
}
