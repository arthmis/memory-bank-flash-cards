package main

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"log"
	"memorybank/api"
	"memorybank/database"
	"memorybank/queries"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// func (env *Env) Login(c echo.Context) error {
// 	component := views.Login()
// 	return html(c, http.StatusOK, component)
// }

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

	handlers := api.NewStrictHandler(env, nil)
	e := echo.New()
	api.RegisterHandlers(e, handlers)
	// e.Static("/dashboard", "app/dist")
	// // e.GET("/dashboard", env.dashboard, cookiesToAuth, handleAuth)
	// // protectedHandler := http.HandlerFunc(clerkAuth)
	// // headerAuthorization := clerkhttp.WithHeaderAuthorization()(protectedHandler)
	// // headerAuthorization := clerkhttp.WithHeaderAuthorization()
	// // authorization := echo.WrapMiddleware(headerAuthorization)
	// // e.GET("/dashboard", env.dashboard, cookiesToAuth, authorization, clerkAuth)
	e.Logger.Fatal(e.Start(":8000"))
}

// func (env *Env) dashboard(c echo.Context) error {
// 	component := views.Dashboard()
// 	return html(c, http.StatusOK, component)
// }

type CardInput struct {
	question string `json:"question"`
	answer   string `json:"answer"`
}

func (env *Env) Cards(c echo.Context) error {
	cardInput := CardInput{}
	err := c.Bind(&cardInput)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	cardParams := queries.CreateCardParams{
		Question: cardInput.question,
		Answer:   cardInput.answer,
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

func (env *Env) CreateDeck(ctx context.Context, request api.CreateDeckRequestObject) (api.CreateDeckResponseObject, error) {
	input := request.Body
	deck, err := env.decks.Queries.CreateDeck(ctx, input.Name)
	if err != nil {
		return api.CreateDeck201JSONResponse{}, err
	}
	return api.CreateDeck201JSONResponse{
		Id:   int(deck.ID),
		Name: deck.Name,
	}, nil
}

func (env *Env) GetDeckById(ctx context.Context, request api.GetDeckByIdRequestObject) (api.GetDeckByIdResponseObject, error) {
	return api.GetDeckById200JSONResponse{}, nil
}

func Map[T, U any](seq iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for a := range seq {
			if !yield(f(a)) {
				return
			}
		}
	}
}

func (env *Env) GetCardsByDeckId(ctx context.Context, request api.GetCardsByDeckIdRequestObject) (api.GetCardsByDeckIdResponseObject, error) {
	fmt.Println("getting cards")
	parsedDeckId := request.DeckId
	deckId := int32(parsedDeckId)
	fmt.Println(deckId)

	cards, err := env.decks.Queries.ListCards(ctx, deckId)
	if err != nil {
		return api.GetCardsByDeckId200JSONResponse{}, nil
	}

	cardsResponse := []api.Card{}
	for _, c := range cards {
		card := api.Card{
			Id:       c.ID,
			Question: c.Question,
			Answer:   c.Answer,
			DeckId:   c.DeckID,
		}
		cardsResponse = append(cardsResponse, card)
	}

	return api.GetCardsByDeckId200JSONResponse{Cards: cardsResponse}, nil
}

func (env *Env) CreateCard(ctx context.Context, request api.CreateCardRequestObject) (api.CreateCardResponseObject, error) {
	deckId := request.DeckId
	cardInput := request.Body

	cardParams := queries.CreateCardParams{
		DeckID:   deckId,
		Question: cardInput.Question,
		Answer:   cardInput.Answer,
	}

	card, err := env.decks.Queries.CreateCard(ctx, cardParams)
	if err != nil {
		fmt.Printf("card %s\n", err)
		return api.CreateCard200JSONResponse{}, nil
	}

	output := api.Card{
		Id:       card.ID,
		Question: card.Question,
		Answer:   card.Answer,
		DeckId:   card.DeckID,
	}
	return api.CreateCard200JSONResponse(output), nil
}
