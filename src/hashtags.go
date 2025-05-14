package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Load environment variables from .env file
func loadEnvFile() map[string]string {
	env := make(map[string]string)
	data, err := os.ReadFile(".env")
	if err != nil {
		return env
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// Remove quotes if present
		value = strings.Trim(value, "\"'")
		env[key] = value
	}
	return env
}

var hashtagRE = regexp.MustCompile(`#([a-zA-Z][\w\-/]{1,60})`)

func isValidHashtag(tag string) bool {
	if matched, _ := regexp.MatchString(`^[0-9a-fA-F]{3}([0-9a-fA-F]{3})?$`, tag); matched {
		return false
	}
	if strings.HasPrefix(tag, "L") && len(tag) > 1 && isNumeric(tag[1:]) {
		return false
	}
	if len(tag) == 1 {
		return false
	}
	return true
}

func isNumeric(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func stripCodeAndStrings(src string) string {
	var out strings.Builder
	inCodeBlock := false

	for _, line := range strings.Split(src, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "```") {
			inCodeBlock = !inCodeBlock
			out.WriteByte('\n')
			continue
		}
		if inCodeBlock {
			out.WriteByte('\n')
			continue
		}

		inSingle, inDouble, inBacktick := false, false, false
		var buf strings.Builder

		for _, ch := range line {
			switch ch {
			case '"':
				if !inSingle && !inBacktick {
					inDouble = !inDouble
					buf.WriteRune(' ')
					continue
				}
			case '\'':
				if !inDouble && !inBacktick {
					inSingle = !inSingle
					buf.WriteRune(' ')
					continue
				}
			case '`':
				if !inSingle && !inDouble {
					inBacktick = !inBacktick
					buf.WriteRune(' ')
					continue
				}
			}
			if inSingle || inDouble || inBacktick {
				buf.WriteRune(' ')
			} else {
				buf.WriteRune(ch)
			}
		}
		out.WriteString(buf.String())
		out.WriteByte('\n')
	}
	return out.String()
}

func main() {
	// Define command line flags
	pathFlag := flag.String("path", "", "Full path to notes directory")
	flag.Parse()

	// Determine notes directory path from various sources
	var notesPath string

	// 1. Check command line argument
	if *pathFlag != "" {
		notesPath = *pathFlag
	} else {
		// 2. Check environment variables from .env file
		env := loadEnvFile()
		if envPath, ok := env["NOTES_PATH"]; ok && envPath != "" {
			notesPath = envPath
		} else {
			// 3. Fallback to default path
			home, _ := os.UserHomeDir()
			notesPath = filepath.Join(home, "notes")
			
			// Error out if not provided explicitly and default doesn't exist
			if _, err := os.Stat(notesPath); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: Notes path not found\n")
				fmt.Fprintf(os.Stderr, "Please specify path using one of these methods:\n")
				fmt.Fprintf(os.Stderr, "1. Command line: -path=/path/to/notes\n")
				fmt.Fprintf(os.Stderr, "2. Environment: Create .env file with NOTES_PATH=/path/to/notes\n")
				os.Exit(1)
			}
		}
	}

	fmt.Printf("Processing notes from: %s\n", notesPath)

	globalFreq := map[string]int{}
	tagToFiles := map[string]map[string]struct{}{}
	var filesToTags []FileData

	_ = filepath.WalkDir(notesPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || strings.ToLower(filepath.Ext(path)) != ".md" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		cleaned := stripCodeAndStrings(string(data))
		matches := hashtagRE.FindAllStringSubmatch(cleaned, -1)

		tagSet := map[string]struct{}{}
		for _, m := range matches {
			tag := m[1]
			if isValidHashtag(tag) {
				globalFreq[tag]++
				tagSet[tag] = struct{}{}
				if _, ok := tagToFiles[tag]; !ok {
					tagToFiles[tag] = make(map[string]struct{})
				}
				tagToFiles[tag][path] = struct{}{}
			}
		}

		if len(tagSet) == 0 {
			return nil
		}

		tags := make([]string, 0, len(tagSet))
		for tag := range tagSet {
			tags = append(tags, tag)
		}
		sort.Strings(tags)

		filesToTags = append(filesToTags, FileData{FilePath: path, Tags: tags})
		return nil
	})

	_ = os.MkdirAll("./dist", 0755)

	writeJSON(sortMapByValue(globalFreq), "./dist/frequency.json")
	writeJSON(convertTagToFiles(tagToFiles), "./dist/tags-files.json")
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

func convertTagToFiles(m map[string]map[string]struct{}) map[string][]string {
	out := make(map[string][]string, len(m))
	for tag, fileSet := range m {
		files := make([]string, 0, len(fileSet))
		for f := range fileSet {
			files = append(files, f)
		}
		sort.Strings(files)
		out[tag] = files
	}
	return out
}
