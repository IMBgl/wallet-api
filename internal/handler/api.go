package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IMBgl/go-wallet-api/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
	"strconv"

	"github.com/go-chi/render"
	"net/http"
)

type apiHandler struct {
	service service.Service
}

var (
	rValidator *requestValidator
)

func ApiHandler(s service.Service) *apiHandler {
	return &apiHandler{service: s}
}

type requestValidator struct {
	en  locales.Translator
	uni *ut.UniversalTranslator
	vd  *validator.Validate
}

func (h *apiHandler) Routes() *chi.Mux {
	mv := NewApiMiddleware(h.service)
	userHandler := &UserHandler{userService: h.service.User()}
	walletHandler := &WalletHandler{walletService: h.service.Wallet(), middleware: mv}
	categoryHandler := &CategoryHandler{categoryService: h.service.Category(), middleware: mv}
	transactionHandler := &TransactionHandler{transactionService: h.service.Transaction(), middleware: mv}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			render.JSON(w, r, map[string]string{"status": "ok"})
		})

		r.Mount("/user", userHandler.Routes())
		r.Mount("/wallet", walletHandler.Routes())
		r.Mount("/category", categoryHandler.Routes())
		r.Mount("/transaction", transactionHandler.Routes())
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, 404)
		render.JSON(w, r, map[string]string{"error": "not found"})
	})

	return r
}

func DateTimeFormat() string {
	return fmt.Sprintf("2006-01-02 15:04:05")
}

func unmarshallRequest(r *http.Request, data interface{}) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
			return errors.New(fmt.Sprintf("field '%s' must be %s", typeErr.Field, typeErr.Type))
		}

		return errors.New("invalid json format")
	}

	return nil
}

func getValidator() (*requestValidator, error) {
	if rValidator == nil {
		rValidator = &requestValidator{}
		rValidator.en = en.New()
		rValidator.uni = ut.New(rValidator.en, rValidator.en)

		trans, _ := rValidator.uni.GetTranslator("en")

		rValidator.vd = validator.New()
		err := en_translations.RegisterDefaultTranslations(rValidator.vd, trans)
		if err != nil {
			return nil, err
		}

		err = rValidator.vd.RegisterTranslation("alphanumunicode", trans, func(ut ut.Translator) error {
			return ut.Add("required", "{0} can only contain alphabetic characters, numbers or unicode", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())

			return t
		})
		if err != nil {
			return nil, err
		}
	}

	return rValidator, nil
}

func validateRequest(request interface{}) (err error, validationErrors []string) {
	rv, err := getValidator()
	if err != nil {
		return
	}

	verr := rv.vd.Struct(request)
	trans, _ := rv.uni.GetTranslator("en")
	if verr != nil {
		errs := verr.(validator.ValidationErrors)

		for _, ve := range errs.Translate(trans) {
			validationErrors = append(validationErrors, ve)
		}
	}

	return
}

func validateFloat32(value interface{}, name string) (float32, error) {
	switch t := value.(type) {
	case string:
		floatVal, err := strconv.ParseFloat(t, 32)
		if err != nil {
			return 0.0, errors.New(fmt.Sprintf("%s value must be float", name))
		}
		return float32(floatVal), nil
	case float32:
		return t, nil
	case float64:
		return float32(t), nil
	}

	return 0.0, errors.New(fmt.Sprintf("%s value must be float", name))
}

func retrieveTokenOrFail(w http.ResponseWriter, r *http.Request) *service.UserToken {
	ctx := r.Context()
	token, ok := ctx.Value("token").(*service.UserToken)
	if !ok {
		render.Render(w, r, ErrNotFound)
		return nil
	}

	return token
}

func retrieveUuidOrFail(w http.ResponseWriter, r *http.Request, paramName string) (uuidVal uuid.UUID) {
	param := chi.URLParam(r, paramName)
	uuidVal, err := uuid.Parse(param)
	if err != nil {
		render.Render(w, r, ErrNotFound)
		return
	}
	return
}
