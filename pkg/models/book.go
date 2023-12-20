package models

import (
	"errors"
	"fmt"

	"github.com/Alb3G/go-bookstore/pkg/config"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model
	Name        string `gorm:"" json:"name"`
	Author      string `json:"author"`
	Publication string `json:"aublication"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Book{})
}

// Asi lo que estamos haciendo es indicarle a Go que esto es un metodo de Book
func (b *Book) CreateBook() *Book {
	db.NewRecord(b)
	db.Create(&b)
	return b
}

func GetAllBooks() []Book {
	var Books []Book
	db.Find(&Books)
	return Books
}

func GetBookById(Id int64) (*Book, *gorm.DB, error) {
	var book Book
	db := db.Where("Id=?", Id).Find(&book)
	// Si db.Error fuera nil solo indica que la consulta se hizo con exito no significa que
	// traiga un libro valido.
	if db.Error != nil {
		// gorm.ErrRecordNotFound error predefinido de Gorm cuando no se encuentra lo que queremos
		// en la consulta realizada.
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			return nil, db, fmt.Errorf("book not found with id: %d", Id)
		}
		return nil, db, db.Error
	}
	return &book, db, nil
}

func DeleteBook(Id int64) (*Book, error) {
	var book Book
	result := db.Where("Id=?", Id).Delete(&book)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("book not found with id: %d", Id)
	}
	return &book, nil
}
