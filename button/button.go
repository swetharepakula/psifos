package button

import (
	"github.com/icza/gowut/gwu"
	"github.com/swetharepakula/psifos/server"
)

type ButtonHandler struct {
	s        *server.PsifosServer
	category string
}

func NewButtonHandler(s *server.PsifosServer, category string) *ButtonHandler {
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
