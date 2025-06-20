package api

import (
	"context"
	"errors"
	"fmt"
	"memorybank/database"
	"memorybank/queries"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/clerk/clerk-sdk-go/v2"
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
	Decks database.DeckModel
}

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
	card, err := env.Decks.Queries.CreateCard(context.Background(), cardParams)
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

func (env *Env) CreateDeck(ctx context.Context, request CreateDeckRequestObject) (CreateDeckResponseObject, error) {
	input := request.Body
	deck, err := env.Decks.Queries.CreateDeck(ctx, input.Name)
	if err != nil {
		return CreateDeck500Response{}, err
	}
	return CreateDeck201JSONResponse{
		Id:   int(deck.ID),
		Name: deck.Name,
	}, nil
}

func (env *Env) GetDeckById(ctx context.Context, request GetDeckByIdRequestObject) (GetDeckByIdResponseObject, error) {
	return GetDeckById200JSONResponse{}, nil
}

func (env *Env) GetCardsByDeckId(ctx context.Context, request GetCardsByDeckIdRequestObject) (GetCardsByDeckIdResponseObject, error) {
	fmt.Println("getting cards")
	parsedDeckId := request.DeckId
	deckId := int32(parsedDeckId)
	fmt.Println(deckId)

	cards, err := env.Decks.Queries.ListCards(ctx, deckId)
	if err != nil {
		return GetCardsByDeckId500Response{}, nil
	}

	cardsResponse := []Card{}
	for _, c := range cards {
		card := Card{
			Id:       c.ID,
			Question: c.Question,
			Answer:   c.Answer,
			DeckId:   c.DeckID,
		}
		cardsResponse = append(cardsResponse, card)
	}

	return GetCardsByDeckId200JSONResponse{Cards: cardsResponse}, nil
}

func (env *Env) CreateCard(ctx context.Context, request CreateCardRequestObject) (CreateCardResponseObject, error) {
	deckId := request.DeckId
	cardInput := request.Body

	cardParams := queries.CreateCardParams{
		DeckID:   deckId,
		Question: cardInput.Question,
		Answer:   cardInput.Answer,
	}

	card, err := env.Decks.Queries.CreateCard(ctx, cardParams)
	if err != nil {
		fmt.Printf("card %s\n", err)
		return CreateCard500Response{}, nil
	}

	output := Card{
		Id:       card.ID,
		Question: card.Question,
		Answer:   card.Answer,
		DeckId:   card.DeckID,
	}
	return CreateCard201JSONResponse(output), nil
}
