package main

import (
	"fmt"
	"github.com/MealeyAU/schema/internal/config"
	"github.com/MealeyAU/schema/internal/file"
	"github.com/MealeyAU/schema/internal/printer"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type execErr struct {
	base   error
	detail string
}

func main() {
	logger := printer.Printer{}

	logger.Stringf("Protobuf binding generator")
	cfg := config.Config{}
	cfg.Init()
	logger.Stringf("Outputs: %v", cfg.EnabledOutputsStrings())
	logger.Separator(printer.SeparatorLong)

	paths, err := findFiles("./proto")
	if err != nil {
		log.Fatalf("failed to search for files: %v", err)
	}

	fmt.Println(fmt.Sprintf("%v", paths))

	err = cleanDirectory("./output")
	if err != nil {
		log.Fatalf("failed to clean existing directory: %v", err)
	}

	err = createOutput("./output")
	if err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	if cfg.GoOutput {
		err = generateGoBindings("./output", paths)
		if err != nil {
			log.Fatalf("failed to generate go bindings: %v", err)
		}
	}
	if cfg.WebOutput {
		err = generateWebBindings("./output", paths)
		if err != nil {
			log.Fatalf("failed to generate js bindings: %v", err)
		}
	}
}

func createCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(
		"bash",
		"-c",
		name+" "+strings.Join(args, " "))

	fmt.Println(fmt.Sprintf("%v %v", name, strings.Join(args, " ")))
	return cmd
}

func findFiles(src string) ([]file.Path, error) {
	// Use a map of paths to efficiently de-dupe parent paths
	pathMap := map[file.Path]struct{}{}
	err := filepath.Walk(src,
		func(pathStr string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			path := file.Path(pathStr)
			if !info.IsDir() && path.Extension() == "proto" {
				pathMap[path.Parent()] = struct{}{}
			}
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to walk path: %v", err)
	}

	var paths []file.Path
	for path := range pathMap {
		paths = append(paths, path)
	}
	return paths, nil
}

func createOutput(dest string) error {
	dir := createCommand("mkdir", dest)
	output, err := dir.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error encountered: %v", string(output))
	}
	return nil
}

func cleanDirectory(dest string) error {
	rm := createCommand("rm", "-rf", fmt.Sprintf("%s/*", dest))
	output, err := rm.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error encountered: %v", string(output))
	}
	return nil
}

func generateProtoCommand(file file.Path, opts ...string) *exec.Cmd {
	args := append([]string{"-Iproto", "-Ithird_party"}, opts...)
	args = append(args, fmt.Sprintf("./%s/*.proto", file))

	return createCommand("protoc", args...)
}

func generateGoBindings(dest string, paths []file.Path) error {
	outputDir := fmt.Sprintf("%s/schema-go", dest)
	err := createOutput(outputDir)
	if err != nil {
		return fmt.Errorf("failed to create output dir: %v", err)
	}
	for _, path := range paths {
		protoc := generateProtoCommand(path,
			fmt.Sprintf("--go_out=%s", outputDir),
			"--go_opt=paths=source_relative",
			fmt.Sprintf("--go-grpc_out=%s", outputDir),
			fmt.Sprintf("--go-grpc_opt=paths=source_relative"))

		output, err := protoc.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error encountered running proto command: %v, %v", string(output), err)
		}
	}
	return nil
}

func generateWebBindings(dest string, paths []file.Path) error {
	outputDir := fmt.Sprintf("%s/schema-web", dest)
	err := createOutput(outputDir)
	if err != nil {
		return fmt.Errorf("failed to create output dir: %v", err)
	}
	for _, path := range paths {
		protoc := generateProtoCommand(path,
			"--js_out=import_style=commonjs:output/schema-web",
			"--grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:output/schema-web",
		)

		output, err := protoc.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error encountered running proto command: %v, %v", string(output), err)
		}
	}
	return nil
}
