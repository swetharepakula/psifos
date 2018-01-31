package server

import (
	"github.com/icza/gowut/gwu"
)

type ButtonHandler struct {
	s        *PsifosServer
	category string
}

func NewButtonHandler(s *PsifosServer, category string) *ButtonHandler {
	return &ButtonHandler{
		s:        s,
		category: category,
	}
}

func (h *ButtonHandler) HandleEvent(gwu.Event) {

	switch h.category {
	case "categoryOne":
		h.s.PutCategoryOne()
	case "categoryTwo":
		h.s.PutCategoryTwo()
	case "setup":
		h.s.SetupDatabase()

	}
}
