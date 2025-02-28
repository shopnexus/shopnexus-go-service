package http

import (
	"net/http"
	"shopnexus-go-service/gen/pb"
	"shopnexus-go-service/internal/http/response"

	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
)

type AccountHandler struct {
	client pb.AccountClient
}

func NewAccountHandler(client pb.AccountClient) http.Handler {
	h := &AccountHandler{client: client}

	r := chi.NewRouter()
	r.Post("/login/user", h.LoginUser)
	r.Post("/login/admin", h.LoginAdmin)
	r.Post("/register/user", h.RegisterUser)
	r.Post("/register/admin", h.RegisterAdmin)

	return r
}

func (h *AccountHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username *string `json:"username,omitempty"`
		Email    *string `json:"email,omitempty"`
		Phone    *string `json:"phone,omitempty"`
		Password string  `json:"password"`
	}

	if err := sonic.ConfigFastest.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Username == nil && req.Email == nil && req.Phone == nil {
		response.FromMessage(w, http.StatusBadRequest, "Must provide username, email, or phone")
		return
	}

	resp, err := h.client.LoginUser(r.Context(), &pb.LoginUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		response.FromMessage(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.FromDTO(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: resp.Token,
	})
}

func (h *AccountHandler) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := sonic.ConfigFastest.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.client.LoginAdmin(r.Context(), &pb.LoginAdminRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		response.FromMessage(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.FromDTO(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: resp.Token,
	})
}

func (h *AccountHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Gender   string `json:"gender"`
		FullName string `json:"full_name"`
	}

	if err := sonic.ConfigFastest.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.client.RegisterUser(r.Context(), &pb.RegisterUserRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Gender:   pb.Gender(pb.Gender_value[req.Gender]),
		FullName: req.FullName,
	})
	if err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	response.FromDTO(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: resp.Token,
	})
}

func (h *AccountHandler) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := sonic.ConfigFastest.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.client.RegisterAdmin(r.Context(), &pb.RegisterAdminRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		response.FromMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	response.FromDTO(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: resp.Token,
	})
}
