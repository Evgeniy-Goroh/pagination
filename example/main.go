package main

import (
	"fmt"
	"github.com/Evgeniy-Goroh/pagination"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("example.html", "example.html")
	router.GET("/example", func(c *gin.Context) {
		p:= Paginator.New(50)
		fmt.Println(p)

		/*
		Get data on paginator

		list := Paging(&Param{
			DB:      orm.DB,
			Page:    page,
			Limit:   limit,
			OrderBy: []string{"id asc"},
			ShowSQL: true,
		}, &LinkListStruct)
		*/

		c.HTML(http.StatusOK, "example.html", gin.H{
			"page": p,
		})
	})
	router.Run(":8080")
}
