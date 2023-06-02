package go_server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	inits "go_server/initializers"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

var db *sql.DB
var rdb *redis.ClusterClient
var ctx context.Context

var (
	id                   int
	event_name           string
	event_type           string
	user_id              string
	wallet_address       string
	insert_id            string
	timestamp            string
	app_version          string
	event_properties     string
	browser_type         string
	os_name              string
	os_version           string
	device_type          string
	auto_track_timestamp string
	url                  string
	referrer             string
	country              string
	language             string
	user_properties      string
	enviroment_type      string
)

type AuthProject struct {
	TenantId   int    `json:"tenant_id"`
	TenantName string `json:"tenant_name"`
	ApiKey     string `json:"api_key"`
}

type UserIds struct {
	UserId        int    `json:"user_id"`
	Walletaddress string `json:"wallet_address"`
}

type OtherAdmVars struct {
	InsertId   int    `json:"insert_id"`
	Timestamp  int    `json:"timestamp"`
	AppVersion string `json:"app_version"`
}

type AutoTrack struct {
	BrowserType string `json:"browser_type"`
	OSName      string `json:"os_name"`
	OSVersion   string `json:"os_version"`
	DeviceType  string `json:"device_type"`
	Timestamp   int    `json:"timestamp"`
	URL         string `json:"url"`
	Referrer    string `json:"referrer"`
	Country     string `json:"country"`
	Language    string `json:"language"`
}

type DbInfo struct {
	EventName       string                 `json:"event_name"`
	EventType       string                 `json:"event_type"`
	AuthProject     AuthProject            `json:"auth_project"`
	UserIds         UserIds                `json:"user_ids"`
	OtherAmdVars    OtherAdmVars           `json:"other_admin_variables"`
	EventProperties map[string]interface{} `json:"event_properties"`
	AutoTrack       AutoTrack              `json:"automatically_tracked"`
	UserProperties  map[string]interface{} `json:"user_properties"`
	EnviornmentType string                 `json:"environment_type"`
}

type DbInfo2 struct { // By Python
	EventName string `json:"event_name"`
	EventType string `json:"event_type"`
	// AuthProject     AuthProject       `json:"auth_project"`
	UserIds         UserIds           `json:"user_ids"`              // отдельное поле +
	OtherAdmVars    OtherAdmVars      `json:"other_admin_variables"` // отдельное поле +
	EventProperties map[string]string `json:"event_properties"`      // отдельный json
	AutoTrack       AutoTrack         `json:"automatically_tracked"` // отдельное поле
	UserProperties  map[string]string `json:"user_properties"`       // отдельный json
}

type GetForDb struct {
	TenantName    string                 `form:"tenant_name"`
	ApyKey        string                 `form:"api_key"`
	To            int                    `form:"to"`
	From          int                    `form:"from"`
	EventName     string                 `form:"event_name"`
	EventType     string                 `form:"event_type"`
	WalletAddress string                 `form:"wallet_address"`
	BrowserType   string                 `form:"browser_type"`
	UserProp      map[string]interface{} `form:"user_prop"`
	EventProp     map[string]interface{} `form:"event_prop"`
}

type GetResponse struct {
	Id                   int
	Event_name           string
	Event_type           string
	User_id              string
	Wallet_address       string
	Insert_id            string
	Timestamp            string
	App_version          string
	Event_properties     map[string]interface{}
	Browser_type         string
	Os_name              string
	Os_version           string
	Device_type          string
	Auto_track_timestamp string
	Url                  string
	Referrer             string
	Country              string
	Language             string
	User_properties      map[string]interface{}
	EnviornmentType      string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetEvent(c *gin.Context) {
	var getData GetForDb
	ServerApiKey := os.Getenv("API_KEY")

	if err := c.Bind(&getData); err != nil {

		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Can not find data"})
		return
	}
	if getData.ApyKey != ServerApiKey {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Bad api key"})
		return
	}
	var columns = []GetResponse{}

	db, ero := inits.GetDb(getData.TenantName, getData.ApyKey, false)
	if ero != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Wrong tenant name"})
		return
	}
	defer db.Close()

	var query string

	// dbstrings := make(map[string]interface{})
	FirstQuery := true
	EndQuery := " ORDER BY id DESC limit 100"

	if getData.To != 0 && getData.From != 0 {
		FirstQuery = false
		query = fmt.Sprintf("SELECT * FROM events WHERE auto_track_timestamp >= %+v and auto_track_timestamp <= %+v", getData.From, getData.To)
	} else {
		query = "SELECT * FROM events "
	}
	if getData.EventName != "" {
		if !FirstQuery {
			query += " AND lower(event_name) LIKE lower('%" + getData.EventName + "%')"
		} else {
			query += " WHERE lower(event_name) LIKE lower('%" + getData.EventName + "%')"
			FirstQuery = false
		}
	}
	if getData.EventType != "" {
		if !FirstQuery {
			query += " AND lower(event_type) LIKE lower('%" + getData.EventType + "%')"
		} else {
			query += " WHERE lower(event_type) LIKE lower('%" + getData.EventType + "%')"
			FirstQuery = false
		}
	}
	if getData.BrowserType != "" {
		if !FirstQuery {
			query += " AND lower(browser_type) LIKE lower('%" + getData.BrowserType + "%')"
		} else {
			query += " WHERE lower(browser_type) LIKE lower('%" + getData.BrowserType + "%')"
			FirstQuery = false
		}
	}
	if getData.WalletAddress != "" {
		if !FirstQuery {
			query += " AND lower(wallet_address) LIKE lower('%" + getData.WalletAddress + "%')"

		} else {
			query += " WHERE lower(wallet_address) LIKE lower('%" + getData.WalletAddress + "%')"
			FirstQuery = false
		}
	}

	if getData.UserProp != nil {
		for key, value := range getData.UserProp {
			if !FirstQuery {
				query += " AND user_properties->>'" + key + "' = '" + value.(string) + "' "

			} else {
				query += " WHERE user_properties->>'" + key + "' = '" + value.(string) + "' "
				FirstQuery = false
			}
		}

	}

	if getData.EventProp != nil {
		for key, value := range getData.EventProp {
			if !FirstQuery {
				query += " AND event_properties->>'" + key + "' = '" + value.(string) + "' "

			} else {
				query += " WHERE event_properties->>'" + key + "' = '" + value.(string) + "' "
				FirstQuery = false
			}
		}

	}

	query += EndQuery
	rows, err := db.Query(query)
	if err != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": err.Error()})
		return
	}
	for rows.Next() {
		column_count, _ := rows.Columns()

		if len(column_count) == 20 {
			er := rows.Scan(
				&id,
				&event_name,
				&event_type,
				&user_id,
				&wallet_address,
				&insert_id,
				&timestamp,
				&app_version,
				&event_properties,
				&browser_type,
				&os_name,
				&os_version,
				&device_type,
				&auto_track_timestamp,
				&url,
				&referrer,
				&country,
				&language,
				&user_properties,
				&enviroment_type,
			)
			if er != nil {
				log.Fatal(er)
			}
		} else {
			er := rows.Scan(
				&id,
				&event_name,
				&event_type,
				&user_id,
				&wallet_address,
				&insert_id,
				&timestamp,
				&app_version,
				&event_properties,
				&browser_type,
				&os_name,
				&os_version,
				&device_type,
				&auto_track_timestamp,
				&url,
				&referrer,
				&country,
				&language,
				&user_properties,
			)
			if er != nil {
				log.Fatal(er)
			}
		}

		mp := make(map[string]interface{})
		erro := json.Unmarshal([]byte(user_properties), &mp)
		if erro != nil {
			panic(erro)
		}
		mp_2 := make(map[string]interface{})
		erro_2 := json.Unmarshal([]byte(event_properties), &mp_2)
		if erro_2 != nil {
			panic(erro)
		}
		columns = append(columns, GetResponse{
			Id:                   id,
			Event_name:           event_name,
			Event_type:           event_type,
			User_id:              user_id,
			Wallet_address:       wallet_address,
			Insert_id:            insert_id,
			Timestamp:            timestamp,
			App_version:          app_version,
			Event_properties:     mp_2,
			Browser_type:         browser_type,
			Os_name:              os_name,
			Os_version:           os_version,
			Device_type:          device_type,
			Auto_track_timestamp: auto_track_timestamp,
			Url:                  url,
			Referrer:             referrer,
			Country:              country,
			Language:             language,
			User_properties:      mp,
			EnviornmentType:      enviroment_type,
		})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	c.IndentedJSON(http.StatusOK, columns)
}

func PostEvent(c *gin.Context) {
	var Info DbInfo
	// ServerApiKey := os.Getenv("API_KEY")

	if err := c.BindJSON(&Info); err != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Can not find data"})
		return
	}

	_, check_tenant := inits.GetDb(Info.AuthProject.TenantName, Info.AuthProject.ApiKey, true)
	if check_tenant != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Bad api key"})
		return
	}

	key := Info.AuthProject.TenantName + "_counter_" + RandStringRunes(10)
	value, _ := json.Marshal(Info)

	rdb = inits.GetRedis()

	err := rdb.Set(key, value, 0).Err()
	if err != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Can not add to redis", "error": string(err.Error())})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "PUT, POST, GET, DELETE, OPTIONS")
	c.IndentedJSON(http.StatusCreated, map[string]string{"message": "Successfully added"})
}

func GetUnicEventType(c *gin.Context) {
	var getData GetForDb
	ServerApiKey := os.Getenv("API_KEY")

	if err := c.Bind(&getData); err != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Can not find data"})
		return
	}
	if getData.ApyKey != ServerApiKey {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Bad api key"})
		return
	}
	db, ero := inits.GetDb(getData.TenantName, getData.ApyKey, false)
	if ero != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Wrong tenant name"})
		return
	}
	defer db.Close()

	query := "SELECT DISTINCT event_type FROM events"
	rows, err := db.Query(query)
	if err != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Can not find tenant_name"})
		return
	}

	var EventTypes []string

	for rows.Next() {
		err := rows.Scan(
			&event_type,
		)

		if err != nil {
			log.Fatal(err)
		}
		EventTypes = append(EventTypes, event_type)
	}
	column := make(map[string]interface{})
	column["event_type"] = EventTypes

	defer rows.Close()

	c.IndentedJSON(http.StatusOK, column)

}

func GetUnicEventName(c *gin.Context) {
	var getData GetForDb
	ServerApiKey := os.Getenv("API_KEY")

	if err := c.Bind(&getData); err != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Can not find data"})
		return
	}
	if getData.ApyKey != ServerApiKey {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Bad api key"})
		return
	}
	db, ero := inits.GetDb(getData.TenantName, getData.ApyKey, false)
	if ero != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Wrong tenant name"})
		return
	}
	defer db.Close()

	query := "SELECT DISTINCT event_name FROM events"
	rows, err := db.Query(query)
	if err != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Can not find tenant_name"})
		return
	}

	var EventNames []string

	for rows.Next() {
		err := rows.Scan(
			&event_name,
		)

		if err != nil {
			log.Fatal(err)
		}
		EventNames = append(EventNames, event_name)
	}
	column := make(map[string]interface{})
	column["event_name"] = EventNames

	defer rows.Close()

	c.IndentedJSON(http.StatusOK, column)

}

func GetUnicBrowserType(c *gin.Context) {
	var getData GetForDb
	ServerApiKey := os.Getenv("API_KEY")

	if err := c.Bind(&getData); err != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Can not find data"})
		return
	}
	if getData.ApyKey != ServerApiKey {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Bad api key"})
		return
	}
	db, ero := inits.GetDb(getData.TenantName, getData.ApyKey, false)
	if ero != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Wrong tenant name"})
		return
	}
	defer db.Close()

	query := "SELECT DISTINCT browser_type FROM events"
	rows, err := db.Query(query)
	if err != nil {
		c.IndentedJSON(http.StatusOK, map[string]string{"message": "Can not find tenant_name"})
		return
	}

	var BrowserTypes []string

	for rows.Next() {
		err := rows.Scan(
			&browser_type,
		)

		if err != nil {
			log.Fatal(err)
		}
		BrowserTypes = append(BrowserTypes, browser_type)
	}
	column := make(map[string]interface{})
	column["browser_type"] = BrowserTypes

	defer rows.Close()

	c.IndentedJSON(http.StatusOK, column)

}
