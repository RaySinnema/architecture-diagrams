package main

import "os"

type TextExporter interface {
	export(model ArchitectureModel, printer *Printer)
}

func Export(model ArchitectureModel, exporter TextExporter, fileName string) error {
	printer := NewPrinter()
	exporter.export(model, printer)
	data := []byte(printer.String())
	return os.WriteFile(fileName, data, 0666)
}
