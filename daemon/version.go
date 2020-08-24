package daemon

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)
type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

func (v *Version) ToString() string{
	strs := []string {
		strconv.FormatInt(int64(v.Major), 10),
		strconv.FormatInt(int64(v.Minor), 10),
		strconv.FormatInt(int64(v.Patch), 10),
	}
	return strings.Join(strs, ".")
}

// This function is used to parse the value from version.txt file
/*
	ex: VERSION_MAJOR=1 will return 1.
 */
func getValue( line string ) string{
	strs := strings.Split(line, "=")
	return strs[1]
}
// read a version from a file, then override the config flag.
// return a version struct, release date, and error if any.
func ReadVersionFromFile( file string ) (Version, string, error) {
	data, err := os.Open( file )
	if err != nil {
		logrus.Errorf("Failed to open %s %s", file, err)
		return Version{}, "", err
	}
	ver := Version{}
	var release string
	scanner := bufio.NewScanner(data)
	/*
	   VERSION_MAJOR=1
	   VERSION_MINOR=0
	   VERSION_PATCH=0
	   RELEASE_DATE=2020.08.30
	*/
	// 4 lines.
	for i := 0; i < 4; i++ {
		if scanner.Scan() {
			value := getValue( scanner.Text())
			if i != 3{
				rawUint8, err := strconv.Atoi(value)
				if  err != nil {
					logrus.Errorf("Failed to convert take version %s", err)
					return Version{}, "", err
				}
				if i == 0 {
					ver.Major = uint8(rawUint8)
				} else if i == 1 {
					ver.Minor = uint8(rawUint8)
				} else {
					ver.Patch = uint8(rawUint8)
				}
			} else {
				release = value
			}
		}
	}
	return ver, release, nil
}
