package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/satori/uuid"

	"google.golang.org/grpc"

	pb "github.com/Javlon2000/gRPC-CRUD-App/proto"
)

const (
	HOST = "localhost"
	PORT = 5432
	USER = ""
	PASSWORD = ""
	DBNAME = "go_grpc"
)

type server struct {
	conn *sql.DB
	pb.UnimplementedToDoServiceServer
}

func (s *server) Create(ctx context.Context, req *pb.CreateRequest) (*pb.ToDo, error) {

	db := s.conn
	id := uuid.NewV4()
	req.ToDo.Id = id.String()

	title := req.GetToDo().GetTitle()
	description := req.GetToDo().GetDescription()
	// completed := req.GetToDo().GetCompleted()

	sqlInsert := `insert into toDo (id, title, description) values ($1, $2, $3);`

	if _, err := db.Exec(sqlInsert, id, title, description); err != nil {
		return nil, errors.Wrapf(err, "To-Do couldn't be inserted")
	}

	return req.ToDo, nil	
}

func (s *server) Read(ctx context.Context, req *pb.ReadRequest) (*pb.ToDo, error) {
	
	db := s.conn
	todo_id := req.Id

	var title, description, id string
	var completed bool

	var object pb.ToDo

	sqlStatement := `select * from toDo where id = $1`
	// var books []*pb.Book
	err := db.QueryRow(sqlStatement, todo_id).Scan(
		&object.Id, 
		&object.Title, 
		&object.Description,
		&object.Completed,
	)

	if err != nil {
		errors.Wrapf(err, "To-Do list can not obtained")
	}

	fmt.Println(id, title, description, completed)


	res := &pb.ToDo{
		Id:        		object.Id,
		Title: 			object.Title,
		Description:  	object.Description,
		Completed:		object.Completed,
	}

	fmt.Println(todo_id)

	return res, nil
}

func (s *server) Update(ctx context.Context, req *pb.UpdateRequest) (*pb.ToDo, error) {
	
	db := s.conn
	
	sqlStatement := `update toDo set title=$2, description=$3, completed=$4 where id=$1`
	
	if _, err := db.Exec(sqlStatement, req.ToDo.Id, req.ToDo.Title, req.ToDo.Description, req.ToDo.Completed); err != nil {
		return nil, err
	}

	return req.ToDo, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.Empty, error) {
	
	db := s.conn
	
	sqlStatement := `delete from toDo where id = $1`

	id := req.GetId()

	if _, err := db.Exec(sqlStatement, id); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *server) ReadAll(ctx context.Context, req *pb.ReadAllRequest) (*pb.ReadAllResponse, error) {
	
	db := s.conn
	
	sqlStatement := `select id, title, description, completed from toDo`
	
	result, err := db.Query(sqlStatement)

	if err != nil {
		fmt.Println(err)
	}

	defer result.Close()

	res := []*pb.ToDo{}

	for result.Next() {

		var id, title, description string
		var  completed bool

		if err = result.Scan(&id, &title, &description, &completed); err != nil {
			
			errors.Wrap(err, "ToDos couldn't be read")
		}

		todo := pb.ToDo {
			Id:			  id,
			Title:        title,
			Description:  description,
			Completed:    completed,
		}

		fmt.Println(todo)



		res = append(res, &todo)
	}

	to_dos := pb.ReadAllResponse{ Todos: res }
	
	return &to_dos, nil
}

func main() {

	fmt.Println("Welcome to the server")

	lis, err := net.Listen("tcp", ":9500")

	if err != nil {
		errors.Wrap(err, "To-Do list can not obtained")
	}

	s := grpc.NewServer()

	conn, err := grpc.Dial("localhost: 9000", grpc.WithInsecure())
	
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	defer conn.Close()

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, USER, PASSWORD, DBNAME)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		errors.Wrap(err, "To-Do list can not obtained")
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		errors.Wrap(err, "To-Do can not be listed")
	}

	pb.RegisterToDoServiceServer(s, &server{conn: db})

	if err := s.Serve(lis); err != nil {
		errors.Wrap(err, "To-Do list can not obtained")
	}
}
