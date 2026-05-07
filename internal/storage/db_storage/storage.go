package db_storage

import (
	"context"
	"rest_library/internal/domain/models"
	"time"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage() (*Storage, error) {
	conn, err := pgx.Connect(context.Background(), "conn str") //TODO Добавить строку подключения
	if err != nil {
		return nil, err
	}
	return &Storage{
		conn: conn,
	}, nil
}

func (s *Storage) GetBooksList() ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	rows, err := s.conn.Query(ctx, "SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(&book.BookID, &book.Author, &book.Lable, &book.Description, &book.Genre, &book.WritedAt, &book.Count)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, nil
}

func (s *Storage) SaveBook(book models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.conn.Exec(ctx, "INSERT INTO books (id, author, lable, description, genre, writed_at, count) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		book.BookID, book.Author, book.Lable, book.Description, book.Genre, book.WritedAt, book.Count)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetUser(user models.UserLogin) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var userM models.User

	row := s.conn.QueryRow(ctx, "SELECT * FRPM users WHERE email = $1", user.Email)

	err := row.Scan(&userM.UserID, &userM.Name, &userM.Age, &userM.Email, &userM.Password)
	if err != nil {
		return models.User{}, err
	}
	return userM, nil
}

func (s *Storage) RegisterUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := s.conn.Exec(ctx, "INSERT INTO user (id, name, age, email, password) VALUES ($1,$2,$3,$4,$5)",
		user.UserID, user.Name, user.Age, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) PrintInfoBook(bookID string) (models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var book models.Book

	row := s.conn.QueryRow(ctx, "SELECT * FROM books WHERE id=$1", bookID)

	err := row.Scan(&book.BookID, &book.Author, &book.Lable, &book.Description, &book.Genre, &book.WritedAt, &book.Count)
	if err != nil {
		return models.Book{}, err
	}
	return book, nil
}
