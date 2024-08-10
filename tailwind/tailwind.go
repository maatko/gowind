package tailwind

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
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

func GetBinary(path string) (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get the current directory: %s", err)
	}

	tailwindDir := fmt.Sprintf("%s/.tailwind", workingDir)

	_, err = os.Stat(tailwindDir)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(tailwindDir, 0755)
		}
	}

	file := fmt.Sprintf("%s/%s", tailwindDir, path)
	if runtime.GOOS == "windows" {
		file += ".exe"
	}

	return file, nil
}

func DownloadMissing(path string, version string) (string, error) {
	binaryPath, err := GetBinary(path)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %s", err)
	}

	_, err = os.Stat(binaryPath)
	if !os.IsNotExist(err) {
		return binaryPath, nil
	}

	log.Print("Downloading TailwindCSS binary")

	resp, err := http.Get(GetVersion(version))
	if err != nil {
		return "", fmt.Errorf("failed to make request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned non-200 status: %s", resp.Status)
	}

	out, err := os.Create(binaryPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %s", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to copy response body: %w", err)
	}

	os.Chmod(binaryPath, 0755)
	return binaryPath, nil
}

func Execute(arguments ...string) error {
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to access the working directory: %s", err)
	}

	binary, err := DownloadMissing("tailwindcss", "latest")
	if err != nil {
		return fmt.Errorf("failed to access the tailwind binary: %s", binary)
	}

	args := make([]string, len(arguments)+1)
	args[0] = binary

	copy(args[1:], arguments)

	process, err := os.StartProcess(binary, args, &os.ProcAttr{
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Dir: workingDir,
	})

	if err != nil {
		return fmt.Errorf("failed to execute tailwind subprocess: %s", err)
	}

	process.Wait()
	return nil
}
