package main

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var hashtagRE = regexp.MustCompile(`#([a-zA-Z][\w\-/]*)`)

func main() {
	home, _ := os.UserHomeDir()
	root := filepath.Join(home, "notes")

	globalFreq := map[string]int{}
	tagToFiles := map[string][]string{}
	var filesToTags []FileData

	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || strings.ToLower(filepath.Ext(path)) != ".md" {
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

		tagSet := map[string]struct{}{}
		for _, m := range matches {
			tag := m[1]
			globalFreq[tag]++
			tagToFiles[tag] = append(tagToFiles[tag], path)
			tagSet[tag] = struct{}{}
		}

		tags := make([]string, 0, len(tagSet))
		for tag := range tagSet {
			tags = append(tags, tag)
		}
		filesToTags = append(filesToTags, FileData{FilePath: path, Tags: tags})

		return nil
	})

	_ = os.MkdirAll("./dist", 0755)

	writeJSON(sortMapByValue(globalFreq), "./dist/frequency.json")
	writeJSON(tagToFiles, "./dist/tags-files.json")
	writeJSON(filesToTags, "./dist/files.tags.json")
}

type FileData struct {
	FilePath string   `json:"file"`
	Tags     []string `json:"tags"`
}

func writeJSON(data any, path string) {
	b, _ := json.MarshalIndent(data, "", "  ")
	_ = os.WriteFile(path, b, 0644)
}

func sortMapByValue(m map[string]int) map[string]int {
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range m {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	out := make(map[string]int, len(sorted))
	for _, kv := range sorted {
		out[kv.Key] = kv.Value
	}
	return out
}
