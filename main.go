package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	app := iris.New()
	app.Get("/api/{apiCall:path}", func(ctx iris.Context) {
		endpoint := "https://jsonplaceholder.typicode.com/"

		apiCall := ctx.Params().Get("apiCall")
		if len(apiCall) > 0 {
			endpoint = fmt.Sprintf("%s%s", endpoint, apiCall)
		}

		rawQuery := ctx.Request().URL.RawQuery
		if len(rawQuery) > 0 {
			endpoint = fmt.Sprintf("%s?%s", endpoint, rawQuery)
		}
		
		/*
			https://jsonplaceholder.typicode.com/todos(/1)
			https://jsonplaceholder.typicode.com/posts(/1)
		*/
		resp, err := http.Get(endpoint)
		if err != nil {
			ctx.JSON(iris.Map{"success": false, "error_message": err.Error()})
			return
		}
		defer resp.Body.Close()
		if (resp.StatusCode == 200 || resp.StatusCode == 201) {
			body, _ := ioutil.ReadAll(resp.Body)
			ctx.ContentType(context.ContentJSONHeaderValue)
			ctx.Write(body)
		} else {
			ctx.JSON(iris.Map{"success": false, "error_message": "target not ok"})
		}
	})
	err := app.Run(
		// Start the web server at localhost:8080
		iris.Addr(":9000"),
		// skip err server closed when CTRL/CMD+C pressed:
		iris.WithoutServerError(iris.ErrServerClosed),
		// enables faster json serialization and more:
		iris.WithOptimizations,
	)

	if err != nil {
		log.Println(err.Error())
	}
}