# Architecture
The task does not specify which type of DB to consume or what kind of API framework to use. For that reason, it will focus on the loose coupling to make it easy in the future to change types of API and DB very quickly.

For the current project, I will use only HTTP and SQL, but the company will be able to change them with other technologies like gRPC, GraphQL, MongoDB etc., by creating drivers for the specific technologies.

It will be like this simple code, where it represents a brain when it says hello and asks for a name. Look how seamlessly corporates the drivers with the logic.

```go
package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/*
curl -X GET http://localhost:1323/hello
curl -X POST http://localhost:1323/whatisyourname -H "Content-Type: application/json" -d '{"name":"manos"}'
*/

func main() {
	rd := NewRamDriver()
	bc := NewBrain(rd)
	ec := NewEchoDriver(bc)
	ec.Run()
}

type BrainInterface interface {
	SayHello(context context.Context) string
	WhatIsYourName(context context.Context, name string) (string, error)
}

type MemoryRepository interface {
	RememberName(context context.Context, name string) error
}

type BrainController struct {
	memory MemoryRepository
}

func NewBrain(memory MemoryRepository) *BrainController {
	return &BrainController{
		memory: memory,
	}
}

func (bc *BrainController) SayHello(context context.Context) string {
	return "hello"
}

func (bc *BrainController) WhatIsYourName(context context.Context, name string) (string, error) {
	if len(name) == 0 {
		return "", errors.New("empty name")
	}
	err := bc.memory.RememberName(context, name)
	if err != nil {
		return "", err
	}
	say := "Nice to meet you, " + name
	return say, nil
}

type EchoDriver struct {
	brain BrainInterface
}

func NewEchoDriver(brain BrainInterface) *EchoDriver {
	return &EchoDriver{
		brain: brain,
	}
}

func (ed *EchoDriver) Run() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/hello", ed.hello)
	e.POST("/whatisyourname", ed.whatIsYourName)

	e.Logger.Fatal(e.Start(":1323"))
}

func (ed *EchoDriver) hello(c echo.Context) error {
	say := ed.brain.SayHello(c.Request().Context())
	return c.String(http.StatusOK, say)
}

type WhatIsYourNameHttpBody struct {
	Name string `json:"name"`
}

func (ed *EchoDriver) whatIsYourName(c echo.Context) error {
	body := WhatIsYourNameHttpBody{}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	reply, err := ed.brain.WhatIsYourName(c.Request().Context(), body.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusOK, reply)
}

type RamDriver struct {
	names []string
}

func NewRamDriver() *RamDriver {
	return &RamDriver{}
}

func (rd *RamDriver) RememberName(ctx context.Context, name string) error {
	if len(name) > 5 {
		return errors.New("not enough space")
	}
	rd.names = append(rd.names, name)
	return nil
}

```

The same way of corporation we will follow for the current project.
We will have an interface for the users and assets to ask for these methods.
* AddUser
* LoginUser
* AddAsset
* DeleteAsset
* UpdateAsset
* GetAsset
* ListAssets
* FavourAnAsset


For simplicity, anyone will be able to add a user. But only admin users can add/update/delete assets.
POST /api/v1/users
POST /api/v1/auth/login
POST /api/v1/auth/logout

POST 	/api/v1/auth/admin/assets
PUT 	/api/v1/auth/admin/assets
DELETE 	/api/v1/auth/admin/assets

GET /api/v1/auth/assets/:id
GET /api/v1/auth/admin/users/:user_id/assets
GET /api/v1/auth/assets
POST /api/v1/auth/assets/:id/favourite

