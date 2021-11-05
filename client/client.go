package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/Javlon2000/gRPC-CRUD-App/proto"
)

var userConnection pb.ToDoServiceClient
var readConnection pb.ReadAllRequest

func CreateToDo(c *gin.Context) {

	var newToDoList pb.ToDo

	if err := c.ShouldBindJSON(&newToDoList); err != nil {
		return
	}

	req := &pb.CreateRequest{
		ToDo: &pb.ToDo{
			Id:        		newToDoList.Id,
			Title: 			newToDoList.Title,
			Description:  	newToDoList.Description,
		},
	}
	res, err := userConnection.Create(context.Background(), req)
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusCreated, res)
}

func Read(c *gin.Context) {

	id := c.Param("id")

	req := &pb.ReadRequest{
		Id: id,
	}

	res, err := userConnection.Read(context.Background(), req)
	
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, res)
}

func UpdateToDo(c *gin.Context) {

	var update pb.ToDo

	if err := c.ShouldBindJSON(&update); err != nil {
		return
	}

	req := &pb.UpdateRequest{
		ToDo: &pb.ToDo{
			Id: 				update.Id,
			Title:  			update.Title,
			Description:        update.Description,
			Completed:       	update.Completed,
		},
	}

	res, err := userConnection.Update(context.Background(), req)
	
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, res)
}

func Delete(c *gin.Context) {

	id := c.Param("id")
	
	req := &pb.DeleteRequest{
		Id: id,
	}
	
	res, err := userConnection.Delete(context.Background(), req)
	
	if err != nil {
		panic(err)
	}
	
	c.IndentedJSON(http.StatusOK, res)
}

func ReadAll(c *gin.Context) {

	req := &pb.ReadAllRequest{}
	
	res, err := userConnection.ReadAll(context.Background(), req)

	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, res)
}


func main() {

	fmt.Println("Welcome Client")

	conn, err := grpc.Dial("localhost:9500", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Couldn't connect %v", err)
	}

	defer conn.Close()

	userConnection = pb.NewToDoServiceClient(conn)

	if err != nil {
		log.Fatalf("Couldn't connect %v", err)
	}

	router := gin.Default()

	router.POST("/todo", CreateToDo)
	router.GET("/todo/:id", Read)
	router.PUT("/todo", UpdateToDo)
	router.DELETE("/todo/:id", Delete)
	router.GET("/todos", ReadAll)

	router.Run("localhost:5000")
}