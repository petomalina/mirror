package bundle

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

// Bundle is a set of templates and logic packed for a purpose of
// code generation that is generalized (into the bundle)
type Bundle struct {
	RunFunc BundleRunFunc
}

type BundleRunFunc func(outDir string, models []interface{}) error

func (b *Bundle) Run(pkgName string, symbols []string, outDir string, generateSymbols, preserveCache bool) error {
	L.Method("Bundle", "Run").Trace("Invoked, should load: ", symbols, "from: ", pkgName)

	pkg, err := findPackage(pkgName)
	if err != nil {
		return err
	}

	// copy to the cache dir
	pkgCacheDir, err := copyPackageToCache(pkg)
	if err != nil {
		return err
	}

	// automatically generate all symbols, creating e.g. XUser from User etc.
	if generateSymbols {
		err = generateSymbolsForModels(symbols, pkgCacheDir)
		if err != nil {
			return err
		}

		// mutate to match the symbol prefix
		for i := range symbols {
			symbols[i] = "X" + symbols[i]
		}
	}

	// remove the cache dir once we are done
	// only if user doesn't want to preserve it
	if !preserveCache {
		L.Method("Bundle", "Run").Trace("Cache will be deleted, as preserveCache is not set")
		defer func(dir string) {
			L.Method("Bundle", "Run").Trace("Removing cache dir: ", dir)
			if err := os.RemoveAll(dir); err != nil {
				L.Method("Bundle", "Run").Warn("An error occurred when removing cache dir: " + dir)
			}
		}(pkgCacheDir)
	}

	L.Method("Bundle", "Run").Trace("Building the plugin: ", pkgCacheDir)
	objPath, err := BuildPlugin(pkgCacheDir)
	if err != nil {
		return err
	}

	// remove the object model when we are done
	// this object resides in the .mirror folder, that's why it's not
	// cleared by the copy cleaner
	defer func(oPath string) {
		L.Method("Bundle", "Run").Trace("Removing object model: ", oPath)
		if err := os.Remove(oPath); err != nil {
			L.Method("Bundle", "Run").Warn("An error occurred when removing: " + oPath)
		}
	}(objPath)

	L.Method("Bundle", "Run").Trace("Opening the plugin: ", objPath)
	models, err := LoadPluginSymbols(objPath, symbols)
	if err != nil {
		return err
	}
	L.Method("Bundle", "Run").Trace("Loaded symbols: ", models)

	// remove the file that was generated by the plugin build
	return b.RunFunc(outDir, models)
}

// CreateDefaultApp returns default flag configuration for bundled apps
func (b *Bundle) CreateDefaultApp(name string) *cli.App {
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

		err = b.Run(
			c.String("pkg"),
			c.StringSlice("models"),
			c.String("out"),
			c.Bool("generateSymbols"),
			c.Bool("preserveCache"),
		)

		if err != nil {
			L.
				Method("Bundle", "CreateDefaultApp").
				Errorln("An error occurred when running the generator: ", err.Error())
		}

		return err
	}

	return app
}

// RunDefaultApp will automatically run the defaultly bundled application
func (b *Bundle) RunDefaultApp(name string) error {
	L.Method("Bundle", "RunDefaultApp").Trace("Invoked  with os args: ", os.Args)
	return b.CreateDefaultApp(name).Run(os.Args)
}
