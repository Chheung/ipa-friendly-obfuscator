package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// type Replacement struct {
// 	PathToFileName  string              `json:"path_to_file_name"`
// 	ReplaceAllInDir bool                `json:"replace_all_in_dir,omitempty"`
// 	Replacements    map[string][]string `json:"replacements"`
// 	IsGlobal        bool                `json:"is_global"`
// }

// type ImageInfo struct {
// 	Size         string `json:"size"`
// 	ExpectedSize string `json:"expected-size"`
// 	Filename     string `json:"filename"`
// 	Folder       string `json:"folder"`
// 	Idiom        string `json:"idiom"`
// 	Scale        string `json:"scale"`
// }

// type AppIconContents struct {
// 	Images []ImageInfo `json:"images"`
// }

// type LaunchImageContents struct {
// 	Images []struct {
// 		Filename string `json:"filename"`
// 		Idiom    string `json:"idiom"`
// 		Scale    string `json:"scale"`
// 	} `json:"images"`
// 	Info struct {
// 		Author  string `json:"author"`
// 		Version int    `json:"version"`
// 	} `json:"info"`
// }

// func main() {
// 	// Remove the folder "AssetsOutput" if it exists
// 	err := utils.RemoveDirIfExists("AssetsOutput")
// 	if err != nil {
// 		log.Fatalf("Failed to remove existing AssetsOutput folder: %v", err)
// 	}

// 	// Remove the folder "output" if it exists
// 	err = utils.RemoveDirIfExists("output")
// 	if err != nil {
// 		log.Fatalf("Failed to remove existing output folder: %v", err)
// 	}

// 	configFile := "./config.json"
// 	// Copy the folder
// 	srcFolder := "input"
// 	destFolder := "output"

// 	err = utils.CopyDir(srcFolder, destFolder)
// 	if err != nil {
// 		log.Fatalf("Failed to copy folder: %v", err)
// 	}

// 	// Read and parse the JSON configuration file
// 	jsonData, err := os.ReadFile(configFile)
// 	if err != nil {
// 		log.Fatalf("Failed to read JSON config file: %v", err)
// 	}

// 	var replacements []Replacement
// 	err = json.Unmarshal(jsonData, &replacements)
// 	if err != nil {
// 		log.Fatalf("Failed to parse JSON config file: %v", err)
// 	}

// 	globalReplacements := map[string][]string{}
// 	for _, replacement := range replacements {
// 		if replacement.IsGlobal {
// 			globalReplacements = mergeMap(globalReplacements, replacement.Replacements)
// 		}
// 	}

// 	// Process each file and its replacements
// 	for _, replacement := range replacements {
// 		if replacement.IsGlobal {
// 			continue
// 		}

// 		utils.ValidateReplacements(replacement.Replacements)

// 		x := mergeMap(globalReplacements, replacement.Replacements)

// 		if replacement.ReplaceAllInDir {
// 			// Process all files in the specified directory
// 			dirPath := filepath.Dir(replacement.PathToFileName)
// 			err := utils.ProcessDirectory(dirPath, x)
// 			if err != nil {
// 				log.Printf("Failed to process directory %s: %v", dirPath, err)
// 			}
// 		} else {
// 			// Process the single specified file
// 			err := utils.ProcessFile(replacement.PathToFileName, x)
// 			if err != nil {
// 				log.Printf("Failed to process file %s: %v", replacement.PathToFileName, err)
// 			}
// 		}
// 	}

// 	// Example directory to process
// 	directory := "./output/Payload/Runner.app/"
// 	err = utils.ProcessImagesInDirectory(directory)
// 	if err != nil {
// 		log.Fatalf("Failed to process images in directory: %v", err)
// 	}

// 	// // Extract the .car file
// 	// err = utils.ExtractCARFile()
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to extract Assets.car: %v", err)
// 	// }

// 	// extractDir := "./Assets.xcassets"
// 	// // Process all images in the extracted directory
// 	// err = filepath.Walk(extractDir, utils.ProcessImage)
// 	// if err != nil {
// 	// 	fmt.Printf("error walking the path %q: %v\n", extractDir, err)
// 	// }

// 	// err = utils.GroupImagesByFolder(extractDir)
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to group images by folder: %v", err)
// 	// }

// 	icon1Path := "./output/Payload/Runner.app/AppIcon60x60@2x.png"
// 	err = utils.ConvertCgBIToPng(icon1Path, icon1Path)
// 	if err != nil {
// 		log.Fatalf("Error converting CgBI to PNG: %v", err)
// 	}

// 	icon1, err := os.Stat(icon1Path)
// 	if err != nil {
// 		log.Fatalf("Error getting image stats of icon1: %v", err)
// 	}

// 	icon2Path := "./output/Payload/Runner.app/AppIcon76x76@2x~ipad.png"
// 	err = utils.ConvertCgBIToPng(icon2Path, icon2Path)
// 	if err != nil {
// 		log.Fatalf("Error converting CgBI to PNG: %v", err)
// 	}

// 	icon2, err := os.Stat(icon2Path)
// 	if err != nil {
// 		log.Fatalf("Error getting image stats of icon2: %v", err)
// 	}

// 	err = utils.ProcessImage(icon1Path, icon1, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to process images in directory: %v", err)
// 	}

// 	err = utils.ProcessImage(icon2Path, icon2, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to process images in directory: %v", err)
// 	}

// 	// // Repack the .car file
// 	// err = utils.RepackCARFile()
// 	// if err != nil {
// 	// 	log.Fatalf("Failed to repackage Assets.car: %v", err)
// 	// }

// 	// Create the final zip file
// 	zipFileName := "final.zip"
// 	err = utils.CreateZipFromFolder("output", zipFileName)
// 	if err != nil {
// 		log.Fatalf("Failed to create zip file: %v", err)
// 	}

// 	// Rename the zip file to .ipa
// 	err = os.Rename(zipFileName, "final.ipa")
// 	if err != nil {
// 		log.Fatalf("Failed to rename zip file to ipa: %v", err)
// 	}

// 	fmt.Println("Finished compressing to final.ipa")

// 	fmt.Println("🚀🚀🚀🚀🚀 Finally done 🚀🚀🚀🚀🚀")
// }

// func mergeMap(m1 map[string][]string, m2 map[string][]string) map[string][]string {
// 	m := map[string][]string{}
// 	for k, v := range m1 {
// 		m[k] = v
// 	}
// 	for k, v := range m2 {
// 		m[k] = v
// 	}

// 	return m
// }

func main() {
	// Specify the directory to search
	dir := "./final"
	targetString := []byte("fhamprjghsjrffwxfqjlapqofhsq")

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Read the file
			data, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Printf("Failed to read file %s: %v\n", path, err)
				return nil
			}

			// Check if the file contains the target string
			if bytes.Contains(data, targetString) {
				fmt.Printf("File contains target string: %s\n", path)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", dir, err)
		return
	}

	fmt.Println("DUNE")
}
