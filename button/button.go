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

func (h *ButtonHandler) HandleEvent(e gwu.Event) {

	switch h.category {
	case "categoryOne":
		h.s.PutCategoryOne()
	case "categoryTwo":
		h.s.PutCategoryTwo()
	case "setup":
		h.s.SetupDatabase()
	}
}

func CreateButton(text, category string, s *server.PsifosServer) gwu.Button {
	btn := gwu.NewButton(text)

	btn.Style().SetHeight("70%").SetWidth("30%").SetFontSize("40pt").SetCursor("pointer").Set("border-radius", "15px").SetMarginRight("4%").SetBackground("#8ac2ea").SetMarginTop("2%")

	btn.AddEHandler(NewButtonHandler(s, category), gwu.ETypeClick)

	return btn
}
