// Code generated by go-swagger; DO NOT EDIT.

package main

import (
	"log"
	"os"

	"github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"

	internal "github.com/kraken-hpc/imageapi/internal/api"
	"github.com/kraken-hpc/imageapi/restapi"
	"github.com/kraken-hpc/imageapi/restapi/operations"
)

// This file was generated by the swagger tool.
// Make sure not to overwrite this file after you generated it because all your edits would be lost!

func main() {
	////////////////////////
	// BEGIN CUSTOM CODE //
	//////////////////////

	internal.ForkInit()

	//////////////////////
	// END CUSTOM CODE //
	////////////////////

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewImageapiAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "Image API"
	parser.LongDescription = "This API specification describes a service for attaching, mounting and preparing container images and manipulating those containers.\n\nIn general, higher level objects can either reference lower level objects (e.g. a mount referencing an attachment point) by a reference ID, \nor, they can contain the full specification of those lower objects.\n\nIf an object references another by ID, deletion of that object does not effect the underlying object.\n\nIf an object defines a lower level object, that lower level object will automatically be deleted on deletion of the higher level object.\n\nFor instance, if a container contains all of the defintions for all mount points and attachments, deletion of the container will automatically unmount\nand detach those lower objects.\n"
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}
