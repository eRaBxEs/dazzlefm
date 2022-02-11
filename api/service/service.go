package service

import (
	"dazzlefm/service/form"
	"dazzlefm/service/model"
	"encoding/json"
	"mime/multipart"
)

// IOffice an interface that manages the following models:
// PracticeArea, People
type IOffice interface {
	GetSchedule(int64) (form.Schedule, error)
	ListSchedule(int, int, int) ([]form.Schedule, error)
	CountSchedule() (int, error)
	ListASchedule(int64) ([]form.Schedule, error)
	AddSchedule(*form.Schedule) (bool, form.Schedule, error)
	SaveSchedule(int64, *form.Schedule) (bool, form.Schedule, error)
	DeleteSchedule(int64) error
	GetLiveSchedule(int, string) (form.Schedule, error)
	GetScheduleSummary(int, string) ([]form.ScheduleHighlight, error)
	GetAllSchedules(int) ([]form.ScheduleHighlight, error)

	GetEvent(int64) (form.Event, error)
	ListEvent(int, int) ([]form.Event, error)
	CountEvent() (int, error)
	AddEvent(*form.Event) error
	SaveEvent(int64, *form.Event) error
	DeleteEvent(int64) error
	GetLiveEvent() (form.Event, error)
	GetOtherEvents() ([]form.Event, error)

	GetPresenter(int64) (form.Presenter, error)
	ListPresenter(int, int, int) ([]form.Presenter, error)
	CountPresenter() (int, error)
	AddPresenter(*form.Presenter) error
	SavePresenter(int64, *form.Presenter) error
	DeletePresenter(int64) error

	GetContact(int64) (form.Contact, error)
	ListContact(int, int) ([]form.Contact, error)
	CountContact() (int, error)
	AddContact(*form.Contact) error
	SaveContact(int64, *form.Contact) error
	DeleteContact(int64) error
}

// IContent an interface that manages the website content models:
// News, Resources
type IContent interface {
	GetNews(int64) (form.News, error)
	ListNews(int, int) ([]form.News, error)
	CountNews() (int, error)
	AddNews(*form.News) error
	SaveNews(int64, *form.News) error
	DeleteNews(int64) error

	GetGallery(int64) (form.Gallery, error)
	ListGallery(int, int) ([]form.Gallery, error)
	CountGallery() (int, error)
	AddGallery(*form.Gallery) error
	SaveGallery(int64, *form.Gallery) error
	DeleteGallery(int64) error

	GetPage(int64) (form.Page, error)
	ListPage() ([]form.Page, error)
	AddPage(*form.Page) error
	SavePage(int64, *form.Page) error
	DeletePage(int64) error

	GetMusic(int64) (form.Music, error)
	ListMusic() ([]form.Music, error)
	AddMusicFile(*multipart.FileHeader) (*model.Upload, error)
	AddMusicFiles(json.RawMessage) (string, error)
	AddMusic(*form.Music) error
	SaveMusic(int64, *form.Music) error
	DeleteMusic(int64) error

	GetSnippet(int64) (form.Script, error)
	ListSnippet() ([]form.Script, error)
	AddSnippet(*form.Script) error
	SaveSnippet(int64, *form.Script) error
	DeleteSnippet(int64) error

	ListSettings() (form.Settings, error)
	AddSettings(*form.Settings) error
}

// IAuthenticate an interface that manages users and user authentication:
// User
type IAuthenticate interface {
	FindUserByUserName(string) (form.User, error)
	GetUser(int64) (form.User, error)
	ListUser() ([]form.User, error)
	AddUser(*form.User) error
	SaveUser(int64, *form.User) error
	DeleteUser(int64) error
}
