package main

import (
	"encoding/json"
	"fmt"
	"hello/utils"
	"log"
	"os"
	"path/filepath"
)

type Replacement struct {
	PathToFileName  string            `json:"path_to_file_name"`
	ReplaceAllInDir bool              `json:"replace_all_in_dir,omitempty"`
	Replacements    map[string]string `json:"replacements"`
	IsGlobal        bool              `json:"is_global"`
}

type ImageInfo struct {
	Size         string `json:"size"`
	ExpectedSize string `json:"expected-size"`
	Filename     string `json:"filename"`
	Folder       string `json:"folder"`
	Idiom        string `json:"idiom"`
	Scale        string `json:"scale"`
}

type AppIconContents struct {
	Images []ImageInfo `json:"images"`
}

type LaunchImageContents struct {
	Images []struct {
		Filename string `json:"filename"`
		Idiom    string `json:"idiom"`
		Scale    string `json:"scale"`
	} `json:"images"`
	Info struct {
		Author  string `json:"author"`
		Version int    `json:"version"`
	} `json:"info"`
}

func main() {
	// Remove the folder "AssetsOutput" if it exists
	err := utils.RemoveDirIfExists("AssetsOutput")
	if err != nil {
		log.Fatalf("Failed to remove existing AssetsOutput folder: %v", err)
	}

	// Remove the folder "caomei_tf_clone" if it exists
	err = utils.RemoveDirIfExists("caomei_tf_clone")
	if err != nil {
		log.Fatalf("Failed to remove existing caomei_tf_clone folder: %v", err)
	}

	configFile := "./config.json"
	// Copy the folder
	srcFolder := "caomei_tf"
	destFolder := "caomei_tf_clone"

	err = utils.CopyDir(srcFolder, destFolder)
	if err != nil {
		log.Fatalf("Failed to copy folder: %v", err)
	}

	// Read and parse the JSON configuration file
	jsonData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read JSON config file: %v", err)
	}

	var replacements []Replacement
	err = json.Unmarshal(jsonData, &replacements)
	if err != nil {
		log.Fatalf("Failed to parse JSON config file: %v", err)
	}

	globalReplacements := map[string]string{}
	for _, replacement := range replacements {
		if replacement.IsGlobal {
			globalReplacements = mergeMap(globalReplacements, replacement.Replacements)
		}
	}

	// Process each file and its replacements
	for _, replacement := range replacements {
		if replacement.IsGlobal {
			continue
		}

		utils.ValidateReplacements(replacement.Replacements)

		x := mergeMap(globalReplacements, replacement.Replacements)

		if replacement.ReplaceAllInDir {
			// Process all files in the specified directory
			dirPath := filepath.Dir(replacement.PathToFileName)
			err := utils.ProcessDirectory(dirPath, x)
			if err != nil {
				log.Printf("Failed to process directory %s: %v", dirPath, err)
			}
		} else {
			// Process the single specified file
			err := utils.ProcessFile(replacement.PathToFileName, x)
			if err != nil {
				log.Printf("Failed to process file %s: %v", replacement.PathToFileName, err)
			}
		}
	}

	// Example directory to process
	directory := "./caomei_tf_clone/Payload/Runner.app/"
	err = utils.ProcessImagesInDirectory(directory)
	if err != nil {
		log.Fatalf("Failed to process images in directory: %v", err)
	}

	// Extract the .car file
	err = utils.ExtractCARFile()
	if err != nil {
		log.Fatalf("Failed to extract Assets.car: %v", err)
	}

	extractDir := "./Assets.xcassets"
	// Process all images in the extracted directory
	err = filepath.Walk(extractDir, utils.ProcessImage)
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", extractDir, err)
	}

	err = utils.GroupImagesByFolder(extractDir)
	if err != nil {
		log.Fatalf("Failed to group images by folder: %v", err)
	}

	err = utils.CopyFile("./Assets.xcassets/AppIcon.appiconset/90.png", "./caomei_tf_clone/Payload/Runner.app/AppIcon60x60@2x.png")
	if err != nil {
		log.Fatalf("Failed to update app icon: %v", err)
	}

	// Repack the .car file
	err = utils.RepackCARFile()
	if err != nil {
		log.Fatalf("Failed to repackage Assets.car: %v", err)
	}

	// Create the final zip file
	zipFileName := "final.zip"
	err = utils.CreateZipFromFolder("caomei_tf_clone", zipFileName)
	if err != nil {
		log.Fatalf("Failed to create zip file: %v", err)
	}

	// Rename the zip file to .ipa
	err = os.Rename(zipFileName, "final.ipa")
	if err != nil {
		log.Fatalf("Failed to rename zip file to ipa: %v", err)
	}

	// fmt.Println("Finished compressing to final.ipa")

	fmt.Println("🚀🚀🚀🚀🚀 Finally done 🚀🚀🚀🚀🚀")
}

func mergeMap(m1 map[string]string, m2 map[string]string) map[string]string {
	m := map[string]string{}
	for k, v := range m1 {
		m[k] = v
	}
	for k, v := range m2 {
		m[k] = v
	}

	return m
}
