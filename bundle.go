package mirror

import (
	"errors"
	"github.com/petomalina/mirror/pkg/bundle"
	"github.com/petomalina/mirror/pkg/plugins"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/tools/go/packages"
	"os"
	"os/signal"
	"syscall"

	. "github.com/petomalina/mirror/pkg/logger"
)

// RunFunc is a callback that will be called when the app bootstrap finishes
type RunFunc func(StructSlice, *Writer, *packages.Package) error

// Writer is an alias for the underlying bundle.Writer type, hidden with its implementation details
type Writer = bundle.Writer

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
		cli.BoolFlag{
			Name:  "watch, w",
			Usage: "(experimental) Watches for file changes in the input directory and triggers the generator",
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

		pkg, err := plugins.FindPackage(c.String("pkg"))
		if err != nil {
			return err
		}

		loader := plugins.Loader{
			TargetPath:      c.String("pkg"),
			GenerateSymbols: c.Bool("generateSymbols"),
			PreserveCache:   c.Bool("preserveCache"),
		}

		// one-shot load
		if !c.Bool("watch") {
			models, err := loader.Load(c.StringSlice("models"))

			if err != nil {
				L.
					Method("Bundle", "CreateDefaultApp").
					Errorln("An error occurred when running the generator: ", err.Error())
			}

			return runFunc(
				ReflectStructs(models...).Each(func(s *Struct) {
					s.OriginalPackage = pkg.PkgPath
				}),
				bundle.NewWriter(c.String("out")),
				pkg,
			)
		} else { // watch for changes
			done := make(chan bool)
			modelsChan, errChan := loader.Watch(c.StringSlice("models"), done)

			// create a channel for notifications from the console
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

			L.
				Method("Bundle", "CreateDefaultApp").
				Info("Starting fsnotify to watch and rebuild")

		watchLoop:
			for {
				select {
				case models, ok := <-modelsChan:
					if !ok {
						break watchLoop
					}

					err := runFunc(
						ReflectStructs(models...).Each(func(s *Struct) {
							s.OriginalPackage = pkg.PkgPath
						}),
						bundle.NewWriter(c.String("out")),
						pkg,
					)

					if err != nil {
						L.
							Method("Bundle", "CreateDefaultApp").
							Errorln("An error occurred when running the generator: ", err.Error())

						done <- true
						return err
					}

				case err := <-errChan:
					done <- true
					return err

				case <-sigs:
					L.
						Method("Bundle", "CreateDefaultApp").
						Infoln("Captured exit signal, signaling watcher to stop")
					// signal the watcher to stop watching for changes
					done <- true
				}
			}
		}

		return nil
	}

	return app
}

// RunDefaultApp will automatically run the defaultly bundled application
func RunDefaultApp(name string, runFunc RunFunc) error {
	L.Method("Bundle", "RunDefaultApp").Trace("Invoked  with os args: ", os.Args)
	return CreateDefaultApp(name, runFunc).Run(os.Args)
}
