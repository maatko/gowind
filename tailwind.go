package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

const (
	TAILWIND_URL = "https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-%s-%s"
)

func GetDownloadLink() string {
	os := runtime.GOOS
	if os == "darwin" {
		os = "macos"
	}
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	}

	link := fmt.Sprintf(TAILWIND_URL, os, arch)
	if os == "windows" {
		link += ".exe"
	}

	fmt.Println(link)
	return link
}

func GetCachedFile(path string) (string, error) {
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

func DownloadFile(filePath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status: %s", resp.Status)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy response body: %w", err)
	}

	return nil
}

func GetBinary() (string, error) {
	tailwind, err := GetCachedFile("tailwindcss")
	if err != nil {
		return "", fmt.Errorf("failed to access the tailwind binary: %s", err)
	}

	_, err = os.Stat(tailwind)
	if os.IsNotExist(err) {
		fmt.Println("> Downloading the TailWind binary")

		err = DownloadFile(tailwind, GetDownloadLink())
		if err != nil {
			return "", fmt.Errorf("failed to download the tailwind binary: %s", err)
		}

		os.Chmod(tailwind, 0755)
	}

	return tailwind, nil
}

func main() {
	systemArguments := os.Args[1:]

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to access the working directory", err)
	}

	binary, err := GetBinary()
	if err != nil {
		log.Fatal("failed to access the tailwind binary", binary)
	}

	args := make([]string, len(systemArguments)+1)
	args[0] = binary

	copy(args[1:], systemArguments)

	process, err := os.StartProcess(binary, args, &os.ProcAttr{
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Dir: workingDir,
	})

	if err != nil {
		log.Fatal("failed to execute tailwind subprocess", err)
	}

	process.Wait()
}
