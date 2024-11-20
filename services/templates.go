package services

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hippo-an/goranchise/config"
	"github.com/hippo-an/goranchise/funcmap"
	"html/template"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

type TemplateRenderer struct {
	templateCache sync.Map
	funcMap       template.FuncMap
	templatesPath string
	config        *config.Config
}

func NewTemplateRenderer(cfg *config.Config) *TemplateRenderer {
	t := &TemplateRenderer{
		templateCache: sync.Map{},
		funcMap:       funcmap.GetFuncMap(),
		config:        cfg,
	}

	t.templatesPath = getTemplatesDirectoryPath()
	return t
}

func (t *TemplateRenderer) GetTemplatesPath() string {
	return t.templatesPath
}

func (t *TemplateRenderer) Load(module, key string) (*template.Template, error) {
	load, ok := t.templateCache.Load(t.getCacheKey(module, key))
	if !ok {
		return nil, errors.New("uncached page template requested")
	}

	templ, ok := load.(*template.Template)
	if !ok {
		return nil, errors.New("unable to cast cached template")
	}

	return templ, nil
}

func (t *TemplateRenderer) Execute(module, key, layoutName string, data interface{}) (*bytes.Buffer, error) {
	tmpl, err := t.Load(module, key)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(buf, layoutName+config.TemplateExt, data)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (t *TemplateRenderer) Parse(module, key, name string, files []string, directories []string) error {
	cacheKey := t.getCacheKey(module, key)
	if _, err := t.Load(module, key); err != nil || t.config.App.Environment == config.EnvironmentProd {
		parsed := template.New(name + config.TemplateExt).
			Funcs(t.funcMap)

		if len(files) > 0 {
			for idx, v := range files {
				files[idx] = fmt.Sprintf("%s/%s%s", t.templatesPath, v, config.TemplateExt)
			}

			parsed, err = parsed.ParseFiles(files...)
			if err != nil {
				return err
			}
		}

		for _, dir := range directories {
			dir = fmt.Sprintf("%s/%s/*%s", t.templatesPath, dir, config.TemplateExt)
			parsed, err = parsed.ParseGlob(dir)
			if err != nil {
				return err
			}
		}

		t.templateCache.Store(cacheKey, parsed)
	}

	return nil
}

func (t *TemplateRenderer) getCacheKey(module, key string) string {
	return fmt.Sprintf("%s:%s", module, key)
}

// getTemplatesDirectoryPath gets the templates directory path
// This is needed in case this is called from a package outside of main, such as testing
func getTemplatesDirectoryPath() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Join(filepath.Dir(d), config.TemplateDir)
}
