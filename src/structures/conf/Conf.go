package conf

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/kirtfieldk/astella/src/api"
	"github.com/kirtfieldk/astella/src/constants/routes"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v3"
)

type Conf struct {
	Aws      AwsConfig `yaml:"aws"`
	Database DbConfig  `yaml:"database"`
	BaseUrl  string    `yaml:"base_url"`
	Port     int       `yaml:"port"`
}

func (c *Conf) GetConf() *Conf {
	yamlFile, err := ioutil.ReadFile("local.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func (c *Conf) BuildApi() {
	var dbConnection = c.CreateDatabaseConnection()
	var baseHandler = api.NewBaseHandler(dbConnection, CreateSession(c.Aws))
	router := gin.Default()
	router.Use(requestIdMiddleware())

	router.GET(routes.GET_EVENT_BY_ID, baseHandler.GetEvent)
	router.GET(routes.GET_EVENTS_MEMBER_OF, baseHandler.GeteventsMemberOf)
	router.GET(routes.GET_EVENTS_MEMBERS, baseHandler.GeteventsMembers)
	router.GET(routes.GET_MESSAGE_IN_EVENT, baseHandler.GetMessagesInEvent)
	router.GET(routes.GET_PIN_MESSAGE, baseHandler.GetPinnedMessaged)
	router.GET(routes.GET_MESSAGE_THREAD, baseHandler.FetchMessageThread)
	router.GET(routes.GET_USRS_LIKE_MESSAGE, baseHandler.GetUserUpvotes)

	router.POST(routes.CREATE_EVENT, baseHandler.CreateEvent)
	router.POST(routes.PIN_MESSAGE, baseHandler.PinMessage)
	router.POST(routes.GET_EVENT_BY_CITY, baseHandler.GetEventByCity)
	router.POST(routes.ADD_USER_TO_EVENT, baseHandler.AddUserToEvent)
	router.POST(routes.POST_MESSAGE_TO_EVENT, baseHandler.PostMessageToEvent)
	router.POST(routes.LIKE_MESSAGE_IN_EVENT, baseHandler.UpvoteMessage)

	router.DELETE(routes.UNPIN_MESSAGE, baseHandler.UnpinMessage)
	router.DELETE(routes.UNLIKE_MESSAGE_IN_EVENT, baseHandler.DownvoteMessage)

	router.PUT(routes.UPDATE_USER, baseHandler.UpdateUser)

	router.Run(c.BaseUrl)
}

func requestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewV4().String())
		c.Next()
	}
}

func (c *Conf) CreateDatabaseConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Database.Host, c.Database.Port, c.Database.Username, c.Database.Password, c.Database.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(psqlInfo)
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

func DbClosedError() error {
	return fmt.Errorf("Database connection closed")
}

func CreateSession(awsConfig AwsConfig) *session.Session {
	session := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(awsConfig.Region),
			Credentials: credentials.NewStaticCredentials(
				awsConfig.AwsAccessKeyId,
				awsConfig.AwsAccessKeySecret,
				"",
			),
		},
	))
	return session
}

func CreateS3Session(session *session.Session) *s3.S3 {
	s3Session := s3.New(session)
	return s3Session
}
