package main

import (
	"spf13/cobra"
	"fmt"
	"os"
	"os/exec"
	"log"
	"strings"
	"bufio"
	"io"
	"bytes"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"flag"
	"strconv"
)

type Kernel struct {
	KernelVersion string `json:"kernel_version"`
}

type KernelBuild struct {
	KernelVersion string `json:"kernel_version"`
	BuildStatus string `json:"build_status"`
}

type DockerImage struct {
	KernelVersion string `json:"kernel_version"`
	ImageStatus string `json:"image_status"`
}

type Pod struct {
	KernelVersion string `json:"kernel_version"`
	PodStatus string `json:"pod_status"`
}

type Api struct {
	KernelVersion string `json:"kernel_version"`
	ApiStatus string `json:"api_status"`
}

type Speed struct {
	KernelVersion string `json:"kernel_version"`
	SpeedStatus string `json:"speed_status"`
}

func main() {
	// Parse command line arguments
	kernelVersion := flag.String("kernel", "5.3.18", "Kernel version to build")
	flag.Parse()

	// Build the kernel
	kernel := Kernel{KernelVersion: *kernelVersion}
	buildKernel(kernel)

	// Build the docker image
	dockerImage := DockerImage{KernelVersion: *kernelVersion}
	buildDockerImage(dockerImage)

	// Run the docker image
	pod := Pod{KernelVersion: *kernelVersion}
	runDockerImage(pod)

	// Check if the API is accessible
	api := Api{KernelVersion: *kernelVersion}
	checkApi(api)

	// Check the speed of the API
	speed := Speed{KernelVersion: *kernelVersion}
	checkSpeed(speed)
	buildfromCobra()
	buildContainerfromCobra()
	runContainerfromCobra()
	checkApifromCobra()
	checkSpeedfromCobra()
}

func buildKernel(kernel Kernel) {
	// Print the kernel version
	fmt.Println("Building kernel version: " + kernel.KernelVersion)

	// Run the build script
	cmd := exec.Command("sh", "build_kernel.sh", kernel.KernelVersion)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Print the build status
	buildStatus := out.String()
	fmt.Println(buildStatus)

	// Write the build status to a file
	kernelBuild := KernelBuild{KernelVersion: kernel.KernelVersion, BuildStatus: buildStatus}
	writeJson("kernel_build.json", kernelBuild)
}

func buildDockerImage(dockerImage DockerImage) {
	// Print the kernel version
	fmt.Println("Building docker image with kernel version: " + dockerImage.KernelVersion)

	// Run the build script
	cmd := exec.Command("sh", "build_docker_image.sh", dockerImage.KernelVersion)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Print the image status
	imageStatus := out.String()
	fmt.Println(imageStatus)

	// Write the image status to a file
	dockerImage.ImageStatus = imageStatus
	writeJson("docker_image.json", dockerImage)
}

func runDockerImage(pod Pod) {
	// Print the kernel version
	fmt.Println("Running docker image with kernel version: " + pod.KernelVersion)

	// Run the build script
	cmd := exec.Command("sh", "run_docker_image.sh", pod.KernelVersion)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Print the pod status
	podStatus := out.String()
	fmt.Println(podStatus)

	// Write the pod status to a file
	pod.PodStatus = podStatus
	writeJson("pod.json", pod)
}

func checkApi(api Api) {
	// Print the kernel version
	fmt.Println("Checking if the API is accessible with kernel version: " + api.KernelVersion)

	// Run the check script
	cmd := exec.Command("sh", "check_api.sh", api.KernelVersion)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Print the API status
	apiStatus := out.String()
	fmt.Println(apiStatus)

	// Write the API status to a file
	api.ApiStatus = apiStatus
	writeJson("api.json", api)
}

func checkSpeed(speed Speed) {
	// Print the kernel version
	fmt.Println("Checking the speed of the API with kernel version: " + speed.KernelVersion)

	// Run the check script
	cmd := exec.Command("sh", "check_speed.sh", speed.KernelVersion)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Print the speed status
	speedStatus := out.String()
	fmt.Println(speedStatus)

	// Write the speed status to a file
	speed.SpeedStatus = speedStatus
	writeJson("speed.json", speed)
}

func writeJson(filename string, data interface{}) {
	// Convert the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	// Write the JSON to a file
	err = ioutil
	.WriteFile
	(filename, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func readJson(filename string, data interface{}) {
	// Read the JSON from a file
	jsonData, err := ioutil.ReadFile
	(filename)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the JSON to data
	err = json.Unmarshal(jsonData, data)
	if err != nil {
		log.Fatal(err)
	}
}

func executeCommand(command string) string {
	// Split the command into parts
	parts := strings.Fields(command)
	head := parts[0]
	parts = parts[1:len(parts)]

	// Create the command
	cmd := exec.Command(head, parts...)

	// Execute the command
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}

func getSpeed(url string) float64 {
	// Send a GET request
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate the speed
	elapsed := time.Since(start)
	speed := float64(len(resp.Header) + len(resp.Status)) / elapsed.Seconds()

	return speed
}

//sfcobra arrangement
//1. build_kernel.sh
//2. build_docker_image.sh
//3. run_docker_image.sh
//4. check_api.s


func buildfromCobra() {
	var rootCmd = &cobra.Command{Use: "kernel"}
	var buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build the kernel",
		Run: func(cmd *cobra.Command, args []string) {
			kernel := Kernel{KernelVersion: *kernelVersion}
			buildKernel(kernel)
		},
	}
	rootCmd.AddCommand(buildCmd)
	rootCmd.Execute()
}

func buildContainerfromCobra() {
	var rootCmd = &cobra.Command{Use: "docker"}
	var buildCmd = &cobra.Command{
		Use:   "build",
		Short: "Build the docker image",
		Run: func(cmd *cobra.Command, args []string) {
			dockerImage := DockerImage{KernelVersion: *kernelVersion}
			buildDockerImage(dockerImage)
		},
	}
	rootCmd.AddCommand(buildCmd)
	rootCmd.Execute()
}

func runContainerfromCobra() {
	var rootCmd = &cobra.Command{Use: "pod"}
	var buildCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the docker image",
		Run: func(cmd *cobra.Command, args []string) {
			pod := Pod{KernelVersion: *kernelVersion}
			runDockerImage(pod)
		},
	}
	rootCmd.AddCommand(buildCmd)
	rootCmd.Execute()
}

func checkApifromCobra() {
	var rootCmd = &cobra.Command{Use: "api"}
	var buildCmd = &cobra.Command{
		Use:   "check",
		Short: "Check if the API is accessible",
		Run: func(cmd *cobra.Command, args []string) {
			api := Api{KernelVersion: *kernelVersion}
			checkApi(api)
		},
	}
	rootCmd.AddCommand(buildCmd)
	rootCmd.Execute()
}

func checkSpeedfromCobra() {
	var rootCmd = &cobra.Command{Use: "speed"}
	var buildCmd = &cobra.Command{
		Use:   "check",
		Short: "Check the speed of the API",
		Run: func(cmd *cobra.Command, args []string) {
			speed := Speed{KernelVersion: *kernelVersion}
			checkSpeed(speed)
		},
	}
	rootCmd.AddCommand(buildCmd)
	rootCmd.Execute()
}



