package tailwind

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const (
	TAILWIND_ENDPOINT = "https://github.com/tailwindlabs/tailwindcss/releases"
	TAILWIND_LATEST   = TAILWIND_ENDPOINT + "/latest/download/tailwindcss-%s-%s"
	TAILWIND_VERSION  = TAILWIND_ENDPOINT + "/download/%s/tailwindcss-%s-%s"
)

func GetVersion(version string) string {
	os := runtime.GOOS
	if os == "darwin" {
		os = "macos"
	}
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}

	var link string
	if version != "latest" {
		link = fmt.Sprintf(TAILWIND_VERSION, version, os, arch)
	} else {
		link = fmt.Sprintf(TAILWIND_LATEST, os, arch)
	}

	if os == "windows" {
		link += ".exe"
	}

	return link
}

func GetBinary(path string) string {
	cmd := exec.Command("go", "env", "GOPATH")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	goPath := strings.TrimSpace(string(output))
	goPathDir := fmt.Sprintf("%s/bin", goPath)

	_, err = os.Stat(goPathDir)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(goPathDir, 0755)
		}
	}

	file := fmt.Sprintf("%s/%s", goPathDir, path)
	if runtime.GOOS == "windows" {
		file += ".exe"
	}

	return file
}

func Download(path string, version string) error {
	binaryPath := GetBinary(path)

	log.Print("Downloading TailwindCSS binary")

	resp, err := http.Get(GetVersion(version))
	if err != nil {
		return fmt.Errorf("failed to make request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status: %s", resp.Status)
	}

	out, err := os.Create(binaryPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy response body: %w", err)
	}

	os.Chmod(binaryPath, 0755)
	return nil
}
