package restoauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/RangelReale/osin"
	"gopkg.in/mgo.v2/bson"
)

//OAuthHandler type
type OAuthHandler struct {
	sconfig *osin.ServerConfig
	server  *osin.Server
	Storage *TestStorage
}

//UserData bson store
type UserData bson.M

// AuthorizeClient is the Authorization code endpoint
func (oauth *OAuthHandler) AuthorizeClient(w http.ResponseWriter, r *http.Request) {
	server := oauth.server
	resp := server.NewResponse()
	defer resp.Close()
	if ar := server.HandleAuthorizeRequest(resp, r); ar != nil {
		if !HandleLoginPage(ar, w, r, true) {
			return
		}
		ar.UserData = UserData{"Login": "test"}
		ar.Authorized = true
		server.FinishAuthorizeRequest(resp, r, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		resp.Output["scope"] = "everything"
	}
	//w.Header().Add("Location", "http://www.baidu.com")
	//w.WriteHeader(302)
	osin.OutputJSON(resp, w, r)
}

//GenerateToken  Access token endpoint
func (oauth *OAuthHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	server := oauth.server
	resp := server.NewResponse()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var fbody interface{}
	err = json.Unmarshal(body, &fbody)
	m := fbody.(map[string]interface{})

	var formString bytes.Buffer
	for k, v := range m {
		formString.WriteString("&")
		formString.WriteString(k)
		formString.WriteString("=")
		formString.WriteString(v.(string))
	}

	r.Body = ioutil.NopCloser(bytes.NewReader([]byte(string(formString.String()))))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if ar := server.HandleAccessRequest(resp, r); ar != nil {

		switch ar.Type {
		case osin.AUTHORIZATION_CODE:
			ar.Authorized = true
		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		case osin.CLIENT_CREDENTIALS:
			ar.Authorized = true
		}
		server.FinishAccessRequest(resp, r, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	if !resp.IsError {
		resp.Output["custom_parameter"] = 19923
	}

	defer resp.Close()
	osin.OutputJSON(resp, w, r)
}

//HandleInfo Information endpoint
func (oauth *OAuthHandler) HandleInfo(w http.ResponseWriter, r *http.Request) {
	server := oauth.server
	resp := server.NewResponse()
	if ir := server.HandleInfoRequest(resp, r); ir != nil {
		server.FinishInfoRequest(resp, r, ir)
	}
	if resp.IsError && resp.InternalError != nil {
		fmt.Printf("ERROR: %s\n", resp.InternalError)
	}
	defer resp.Close()
	osin.OutputJSON(resp, w, r)
}

//NewOAuthHandler new the oauth handler
func NewOAuthHandler(session string, dbName string) *OAuthHandler {
	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.CODE, osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN, osin.PASSWORD, osin.CLIENT_CREDENTIALS, osin.ASSERTION}

	sconfig.AllowClientSecretInParams = true
	sconfig.AllowGetAccessRequest = true
	storage := NewTestStorage()
	server := osin.NewServer(sconfig, storage)
	return &OAuthHandler{sconfig, server, storage}

}
