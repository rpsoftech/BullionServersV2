package firebase

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/rpsoftech/bullion-server/src/env"
	"google.golang.org/api/option"
)

var firebaseApp *firebase.App
var firebaseCtx context.Context

func Init() {
	if env.Env.APP_ENV == env.APP_ENV_DEVELOPE {
		return
	}
	firebaseCtx = context.Background()
	opt := option.WithCredentialsJSON([]byte(env.Env.FIREBASE_JSON_STRING))

	app, err := firebase.NewApp(firebaseCtx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	firebaseApp = app
	firebasedb, _ := firebaseApp.Database(firebaseCtx)
	ref := firebasedb.NewRef("/")
	c1 := ref.Child(fmt.Sprintf("%d", time.Now().Unix()))
	c1.Set(firebaseCtx, map[string]string{
		"AAAAA": "Ohhyes",
	})
}

// ctx := context.Background()
// conf := &firebase.Config{
//         DatabaseURL: "https://databaseName.firebaseio.com",
// }
// // Fetch the service account key JSON file contents
// opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")

// // Initialize the app with a service account, granting admin privileges
// app, err := firebase.NewApp(ctx, conf, opt)
// if err != nil {
//         log.Fatalln("Error initializing app:", err)
// }

// client, err := app.Database(ctx)
// if err != nil {
//         log.Fatalln("Error initializing database client:", err)
// }

// // As an admin, the app has access to read and write all data, regradless of Security Rules
// ref := client.NewRef("restricted_access/secret_document")
// var data map[string]interface{}
// if err := ref.Get(ctx, &data); err != nil {
//         log.Fatalln("Error reading from database:", err)
// }
// fmt.Println(data)
