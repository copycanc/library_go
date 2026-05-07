package in_memory

import (
	"fmt"
	"rest_library/internal/domain/models"

	"golang.org/x/crypto/bcrypt"
)
import "rest_library/internal/domain/errors"
import "github.com/google/uuid"

type Storage struct {
	books map[string]models.Book
	user  map[string]models.User
}

func NewStorage() *Storage {
	return &Storage{
		books: make(map[string]models.Book),
		user:  make(map[string]models.User),
	}
}

func (s *Storage) GetBooksList() ([]models.Book, error) {
	var bookList []models.Book
	if len(s.books) == 0 {
		return nil, errors.ErrBooksListIsEmpty
	}
	for _, book := range s.books {
		bookList = append(bookList, book)
	}
	return bookList, nil
}

func (s *Storage) SaveBook(book models.Book) error {
	for key, b := range s.books {
		if b.Author == book.Author && b.Lable == book.Lable {
			mBook := s.books[key]
			mBook.Count++
			s.books[key] = mBook
			return nil
		}

	}
	bookID := uuid.New().String()
	book.Count = 1
	book.BookID = bookID
	s.books[bookID] = book
	return nil
}

func (s *Storage) RegisterUser(user models.User) error {
	for _, u := range s.user {
		if u.Email == user.Email {
			return errors.ErrUserExist
		}
	}
	userID := uuid.New().String()
	user.UserID = userID
	s.user[userID] = user
	return nil
}

func (s *Storage) GetUser(user models.UserLogin) (models.User, error) {
	for _, dbUser := range s.user {
		if dbUser.Email == user.Email {
			if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
				return models.User{}, errors.ErrInvalidCreds
			}
			return dbUser, nil
		}
	}
	return models.User{}, errors.ErrInvalidCreds
}

func (s *Storage) PrintInfoBook(bookID string) (models.Book, error) {
	fmt.Println(s.books)
	book, exists := s.books[bookID]
	if !exists {
		return models.Book{}, errors.ErrBookNotFound
	}
	return book, nil
}
