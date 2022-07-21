package main

import (
	"fmt"
	"os"

	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
)

func main() {

	report, err := runner.Run(
		"trtl.v1.Trtl.Get",
		"localhost:4436",
		runner.WithProtoFile("proto/trtl/v1/trtl.proto", []string{}),
		runner.WithDataFromJSON(`{"key": "MjE1aktiVFpheGhpWVRsRmcyT2FyNkd0VFJv","namespace": "people"}`),
		runner.WithInsecure(true),
		runner.WithTotalRequests(100),
		runner.WithConcurrency(100),
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	printer := printer.ReportPrinter{
		Out:    os.Stdout,
		Report: report,
	}

	printer.Print("summary")

}
