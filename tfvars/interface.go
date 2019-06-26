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
func New(gen Generator) *Service {
	return &Service{
		generator: gen,
		tmpl:      template.Must(template.New("").Parse(gen.Template())),
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
