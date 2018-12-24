package bundle

import "golang.org/x/tools/go/packages"
import . "github.com/petomalina/mirror/pkg/logger"

// Out is an encapsulation of methods used to write to the output directory
// which must represent the package
type Out struct {
	pkgPath string
	pkg     *packages.Package
}

// NewOut creates a new Out wrapper
func NewOut(pkgPath string) *Out {
	pkg, err := packages.Load(&packages.Config{
		Mode:  packages.LoadFiles,
		Tests: false,
	}, pkgPath)

	if err != nil {
		L.Method("Out", "NewOut").Warnln("An error occurred when loading output package:", err)
	}

	return &Out{
		pkgPath: pkgPath,
		pkg:     pkg[0],
	}
}
