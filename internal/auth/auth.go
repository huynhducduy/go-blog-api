package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	firebase "firebase.google.com/go"
	"github.com/go-sql-driver/mysql"
	"go-blog/internal/db"
	"go-blog/internal/user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go-blog/internal/config"
	"go-blog/pkg/utils"
	//"firebase.google.com/go/auth"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Claims struct {
	Id        int64 `json:"id"`
	ExpiresAt int64 `json:"expires_at"`
	jwt.StandardClaims
}

type Token struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
	Username  string `json:"username"`
}

func getToken(r *http.Request) (string, error) {
	if r.Header["Authorization"] != nil && len(strings.Split(r.Header["Authorization"][0], " ")) == 2 {
		return strings.Split(r.Header["Authorization"][0], " ")[1], nil
	} else {
		return "", errors.New("No bearer token.")
	}
}

func generateToken(id int64, username string) Token {
	expAt := time.Now().Unix() + 604800 // 1 week

	payload := Claims{Id: id, ExpiresAt: expAt}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), payload)
	tokenString, _ := token.SignedString([]byte(config.GetConfig().SECRET))

	return Token{
		Token:     tokenString,
		ExpiresAt: expAt,
		Username: username,
	}
}

func hashAndSalt(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashed string, plain string) bool {
	byteHashed := []byte(hashed)
	bytePlain := []byte(plain)
	err := bcrypt.CompareHashAndPassword(byteHashed, bytePlain)
	if err != nil {
		return false
	}
	return true
}

var app *firebase.App

func InitFirebase() {

	var err error

	opt := option.WithCredentialsFile(config.GetConfig().FIREBASE_PRIVATEKEY)
	config := &firebase.Config{ProjectID: "swd391"}
	app, err = firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}
}

func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		empToken, err := getToken(r)
		if err != nil {
			utils.ResponseMessage(w, http.StatusUnauthorized, err.Error())
			return
		}

		token, err := jwt.ParseWithClaims(empToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetConfig().SECRET), nil
		})

		if err == nil && token.Valid {
			var emp *user.User
			emp, err = user.Read(token.Claims.(*Claims).Id)
			if err != nil {
				utils.ResponseInternalError(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), "user", emp)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		utils.ResponseMessage(w, http.StatusUnauthorized, "Invalid token!")
	})
}

func Login(w http.ResponseWriter, r *http.Request) {

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	var credential Credential
	json.Unmarshal(reqBody, &credential)

	if credential.Token != "" {

		ctx := context.Background()

		client, err := app.Auth(ctx)
		if err != nil {
			utils.ResponseInternalError(w, err)
			return
		}

		token, err := client.VerifyIDToken(ctx, credential.Token)
		if err != nil {
			utils.ResponseInternalError(w, err)
			return
		}

		log.Printf("Verified ID token: %v\n", token)

		email := token.Firebase.Identities["email"].([]interface{})[0].(string)
		var id int64
		var username string

		db := db.GetConnection()

		results := db.QueryRow("SELECT `id`, `username` FROM `users` where `email` = ?", email)
		err = results.Scan(&id, &username)
		if err == sql.ErrNoRows {
			utils.ResponseMessage(w, http.StatusNotFound, "Email is not registered!")
			return
		} else if err != nil {
			utils.ResponseInternalError(w, err)
			return
		}

		userToken := generateToken(id, username)

		utils.Response(w, http.StatusOK, userToken)

	} else {
		if credential.Username == "" || credential.Password == "" {
			utils.ResponseMessage(w, http.StatusBadRequest, "Username and password must not be empty!")
			return
		}

		var id int64
		var pass string
		db := db.GetConnection()

		results := db.QueryRow("SELECT `id`, `password` FROM `users` where `username` = ?", credential.Username)
		err = results.Scan(&id, &pass)
		if err == sql.ErrNoRows {
			utils.ResponseMessage(w, http.StatusNotFound, "Username and password is incorrect!")
			return
		} else if err != nil {
			utils.ResponseInternalError(w, err)
			return
		}

		if !comparePasswords(pass, credential.Password) {
			utils.ResponseMessage(w, http.StatusNotFound, "Username and password is incorrect!")
			return
		}

		token := generateToken(id, credential.Username)

		utils.Response(w, http.StatusOK, token)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	var newUser user.User
	json.Unmarshal(reqBody, &newUser)

	if newUser.Username == nil || *newUser.Username == "" {
		utils.ResponseMessage(w, http.StatusBadRequest, "Username cannot be empty!")
		return
	}

	if newUser.Name == nil || *newUser.Name == "" {
		utils.ResponseMessage(w, http.StatusBadRequest, "Name cannot be empty!")
		return
	}

	if newUser.Email == nil || *newUser.Email == "" {
		utils.ResponseMessage(w, http.StatusBadRequest, "Email cannot be empty!")
		return
	}

	if newUser.Password == nil || *newUser.Password == "" {
		utils.ResponseMessage(w, http.StatusBadRequest, "Password cannot be empty!")
		return
	}



	newUser.Role = nil
	newPassword := hashAndSalt(*newUser.Password)
	newUser.Password = &newPassword

	id, err := user.Create(newUser)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			utils.ResponseMessage(w, http.StatusBadRequest, "Email or username is already used!")
			return
		}
		utils.ResponseInternalError(w, err)
		return
	}

	utils.ResponseCreated(w, id)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(*user.User)
	utils.Response(w, 200, usr)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(*user.User)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	json.Unmarshal(reqBody, &usr)

	err = user.Update(*usr)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			utils.ResponseMessage(w, http.StatusBadRequest, "Email or username is already used!")
			return
		}
		utils.ResponseInternalError(w, err)
		return
	}

	utils.ResponseMessage(w, http.StatusOK, "Update profile success!")
}


func ChangeProfilePassword(w http.ResponseWriter, r *http.Request) {
	usr := r.Context().Value("user").(*user.User)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	json.Unmarshal(reqBody, &usr)
	password := hashAndSalt(*usr.Password)
	usr.Password = &password

	oldPasswordFromDB, err := user.ReadPassword(*usr.Id)
	if err != nil {
		utils.ResponseInternalError(w, err)
		return
	}

	if !comparePasswords(oldPasswordFromDB, *usr.OldPassword) {
		utils.ResponseMessage(w, http.StatusNotFound, "Old password is incorrect")
		return
	}

	err = user.UpdatePassword(*usr)

	utils.ResponseMessage(w, http.StatusOK, "Update password success!")
}

//func GetPwd(w http.ResponseWriter, r *http.Request) {
//	utils.ResponseMessage(w, 200, hashAndSalt("password123"))
//}