package button

import (
	"fmt"

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

	var err error
	switch h.category {
	case "categoryOne":
		err = h.s.PutCategoryOne()
	case "categoryTwo":
		err = h.s.PutCategoryTwo()
	case "setup":
		err = h.s.SetupDatabase()
	}

	if err != nil {
		fmt.Println("Received error: ", err)
	}
}

func CreateButton(text, category string, s *server.PsifosServer) gwu.Button {
	btn := gwu.NewButton(text)

	btn.Style().SetHeight("70%").SetWidth("30%").SetFontSize("40pt").SetCursor("pointer").Set("border-radius", "15px").SetMarginRight("4%").SetBackground("#8ac2ea").SetMarginTop("2%")

	btn.AddEHandler(NewButtonHandler(s, category), gwu.ETypeClick)

	return btn
}
