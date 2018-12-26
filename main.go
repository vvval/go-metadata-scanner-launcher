package launcher

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const scanDir = "scanner"
const scanFile = "scanner"

func main() {
	pwd, err := pwd()
	if err != nil {
		fmt.Printf("error: '%s'\nPress <Enter> to exit", err.Error())

		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(0)
	}

	fmt.Println("Enter csv filename...")

	var filename, directory string
	for {
		f, err, repeat := readFilename()
		if err != nil {
			if repeat {
				fmt.Printf("error: '%s', Please retry...\n", err.Error())

				continue
			} else {
				fmt.Printf("error: %s\nPress <Enter> to exit", err.Error())

				bufio.NewReader(os.Stdin).ReadBytes('\n')
				os.Exit(0)
			}
		}

		filename = f

		break
	}

	fmt.Print("Working with filename '" + makePath(pwd, filename) + "'\nEnter directory...\n")

	directory, err = readDir(filename)
	if err != nil {
		fmt.Printf("error: %s\nPress <Enter> to exit", err.Error())

		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(0)
	}

	fmt.Print("Working with directory '" + makePath(pwd, directory) + "'\n")

	err = run(pwd, makePath(pwd, filename), makePath(pwd, directory))
	if err != nil {
		fmt.Printf("error: %s\nPress <Enter> to exit", err.Error())

		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(0)
	}

	fmt.Println("Press <Enter> to exit...")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func readFilename() (string, error, bool) {
	val, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		return "", err, false
	}

	filename := strings.Trim(string(val), " \t\r\n")
	if len(filename) == 0 {
		return "", errors.New("filename is empty"), true
	}

	return filename, nil, false
}

func readDir(filename string) (string, error) {
	val, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		return "", err
	}

	dir := strings.Trim(string(val), " \t\r\n")
	if len(dir) == 0 {
		return filepath.Dir(filename), nil
	}

	return dir, nil
}

func pwd() (string, error) {
	//for debug mode:
	//return ".", nil
	path, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(path), nil
}

func makePath(pwd, filename string) string {
	return filepath.Join(filepath.Dir(pwd), filename)
}

func run(pwd, filename, directory string) error {
	command := exec.Command(cmd(pwd), "write", "-f", filename, "-d", directory, "--v")
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

func cmd(pwd string) string {
	return filepath.Join(pwd, scanDir, scanFile+platformPostfix())
}

func platformPostfix() string {
	switch runtime.GOOS {
	case "windows":
		return "-windows.exe"

	case "darwin":
		return "-mac"

	case "linux":
		return "-linux"

	default:
		return ""
	}
}
