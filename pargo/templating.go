package main

import (
	"fmt"
	"html/template"
	"os"
)

type Template struct {
	OutputFilePath   string
	TemplateFilePath string
	TemplateFileName string
}

func (t *Template) Generate(
	outputFileName string,
	templateFileName string,
	data interface{}) (err error) {

	//parse template
	templateFile := fmt.Sprintf("%s\\%s", t.TemplateFilePath, templateFileName)
	tempTemplate, err := template.ParseFiles(templateFile)
	if err != nil {
		return &os.PathError{}
	}
	//create output directory
	_, err = os.Stat(t.OutputFilePath)
	if os.IsNotExist(err) {
		err = os.Mkdir(t.OutputFilePath, os.ModeDir)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		return err
	}

	// create generated file
	genFile := fmt.Sprintf("%s\\%s", t.OutputFilePath, outputFileName)
	file, err := os.Create(genFile)
	if err != nil {
		return err
	}
	//execture template to the generated file
	err = tempTemplate.Execute(file, data)
	if err != nil {
		return err
	}
	return nil
}
