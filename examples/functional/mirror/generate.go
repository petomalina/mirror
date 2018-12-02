package main

import (
	"github.com/petomalina/mirror"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const mapTemplate = `type _T_Slice []*_T_

type _T_MapCallback func(*_T_) *_T_

func (us _T_Slice) Map(cb _T_MapCallback) _T_Slice {
	newSlice := _T_Slice{}
	for _, o := range us {
		newSlice = append(newSlice, cb(o))
	}
	
	return newSlice
}
`

func main() {
	app := cli.NewApp()
	app.Name = "functional-generator"
	app.Usage = "generates map/filter/reduce for input models"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "pkg, p",
			Value: ".",
			Usage: "Package to be used for model determination",
		},
		cli.StringSliceFlag{
			Name:  "models, m",
			Usage: "Models that should be considered when generating",
		},
		cli.StringFlag{
			Name:  "out, o",
			Value: ".",
			Usage: "Directory for the generated files to be saved in",
		},
	}
	app.Action = func(c *cli.Context) error {
		bundle := &mirror.Bundle{}
		return bundle.Run("", []string{}, "", func(outDir string, models []interface{}) error {
			out := mirror.File(filepath.Join(outDir, "functional.go"))
			blocks := []string{}

			for _, m := range models {
				str := mirror.ReflectStruct(m)

				blocks = append(blocks, strings.Replace(mapTemplate, "_T_", str.Name(), -1))
			}

			return out.Write(blocks...)
		})
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
