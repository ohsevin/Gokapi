package webserver

import (
	"Gokapi/internal/configuration"
	testconfiguration "Gokapi/internal/test"
	testconfiguration2 "Gokapi/internal/test/testconfiguration"
	"html/template"
	"io/fs"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	testconfiguration2.Create(true)
	configuration.Load()
	go Start()
	time.Sleep(1 * time.Second)
	exitVal := m.Run()
	testconfiguration2.Delete()
	os.Exit(exitVal)
}

func TestEmbedFs(t *testing.T) {
	templates, err := template.ParseFS(templateFolderEmbedded, "web/templates/*.tmpl")
	if err != nil {
		t.Error("Unable to read templates")
	}
	if !strings.Contains(templates.DefinedTemplates(), "app_name") {
		t.Error("Unable to parse templates")
	}
	_, err = fs.Stat(staticFolderEmbedded, "web/static/expired.png")
	if err != nil {
		t.Error("Static webdir incomplete")
	}
}

func TestIndexRedirect(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/",
		RequiredContent: []string{"<html><head><meta http-equiv=\"Refresh\" content=\"0; URL=./index\"></head></html>"},
		IsHtml:          true,
	})
}
func TestIndexFile(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/index",
		RequiredContent: []string{configuration.ServerSettings.RedirectUrl},
		IsHtml:          true,
	})
}
func TestStaticDirs(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/css/cover.css",
		RequiredContent: []string{".btn-secondary:hover"},
	})
}
func TestLogin(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/login",
		RequiredContent: []string{"id=\"uname_hidden\""},
		IsHtml:          true,
	})
}
func TestAdminNoAuth(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/admin",
		RequiredContent: []string{"URL=./login\""},
		IsHtml:          true,
	})
}
func TestAdminAuth(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/admin",
		RequiredContent: []string{"Downloads remaining"},
		IsHtml:          true,
		Cookies: []testconfiguration.Cookie{{
			Name:  "session_token",
			Value: "validsession",
		}},
	})
}
func TestAdminExpiredAuth(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/admin",
		RequiredContent: []string{"URL=./login\""},
		IsHtml:          true,
		Cookies: []testconfiguration.Cookie{{
			Name:  "session_token",
			Value: "expiredsession",
		}},
	})
}

func TestAdminRenewalAuth(t *testing.T) {
	t.Parallel()
	cookies := testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/admin",
		RequiredContent: []string{"Downloads remaining"},
		IsHtml:          true,
		Cookies: []testconfiguration.Cookie{{
			Name:  "session_token",
			Value: "needsRenewal",
		}},
	})
	sessionCookie := "needsRenewal"
	for _, cookie := range cookies {
		if (*cookie).Name == "session_token" {
			sessionCookie = (*cookie).Value
			break
		}
	}
	if sessionCookie == "needsRenewal" {
		t.Error("Session not renewed")
	}
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/admin",
		RequiredContent: []string{"Downloads remaining"},
		IsHtml:          true,
		Cookies: []testconfiguration.Cookie{{
			Name:  "session_token",
			Value: sessionCookie,
		}},
	})
}

func TestAdminInvalidAuth(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/admin",
		RequiredContent: []string{"URL=./login\""},
		IsHtml:          true,
		Cookies: []testconfiguration.Cookie{{
			Name:  "session_token",
			Value: "invalid",
		}},
	})
}

func TestInvalidLink(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/d?id=123",
		RequiredContent: []string{"URL=./error\""},
		IsHtml:          true,
	})
}
func TestError(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/error",
		RequiredContent: []string{"this file cannot be found"},
		IsHtml:          true,
	})
}
func TestForgotPw(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/forgotpw",
		RequiredContent: []string{"--reset-pw"},
		IsHtml:          true,
	})
}
func TestLoginCorrect(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/login",
		RequiredContent: []string{"URL=./admin\""},
		IsHtml:          true,
		Method:          "POST",
		PostValues:      []testconfiguration.PostBody{{"username", "test"}, {"password", "testtest"}},
	})
}

func TestLoginIncorrectPassword(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/login",
		RequiredContent: []string{"Incorrect username or password"},
		IsHtml:          true,
		Method:          "POST",
		PostValues:      []testconfiguration.PostBody{{"username", "test"}, {"password", "incorrect"}},
	})
}
func TestLoginIncorrectUsername(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/login",
		RequiredContent: []string{"Incorrect username or password"},
		IsHtml:          true,
		Method:          "POST",
		PostValues:      []testconfiguration.PostBody{{"username", "incorrect"}, {"password", "incorrect"}},
	})
}

func TestLogout(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/admin",
		RequiredContent: []string{"Downloads remaining"},
		IsHtml:          true,
		Cookies: []testconfiguration.Cookie{{
			Name:  "session_token",
			Value: "logoutsession",
		}},
	})
	// Logout
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/logout",
		RequiredContent: []string{"URL=./login\""},
		IsHtml:          true,
		Cookies: []testconfiguration.Cookie{{
			Name:  "session_token",
			Value: "logoutsession",
		}},
	})
	// Admin after logout
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://localhost:53843/admin",
		RequiredContent: []string{"URL=./login\""},
		IsHtml:          true,
		Cookies: []testconfiguration.Cookie{{
			Name:  "session_token",
			Value: "logoutsession",
		}},
	})
}

func TestDownloadHotlink(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/hotlink/PhSs6mFtf8O5YGlLMfNw9rYXx9XRNkzCnJZpQBi7inunv3Z4A.jpg",
		RequiredContent: []string{"123"},
	})
	// Download expired hotlink
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/hotlink/PhSs6mFtf8O5YGlLMfNw9rYXx9XRNkzCnJZpQBi7inunv3Z4A.jpg",
		RequiredContent: []string{"Created with GIMP"},
	})
}

func TestDownloadNoPassword(t *testing.T) {
	t.Parallel()
	// Show download page
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/d?id=Wzol7LyY2QVczXynJtVo",
		IsHtml:          true,
		RequiredContent: []string{"smallfile2"},
	})
	// Download
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/downloadFile?id=Wzol7LyY2QVczXynJtVo",
		RequiredContent: []string{"789"},
	})
	// Show download page expired file
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/d?id=Wzol7LyY2QVczXynJtVo",
		IsHtml:          true,
		RequiredContent: []string{"URL=./error\""},
	})
	// Download expired file
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/downloadFile?id=Wzol7LyY2QVczXynJtVo",
		IsHtml:          true,
		RequiredContent: []string{"URL=./error\""},
	})
}

func TestDownloadPagePassword(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/d?id=jpLXGJKigM4hjtA6T6sN",
		IsHtml:          true,
		RequiredContent: []string{"Password required"},
	})
}
func TestDownloadPageIncorrectPassword(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/d?id=jpLXGJKigM4hjtA6T6sN",
		IsHtml:          true,
		RequiredContent: []string{"Incorrect password!"},
		Method:          "POST",
		PostValues:      []testconfiguration.PostBody{{"password", "incorrect"}},
	})
}

func TestDownloadIncorrectPasswordCookie(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/d?id=jpLXGJKigM4hjtA6T6sN",
		IsHtml:          true,
		RequiredContent: []string{"Password required"},
		Cookies:         []testconfiguration.Cookie{{"pjpLXGJKigM4hjtA6T6sN", "invalid"}},
	})
}

func TestDownloadIncorrectPassword(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/downloadFile?id=jpLXGJKigM4hjtA6T6sN",
		IsHtml:          true,
		RequiredContent: []string{"URL=./d?id=jpLXGJKigM4hjtA6T6sN"},
		Cookies:         []testconfiguration.Cookie{{"pjpLXGJKigM4hjtA6T6sN", "invalid"}},
	})
}

func TestDownloadCorrectPassword(t *testing.T) {
	t.Parallel()
	// Submit download page correct password
	cookies := testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/d?id=jpLXGJKigM4hjtA6T6sN2",
		IsHtml:          true,
		RequiredContent: []string{"URL=./d?id=jpLXGJKigM4hjtA6T6sN2"},
		Method:          "POST",
		PostValues:      []testconfiguration.PostBody{{"password", "123"}},
	})
	pwCookie := ""
	for _, cookie := range cookies {
		if (*cookie).Name == "pjpLXGJKigM4hjtA6T6sN2" {
			pwCookie = (*cookie).Value
			break
		}
	}
	if pwCookie == "" {
		t.Error("Cookie not set")
	}
	// Show download page correct password
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/d?id=jpLXGJKigM4hjtA6T6sN2",
		IsHtml:          true,
		RequiredContent: []string{"smallfile"},
		Cookies:         []testconfiguration.Cookie{{"pjpLXGJKigM4hjtA6T6sN2", pwCookie}},
	})
	// Download correct password
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/downloadFile?id=jpLXGJKigM4hjtA6T6sN2",
		RequiredContent: []string{"456"},
		Cookies:         []testconfiguration.Cookie{{"pjpLXGJKigM4hjtA6T6sN2", pwCookie}},
	})
}

func TestDeleteFileNonAuth(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/delete?id=e4TjE7CokWK0giiLNxDL",
		IsHtml:          true,
		RequiredContent: []string{"URL=./login"},
	})
}

func TestDeleteFileInvalidKey(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/delete",
		IsHtml:          true,
		RequiredContent: []string{"URL=./admin"},
		Cookies: []testconfiguration.Cookie{{
			Name:  "session_token",
			Value: "validsession",
		}},
	})
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/delete?id=",
		IsHtml:          true,
		RequiredContent: []string{"URL=./admin"},
		Cookies: []testconfiguration.Cookie{{
			Name:  "session_token",
			Value: "validsession",
		}},
	})
}

func TestDeleteFile(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPageResult(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/delete?id=e4TjE7CokWK0giiLNxDL",
		IsHtml:          true,
		RequiredContent: []string{"URL=./admin"},
		Cookies: []testconfiguration.Cookie{{
			Name:  "session_token",
			Value: "validsession",
		}},
	})
}

func TestPostUploadNoAuth(t *testing.T) {
	t.Parallel()
	testconfiguration.HttpPostRequest(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/upload",
		UploadFileName:  "test/fileupload.jpg",
		UploadFieldName: "file",
		RequiredContent: []string{"{\"Result\":\"error\",\"ErrorMessage\":\"Not authenticated\"}"},
	})
}
func TestPostUpload(t *testing.T) {
	testconfiguration.HttpPostRequest(t, testconfiguration.HttpTestConfig{
		Url:             "http://127.0.0.1:53843/upload",
		UploadFileName:  "test/fileupload.jpg",
		UploadFieldName: "file",
		RequiredContent: []string{"{\"Result\":\"OK\"", "\"Name\":\"fileupload.jpg\"", "\"SHA256\":\"a9993e364706816aba3e25717850c26c9cd0d89d\"", "DownloadsRemaining\":3"},
		ExcludedContent: []string{"\"Id\":\"\"", "HotlinkId\":\"\""},
		Cookies: []testconfiguration.Cookie{{
			Name:  "session_token",
			Value: "validsession",
		}},
	})
}