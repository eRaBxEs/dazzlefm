package model

import (
	"encoding/json"
	"time"
)

// Presenter ...
type Presenter struct {
	ID          int64           `json:"id"`
	Image       json.RawMessage `json:"image"`
	Name        string          `json:"name"`
	Email       string          `json:"email"`
	Description string          `json:"description"`
	Status 		int 			`json:"status"`
}

// PersonRole ...
type PersonRole struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// News ...
type News struct {
	ID      int64           `json:"id"`
	Title   string          `json:"title"`
	Date    time.Time       `json:"date"`
	Author  string          `json:"author"`
	Content string          `json:"content"`
	Image   json.RawMessage `json:"image"`
}

// Schedule ...
type Schedule struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	PresenterID   int64     `json:"presenter_id"`
	CoPresenterID int64     `json:"copresenter_id" gorm:"column:copresenter_id"`
	Day           int       `json:"day"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
}

// User ...
type User struct {
	ID         int64  `json:"id"`
	Role       int    `json:"role"`
	UserName   string `json:"user_name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	Surname    string `json:"surname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

// Page ...
type Page struct {
	ID         int64           `json:"id"`
	Name       string          `json:"name"`
	Route      string          `json:"route"`
	Attributes json.RawMessage `json:"attributes"`
}

// Login sent during login
type Login struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// Contact ...
type Contact struct {
	ID      int64     `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Subject string    `json:"subject"`
	Date    time.Time `json:"date"`
	Message string    `json:"message"`
}

// Gallery ...
type Gallery struct {
	ID    int64           `json:"id"`
	Image json.RawMessage `json:"image"`
}

// Music ...
type Music struct {
	ID        int64           `json:"id"`
	Title     string          `json:"title"`
	Artist    string          `json:"artist"`
	UploadID  int64           `json:"upload_id"`
	MusicFile json.RawMessage `json:"music_file"`
	Rank      int             `json:"rank"`
}

// Upload ...
type Upload struct {
	ID   int64  `json:"id"`
	Path string `json:"path"`
	Size int64  `json:"size"`
}

// Event ...
type Event struct {
	ID     int64     `json:"id"`
	Title  string    `json:"title"`
	Venue  string    `json:"venue"`
	Status int       `json:"status"`
	Date   time.Time `json:"date"`
	Link   string    `json:"link"`
}

// Script ...
type Script struct {
	ID        int64  `json:"id"`
	Source    string `json:"src"`
	Async     bool   `json:"async"`
	Defer     bool   `json:"defer"`
	InnerHTML string `json:"innerHTML"`
}

// Settings ...
type Settings struct {
	StreamURL string `json:"stream_url"`
}
