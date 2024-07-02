package utils

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
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

func ProcessImage(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	if !isImageFile(path) {
		return nil
	}

	img, err := imaging.Open(path)
	if err != nil {
		fmt.Printf("failed to open image: %v\n", err)
		return nil
	}

	// Convert to RGBA to allow modification
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	dotColor := color.RGBA{0, 0, 0, 255} // Black color
	dotRadius := 1                       // Radius of the dot

	bounds := rgba.Bounds()
	centerX, centerY := bounds.Dx()/2, bounds.Dy()/2

	for y := -dotRadius; y <= dotRadius; y++ {
		for x := -dotRadius; x <= dotRadius; x++ {
			if x*x+y*y <= dotRadius*dotRadius {
				rgba.Set(centerX+x, centerY+y, dotColor)
			}
		}
	}

	err = imaging.Save(rgba, path)
	if err != nil {
		fmt.Printf("failed to save image: %v\n", err)
	}

	fmt.Printf("Processed image: %s\n", path)
	return nil
}

func isImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp" || ext == ".tiff"
}

func ExtractCARFile() error {
	cmd := exec.Command("mkdir", "-p", "./AssetsOutput")
	cmd2 := exec.Command("./acextract", "-i", "./caomei_tf_clone/Payload/Runner.app/Assets.car", "-o", "./AssetsOutput")
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	_, err = cmd2.Output()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func RepackCARFile() error {
	cmd := exec.Command("actool", "--output-format", "human-readable-text", "--notices", "--warnings", "--platform", "iphoneos", "--minimum-deployment-target", "12.0", "--target-device", "iphone", "--target-device", "ipad", "--compile", "./caomei_tf_clone/Payload/Runner.app", "./AssetsOutput")
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}

func GroupImagesByFolder(basePath string) error {
	imageGroups := make(map[string][]ImageInfo)
	imageRegex := regexp.MustCompile(`(AppIcon|LaunchImage)(\d+x\d+)?(@(\dx))?(~(\w+))?\.png`)

	err := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !isImageFile(path) {
			return nil
		}

		filename := filepath.Base(path)
		match := imageRegex.FindStringSubmatch(filename)
		if match == nil {
			return nil
		}

		baseName := match[1]
		size := match[2]
		scale := match[4]
		if scale == "" {
			scale = "1x"
		}
		idiom := match[6]
		if idiom == "" {
			idiom = "universal"
		}

		expectedSize := size
		if size != "" && scale != "" {
			width := strings.Split(size, "x")[0]
			widthInt, err := strconv.Atoi(width)
			if err != nil {
				return err
			}
			scaleFactor, err := strconv.Atoi(string(scale[0]))
			if err != nil {
				return err
			}
			expectedSize = strconv.Itoa(scaleFactor * widthInt)
		}

		var folderName string
		if baseName == "AppIcon" {
			folderName = fmt.Sprintf("%s.iconimageset", baseName)
		} else if baseName == "LaunchImage" {
			folderName = fmt.Sprintf("%s.imageset", baseName)
		}
		folderPath := filepath.Join(basePath, folderName)

		var destFilename string
		if baseName == "AppIcon" && size != "" && scale != "" {
			destFilename = fmt.Sprintf("icon-%s@%s.png", size, scale)
		} else {
			destFilename = filename
		}

		imageInfo := ImageInfo{
			Size:         size,
			ExpectedSize: expectedSize,
			Filename:     destFilename,
			Folder:       folderPath,
			Idiom:        idiom,
			Scale:        scale,
		}

		imageGroups[folderName] = append(imageGroups[folderName], imageInfo)

		// Create folder if it doesn't exist
		err = os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return err
		}

		// Move and rename file to the new folder
		destPath := filepath.Join(folderPath, destFilename)
		err = os.Rename(path, destPath)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Generate Contents.json files for each folder
	for folder, images := range imageGroups {
		if strings.Contains(folder, "AppIcon") {
			// Update folder path for AppIcon images
			for i := range images {
				images[i].Folder = filepath.Join(basePath, folder)
			}
			contents := AppIconContents{Images: images}
			contentsPath := filepath.Join(basePath, folder, "Contents.json")
			file, err := os.Create(contentsPath)
			if err != nil {
				return err
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			err = encoder.Encode(contents)
			if err != nil {
				return err
			}
		} else if strings.Contains(folder, "LaunchImage") {
			launchImages := []struct {
				Filename string `json:"filename"`
				Idiom    string `json:"idiom"`
				Scale    string `json:"scale"`
			}{}
			for _, img := range images {
				launchImages = append(launchImages, struct {
					Filename string `json:"filename"`
					Idiom    string `json:"idiom"`
					Scale    string `json:"scale"`
				}{
					Filename: img.Filename,
					Idiom:    img.Idiom,
					Scale:    img.Scale,
				})
			}
			contents := LaunchImageContents{
				Images: launchImages,
				Info: struct {
					Author  string `json:"author"`
					Version int    `json:"version"`
				}{
					Author:  "xcode",
					Version: 1,
				},
			}
			contentsPath := filepath.Join(basePath, folder, "Contents.json")
			file, err := os.Create(contentsPath)
			if err != nil {
				return err
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			err = encoder.Encode(contents)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ProcessImagesInDirectory processes all images in a directory and its subdirectories.
func ProcessImagesInDirectory(root string) error {
	return filepath.Walk(root, ProcessImage)
}
