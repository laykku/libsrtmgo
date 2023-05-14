# libsrtmgo

Load and parse the SRTM HGT data into array of lat/lon/elevation floats

##### Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/mr-marsh/libsrtmgo/srtm"
)

func main() {
	srtm.Init("https://step.esa.int/auxdata/dem/SRTMGL1/", srtm.SRTMGL1) // one arc-second srtm data

	if points, err := srtm.LoadTile(40.7128, -74.0060); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(points[0])
	}
}
```