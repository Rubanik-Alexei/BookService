package main

import (
	"BookService/protobuff"
	"context"
	"log"
	"net"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

//mocking real grpc connection
func init() {
	lis = bufconn.Listen(bufSize)
	log := logrus.New()
	ms := NewMyServer(log)
	s := grpc.NewServer()
	protobuff.RegisterBookServiceServer(s, ms)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

//testing successful SearchBook request
func TestSearchBook(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := protobuff.NewBookServiceClient(conn)
	var mock sqlmock.Sqlmock
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	prep := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT books.book_id, books.name, authors.name
	FROM books
		INNER JOIN book_author
		ON books.book_id = book_author.book_id
		INNER JOIN authors
		ON authors.author_id = book_author.author_id
        where books.name = ?
        group by books.book_id,authors.name;`))
	prep.ExpectQuery().WithArgs("Designing Data-Intensive Applications").WillReturnRows(sqlmock.NewRows([]string{"2", "Designing Data-Intensive Applications", "Martin Kleppman"}))
	resp := &protobuff.Books{}
	if resp, err = client.SearchBook(ctx, &protobuff.BookName{BookName: "Designing Data-Intensive Applications"}); err != nil {
		t.Fatalf("SearchBook failed: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {

		t.Errorf("there were unfulfilled expectations: %s", err)

	}
	log.Printf("Response: %+v", resp)
}

//testing successful SearchAuthor request
func TestSearchAuthor(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := protobuff.NewBookServiceClient(conn)
	var mock sqlmock.Sqlmock
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	prep := mock.ExpectPrepare(regexp.QuoteMeta(`SELECT books.book_id, books.name, authors.name
	FROM books
		INNER JOIN book_author
		ON books.book_id = book_author.book_id
		INNER JOIN authors
		ON authors.author_id = book_author.author_id
        where authors.name = ?
        group by books.book_id,authors.name;`))
	prep.ExpectQuery().WithArgs("Rob Pike").WillReturnRows(sqlmock.NewRows([]string{"1", "Golang", "Rob Pike"}))
	resp := &protobuff.Books{}
	if resp, err = client.SearchAuthor(ctx, &protobuff.SearchAuthorRequest{Authors: "Rob Pike"}); err != nil {
		t.Fatalf("SearchAuthor failed: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {

		t.Errorf("there were unfulfilled expectations: %s", err)

	}
	log.Printf("Response: %+v", resp)
}
