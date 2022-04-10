package main

import (
	"BookService/protobuff"
	"context"

	"github.com/sirupsen/logrus"
)

//grpcurl --plaintext -d '{"Authors":"Dave Cheney"}' localhost:9092 BookService.SearchAuthor
//grpcurl --plaintext -d '{"BookName":"Golang"}' localhost:9092 BookService.SearchBook

//entity for Database
type BookDB struct {
	Id      int
	Name    string
	Authors string
}

//service's entity
type Book struct {
	Name    string
	Authors []string
}

type Server struct {
	log *logrus.Logger
	protobuff.UnimplementedBookServiceServer
}

func NewMyServer(l *logrus.Logger) *Server {
	return &Server{l, protobuff.UnimplementedBookServiceServer{}}
}

//function for converting db's Book entity to service's Book entity
func BookConvert(bookdb BookDB) (book Book) {
	return Book{Name: bookdb.Name, Authors: []string{bookdb.Authors}}
}

//Searching authors by book's name
func (s *Server) SearchBook(ctx context.Context, req *protobuff.BookName) (*protobuff.Books, error) {
	book := req.GetBookName()
	//preparing query for sanitizing purposes
	stmt, err := db.Prepare(`SELECT books.book_id, books.name, authors.name
	FROM books
		INNER JOIN book_author
		ON books.book_id = book_author.book_id
		INNER JOIN authors
		ON authors.author_id = book_author.author_id
        where books.name = ?
        group by books.book_id,authors.name;`)
	//closing to free resources
	defer stmt.Close()
	if err != nil {
		s.log.Error("Error with SearchBook request:" + err.Error())
		return nil, err
	}
	rows, err := stmt.Query(book)
	if err != nil {
		s.log.Error("Error with SearchBook request:" + err.Error())
		return nil, err
	}
	//closing to free resources
	defer rows.Close()
	var id int
	var tmpBook = BookDB{}
	var respBook = Book{}
	var books = []Book{}
	//looping through returned rows
	for rows.Next() {
		rows.Scan(&tmpBook.Id, &tmpBook.Name, &tmpBook.Authors)
		//checking to determine if it's first row or not
		if id == 0 {
			respBook = BookConvert(tmpBook)
		} else {
			// we got new book with different id - need to add book into response
			if id != tmpBook.Id {
				books = append(books, respBook)
				respBook = BookConvert(tmpBook)
			} else {
				//id stays the same - just adding authors of the book
				respBook.Authors = append(respBook.Authors, tmpBook.Authors)
			}
		}
		//updating id
		id = tmpBook.Id
	}
	//adding last found book
	books = append(books, respBook)
	//manually closing for better assurance
	err = rows.Close()
	if err != nil {
		s.log.Error("Error with SearchBook request:" + err.Error())
		return nil, err
	}
	//creating proper response struct
	response := &protobuff.Books{}
	for _, v := range books {
		response.Book = append(response.Book, &protobuff.BookInfo{Name: v.Name, Authors: v.Authors})
	}
	s.log.Info("Successful SearchBook response")
	return response, nil
}

//searching books by author's name
func (s *Server) SearchAuthor(ctx context.Context, req *protobuff.SearchAuthorRequest) (*protobuff.Books, error) {
	author := req.GetAuthors()
	//preparing query for sanitizing purposes
	stmt, err := db.Prepare(`SELECT books.book_id, books.name, authors.name
	FROM books
		INNER JOIN book_author
		ON books.book_id = book_author.book_id
		INNER JOIN authors
		ON authors.author_id = book_author.author_id
        where authors.name = ?
        group by books.book_id,authors.name;`)
	//closing to free resources
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(author)
	if err != nil {
		return nil, err
	}
	var id int
	var tmpBook = BookDB{}
	var respBook = Book{}
	var books = []Book{}
	//closing to free resources
	defer rows.Close()
	//looping through returned rows
	for rows.Next() {
		rows.Scan(&tmpBook.Id, &tmpBook.Name, &tmpBook.Authors)
		//checking to determine if it's first row or not
		if id == 0 {
			respBook = BookConvert(tmpBook)
		} else {
			// we got new book with different id - need to add book into response
			if id != tmpBook.Id {
				books = append(books, respBook)
				respBook = BookConvert(tmpBook)
			} else {
				//id stays the same - just adding authors of the book
				respBook.Authors = append(respBook.Authors, tmpBook.Authors)
			}
		}
		//updating id
		id = tmpBook.Id
	}
	//adding last found book
	books = append(books, respBook)
	//manually closing for better assurance
	err = rows.Close()
	if err != nil {
		return nil, err
	}
	//creating proper response struct
	response := &protobuff.Books{}
	for _, v := range books {
		response.Book = append(response.Book, &protobuff.BookInfo{Name: v.Name, Authors: v.Authors})
	}
	s.log.Info("Successful SearchAuthor response")
	return response, nil
}
