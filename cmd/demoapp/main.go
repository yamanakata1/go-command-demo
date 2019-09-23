package main

import (
	"fmt"
	"os"
	"sync"
	"runtime"
	"os/exec"
	"strings"
	"io/ioutil"
	"path/filepath"
)

func main() {
	dir := "./web/img"
	paths := getImageFilePaths(dir)

	var wg sync.WaitGroup
	cpus := runtime.NumCPU()
	fmt.Println("cpus: ", cpus)
	semaphore := make(chan int, cpus)
	for i, _ := range paths {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			semaphore <- 1
			out, err := exec.Command("sh", "-c", "guetzli --quality 84 " + paths[i] + " " + dir + "/" + getFileNameWithoutExt(paths[i]) + "_compressed.jpg").Output()
			if err != nil {
				fmt.Println(err)
			}
			if out != nil {
				fmt.Println(string(out))
				fmt.Printf("%s Done.\n", paths[i])
			}
			<-semaphore
		}(i)
	}
	wg.Wait()
	fmt.Println("Finish.")
}

func getImageFilePaths(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), "jpg") &&
		!strings.HasSuffix(file.Name(), "jpeg") &&
		!strings.HasSuffix(file.Name(), "png") &&
		!strings.HasSuffix(file.Name(), "gif") {
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}
	return paths
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
