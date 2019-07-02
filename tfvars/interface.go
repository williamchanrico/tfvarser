package tfvars

import (
	"io"
	"os"
	"path"
	"text/template"
)

// Generator  for packages that want to produce tfvars file
type Generator interface {
	Name() string
	Kind() string
	Template() string
	Execute(io.Writer, *template.Template) error
}

// Service struct maps provider+object to its corresponding generator/template
type Service struct {
	generator Generator
	tmpl      *template.Template
}

// New return new tfvars generator
func New(gen Generator, funcMap ...template.FuncMap) *Service {
	var tmpl *template.Template
	if len(funcMap) > 0 {
		tmpl = template.New(gen.Name()).Funcs(funcMap[0])
	} else {
		tmpl = template.New(gen.Name())
	}

	return &Service{
		generator: gen,
		tmpl:      template.Must(tmpl.Parse(gen.Template())),
	}
}

// Generate a tmpl from initialized Generator (name is a concat of provider+object)
func (s *Service) Generate(dir, filename string) error {
	err := makeDirIfNotExists(dir)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(dir, filename))
	defer f.Close()
	if err != nil {
		return err
	}

	return s.generator.Execute(f, s.tmpl)
}
