package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var hashtagRE = regexp.MustCompile(`#([a-zA-Z][\w\-/]*)`)
var globalHashtags = make(map[string]int)

type FileData struct {
	FilePath string   `json:"file"`
	Tags     []string `json:"tags"`
}

type TagFiles struct {
	Tag   string   `json:"tag"`
	Files []string `json:"files"`
}

func main() {
	home, _ := os.UserHomeDir()
	root := filepath.Join(home, "notes")

	var filesData []FileData

	tagToFiles := make(map[string][]string)

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		matches := hashtagRE.FindAllStringSubmatch(string(data), -1)
		if len(matches) == 0 {
			return nil
		}
		counts := map[string]int{}
		for _, m := range matches {
			counts[m[1]]++

			globalHashtags[m[1]]++
		}
		fmt.Println(path)

		fileData := FileData{
			FilePath: path,
			Tags:     []string{},
		}

		for tag := range counts {
			fmt.Printf("  %s: %d\n", tag, counts[tag])
			fileData.Tags = append(fileData.Tags, tag)

			tagToFiles[tag] = append(tagToFiles[tag], path)
		}

		filesData = append(filesData, fileData)

		return nil
	})

	fmt.Println("\n=== Global Hashtags (Sorted by Frequency) ===")
	printSortedHashtags(globalHashtags)

	outputPath1 := filepath.Join("./dist", "files_to_tags.json")
	writeHashtagsToJSON(filesData, outputPath1)
	fmt.Printf("\nFiles to tags data written to: %s\n", outputPath1)

	outputPath2 := filepath.Join("./dist", "tags_to_files.json")
	writeTagsToFilesJSON(tagToFiles, outputPath2)
	fmt.Printf("Tags to files data written to: %s\n", outputPath2)
}

func printSortedHashtags(hashtags map[string]int) {
	type tagFreq struct {
		tag   string
		count int
	}

	pairs := make([]tagFreq, 0, len(hashtags))

	for tag, count := range hashtags {
		pairs = append(pairs, tagFreq{tag, count})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	for _, pair := range pairs {
		fmt.Printf("  %s: %d\n", pair.tag, pair.count)
	}
}

func writeHashtagsToJSON(filesData []FileData, outputPath string) error {
	jsonData, err := json.MarshalIndent(filesData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, jsonData, 0644)
}

func writeTagsToFilesJSON(tagToFiles map[string][]string, outputPath string) error {
	var tagsData []TagFiles

	for tag, files := range tagToFiles {
		tagsData = append(tagsData, TagFiles{
			Tag:   tag,
			Files: files,
		})
	}

	sort.Slice(tagsData, func(i, j int) bool {
		return tagsData[i].Tag > tagsData[j].Tag
	})

	jsonData, err := json.MarshalIndent(tagsData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, jsonData, 0644)
}
