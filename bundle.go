package mirror

import (
	"errors"
	"github.com/petomalina/mirror/pkg/bundle"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/tools/go/packages"
	"os"

	. "github.com/petomalina/mirror/pkg/logger"
)

// RunFunc is a callback that will be called when the app bootstrap finishes
type RunFunc func(StructSlice, *Out, *packages.Package) error

// Out is an alias for the underlying bundle.Out type, hidden with its implementation details
type Out = bundle.Out

// CreateDefaultApp returns default flag configuration for bundled apps
func CreateDefaultApp(name string, runFunc RunFunc) *cli.App {
	// override the version flag so we can use -v for verbosity
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version",
		Usage: "Prints the version of the cli",
	}

	app := cli.NewApp()
	app.Name = name
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "pkg, p",
			Value: ".",
			Usage: "Package to be used for model determination",
		},
		cli.StringSliceFlag{
			Name:   "models, m",
			Usage:  "Models that should be considered when generating, 'all' for all that can be found",
			EnvVar: "MIRROR_MODELS",
		},
		cli.StringFlag{
			Name:  "out, o",
			Value: ".",
			Usage: "Directory for the generated files to be saved in",
		},
		cli.StringFlag{
			Name:   "verbosity, v",
			Value:  "info",
			Usage:  "Sets the logging level for the bundle",
			EnvVar: "MIRROR_LOG_LEVEL",
		},
		cli.BoolFlag{
			Name:  "generateSymbols, x",
			Usage: "(experimental) Defines if symbols should be generated automatically or not",
		},
		cli.BoolFlag{
			Name:  "preserveCache, c",
			Usage: "(experimental) Preserves the cache after the build for further examination",
		},
	}
	app.Action = func(c *cli.Context) error {
		logLevel, err := logrus.ParseLevel(c.String("verbosity"))
		if err != nil {
			L.
				Method("Bundle", "CreateDefaultApp").
				Errorln("An error occurred when parsing log level: ", err.Error())
			return err
		}
		L.SetLevel(logLevel)

		// validate the generateSymbols flag, so we can warn the user beforehand
		if c.Bool("generateSymbols") {
			if len(c.StringSlice("models")) == 0 {
				return errors.New("can't use generateSymbols (-x) without specifying models")
			}

			if len(c.StringSlice("models")) == 1 && c.StringSlice("models")[0] == "all" {
				return errors.New("using 'all' with generateSymbols (-x) is forbidden - can't find models")
			}
		}

		// create the internal bundle that will be run
		b := bundle.Bundle{}

		symbols, pkg, err := b.Run(
			c.String("pkg"),
			c.StringSlice("models"),
			bundle.RunOptions{
				GenerateSymbols: c.Bool("generateSymbols"),
				PreserveCache:   c.Bool("preserveCache"),
			},
		)

		if err != nil {
			L.
				Method("Bundle", "CreateDefaultApp").
				Errorln("An error occurred when running the generator: ", err.Error())
		}

		return runFunc(
			ReflectStructs(symbols...).Each(func(s *Struct) {
				s.OriginalPackage = pkg.PkgPath
			}),
			bundle.NewOut(c.String("out")),
			pkg,
		)
	}

	return app
}

// RunDefaultApp will automatically run the defaultly bundled application
func RunDefaultApp(name string, runFunc RunFunc) error {
	L.Method("Bundle", "RunDefaultApp").Trace("Invoked  with os args: ", os.Args)
	return CreateDefaultApp(name, runFunc).Run(os.Args)
}
