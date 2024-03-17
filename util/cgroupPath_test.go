package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var subSystemNameIns = []string{
	"memory",
	"cpu",
	"cpuset",
}

func TestFindCgroupRootPath(t *testing.T) {

	for _, subSysSubSystemName := range subSystemNameIns {
		_, err := FindCgroupRootPath(subSysSubSystemName)
		require.NoError(t, err)
	}
}

func Testutil.GetCgroupPath(t *testing.T) {
	for _, subSysSubSystemName := range subSystemNameIns {
		path, err := util.GetCgroupPath(subSysSubSystemName, "myTest", true)
		fmt.Println(path)
		require.NoError(t, err)
	}
}
