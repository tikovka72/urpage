package handlers

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4"
	"net/http"
	"urpage/constants"
	"urpage/jwt_api"
	"urpage/session"
	"urpage/storage"
)

func CreateDoEditPassword(conn *pgx.Conn, rdb *redis.Client) {
	doEditPassword := func(writer http.ResponseWriter, request *http.Request) {
		var (
			userId                              int
			oldPassword, newPassword, CSRFToken string
			jsonAnswer                          []byte
			user                                storage.User
			err                                 error
		)

		if request.Method != "POST" {
			return
		}

		defer func() { SendJson(writer, jsonAnswer) }()

		{ // CSRF check
			_, CSRFToken, err = session.CheckSessionId(writer, request, rdb)

			if err != nil {
				jsonAnswer, _ = json.Marshal(Answer{Err: "no-csrf"})
				return
			}
		}

		{ // check user authed
			userId, err = jwt_api.CheckIfUserAuth(writer, request, rdb)

			if err != nil {
				http.Error(writer, "no jwt", http.StatusForbidden)
				return
			}
		}

		{ // work with form
			CSRFTokenForm := request.FormValue("csrf")

			if CSRFToken != CSRFTokenForm {
				jsonAnswer, _ = json.Marshal(Answer{Err: "no-csrf"})
				return
			}

			oldPassword = request.FormValue("old")
			newPassword = request.FormValue("new")
		}

		{ // get user
			user, err = storage.GetUserViaId(conn, userId)

			if err != nil {
				http.Error(writer, "error getting user", http.StatusInternalServerError)
				return
			}
		}

		{ // check old password
			correct, err := storage.CheckPassword(oldPassword, user.Password)
			if err != nil || !correct {
				jsonAnswer, _ = json.Marshal(Answer{Err: "wrong-password"})
				return
			}
		}

		{ // set new data
			user.ImagePath = user.ImagePath[len(constants.UserImages):]

			user.Password, err = storage.HashPassword(newPassword)

			if err != nil {
				http.Error(writer, "error hashing password", http.StatusInternalServerError)
				return
			}

			err = storage.UpdateUserMainInfo(conn, user)
			if err != nil {
				http.Error(writer, "error updating user", http.StatusInternalServerError)
				return
			}
			jsonAnswer, _ = json.Marshal(Answer{Err: ""})
		}
	}

	http.HandleFunc("/do/edit_password", doEditPassword)
}
