package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	Reset   = "\033[0m"
	Cyan    = "\033[36m"
	Magenta = "\033[35m"
	Yellow  = "\033[33m"
)

// Fungsi untuk mencetak ASCII art
func printAsciiArt() {
	asciiArt := []string{
		"     _ _____ _     _____  __",
		"    | | ____| |   / _ \\ \\/ /",
		" _  | |  _| | |  | | | \\  / ",
		"| |_| | |___| |__| |_| /  \\ ",
		" \\___/|_____|_____\\___/_/\\_\\",
	}

	for _, line := range asciiArt {
		fmt.Println(Cyan + line + Reset)
	}
}

// Fungsi untuk mencetak keterangan
func printUsage() {
	fmt.Println("Usage: go run main.go [OPTIONS]")
	fmt.Println("Options:")
	fmt.Println("  -o FILE  Output file (default: wordlist.txt)")
	fmt.Println("  -d DIR   Directory to process")
	fmt.Println("  -l FILE  List file containing repository URLs")
	fmt.Println("  -u URL   Single repository URL to process")
	fmt.Println("  -h       Show help message")
}

// Fungsi untuk membersihkan URL (menghapus spasi ekstra)
func cleanURL(url string) string {
	return strings.TrimSpace(url)
}

// Fungsi untuk meng-clone repo dan mengambil daftar file
func cloneAndListRepo(repoURL string) ([]string, error) {
	tmpDir := filepath.Join(os.TempDir(), "jelox-repo", "jelo-"+uuid.New().String())

	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error creating temp directory: %v", err)
	}

	cmd := exec.Command("git", "clone", repoURL, tmpDir)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error cloning repo: %v, stderr: %s", err, stderr.String())
	}

	cmd = exec.Command("git", "ls-tree", "-r", "--name-only", "HEAD")
	cmd.Dir = tmpDir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error listing files: %v, stderr: %s", err, stderr.String())
	}

	if err := os.RemoveAll(tmpDir); err != nil {
		return nil, fmt.Errorf("error removing temp repo: %v", err)
	}

	files := strings.Split(out.String(), "\n")
	return files, nil
}

// Fungsi untuk membaca repositori dari input stdin atau file
func readReposFromInputOrFile(listFile string) ([]string, error) {
	var repos []string

	if listFile != "" {
		file, err := os.Open(listFile)
		if err != nil {
			return nil, fmt.Errorf("error opening list file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			repoURL := scanner.Text()
			repos = append(repos, cleanURL(repoURL))
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			repoURL := scanner.Text()
			repos = append(repos, cleanURL(repoURL))
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return repos, nil
}

// Fungsi untuk memproses direktori lokal
func processDirectory(directory string) ([]string, error) {
	var files []string

	// Check if the directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %v", directory)
	}

	// Normalize the path to ensure it's in the correct format for the OS
	absDir, err := filepath.Abs(directory)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute directory path: %v", err)
	}
	fmt.Println("Processing directory:", absDir)

	// Use the appropriate command based on the operating system
	var cmd *exec.Cmd
	if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		cmd = exec.Command("cmd", "/C", "dir", "/B", "/S", directory)
	} else {
		cmd = exec.Command("find", directory, "-type", "f")
	}

	// Capture the output
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error listing files in directory: %v, stderr: %s", err, stderr.String())
	}

	// Parse and clean up the output
	files = strings.Split(out.String(), "\n")

	// Remove empty entries and convert absolute paths to relative paths
	var validFiles []string
	for _, file := range files {
		if file != "" {
			relativePath, err := filepath.Rel(absDir, file)
			if err != nil {
				return nil, fmt.Errorf("error calculating relative path: %v", err)
			}
			validFiles = append(validFiles, relativePath)
		}
	}

	return validFiles, nil
}

func main() {
	printAsciiArt()

	if len(os.Args) == 1 {
		fmt.Println("No arguments provided, showing help...")
		flag.Parse()
		printUsage()
		return
	}

	outPtr := flag.String("o", "wordlist.txt", "Output file")
	dirPtr := flag.String("d", "", "Directory to process")
	listFilePtr := flag.String("l", "", "List file containing repository URLs")
	singleRepoPtr := flag.String("u", "", "Single repository URL to process")
	showHelp := flag.Bool("h", false, "Show help message")
	flag.Parse()

	if *showHelp {
		printUsage()
		return
	}

	var files []string

	if *singleRepoPtr != "" {
		// If the -u flag is set, process a single repository URL
		fmt.Println(Magenta + "Processing single repo: " + *singleRepoPtr + Reset)
		files, _ = cloneAndListRepo(*singleRepoPtr)
	} else if *dirPtr != "" {
		// If the -d flag is set, process a local directory
		fmt.Println(Magenta + "Processing directory: " + *dirPtr + Reset)
		files, _ = processDirectory(*dirPtr)
	} else {
		// If the -l flag is set, process the list of repositories from the file
		repos, err := readReposFromInputOrFile(*listFilePtr)
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		for _, repoURL := range repos {
			fmt.Println(Magenta + "Processing repo: " + repoURL + Reset)

			f, err := cloneAndListRepo(repoURL)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			files = append(files, f...)
		}
	}

	// Print the files to stdout
	if len(files) > 0 {
		for _, file := range files {
			if file != "" {
				fmt.Println(file)
			}
		}
	} else {
		fmt.Println("No files found.")
	}

	// Output the files to a specified file if -o is set
	if *outPtr != "" {
		outputFile, err := os.Create(*outPtr)
		if err != nil {
			fmt.Println("Error creating output file:", err)
			return
		}
		defer outputFile.Close()

		for _, file := range files {
			if file != "" {
				outputFile.WriteString(file + "\n")
			}
		}
	}
}
