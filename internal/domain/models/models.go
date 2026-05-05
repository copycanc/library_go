package models

type Book struct {
	BookID      string `json:"book_id,omitempty"`
	Author      string `json:"author"`
	Lable       string `json:"lable"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
	WritedAt    string `json:"writed_at"`
	Count       int    `json:"count,omitempty"`
}

type User struct {
	UserID   string `json:"user_id,omitempty"`
	Name     string `json:"name"`
	Age      int    `json:"age" validate:"gte=12"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=8"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password"`
}
