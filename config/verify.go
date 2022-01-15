package config

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

var snowflakeVarifyReg = regexp.MustCompile(`\d{18}`)

func VerifySnowflake(id string) (string, bool) {
	strFound := snowflakeVarifyReg.FindString(id)
	ok := len(strFound) != 0
	return strFound, ok
}

//VerifyToken will make a request to the user endpoint with the token as "Authorization" header
//will throw an err if the token is not a user one.
//GET https://discord.com/api/v9/users/@me with a Bot token returns 401
func VerifyToken(token string) error {
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	if err != nil {
		return fmt.Errorf("error checking token: %s", err)
	}
	//Set the necessary headers
	req.Header = http.Header{"Authorization": []string{token}, "Accept": []string{"application/json"}, "User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) discord/1.0.9003 Chrome/91.0.4472.164 Electron/13.4.0 Safari/537.36"}}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error checking token: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 401 {
		return fmt.Errorf("error: the provided token is not a valid user token")
	} else if res.StatusCode > 400 {
		return fmt.Errorf("wrong status code: %d", res.StatusCode)
	}
	return nil
}
