package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"groot/comm/conf"
	"time"
)

var (
	ErrMissingHeader = errors.New("The length of the 'Authorization' header is zero.")
)

//password processing
func secretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		//make sure the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	}
}

//Parse validates the token with the special secret.
// and returns the context if the token is valid.
func Parse(tokenString string, secret string) (*AuthContext, error) {
	ctx := &AuthContext{}

	//Parse the token.
	token, e := jwt.Parse(tokenString, secretFunc(secret))
	//Parse error
	if e != nil {
		return ctx, e
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.UIN = claims["uin"].(string)
		ctx.DeviceID = claims["device_id"].(string)
		ctx.ClientVersion = uint32(claims["client_version"].(float64))
		ctx.DeviceType = claims["device_type"].(string)
		ctx.TokenSecKey = claims["token_sec_key"].(string)
		return ctx, nil
	} else {
		return ctx, e
	}
}

// ParseReqest gets the token from the header end
// pass it to the parse function to parse the token
func ParseRequest(c *gin.Context) (*AuthContext, error) {
	header := c.Request.Header.Get("Authorization")
	// Load the jwt secret from config
	secret := conf.GetAppConfig().AuthToken.JwtSignKeySecret //viper.GetString("jwt_secret")

	if len(header) == 0 {
		return &AuthContext{}, ErrMissingHeader
	}

	var t string
	//parse the header to get the token part.
	fmt.Sscanf(header, "Bearer %s", &t)
	return Parse(t, secret)
}

//Sign signs the context with sepcial secret
func Sign(ctx *gin.Context, c AuthContext, secret string) (tokenString string, e error) {
	//Load the jwt secret from config if secret is nil
	if secret == "" {
		secret = conf.GetAppConfig().AuthToken.JwtSignKeySecret //viper.GetString("jwt_secret")
	}

	//the token content.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uin":            c.UIN,
		"device_id":      c.DeviceID,
		"client_version": c.ClientVersion,
		"device_type":    c.DeviceType,
		"nbf":            time.Now().Unix(),
		"iat":            time.Now().Unix(),
		"token_sec_key":  c.TokenSecKey,
	})

	tokenString, e = token.SignedString([]byte(secret))

	return
}

//Parse validates the token with the special secret.
// and returns the context if the token is valid.
func ParseLongConnAuthToken(tokenString string) (*LongConnAuthContext, error) {
	// Load the jwt secret from config
	secret := conf.GetAppConfig().AuthToken.SessionSignKeySecret//viper.GetString("session_key_secret")

	ctx := &LongConnAuthContext{}

	//Parse the token.
	token, e := jwt.Parse(tokenString, secretFunc(secret))

	//Parse error
	if e != nil {
		return ctx, e
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.UIN = claims["uin"].(string)
		ctx.DeviceID = claims["device_id"].(string)
		ctx.ClientVersion = uint32(claims["client_version"].(float64))
		ctx.DeviceType = claims["device_type"].(string)
		return ctx, nil
	} else {
		return ctx, e
	}
}

//Sign signs the context with sepcial secret
func SignLongConnAuthToken(c LongConnAuthContext, secret string) (tokenString string, e error) {
	//Load the jwt secret from config if secret is nil
	if secret == "" {
		secret =  conf.GetAppConfig().AuthToken.SessionSignKeySecret// viper.GetString("session_key_secret")
	}

	//the token content.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uin":            c.UIN,
		"device_id":      c.DeviceID,
		"client_version": c.ClientVersion,
		"device_type":    c.DeviceType,
		"nbf":            time.Now().Unix(),
		"iat":            time.Now().Unix(),
	})

	tokenString, e = token.SignedString([]byte(secret))

	return
}

//Parse validates the token with the special secret.
// and returns the context if the token is valid.
func ParseFriendToken(tokenString string) (*FriendContext, error) {
	// Load the jwt secret from config
	secret := conf.GetAppConfig().AuthToken.FriendSignKeySecret //viper.GetString("friend_control_key_secret")

	ctx := &FriendContext{}

	//Parse the token.
	token, e := jwt.Parse(tokenString, secretFunc(secret))

	//Parse error
	if e != nil {
		return ctx, e
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.SendUserName = claims["send_user_name"].(string)
		ctx.ReceiveUserName = claims["receive_user_name"].(string)
		ctx.Timestamp = int64(claims["timestamp"].(float64))
		ctx.EventType = int(claims["event_type"].(float64))
		return ctx, nil
	} else {
		return ctx, e
	}
}

//Sign signs the context with sepcial secret
func SignFriendToken(c FriendContext, secret string) (tokenString string, e error) {
	//Load the jwt secret from config if secret is nil
	if secret == "" {
		secret = conf.GetAppConfig().AuthToken.FriendSignKeySecret //viper.GetString("friend_control_key_secret")
	}

	//the token content.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"send_user_name":    c.SendUserName,
		"receive_user_name": c.ReceiveUserName,
		"timestamp":         c.Timestamp,
		"event_type":        c.EventType,
		"nbf":               time.Now().Unix(),
		"iat":               time.Now().Unix(),
	})

	tokenString, e = token.SignedString([]byte(secret))

	return
}

//Parse validates the token with the special secret.
// and returns the context if the token is valid.
func ParseFileToken(tokenString string) (*FileContext, error) {
	// Load the jwt secret from config
	secret := conf.GetAppConfig().AuthToken.FileSignKeySecret //viper.GetString("file_key_secret")

	ctx := &FileContext{}

	//Parse the token.
	token, e := jwt.Parse(tokenString, secretFunc(secret))

	//Parse error
	if e != nil {
		return ctx, e
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.UIN = claims["uin"].(string)
		ctx.DeviceID = claims["device_id"].(string)
		ctx.ClientVersion = uint32(claims["client_version"].(float64))
		ctx.DeviceType = claims["device_type"].(string)
		return ctx, nil
	} else {
		return ctx, e
	}
}

//Sign signs the context with sepcial secret
func SignFileToken(c FileContext, secret string) (tokenString string, e error) {
	//Load the jwt secret from config if secret is nil
	if secret == "" {
		secret = conf.GetAppConfig().AuthToken.FileSignKeySecret //viper.GetString("file_key_secret")
	}

	//the token content.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uin":            c.UIN,
		"device_id":      c.DeviceID,
		"client_version": c.ClientVersion,
		"device_type":    c.DeviceType,
		"nbf":            time.Now().Unix(),
		"iat":            time.Now().Unix(),
	})

	tokenString, e = token.SignedString([]byte(secret))

	return
}

//Parse validates the token with the special secret.
// and returns the context if the token is valid.
func ParseSMSCodeToken(tokenString string) (*SMSCodeContext, error) {
	// Load the jwt secret from config
	secret := conf.GetAppConfig().AuthToken.SmsCodeSignKeySecret // viper.GetString("sms_code_key_secret")

	ctx := &SMSCodeContext{}

	//Parse the token.
	token, e := jwt.Parse(tokenString, secretFunc(secret))

	//Parse error
	if e != nil {
		return ctx, e
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		ctx.UserName = claims["user_name"].(string)
		ctx.DeviceID = claims["device_id"].(string)
		ctx.TimeStamp = int64(claims["time_stamp"].(float64))
		return ctx, nil
	} else {
		return ctx, e
	}
}

//Sign signs the context with sepcial secret
func SignSMSCodeToken(c SMSCodeContext, secret string) (tokenString string, e error) {
	//Load the jwt secret from config if secret is nil
	if secret == "" {
		secret = conf.GetAppConfig().AuthToken.SmsCodeSignKeySecret //viper.GetString("file_key_secret")
	}

	//the token content.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name":  c.UserName,
		"device_id":  c.DeviceID,
		"time_stamp": c.TimeStamp,
	})

	tokenString, e = token.SignedString([]byte(secret))

	return
}
