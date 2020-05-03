package provider

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"github.com/sirupsen/logrus"
)

// LoginHandler catch request from ORY Hydra with login challenge.
func (s Service) LoginHandler(c echo.Context) error {
	ctx := context.Background()
	loginChallenge := c.QueryParam("login_challenge")
	// get login request from Hydra admin API.
	loginResp, err := s.hydra.Admin.GetLoginRequest(
		&admin.GetLoginRequestParams{
			Context:        ctx,
			LoginChallenge: loginChallenge,
		})
	if err != nil {
		logrus.WithError(err).Errorln("getLoginRequest error")
		return err
	}
	// if client already authorized.
	if loginResp.Payload.Skip {
		acceptResp, err := s.hydra.Admin.AcceptLoginRequest(&admin.AcceptLoginRequestParams{
			Body: &models.AcceptLoginRequest{
				Subject:     &loginResp.Payload.Subject,
				Remember:    true,
				RememberFor: 30, // 30 seconds.
			},
			LoginChallenge: loginChallenge,
			Context:        ctx,
		})
		if err != nil {
			logrus.WithError(err).Errorln("acceptLoginRequest error")
			return err
		}
		logrus.WithField("url", acceptResp.Payload.RedirectTo).Infoln("LoginHandler redirect")
		return c.Redirect(http.StatusFound, acceptResp.Payload.RedirectTo)
	}

	return c.Render(http.StatusOK, "login.html", map[string]interface{}{"challenge": loginChallenge})
}

// ConsentHandler catch request from ORY Hydra with consent challenge.
func (s Service) ConsentHandler(c echo.Context) error {
	ctx := context.Background()
	consentChallenge := c.QueryParam("consent_challenge")

	resp, err := s.hydra.Admin.GetConsentRequest(
		&admin.GetConsentRequestParams{
			Context:          ctx,
			ConsentChallenge: consentChallenge,
		})
	if err != nil {
		logrus.WithError(err).Errorln("getConsentRequest error")
		return err
	}

	// account data.
	var (
		accountID        = 10
		accountFirstName = "Elon"
		accountLastName  = "Musk"
	)

	// We want to allow everything that we requested, as this is our application.
	acceptResp, err := s.hydra.Admin.AcceptConsentRequest(
		&admin.AcceptConsentRequestParams{
			Body: &models.AcceptConsentRequest{
				GrantScope:               resp.Payload.RequestedScope,
				GrantAccessTokenAudience: resp.Payload.RequestedAccessTokenAudience,
				Remember:                 true,
				RememberFor:              30,
				Session: &models.ConsentRequestSession{
					// Sets session data for the OpenID Connect ID token.
					IDToken: map[string]interface{}{
						"extra_vars": map[string]interface{}{
							"id":        accountID,
							"firstName": accountFirstName,
							"lastName":  accountLastName,
						},
					},
				},
			},
			ConsentChallenge: consentChallenge,
			Context:          ctx,
		})
	if err != nil {
		logrus.WithError(err).Errorln("AcceptConsentRequest error")
		return err
	}

	logrus.WithField("url", acceptResp.Payload.RedirectTo).Infoln("ConsentHandler redirect")
	return c.Redirect(http.StatusFound, acceptResp.Payload.RedirectTo)
}

// SignInHandler checks the username and password of the user and draws conclusions.
func (s Service) SignInHandler(c echo.Context) error {
	login := c.FormValue("login")
	password := c.FormValue("password")
	loginChallenge := c.FormValue("loginChallenge")
	logrus.
		WithFields(
			logrus.Fields{
				"login":          login,
				"password":       password,
				"loginChallenge": loginChallenge,
			}).
		Infoln("signIn request")
	// Oh no..
	if login != password {
		return errors.New("bad credentials")
	}

	acceptResp, err := s.hydra.Admin.AcceptLoginRequest(&admin.AcceptLoginRequestParams{
		Body: &models.AcceptLoginRequest{
			Subject:     &login,
			Remember:    true,
			RememberFor: 30, // 30 seconds.
		},
		LoginChallenge: loginChallenge,
		Context:        context.Background(),
	})
	if err != nil {
		logrus.WithError(err).Errorln("acceptLoginRequest error")
		return err
	}
	logrus.WithField("url", acceptResp.Payload.RedirectTo).Infoln("SignInHandler redirect")
	return c.Redirect(http.StatusFound, acceptResp.Payload.RedirectTo)
}
