package handler

import (
	"RMS/database"
	"RMS/database/helper"
	"RMS/models"
	"RMS/utilities"
	"cloud.google.com/go/firestore"
	cloud "cloud.google.com/go/storage"
	"database/sql"
	firebase "firebase.google.com/go"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"strconv"

	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type App struct {
	Ctx     context.Context
	Client  *firestore.Client
	Storage *cloud.Client
}

func Upload(r *http.Request) string {
	client := App{}
	client.Ctx = context.Background()
	credentialsFile := option.WithCredentialsJSON([]byte(os.Getenv("firebase_key")))
	// file :=
	var err error

	app, err := firebase.NewApp(client.Ctx, nil, credentialsFile)
	if err != nil {
		logrus.Printf("error is: %v", err)
		return ""
	}
	client.Client, err = app.Firestore(client.Ctx)
	if err != nil {
		logrus.Printf("error is:%v", err)
		return ""
	}
	client.Storage, err = cloud.NewClient(client.Ctx, credentialsFile)
	if err != nil {
		logrus.Printf("error is:%v", err)
		return ""
	}
	file, handler, err := r.FormFile("image")

	if err != nil {
		logrus.Printf("error is:%v", err)
		return ""
	}
	err1 := r.ParseMultipartForm(10 << 20)
	if err1 != nil {
		return ""
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			logrus.Printf("Upload: unable to close file")
			return
		}
	}(file)

	imagePath := handler.Filename

	bucket := "rms-project-ccbd6.appspot.com"

	wc := client.Storage.Bucket(bucket).Object("images/" + imagePath).NewWriter(client.Ctx)
	_, err = io.Copy(wc, file)
	if err != nil {
		logrus.Printf("error is:%v", err)
		return ""
	}
	if err := wc.Close(); err != nil {
		logrus.Printf("error is:%v", err)
		return ""
	}

	URL := "https://storage.cloud.google.com/" + bucket + "/" + "images/" + imagePath

	return URL
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	url := Upload(r)

	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateRestaurant:QueryParam for ID:%v", ok)
		return
	}

	imageType := r.URL.Query().Get("name")

	imageID, err := helper.UploadImage(url, contextValues.ID, imageType)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("UploadImage: not able to upload image:%v", err)
		return
	}
	err = utilities.Encoder(w, imageID)
	if err != nil {
		logrus.Printf("UploadImage: not able to upload image:%v", err)
		return
	}
}

// func BulkInsert(w http.ResponseWriter, r *http.Request) {
//	bulkDishes := make([]models.BulkDishes, 0)
//
//	err := utilities.Decoder(r, &bulkDishes)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		logrus.Printf("Decoder Error:%v", err)
//		return
//	}
//
//	err = helper.BulkDishes(bulkDishes)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		logrus.Printf("BulkInsert: not able to insert bulk dishes:%v", err)
//		return
//	}
// }

// func BulkInsertDishes(w http.ResponseWriter, r *http.Request) {
//	bulkDishes := make([]models.BulkDishes, 0)
//
//	err := utilities.Decoder(r, &bulkDishes)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		logrus.Printf("Decoder Error:%v", err)
//		return
//	}
//
//	err = helper.BulkInsertDishes(bulkDishes)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		logrus.Printf("BulkInsert: not able to insert bulk dishes:%v", err)
//		return
//	}
// }

func BulkInsertDishes(w http.ResponseWriter, r *http.Request) {
	bulkDishes := make([]models.BulkDishes, 0)

	err := utilities.Decoder(r, &bulkDishes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder Error:%v", err)
		return
	}

	if len(bulkDishes) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("BulkInsertDishes: Dishes cannot be empty")
		return
	}

	dishID, err := strconv.Atoi(chi.URLParam(r, "dishID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("parsing error :%v", err)
		return
	}

	err = helper.BulkInsertDishes(bulkDishes, dishID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("BulkInsert: not able to insert bulk dishes:%v", err)
		return
	}
}

func CreateDishes(w http.ResponseWriter, r *http.Request) {
	var dishDetails models.DishDetails

	decoderErr := utilities.Decoder(r, &dishDetails)
	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateDishes:QueryParam for ID:%v", ok)
		return
	}

	restaurantID, err := strconv.Atoi(chi.URLParam(r, "restaurantID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("parsing error :%v", err)
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		dishID, err := helper.CreateDishes(dishDetails, contextValues.ID, restaurantID, tx)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusBadRequest)
				logrus.Printf("CreateDishes:error is:%v", err)
				_, err := w.Write([]byte("Dish already exists in database"))
				if err != nil {
					return err
				}
			} else {
				logrus.Printf("CreateDishes:CreateDishes:%v", err)
				return err
			}
		}
		err = helper.InsertDishImage(dishDetails.ImageID, dishID, tx)
		return err
	})
	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateDishes: InsertDishImage:%v", txErr)
		return
	}
	message := "Created Dish Successfully"

	err = utilities.Encoder(w, message)
	if err != nil {
		logrus.Printf("CreateDishes:%v", err)
		return
	}
}
