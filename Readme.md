# Go SimpleStorage

If you are reading this in github be aware this s mirrored from Gitlab @ gitlab.com/tulpenhaendler/simplestore


This is a Key-Value Store that stores stuff into one single file,

not good for anything complex,
but good if you have a small app and only need to store a little bit of data.


## Installing

```
go get gitlab.com/tulpenhaendler/simplestore
```

## Example:

```
package main

import (
	"fmt"
	"gitlab.com/tulpenhaendler/simplestore"
)

func main(){
	var config *SimpleStorage.Config
	config = nil
	appname := "myapp"
	
	config = &SimpleStorage.Config{

		// Default: $HOME/.$APPNAME/storage
		// Overwritten by env(STORAGE_DIR)
		StorageDir: "./simpstore",

	}
	
	s := SimpleStorage.NewSimpleStorage(appname, config)
	// Example of how to Store and Get basic types:

	s.StoreString("key","value")
	s.StoreUint64("bignumber",9223372036854775805)

	fmt.Println( s.GetFloat64("test") )

	// Example of how to set and Get anything:
	type test struct{
		First string
		Second string
	}
	set := test{
		First:"hi",
		Second: "world",
	}
	s.StoreInterface("structexample",&set)

	get := test{}
        s.GetInterface("structexample",&get)
	fmt.Println( get )
}

```
