package hoard

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	api "github.com/oswines/landset"
)

type httpServer struct {
	Hoard *Hoard
}

// w http.ResponseWriter, r *http.Request
func (s *httpServer) handleInsert(c echo.Context) error {
	var req api.InlayDocument

	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request\n")
	}

	id, err := s.Hoard.Insert(req.Inlay)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot insert into DB\n")
	}

	res := api.IDDocument{ID: id}
	err = json.NewEncoder(c.Response().Writer).Encode(res)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot encode something\n")
	}

	return c.String(http.StatusOK, "Inlay inserted\n")
}

func (s *httpServer) handleGetByID(c echo.Context) error {
	var req api.IDDocument

	err := json.NewDecoder(c.Request().Body).Decode(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request\n")
	}

	inlay, err := s.Hoard.Fetch(req.ID)
	if err == ErrIDNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Non-ID error fetching inlay\n")
	}

	log.Printf("Returning %v\n", inlay)

	res := api.InlayDocument{Inlay: inlay}
	err = json.NewEncoder(c.Response().Writer).Encode(res)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Cannot encode something\n")
	}

	return c.String(http.StatusOK, "Fetched Inlay\n")
}

func NewHTTPServer(addr string) {
	var hrd *Hoard
	var err error
	if hrd, err = NewHoard(); err != nil {
		log.Fatal(err)
	}
	server := &httpServer{
		Hoard: hrd,
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/hoard", server.handleInsert)
	e.GET("/hoard", server.handleGetByID)
	e.GET("/", eala)
	e.GET("/admin", admin)

	e.Logger.Fatal(e.Start(addr))

}

func eala(c echo.Context) error {
	return c.String(http.StatusOK, "Eala, world!")
}

func admin(c echo.Context) error {
	return c.String(http.StatusOK, "Eala, admin!")
}
