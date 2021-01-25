package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hrfee/jfa-go/common"
	"github.com/hrfee/jfa-go/mediabrowser"
)

func (app *appContext) ServeSetup(gc *gin.Context) {
	lang := gc.Query("lang")
	if lang == "" {
		lang = "en-us"
	} else if _, ok := app.storage.lang.Admin[lang]; !ok {
		lang = "en-us"
	}
	emailLang := lang
	if _, ok := app.storage.lang.Email[lang]; !ok {
		emailLang = "en-us"
	}
	gc.HTML(200, "setup2.html", gin.H{
		"lang":      app.storage.lang.Setup[lang],
		"emailLang": app.storage.lang.Email[emailLang],
	})
}

type testReq struct {
	Host     string `json:"jfHost"`
	Username string `json:"jfUser"`
	Password string `json:"jfPassword"`
}

func (app *appContext) TestJF(gc *gin.Context) {
	var req testReq
	gc.BindJSON(&req)
	tempjf, _ := mediabrowser.NewServer(mediabrowser.JellyfinServer, req.Host, "jfa-go-setup", app.version, "auth", "auth", common.NewTimeoutHandler("authJF", req.Host, true), 30)
	_, status, err := tempjf.Authenticate(req.Username, req.Password)
	if !(status == 200 || status == 204) || err != nil {
		app.info.Printf("Auth failed with code %d (%s)", status, err)
		gc.JSON(401, map[string]bool{"success": false})
		return
	}
	gc.JSON(200, map[string]bool{"success": true})
}

func (st *Storage) loadLangSetup() error {
	st.lang.Setup = map[string]setupLang{}
	var english setupLang
	load := func(fname string) error {
		index := strings.TrimSuffix(fname, filepath.Ext(fname))
		lang := setupLang{}
		f, err := ioutil.ReadFile(filepath.Join(st.lang.SetupPath, fname))
		if err != nil {
			return err
		}
		err = json.Unmarshal(f, &lang)
		if err != nil {
			return err
		}
		st.lang.Common.patchCommon(index, &lang.Strings)
		if fname != "en-us.json" {
			patchLang(&english.Strings, &lang.Strings)
			patchLang(&english.StartPage, &lang.StartPage)
			patchLang(&english.EndPage, &lang.EndPage)
			patchLang(&english.Language, &lang.Language)
			patchLang(&english.Login, &lang.Login)
			patchLang(&english.JellyfinEmby, &lang.JellyfinEmby)
			patchLang(&english.Email, &lang.Email)
			patchLang(&english.Notifications, &lang.Notifications)
			patchLang(&english.PasswordResets, &lang.PasswordResets)
			patchLang(&english.InviteEmails, &lang.InviteEmails)
			patchLang(&english.PasswordValidation, &lang.PasswordValidation)
			patchLang(&english.HelpMessages, &lang.HelpMessages)
		}
		st.lang.Setup[index] = lang
		return nil
	}
	err := load("en-us.json")
	if err != nil {
		return err
	}
	english = st.lang.Setup["en-us"]
	files, err := ioutil.ReadDir(st.lang.SetupPath)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.Name() != "en-us.json" {
			err = load(f.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}
