package utils

import (
	"fmt"
	"os/exec"
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

func ExtractCARFile() error {
	cmd := exec.Command("mkdir", "-p", "./Assets.xcassets")
	cmd2 := exec.Command("./acextract", "-i", "./caomei_tf_clone/Payload/Runner.app/Assets.car", "-o", "./Assets.xcassets")
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
	cmd := exec.Command("actool", "--output-format", "human-readable-text", "--notices", "--warnings", "--platform", "iphoneos", "--minimum-deployment-target", "12.0", "--target-device", "iphone", "--target-device", "ipad", "--compile", "./caomei_tf_clone/Payload/Runner.app", "./Assets.xcassets")
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	return err
}
