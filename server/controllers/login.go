package controllers

import (
	"database/sql"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"pmail/db"
	"pmail/dto"
	"pmail/dto/response"
	"pmail/i18n"
	"pmail/models"
	"pmail/session"
	"pmail/utils/password"
)

type loginRequest struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func Login(ctx *dto.Context, w http.ResponseWriter, req *http.Request) {

	reqBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Errorf("%+v", err)
	}
	var reqData loginRequest
	err = json.Unmarshal(reqBytes, &reqData)
	if err != nil {
		log.Errorf("%+v", err)
	}

	var user models.User

	encodePwd := password.Encode(reqData.Password)

	err = db.Instance.Get(&user, db.WithContext(ctx, "select * from user where account =? and password =?"),
		reqData.Account, encodePwd)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("%+v", err)
	}

	if user.ID != 0 {
		userStr, _ := json.Marshal(user)
		session.Instance.Put(req.Context(), "user", string(userStr))
		response.NewSuccessResponse("").FPrint(w)
	} else {
		response.NewErrorResponse(response.ParamsError, i18n.GetText(ctx.Lang, "aperror"), "").FPrint(w)
	}
}
