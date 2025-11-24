package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Bayan2019/ai-hackathon-2025-api/configuration"
	"github.com/Bayan2019/ai-hackathon-2025-api/repositories/database"
	"github.com/Bayan2019/ai-hackathon-2025-api/views"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	mrand "math/rand"

	gomail "gopkg.in/mail.v2"
)

type TokenType string

const (
	// TokenTypeAccess -
	// Set the Issuer to "ozinshe"
	TokenTypeAccess TokenType = "ai-hackathon-access"
)

type AuthHandlers struct {
	DB          *database.Queries
	jwtSecret   string
	email       string
	appPassword string
}

func NewAuthHandlers(config configuration.ApiConfiguration) *AuthHandlers {
	return &AuthHandlers{
		DB:        config.DB,
		jwtSecret: config.JwtSecret,
	}
}

type UserClaims struct {
	// Email  string           `json:"email"`
	// Role database.Roles `json:"role"`
	// 	Last   string           `json:"last"`
	// 	Iat    *jwt.NumericDate `json:"iat"`
	// 	Eat    *jwt.NumericDate `json:"eat"`
	jwt.RegisteredClaims
}

type authedHandler func(http.ResponseWriter, *http.Request, views.User)

func (ah *AuthHandlers) MiddlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := getBearerToken(r.Header)
		if err != nil {
			views.RespondWithError(w, http.StatusUnauthorized, "Couldn't find token", err)
			return
		}

		email, err := validateJWT(jwtToken, ah.jwtSecret)
		if err != nil {
			views.RespondWithError(w, http.StatusInternalServerError, "Couldn't get email from token", err)
			return
		}

		user, err := ah.DB.GetUserByEmail(r.Context(), email)
		if err != nil {
			views.RespondWithError(w, http.StatusNotFound, "Couldn't get user", err)
			return
		}

		_, err = ah.DB.GetRefreshTokenOfUser(r.Context(), user.Email)
		if err != nil {
			views.RespondWithError(w, http.StatusInternalServerError, "Couldn't find refresh_token for user", err)
			return
		}

		handler(w, r, views.User{
			FirstName: user.FirstName,
			Email:     user.Email,
			LastName:  user.LastName,
		})
	}
}

// SignIn godoc
// @Tags Auth
// @Summary      Login to send the code to phone
// @Accept       json
// @Produce      json
// @Param request body views.SignInRequest true "LogIn data"
// @Success      200  {object} views.ResponseMessage "OK"
// @Failure   	 400  {object} views.ErrorResponse "Invalid Data"
// @Failure   	 401  {object} views.ErrorResponse "wrong passwors"
// @Failure   	 404  {object} views.ErrorResponse "Not Present"
// @Failure   	 425  {object} views.ErrorResponse "Too Often"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't query"
// @Router       /v1/auth [post]
func (ah *AuthHandlers) SignIn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var signInReq views.SignInRequest
	err := decoder.Decode(&signInReq)
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Invalid Data", err)
		return
	}

	isPresent, err := ah.DB.IsUserRegistered(r.Context(), signInReq.Email)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't query", err)
		return
	}

	if isPresent == 0 {
		views.RespondWithError(w, http.StatusNotFound, "There is no user with such email", errors.New("not found"))
		return
	}

	user, err := ah.DB.GetUserByEmail(r.Context(), signInReq.Email)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't query", err)
		return
	}

	err = checkPasswordHash(signInReq.Password, user.PasswordHash)
	if err != nil {
		views.RespondWithError(w, http.StatusUnauthorized, "Incorrect password", err)
		return
	}

	codes, err := ah.DB.GetCodesOfUser(r.Context(), signInReq.Email)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't get codes", err)
		return
	}

	if len(codes) >= 1 {
		if codes[0].CreatedAt > time.Now().Add(-70*time.Second).String() {
			views.RespondWithError(w, http.StatusTooEarly, "Wait another 20 seconds", errors.New("too many"))
			return
		}
	}

	code, err := ah.SendCode2Mail(r.Context(), signInReq.Email)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't send code", err)
		return
	}

	views.RespondWithJSON(w, http.StatusOK, views.ResponseMessage{
		Message: fmt.Sprintf("code (%s) is sent", code),
	})
}

// SignInCode godoc
// @Tags Auth
// @Summary      Login to get the code and compare
// @Accept       json
// @Produce      json
// @Param request body views.SignInCodeRequest true "LogIn data"
// @Success      200  {object} views.TokensResponse "OK"
// @Failure   	 404  {object} views.ErrorResponse "Invalid Data"
// @Failure   	 409  {object} views.ErrorResponse "Conflict of Code"
// @Failure   	 500  {object} views.ErrorResponse "Couldn't find the user"
// @Router       /v1/auth [patch]
func (ah *AuthHandlers) SignInCode(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var signInReq views.SignInCodeRequest
	err := decoder.Decode(&signInReq)
	if err != nil {
		views.RespondWithError(w, http.StatusBadRequest, "Invalid Data", err)
		return
	}

	code, err := ah.DB.GetCodeOfUser(r.Context(), signInReq.Email)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't get code", err)
		return
	}

	if code.Code != signInReq.Code {
		views.RespondWithError(w, http.StatusBadRequest, "Wrong code", err)
		return
	}

	err = ah.DB.ConfirmCode(r.Context(), database.ConfirmCodeParams{
		Email:     signInReq.Email,
		Code:      signInReq.Code,
		CreatedAt: code.CreatedAt,
	})
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't save code or refresh token in DataBase", err)
		return
	}

	accessToken, err := makeJWT(
		signInReq.Email,
		ah.jwtSecret,
		time.Hour*24,
	)
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	refreshToken, err := makeRefreshToken()
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	err = ah.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		Email:     signInReq.Code,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60).Format(time.RFC3339),
	})
	if err != nil {
		views.RespondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	views.RespondWithJSON(w, http.StatusOK, views.TokensResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////
//// accommodating functions ////

func (ah *AuthHandlers) SendCode2Mail(ctx context.Context, email string) (string, error) {
	code := generateRandomDigits(6)
	err := ah.DB.CreateCode(ctx, database.CreateCodeParams{
		Email: email,
		Code:  code,
	})
	if err != nil {
		return "", err
	}

	err = ah.sendEmail(email, code)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (ah *AuthHandlers) sendEmail(email string, code string) error {
	// Create a new message
	message := gomail.NewMessage()
	from := ah.email
	// Set email headers
	message.SetHeader("From", from)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Code to confirm")
	// Set email body
	message.SetBody("text/html", fmt.Sprintf("To complete the authentication use the code: %s", code))
	password := ah.appPassword // Consider using application-specific passwords for security

	smtpHost := "smtp.gmail.com" // e.g., smtp.gmail.com
	smtpPort := 587              // or 465 for SMTPS
	// Set up the SMTP dialer
	dialer := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// Send the email
	err := dialer.DialAndSend(message)
	if err != nil {
		return err
	}

	return nil
}

func generateRandomDigits(n int) string {
	// Seed the random number generator using the current time for better randomness.
	// Note: As of Go 1.20, Seed() is deprecated, and a default source is used if not provided.
	// For older Go versions or more explicit control, you might use:
	// rand.NewSource(time.Now().UnixNano())
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))

	var digits = make([]byte, n)
	for i := 0; i < n; i++ {
		// Generate a random digit (0-9) and convert it to its ASCII character representation.
		digits[i] = byte('0' + r.Intn(10))
	}
	return string(digits)
}

func makeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func makeJWT(
	phone string,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	signingKey := []byte(tokenSecret)
	// Use jwt.NewWithClaims to create a new token
	token := jwt.NewWithClaims(
		// Use jwt.SigningMethodHS256 as the signing method.
		jwt.SigningMethodHS256,
		// Use jwt.RegisteredClaims as the claims
		UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer: string(TokenTypeAccess),
				// Set IssuedAt to the current time in UTC
				IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
				// Set ExpiresAt to the current time plus the expiration time (expiresIn)
				ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
				// Set the Subject to a stringified version of the user's email
				Subject: phone,
			},
		})
	// Use token.SignedString to sign the token with the secret key.
	return token.SignedString(signingKey)
}

func validateJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := UserClaims{}
	// Use the jwt.ParseWithClaims function
	// to validate the signature of the JWT
	// and extract the claims into a *jwt.Token struct.
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return "", err
	}

	_, ok := token.Claims.(*UserClaims)
	// _, ok := token.Claims.(*UserClaims)
	if !ok {
		return "", errors.New("couldn't extract claims")
	}
	if !token.Valid {
		return "", errors.New("token not valid")
	}

	email, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}
	if issuer != string(TokenTypeAccess) {
		return "", errors.New("invalid issuer")
	}

	return email, nil
}

func getBearerToken(headers http.Header) (string, error) {
	// Auth information will come into our server
	// in the Authorization header.
	authHeader := headers.Get("Authorization")
	// fmt.Println(authHeader)
	if authHeader == "" {
		// If the header doesn't exist, return an error.
		return "", errors.New("no auth header included in request")
	}
	// stripping off the Bearer prefix and whitespace
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		// If the header doesn't exist, return an error.
		return "", errors.New("malformed authorization header")
	}
	// return the TOKEN_STRING if it exists
	return splitAuth[1], nil
	// return authHeader, nil
}
