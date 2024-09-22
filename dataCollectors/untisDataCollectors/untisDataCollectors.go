package untisDataCollectors

import (
	"fmt"

	"github.com/Mr-Comand/goUntisAPI/structs"
	"github.com/Mr-Comand/goUntisAPI/untisApi"
)

func Init(apiConfig structs.ApiConfig) {
	c := untisApi.NewClient(apiConfig)
	err := c.Authenticate()
	if err != nil {
		fmt.Println("Error authenticating:", err)
	}
	c.Test()
	c.Logout()
}
