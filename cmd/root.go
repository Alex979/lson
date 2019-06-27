package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// File represents a unix file
type File struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

// Directory extends File and contains a splice of files and/or directorie as Children
type Directory struct {
	File
	Children []interface{} `json:"children"`
}

// Setup for the lson command
var rootCmd = &cobra.Command{
	Use:     "lson [file]",
	Short:   "Display directory structure in json format.",
	Version: "1.0.0",
	Args: func(cmd *cobra.Command, args []string) error {
		// Ensure a file argument is passed in
		if len(args) < 1 {
			return errors.New("requires a file argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Check if first argument passed is a valid file or folder
		fileInfo, err := os.Lstat(args[0])
		if err != nil {
			log.Fatal(err)
		}

		// If file is a folder, begin building json structure with a directory as the root
		if fileInfo.IsDir() {

			jsonDir, err := json.MarshalIndent(buildDir(args[0]), "", "  ")
			if err != nil {
				log.Fatal(err)
			}

			// Write json output to stdout
			os.Stdout.Write(append(jsonDir, []byte("\n")...))
		} else {
			// Otherwise, build json structure as a single file for the root

			jsonDir, err := json.MarshalIndent(buildFile(args[0]), "", "  ")
			if err != nil {
				log.Fatal(err)
			}

			// Write json output to stdout
			os.Stdout.Write(append(jsonDir, []byte("\n")...))
		}
	},
}

// Returns a File struct for the given file path
func buildFile(file string) File {
	// Get file info
	fileInfo, err := os.Lstat(file)
	if err != nil {
		log.Fatal(err)
	}

	// Return new File struct
	return File{
		Name: fileInfo.Name(),
		Type: "file",
		Size: fileInfo.Size(),
	}
}

// Returns a Directory struct for the given file path that recursively builds children
func buildDir(directory string) Directory {
	// Initialize Directory fields
	var children []interface{}
	var size int64

	// Read files from directory
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through all the directory's children and build the File/Directory objects
	for _, file := range files {
		if file.IsDir() {
			children = append(children, buildDir(directory+"/"+file.Name()))
		} else {
			children = append(children, buildFile(directory+"/"+file.Name()))
		}
	}

	// Add the size of each child file/directory to this directories size
	for _, child := range children {
		switch child.(type) {
		case File:
			size += child.(File).Size
			break
		case Directory:
			size += child.(Directory).Size
		}
	}

	// Get file info for the directory
	fileInfo, err := os.Lstat(directory)
	if err != nil {
		log.Fatal(err)
	}

	// Return new Directory struct
	return Directory{
		File: File{
			Name: fileInfo.Name(),
			Type: "directory",
			Size: size,
		},
		Children: children,
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
