package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Replacement struct {
	PathToFileName  string            `json:"path_to_file_name"`
	ReplaceAllInDir bool              `json:"replace_all_in_dir,omitempty"`
	Replacements    map[string]string `json:"replacements"`
}

func main() {
	// Remove the folder "caomei_tf_clone" if it exists
	err := removeDirIfExists("caomei_tf_clone")
	if err != nil {
		log.Fatalf("Failed to remove existing caomei_tf_clone folder: %v", err)
	}

	configFile := "./config.json"
	// Copy the folder
	srcFolder := "caomei_tf"
	destFolder := "caomei_tf_clone"

	err = copyDir(srcFolder, destFolder)
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

	// Process each file and its replacements
	for _, replacement := range replacements {
		validateReplacements(replacement.Replacements)

		if replacement.ReplaceAllInDir {
			// Process all files in the specified directory
			dirPath := filepath.Dir(replacement.PathToFileName)
			err := processDirectory(dirPath, replacement.Replacements)
			if err != nil {
				log.Printf("Failed to process directory %s: %v", dirPath, err)
			}
		} else {
			// Process the single specified file
			err := processFile(replacement.PathToFileName, replacement.Replacements)
			if err != nil {
				log.Printf("Failed to process file %s: %v", replacement.PathToFileName, err)
			}
		}
	}

	// idir := "./caomei_tf_clone/Payload/Runner.app/Assets.car"

	// Extract the .car file
	// err = extractCARFile()
	// if err != nil {
	// 	log.Fatalf("Failed to extract Assets.car: %v", err)
	// }

	// extractDir := "./AssetsOutput"
	// // Process all images in the extracted directory
	// err = filepath.Walk(extractDir, processImage)
	// if err != nil {
	// 	fmt.Printf("error walking the path %q: %v\n", extractDir, err)
	// }

	// // Repack the .car file
	// err = repackCARFile()
	// if err != nil {
	// 	log.Fatalf("Failed to repackage Assets.car: %v", err)
	// }

	fmt.Println("🚀🚀🚀🚀🚀 Finally done 🚀🚀🚀🚀🚀")
}

func copyDir(src string, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		} else {
			return copyFile(path, destPath)
		}
	})
}

func copyFile(src string, dest string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dest, input, 0644)
}

func removeDirIfExists(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return os.RemoveAll(dir)
	}
	return nil
}

func processFile(filePath string, replacements map[string]string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	modifiedData := data
	for searchSeq, replaceSeq := range replacements {
		modifiedData = bytes.Replace(modifiedData, []byte(searchSeq), []byte(replaceSeq), -1)
	}

	err = os.WriteFile(filePath, modifiedData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %v", filePath, err)
	}

	fmt.Printf("File %s modified successfully\n", filePath)
	return nil
}

func processDirectory(dirPath string, replacements map[string]string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			err := processFile(path, replacements)
			if err != nil {
				log.Printf("Failed to process file %s: %v", path, err)
			}
		}

		return nil
	})
}

func validateReplacements(replacements map[string]string) {
	for key, value := range replacements {
		if len(key) != len(value) {
			log.Printf("replacement length mismatch: key '%s' (length %d) and value '%s' (length %d) do not match", key, len(key), value, len(value))
			panic("WTF")
		}
	}
}

// func processImage(path string, info os.FileInfo, err error) error {
// 	if err != nil {
// 		return err
// 	}

// 	if info.IsDir() {
// 		return nil
// 	}

// 	if !isImageFile(path) {
// 		return nil
// 	}

// 	img, err := imaging.Open(path)
// 	if err != nil {
// 		fmt.Printf("failed to open image: %v\n", err)
// 		return nil
// 	}

// 	// Convert to RGBA to allow modification
// 	rgba := image.NewRGBA(img.Bounds())
// 	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

// 	dotColor := color.RGBA{0, 0, 0, 255} // Black color
// 	dotRadius := 10                      // Radius of the dot

// 	bounds := rgba.Bounds()
// 	centerX, centerY := bounds.Dx()/2, bounds.Dy()/2

// 	for y := -dotRadius; y <= dotRadius; y++ {
// 		for x := -dotRadius; x <= dotRadius; x++ {
// 			if x*x+y*y <= dotRadius*dotRadius {
// 				rgba.Set(centerX+x, centerY+y, dotColor)
// 			}
// 		}
// 	}

// 	err = imaging.Save(rgba, path)
// 	if err != nil {
// 		fmt.Printf("failed to save image: %v\n", err)
// 	}

// 	fmt.Printf("Processed image: %s\n", path)
// 	return nil
// }

// func isImageFile(path string) bool {
// 	ext := strings.ToLower(filepath.Ext(path))
// 	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp" || ext == ".tiff"
// }

// func extractCARFile() error {
// 	cmd := exec.Command("mkdir", "-p", "./AssetsOutput")
// 	cmd2 := exec.Command("./acextract", "-i", "./caomei_tf_clone/Payload/Runner.app/Assets.car", "-o", "./AssetsOutput")
// 	_, err := cmd.Output()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return err
// 	}

// 	_, err = cmd2.Output()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return err
// 	}

// 	return nil
// }

// func repackCARFile() error {
// 	// cmd := exec.Command("mkdir", "-p", "./AssetsOutput")
// 	return nil
// }
