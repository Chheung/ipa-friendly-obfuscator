package utils

// import (
// 	"archive/zip"
// 	"fmt"
// 	"image"
// 	"image/color"
// 	"image/draw"
// 	"image/png"
// 	"io"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strings"
// )

// func main() {
// 	ipaFile := "Archive.ipa"
// 	extractDir := "./extracted"
// 	assetsDir := filepath.Join(extractDir, "Assets.xcassets")
// 	outputDir := "./output"
// 	assetsCar := filepath.Join(extractDir, "Payload", "Runner.app", "Assets.car")

// 	err := extractIPA(ipaFile, extractDir)
// 	if err != nil {
// 		fmt.Printf("Error extracting IPA: %v\n", err)
// 		return
// 	}

// 	err = os.MkdirAll(assetsDir, os.ModePerm)
// 	if err != nil {
// 		fmt.Printf("Error creating assets directory: %v\n", err)
// 		return
// 	}

// 	err = parseAssetsCar(assetsCar, assetsDir)
// 	if err != nil {
// 		fmt.Printf("Error parsing Assets.car: %v\n", err)
// 		return
// 	}

// 	err = modifyAssets(assetsDir)
// 	if err != nil {
// 		fmt.Printf("Error modifying assets: %v\n", err)
// 		return
// 	}

// 	err = repackageAssets(assetsDir, outputDir)
// 	if err != nil {
// 		fmt.Printf("Error repackaging assets: %v\n", err)
// 		return
// 	}

// 	// err = replaceAssetsCar(extractDir, outputDir)
// 	// if err != nil {
// 	// 	fmt.Printf("Error replacing Assets.car: %v\n", err)
// 	// 	return
// 	// }

// 	// fmt.Println("Completed successfully!")
// }

// func extractIPA(ipaFile, extractDir string) error {
// 	r, err := zip.OpenReader(ipaFile)
// 	if err != nil {
// 		return err
// 	}
// 	defer r.Close()

// 	for _, f := range r.File {
// 		fpath := filepath.Join(extractDir, f.Name)

// 		if f.FileInfo().IsDir() {
// 			os.MkdirAll(fpath, os.ModePerm)
// 			continue
// 		}

// 		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
// 			return err
// 		}

// 		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
// 		if err != nil {
// 			return err
// 		}

// 		rc, err := f.Open()
// 		if err != nil {
// 			return err
// 		}

// 		_, err = io.Copy(outFile, rc)

// 		outFile.Close()
// 		rc.Close()

// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func parseAssetsCar(assetsCar, assetsDir string) error {
// 	cmd := exec.Command("./acextract", "-i", assetsCar, "-o", assetsDir)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	return cmd.Run()
// }

// func modifyAssets(assetsDir string) error {
// 	return filepath.Walk(assetsDir, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		if !info.IsDir() && strings.Contains(path, ".png") {
// 			img, err := addDotsToImage(path)
// 			if err != nil {
// 				return err
// 			}
// 			outFile, err := os.Create(path)
// 			if err != nil {
// 				return err
// 			}
// 			defer outFile.Close()
// 			err = png.Encode(outFile, img)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }

// func addDotsToImage(path string) (image.Image, error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	img, err := png.Decode(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	bounds := img.Bounds()
// 	dottedImg := image.NewRGBA(bounds)
// 	draw.Draw(dottedImg, bounds, img, bounds.Min, draw.Over)

// 	dotColor := color.RGBA{0, 0, 0, 255}
// 	for y := bounds.Min.Y; y < bounds.Max.Y; y += 20 {
// 		for x := bounds.Min.X; x < bounds.Max.X; x += 20 {
// 			drawDot(dottedImg, x, y, dotColor)
// 		}
// 	}

// 	return dottedImg, nil
// }

// func drawDot(img *image.RGBA, x, y int, clr color.Color) {
// 	radius := 5
// 	for dy := -radius; dy <= radius; dy++ {
// 		for dx := -radius; dx <= radius; dx++ {
// 			if dx*dx+dy*dy <= radius*radius {
// 				img.Set(x+dx, y+dy, clr)
// 			}
// 		}
// 	}
// }

// func repackageAssets(assetsDir, outputDir string) error {
// 	err := os.MkdirAll(outputDir, os.ModePerm)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println(outputDir)

// 	cmd := exec.Command("xcrun", "actool",
// 		"--output-format", "human-readable-text",
// 		"--notices",
// 		"--warnings",
// 		"--platform", "iphoneos",
// 		"--minimum-deployment-target", "12.0",
// 		"--target-device", "iphone",
// 		"--target-device", "ipad",
// 		"--compile", outputDir, assetsDir)

// 	fmt.Println(cmd.String())
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	return cmd.Run()
// }

// func replaceAssetsCar(extractDir, outputDir string) error {
// 	src := filepath.Join(outputDir, "Assets.car")
// 	dest := filepath.Join(extractDir, "Payload", "Runner.app", "Assets.car")

// 	fmt.Printf("Moving %s to %s\n", src, dest)

// 	// Ensure the destination file does not exist before renaming
// 	if _, err := os.Stat(dest); err == nil {
// 		if err := os.Remove(dest); err != nil {
// 			return fmt.Errorf("error removing original Assets.car: %v", err)
// 		}
// 	}

// 	// Rename the new Assets.car to the destination
// 	if err := os.Rename(src, dest); err != nil {
// 		return fmt.Errorf("error renaming Assets.car: %v", err)
// 	}

// 	return nil
// }
