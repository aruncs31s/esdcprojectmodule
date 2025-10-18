# Project Module
This is a go module for managing projects in the [ESDC Website](https://esdc.vercel.app/).

## Install 

```bash
go get github.com/aruncs31s/esdcprojectmodule
```

## Example Usage
```go
package main

import (
	project "github.com/aruncs31s/esdcprojectmodule"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	db, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	project.InitProjectModule(r, db)
	project.RegisterPublicProjectRoutes()
	r.Run()
}
```