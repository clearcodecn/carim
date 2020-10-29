package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	protocVersion  = "3.13.0"
	protoURLFormat = "https://github.com/protocolbuffers/protobuf/releases/download/v%s/protoc-%s-%s-x%s.zip"
	saveFile       = "protoc.zip"
	gogoURL        = "https://github.com/gogo/protobuf/archive/v1.3.1.zip"
)

var (
	forceInstall bool
)

func init() {
	flag.BoolVar(&forceInstall, "f", false, "if protoc exist skip")
}

func main() {
	flag.Parse()
	installProtoc()
	setupEnv()
	installGoPlugins()
	buildProtoBuffers()
}

func installProtoc() {
	exist := fileExist(fmt.Sprintf("build/bin/protoc"))
	if exist && !forceInstall {
		return
	}

	arch := "86_32"
	goos := "linux"
	switch runtime.GOARCH {
	case "amd64":
		arch = "86_64"
	default:
		arch = "86_32"
	}
	switch runtime.GOOS {
	case "windows":
		goos = "windows"
	case "darwin":
		goos = "osx"
	case "linux":
		goos = "linux"
	}
	url := fmt.Sprintf(protoURLFormat, protocVersion, protocVersion, goos, arch)

	if !fileExist(saveFile) {
		if err := download(url, saveFile); err != nil {
			log.Println("download failed", err)
			return
		}
	}
	if err := unzip(saveFile, "build", ""); err != nil {
		log.Println("unzip failed", err)
		return
	}

	os.Remove(saveFile)
}

func installGoPlugins() {
	if err := execCommand("go", "install", "-v", "github.com/golang/protobuf/protoc-gen-go"); err != nil {
		log.Println(err)
	}
	if err := execCommand("go", "install", "-v", "github.com/gogo/protobuf/protoc-gen-gogoslick"); err != nil {
		log.Println(err)
	}
	if err := execCommand("go", "install", "-v", "github.com/golang/mock/mockgen"); err != nil {
		log.Println(err)
	}
	// install gogo protobuf.
	if !fileExist("build/include/github.com/gogo/protobuf") {
		if err := download(gogoURL, "gogo.zip"); err != nil {
			return
		}
		defer func() {
			os.RemoveAll("gogo.zip")
		}()
		if err := unzip("gogo.zip", "build/include/github.com/gogo/protobuf", "protobuf-1.3.1"); err != nil {
			log.Println(err)
		}
	}
}

func buildProtoBuffers() {
	args := []string{
		"-Ibuild/include",
		"-Iproto/",
		"--gofast_out=plugins=grpc:./proto",
		"proto/protocol/proto.proto",
	}

	if err := execCommand("protoc", args...); err != nil {
		log.Println(err)
	}
}

func fileExist(f string) bool {
	_, err := os.Stat(f)
	if err != nil {
		return false
	}

	return true
}

func download(url string, save string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}
	f, err := os.OpenFile(save, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func unzip(file string, dir string, trim string) error {
	os.MkdirAll(dir, 0777)
	r, err := zip.OpenReader(file)
	if err != nil {
		return err
	}
	for _, file := range r.File {
		path := strings.TrimPrefix(file.Name, trim)
		path = filepath.Join(dir, path)
		rc, err := file.Open()
		if err != nil {
			return err
		}
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, 0777)
			continue
		} else {
			os.MkdirAll(filepath.Dir(path), 0777)
		}
		dst, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			return err
		}
		_, err = io.Copy(dst, rc)
		if err != nil {
			return err
		}
		rc.Close()
	}
	return nil
}

func execCommand(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func setupEnv() {
	path := os.Getenv("PATH")
	dir, _ := os.Getwd()
	bin := filepath.Join(dir, "build/bin")
	os.Setenv("PATH", fmt.Sprintf("%s:", path, bin))
}
