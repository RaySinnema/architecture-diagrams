package main

import "os"

type TextExporter interface {
	export(model ArchitectureModel, printer *Printer) error
}

func Export(model ArchitectureModel, exporter TextExporter, fileName string) error {
	printer := NewPrinter()
	err := exporter.export(model, printer)
	if err != nil {
		return err
	}
	data := []byte(printer.String())
	return os.WriteFile(fileName, data, 0666)
}
