package mirror

import (
	"fmt"
	"github.com/urfave/cli"
	"math/rand"
	"os"
	"os/exec"
	"plugin"
)

// Bundle is a set of templates and logic packed for a purpose of
// code generation that is generalized (into the bundle)
type Bundle struct {
	RunFunc BundleRunFunc
}

type BundleRunFunc func(outDir string, models []interface{}) error

func (b *Bundle) Run(pkg string, symbols []string, outDir string) error {
	// random file name so we'll get unique loader each time
	uniq := rand.Int()

	objPath := fmt.Sprintf("mirror-%d.so", uniq)

	// create the plugin from the passed package
	err := WithChangedPackage(pkg, "main", func() error {
		cmd := exec.Command("go", "build", "-buildmode=plugin", "-o="+objPath, pkg)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	})
	if err != nil {
		return err
	}

	// remove the object model when we are done
	defer func(oPath string) {
		if err := os.Remove(oPath); err != nil {
			fmt.Println("An error occurred when removing " + oPath)
		}
	}(objPath)

	p, err := plugin.Open(objPath)
	if err != nil {
		return err
	}

	// add model symbols that were loaded from the built plugin
	models := []interface{}{}
	for _, symName := range symbols {
		sym, err := p.Lookup(symName)
		if err != nil {
			return err
		}

		models = append(models, sym)
	}

	// remove the file that was generated by the plugin build
	return b.RunFunc(outDir, models)
}

// CreateDefaultApp returns default flag configuration for bundled apps
func (b *Bundle) CreateDefaultApp(name string) *cli.App {
	app := cli.NewApp()
	app.Name = name
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
		return b.Run(
			c.String("pkg"),
			c.StringSlice("models"),
			c.String("out"),
		)
	}

	return app
}

// RunDefaultApp will automatically run the defaultly bundled application
func (b *Bundle) RunDefaultApp(name string) error {
	return b.CreateDefaultApp(name).Run(os.Args)
}
