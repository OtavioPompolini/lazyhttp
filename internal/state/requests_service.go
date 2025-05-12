package state

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/OtavioPompolini/project-postman/internal/database"
	"github.com/OtavioPompolini/project-postman/internal/types"
	"github.com/OtavioPompolini/project-postman/internal/utils"
)

type RequestStateService struct {
	state       *State
	persistance database.PersistanceAdapter
}

func newRequestStateService(state *State, db database.PersistanceAdapter) RequestStateService {
	return RequestStateService{
		state:       state,
		persistance: db,
	}
}
