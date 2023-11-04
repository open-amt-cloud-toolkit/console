package templates

import (
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// TemplateParseFSRecursive recursively parses all templates in the FS with the given extension.
// File paths are used as template names to support duplicate file names.
// Use nonRootTemplateNames to exclude root directory from template names
// (e.g. index.html instead of templates/index.html)
func TemplateParseFSRecursive(
	templates fs.FS,
	ext string,
	nonRootTemplateNames bool,
	funcMap template.FuncMap) (*template.Template, error) {

	root := template.New("")
	err := fs.WalkDir(templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(path, ext) {
			if err != nil {
				return err
			}
			b, err := fs.ReadFile(templates, path)
			if err != nil {
				return err
			}
			name := ""
			if nonRootTemplateNames {
				//name the template based on the file path (excluding the root)
				parts := strings.Split(filepath.FromSlash(path), string(os.PathSeparator))
				name = strings.Join(parts[1:], "/")
				fmt.Println(name)
			}
			t := root.New(name).Funcs(funcMap)
			_, err = t.Parse(string(b))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return root, err
}
