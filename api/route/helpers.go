package route

import (
	"dazzlefm/util"

	"github.com/labstack/echo"
)

// ModlueEnum the enum type for Module IDs
type ModlueEnum int

const (
	// AllModule ...
	AllModule ModlueEnum = iota
	// PageModule ...
	PageModule
	// PresenterModule ...
	PresenterModule
	// ScheduleModule ...
	ScheduleModule
	// GalleryModule ...
	GalleryModule
	// ContactsModule ...
	ContactsModule
	// MusicModule ...
	MusicModule
	// EventModule ...
	EventModule
	// UserModule ...
	UserModule
	// NewsModule ...
	NewsModule
	// ScriptModule ...
	ScriptModule
	// SettingsModule ...
	SettingsModule
)

// RoleDefinition ...
type RoleDefinition struct {
	ID     int
	Edit   []ModlueEnum
	View   []ModlueEnum
	Delete []ModlueEnum
	Create []ModlueEnum
}

// AccessControl ...
var AccessControl map[int]RoleDefinition

func init() {
	AccessControl = map[int]RoleDefinition{
		1: {1,
			[]ModlueEnum{AllModule},
			[]ModlueEnum{AllModule},
			[]ModlueEnum{AllModule},
			[]ModlueEnum{AllModule},
		},

		2: {2,
			[]ModlueEnum{PageModule, PresenterModule, GalleryModule, ContactsModule, MusicModule, EventModule, NewsModule, ScheduleModule},
			[]ModlueEnum{PageModule, PresenterModule, ScheduleModule, GalleryModule, ContactsModule, MusicModule, EventModule, NewsModule, SettingsModule},
			[]ModlueEnum{PageModule, PresenterModule, GalleryModule, ContactsModule, MusicModule, EventModule, NewsModule, ScheduleModule},
			[]ModlueEnum{PageModule, PresenterModule, GalleryModule, ContactsModule, MusicModule, EventModule, NewsModule, ScheduleModule},
		},

		3: {3,
			[]ModlueEnum{PageModule, PresenterModule, GalleryModule, MusicModule, EventModule, NewsModule},
			[]ModlueEnum{PageModule, PresenterModule, ScheduleModule, GalleryModule, MusicModule, EventModule, NewsModule, SettingsModule},
			[]ModlueEnum{PageModule, PresenterModule, GalleryModule, MusicModule, EventModule, NewsModule},
			[]ModlueEnum{PageModule, PresenterModule, GalleryModule, MusicModule, EventModule, NewsModule},
		},
	}
}

// IsAuthenticated checks if the current user is logged in
func IsAuthenticated(c echo.Context) (ok bool, err error) {
	sMgr, err := util.NewSessionMgr(c, "session")
	if err != nil {
		return
	}

	// sMgr.Set("loggedIn", true)
	// sMgr.Set("userName", record.UserName)
	// sMgr.Set("userID", record.ID)
	// sMgr.Set("userRole", record.Role)

	ok = sMgr.Bool("loggedIn")

	return
}

// IsAuthorized checks if the current user is logged in and has a suitable role definition
// verb: 1: Get / (view), 2: Post / (create), 3 Post /id (edit), 4 Delete / (delete)
func IsAuthorized(c echo.Context, verb int, module ModlueEnum) (ok bool, err error) {
	sMgr, err := util.NewSessionMgr(c, "session")
	if err != nil {
		return
	}

	if !sMgr.Bool("loggedIn") {
		ok = false
		return
	}

	userRole := sMgr.Int("userRole")
	if userRole == 0 {
		ok = false
		return
	}
	roleDef := AccessControl[userRole]
	switch verb {
	case 1:
		ok = findModule(module, roleDef.View)
		return
	case 2:
		ok = findModule(module, roleDef.Create)
		return
	case 3:
		ok = findModule(module, roleDef.Edit)
		return
	case 4:
		ok = findModule(module, roleDef.Delete)
		return
	}

	return
}

func findModule(mdl ModlueEnum, allowedModules []ModlueEnum) bool {
	for _, i := range allowedModules {
		if i == mdl || i == AllModule {
			return true
		}
	}

	return false
}
