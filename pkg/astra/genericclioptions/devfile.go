package genericclioptions

import (
	"fmt"

	"github.com/devfile/library/v2/pkg/devfile/parser"
	dfutil "github.com/devfile/library/v2/pkg/util"
	"github.com/spf13/cobra"

	"github\.com/danielpickens/astra/pkg/component"
	"github\.com/danielpickens/astra/pkg/devfile"
	"github\.com/danielpickens/astra/pkg/devfile/location"
	"github\.com/danielpickens/astra/pkg/devfile/validate"
	"github\.com/danielpickens/astra/pkg/testingutil/filesystem"
	astrautil "github\.com/danielpickens/astra/pkg/util"
)

func getDevfileInfo(cmd *cobra.Command, fsys filesystem.Filesystem, workingDir string, variables map[string]string, imageRegistry string) (
	devfilePath string,
	devfileObj *parser.DevfileObj,
	componentName string,
	err error,
) {
	devfilePath = location.DevfileLocation(fsys, workingDir)
	isDevfile := astrautil.CheckPathExists(fsys, devfilePath)
	if isDevfile {
		devfilePath, err = dfutil.GetAbsPath(devfilePath)
		if err != nil {
			return "", nil, "", err
		}
		// Parse devfile and validate
		var devObj parser.DevfileObj
		devObj, err = devfile.ParseAndValidateFromFileWithVariables(devfilePath, variables, imageRegistry, true)
		if err != nil {
			return "", nil, "", fmt.Errorf("failed to parse the devfile %s: %w", devfilePath, err)
		}
		devfileObj = &devObj
		err = validate.ValidateDevfileData(devfileObj.Data)
		if err != nil {
			return "", nil, "", err
		}

		componentName, err = component.GatherName(workingDir, devfileObj)
		if err != nil {
			return "", nil, "", err
		}
	}

	return devfilePath, devfileObj, componentName, nil
}
