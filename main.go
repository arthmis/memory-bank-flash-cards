package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"memorybank/database"
	"memorybank/queries"
	"memorybank/views"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func (env *Env) Login(c echo.Context) error {
	component := views.Login()
	return html(c, http.StatusOK, component)
}

func html(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

type Env struct {
	decks database.DeckModel
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clerkKey := os.Getenv("clerk_secret_key")

	// Set the API key with your Clerk Secret Key
	clerk.SetKey(clerkKey)

	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, "user=postgres dbname=postgres password=postgres port=7777")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer dbpool.Close()

	deckQueries := queries.New(dbpool)

	env := &Env{
		decks: database.DeckModel{DB: dbpool, Queries: *deckQueries},
	}

	e := echo.New()
	e.Static("/", "static")
	// e.GET("/dashboard", env.dashboard, cookiesToAuth, handleAuth)
	// protectedHandler := http.HandlerFunc(clerkAuth)
	// headerAuthorization := clerkhttp.WithHeaderAuthorization()(protectedHandler)
	headerAuthorization := clerkhttp.WithHeaderAuthorization()
	authorization := echo.WrapMiddleware(headerAuthorization)
	e.GET("/dashboard", env.dashboard, cookiesToAuth, authorization, clerkAuth)
	e.GET("/login", env.Login)
	// e.POST("/api/cards", env.Cards, cookiesToAuth, authorization, clerkAuth)
	e.POST("/api/cards", env.Cards)
	// e.GET("api/login", env.login)
	e.Logger.Fatal(e.Start(":8000"))
}

func (env *Env) dashboard(c echo.Context) error {
	component := views.Dashboard()
	return html(c, http.StatusOK, component)
}

type CardInput struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func (env *Env) Cards(c echo.Context) error {
	cardInput := CardInput{}
	err := c.Bind(&cardInput)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	cardParams := queries.CreateCardParams{
		Question: cardInput.Question,
		Answer:   cardInput.Answer,
	}
	card, err := env.decks.Queries.CreateCard(context.Background(), cardParams)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, card)
}

func clerkAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, ok := clerk.SessionClaimsFromContext(c.Request().Context())
		if !ok {
			c.Response().WriteHeader(http.StatusUnauthorized)
			c.Response().Write([]byte(`{"access": "unauthorized"}`))
			return errors.New("unauthorized")
		}

		// usr, err := user.Get(c.Request().Context(), claims.Subject)
		// if err != nil {
		// 	// handle the error
		// }
		// fmt.Fprintf(c.Response().Writer, `{"user_id": "%s", "user_banned": "%t"}`, usr.ID, usr.Banned)

		next(c)
		return nil
	}
}

func cookiesToAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		cookieHeader := c.Request().Header.Get("Cookie")
		requestCookie, cookieErr := c.Cookie("__session")
		if authHeader != "" {
			fmt.Printf("Auth header ----- %s", authHeader)
			next(c)
			return nil
		} else {
			if cookieErr == nil && requestCookie.Value != "" {
				// fmt.Printf("Request Cookie ----- %s", requestCookie.Value)
				setAuthHeader(c.Request(), requestCookie.Value)
				next(c)
				return nil
			}
			if cookieHeader != "" {
				// fmt.Printf("Cookie header ----- %s", cookieHeader)
				session := getSessionFromCookieHeader(cookieHeader)
				setAuthHeader(c.Request(), session)
				next(c)
				return nil
			}
			next(c)
			return errors.New("couldn't find cookie or auth header for authentication")
		}
	}
}

func getSessionFromCookieHeader(cookie string) string {
	value := strings.Split(cookie, " ")
	var session string
	for _, v := range value {
		if strings.Contains(v, "__session_") {
			val := strings.Split(v, "=")
			if len(val) >= 1 {
				session = val[1]
				session = strings.Replace(session, ";", "", -1)
				fmt.Printf("Session ----- %s", session)
				return session
			}
		}
	}
	return session
}

func setAuthHeader(r *http.Request, value string) {
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", value))
}

func (env *Env) createDeck(c echo.Context) error {
	name := c.FormValue("name")
	env.decks.CreateDeck(database.NewDeck{Name: name})
	return c.NoContent(http.StatusOK)
}
