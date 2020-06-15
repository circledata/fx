package fx

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/rivo/sessions"
	"github.com/gorilla/mux"
)

type WebUser struct {
	sessions.User
}

type WebModule struct {
	Mux    *mux.Router
	sessionCookieName string
	sessionExpiry time.Duration
}

func (m *WebModule) InitializeUserSession(userProviderFunc func(id interface{}) (sessions.User, error)) error {

	sessionCookieName := "fx"

	if m.sessionCookieName != "" {
		sessionCookieName = m.sessionCookieName
	}

	//sessions
	sessions.SessionCookie = sessionCookieName
	sessions.SessionExpiry = 60 * time.Minute
	//sessions.SessionIDExpiry = 15 * time.Minute
	//sessions.SessionIDGracePeriod = 15 * * time.Minute

	sessions.NewSessionCookie = func() *http.Cookie {
		return &http.Cookie{
			Expires:  time.Now().Add(24 * time.Hour),
			MaxAge:   24 * 60 * 60,
			HttpOnly: true,
			//Domain:   "www.example.com",
			Path: "/",
			//Secure: true,
		}
	}

	persistence, persistenceOK := sessions.Persistence.(sessions.ExtendablePersistenceLayer)

	if !persistenceOK {
		return errors.New("error initializing session persistence: persistence is not of type ExtendablePersistenceLayer")
	}

	persistence.LoadUserFunc = userProviderFunc

	return nil
}

func (m *WebModule) StartUserSession(w http.ResponseWriter, r *http.Request, user *WebUser) error {
	session, sessionError := sessions.Start(w, r, true)

	if sessionError != nil {
		return sessionError
	}

	loginErr := session.LogIn(user, false, w)

	if loginErr != nil {
		return loginErr
	}

	return nil
}

func (m *WebModule) EndUserSession(w http.ResponseWriter, r *http.Request) error {

	session, sessionError := sessions.Start(w, r, false)

	if sessionError != nil {
		return sessionError
	}

	if session == nil || session.User() == nil {
		return nil
	}

	if logoutError := session.LogOut(); logoutError != nil {
		return logoutError
	}

	if sessionRegenerateError := session.RegenerateID(w); sessionRegenerateError != nil {
		return sessionRegenerateError
	}

	if destroyErr := session.Destroy(w, r); destroyErr != nil {
		return destroyErr
	}

	return nil
}

func (m *WebModule) GetLoggedInUser(w http.ResponseWriter, r *http.Request) (*WebUser, error) {
	session, sessionError := sessions.Start(w, r, false)

	if sessionError != nil {
		return nil, sessionError
	}

	if session == nil {
		return nil, nil
	}

	sessionUser := session.User().(*WebUser)

	if session == nil || sessionUser == nil {
		return nil, nil
	}

	return sessionUser, nil
}

func (m *WebModule) GetViewsDirectory() string {
	return ""
}

func (m *WebModule) GetContextPath() string {
	return ""
}

/// Views
type ViewData struct {
	ContextPath string
	LoggedInUser       *WebUser
	SystemDate         time.Time
	PageInfoMessage    string
	PageSuccessMessage string
	PageWarningMessage string
	PageErrorMessage   string
	Data               interface{}
}

type LayoutView struct {
	Template *template.Template
	ServerModule       ServerModule
}

func (v *LayoutView) Render(w http.ResponseWriter, r *http.Request, data interface{}) error {
	session, sessionErr := sessions.Start(w, r, false)

	if sessionErr != nil {
		return sessionErr
	}

	pageInfoMessage := session.GetAndDelete("pageInfoMessage", nil)

	if pageInfoMessage == nil {
		pageInfoMessage = ""
	}

	pageSuccessMessage := session.GetAndDelete("pageSuccessMessage", nil)

	if pageSuccessMessage == nil {
		pageSuccessMessage = ""
	}

	pageWarningMessage := session.GetAndDelete("pageWarningMessage", nil)

	if pageWarningMessage == nil {
		pageWarningMessage = ""
	}

	pageErrorMessage := session.GetAndDelete("pageErrorMessage", nil)

	if pageErrorMessage == nil {
		pageErrorMessage = ""
	}

	loggedInUser := session.User().(*WebUser)

	return v.Template.ExecuteTemplate(w, "layout", ViewData{
		ContextPath:       v.ServerModule.GetContextPath(),
		SystemDate:       time.Now(),
		LoggedInUser:       loggedInUser,
		PageInfoMessage:    pageInfoMessage.(string),
		PageSuccessMessage: pageSuccessMessage.(string),
		PageWarningMessage: pageWarningMessage.(string),
		PageErrorMessage:   pageErrorMessage.(string),
		Data:               data,
	})
}

type SimpleView struct {
	Template *template.Template
}

func (v *SimpleView) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, "page", ViewData{
		SystemDate:       time.Now(),
		Data: data,
	})
}

