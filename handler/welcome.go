package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"RMS/database"
	"RMS/database/helper"
	"RMS/models"
	"RMS/utilities"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var JwtKey = []byte("secret_key")

func Login(w http.ResponseWriter, r *http.Request) {
	var userDetails models.UsersLoginDetails
	// err := json.NewDecoder(r.Body).Decode(&userDetails)
	// if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	logrus.Printf("Decoder error:%v", err)
	//	return
	// }
	decoderErr := utilities.Decoder(r, &userDetails)

	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	userDetails.Email = strings.ToLower(userDetails.Email)

	userCredentials, fetchErr := helper.FetchPasswordAndIDANDRole(userDetails.Email, userDetails.Role)

	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("ERROR: Wrong details"))
			if err != nil {
				return
			}

			logrus.Printf("FetchPasswordAndId: not able to get password, id, or role:%v", fetchErr)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if PasswordErr := bcrypt.CompareHashAndPassword([]byte(userCredentials.Password), []byte(userDetails.Password)); PasswordErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logrus.Printf("password misMatch")
		_, err := w.Write([]byte("ERROR: Wrong password"))
		if err != nil {
			return
		}
		return
	}

	expiresAt := time.Now().Add(60 * time.Minute)

	claims := &models.Claims{
		ID:   userCredentials.ID,
		Role: userCredentials.Role,
		StandardClaims: jwt.StandardClaims{

			ExpiresAt: expiresAt.Unix(),
			// Issuer:    userCredentials.Role,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("TokenString: cannot create token string:%v", err)
		return
	}

	userOutboundData := make(map[string]interface{})

	userOutboundData["token"] = tokenString

	err = utilities.Encoder(w, userOutboundData)
	if err != nil {
		logrus.Printf("Login: Not able to login:%v", err)
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	var userDetails models.Users

	decoderErr := utilities.Decoder(r, &userDetails)

	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, err := helper.CreateUser(&userDetails, tx)
		if err != nil {
			logrus.Printf("Register:CreateUser:%v", err)
			return err
		}
		userDetails.ID = userID
		err = helper.CreateAddress(userID, &userDetails, tx)
		if err != nil {
			logrus.Printf("Register:CreateAddress:%v", err)
			return err
		}
		err = helper.CreateRole(userID, models.UserRoleUser, tx)
		return err
	})
	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("Register:%v", txErr)
		return
	}

	userOutboundData := make(map[string]int)

	userOutboundData["Successfully Registered: ID is"] = userDetails.ID

	err := utilities.Encoder(w, userOutboundData)
	if err != nil {
		logrus.Printf("Register:%v", err)
		return
	}
}

// func CreateAdmin(w http.ResponseWriter, r *http.Request) {
//	var userDetails models.Users
//
//	err := json.NewDecoder(r.Body).Decode(&userDetails)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		logrus.Printf("decoder error:%v", err)
//		return
//	}
//	txErr := database.Tx(func(tx *sqlx.Tx) error {
//		userID, err := helper.CreateUser(userDetails, tx)
//		if err != nil {
//
//			logrus.Printf("CreateAdmin:CreateUser:%v", err)
//			return err
//		}
//		err = helper.CreateAddress(userID, userDetails, tx)
//		if err != nil {
//			logrus.Printf("CreateAdmin:CreateAddress:%v", err)
//			return err
//		}
//		err = helper.CreateRoleAdmin(userID, tx)
//		return err
//	})
//	if txErr != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		logrus.Printf("CreateAdmin:%v", txErr)
//		return
//	}
//
// }

func CreateSubAdmin(w http.ResponseWriter, r *http.Request) {
	var userDetails models.Users

	decoderErr := utilities.Decoder(r, &userDetails)

	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, err := helper.CreateUser(&userDetails, tx)

		if err != nil {
			logrus.Printf("CreateSubAdmin:CreateUser:%v", err)
			return err
		}

		userDetails.ID = userID
		err = helper.CreateAddress(userID, &userDetails, tx)

		if err != nil {
			logrus.Printf("CreateSubAdmin:CreateAddress:%v", err)
			return err
		}

		err = helper.CreateRole(userID, models.UserRoleSubAdmin, tx)
		return err
	})
	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateSubAdmin:%v", txErr)
		return
	}

	userOutboundData := make(map[string]int)

	userOutboundData["Successfully Created SubAdmin: ID is"] = userDetails.ID

	err := utilities.Encoder(w, userOutboundData)
	if err != nil {
		logrus.Printf("CreateSubAdmin:%v", err)
		return
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var userDetails models.Users

	decoderErr := utilities.Decoder(r, &userDetails)

	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, err := helper.CreateUser(&userDetails, tx)

		if err != nil {
			logrus.Printf("CreateUser:CreateUser:%v", err)
			return err
		}

		userDetails.ID = userID
		err = helper.CreateAddress(userID, &userDetails, tx)

		if err != nil {
			logrus.Printf("CreateUser:CreateAddress:%v", err)
			return err
		}

		err = helper.CreateRole(userID, models.UserRoleUser, tx)
		return err
	})
	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateUser:%v", txErr)
		return
	}
	userOutboundData := make(map[string]int)

	userOutboundData["Successfully Created User: ID is"] = userDetails.ID

	err := utilities.Encoder(w, userOutboundData)
	if err != nil {
		logrus.Printf("CreateUser:%v", err)
		return
	}
}

func AddAddress(w http.ResponseWriter, r *http.Request) {
	var userDetails models.Users

	decoderErr := utilities.Decoder(r, &userDetails)

	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("AddAddress:QueryParam for ID:%v", ok)
		return
	}

	err := helper.AddAddress(contextValues.ID, &userDetails)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("AddAddress: cannot add another address:%v", err)
		return
	}

	message := "Added Another address successfully"
	err = utilities.Encoder(w, message)
	if err != nil {
		logrus.Printf("AddAddress:%v", err)
		return
	}
}

func CreateRestaurants(w http.ResponseWriter, r *http.Request) {
	var restaurantDetails models.Restaurants

	decoderErr := utilities.Decoder(r, &restaurantDetails)
	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateRestaurant:QueryParam for ID:%v", ok)
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		restaurantID, err := helper.CreateRestaurants(&restaurantDetails, contextValues.ID, tx)
		if err != nil {
			logrus.Printf("CreateRestaurants:CreateRestaurants:%v", err)
			return err
		}
		restaurantDetails.ID = restaurantID
		err = helper.CreateRestaurantAddress(&restaurantDetails, restaurantID, tx)
		return err
	})
	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateRestaurants:%v", txErr)
		return
	}
	userOutboundData := make(map[string]int)

	userOutboundData["Created restaurant with restaurant ID:"] = restaurantDetails.ID

	err := utilities.Encoder(w, userOutboundData)
	if err != nil {
		logrus.Printf("CreateRestaurants:%v", err)
		return
	}
}

// func CreateDishes(w http.ResponseWriter, r *http.Request) {
//	// var route *models.App
//
//	var dishDetails models.DishDetails
//
//	decoderErr := utilities.Decoder(r, &dishDetails)
//	if decoderErr != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		logrus.Printf("Decoder error:%v", decoderErr)
//		return
//	}
//
//	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)
//
//	if !ok {
//		w.WriteHeader(http.StatusInternalServerError)
//		logrus.Printf("CreateRestaurant:QueryParam for ID:%v", ok)
//		return
//	}
//
//	restaurantID, err := strconv.Atoi(chi.URLParam(r, "restaurantID"))
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		logrus.Printf("parsing error :%v", err)
//		return
//	}
//
//	dishID, err := helper.CreateDishes(dishDetails, contextValues.ID, restaurantID)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		logrus.Printf("CreateDishes:%v", err)
//		return
//	}
//
//	userOutboundData := make(map[string]int)
//
//	userOutboundData["Created restaurant with restaurant ID:"] = dishID
//
//	err = utilities.Encoder(w, userOutboundData)
//	if err != nil {
//		logrus.Printf("CreateRestaurants:%v", err)
//		return
//	}
// }

// func SetRole(w http.ResponseWriter, r *http.Request) {
//	id, err := strconv.Atoi(chi.URLParam(r, "ID"))
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		log.Printf("parsing error :%v", err)
//		return
//	}
//	role := chi.URLParam(r, "ROLE")
//
//	err = helper.SetRole(id, role)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("SetRole:error is:%v", err)
//		return
//	}
// }

// func GetAdminDetails(w http.ResponseWriter, r *http.Request) {
//
//	id, ok := r.Context().Value("ID").(int)
//	if !ok {
// 		w.WriteHeader(http.StatusInternalServerError)
//		logrus.Printf("GetAdmin:QueryParam for ID:%v", ok)
//		return
//	}
//
//	adminDetails, GetAdminErr := helper.GetAdminDetails(id)
//	if GetAdminErr != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		logrus.Printf("GetAdmin: error is:%v ", GetAdminErr)
//		return
//	}
//	err := json.NewEncoder(w).Encode(adminDetails)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		logrus.Printf("Encoding error:%v", err)
//		return
//	}
// }

// func AddDishes(w http.ResponseWriter, r *http.Request) {
//	id, ok := r.Context().Value("ID").(int)
//	if !ok {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("AddDishes:QueryParam for ID:%v", ok)
//		return
//	}
//	err := helper.AddDishes(id)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("AddDishes:error is:%v", err)
//		return
//	}
// }

func GetSubAdmin(w http.ResponseWriter, r *http.Request) {
	filterCheck, err := filters(r)
	if err != nil {
		logrus.Printf("GetSubAdmin: filterCheck error:%v", err)
		return
	}

	subAdminDetails, GetSubAdminErr := helper.GetSubAdmin(filterCheck)
	if GetSubAdminErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("GetSubAdmin: not able to get subadmin :%v ", GetSubAdminErr)
		return
	}

	err = utilities.Encoder(w, subAdminDetails)
	if err != nil {
		logrus.Printf("GetSubAdmin:%v", err)
		return
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	filterCheck, err := filters(r)
	if err != nil {
		logrus.Printf("GetUsers:filterCheck:%v", err)
	}

	adminGetUserDetails, AdminGetUserErr := helper.GetUsers(filterCheck)
	if AdminGetUserErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("AdminGetUser: not able to get users :%v ", AdminGetUserErr)
		return
	}

	err = utilities.Encoder(w, adminGetUserDetails)
	if err != nil {
		logrus.Printf("GetUser:%v", err)
		return
	}
}

func filters(r *http.Request) (models.FiltersCheck, error) {
	filtersCheck := models.FiltersCheck{}
	isSearched := false
	searchedName := r.URL.Query().Get("name")
	if searchedName != "" {
		isSearched = true
	}

	var limit int
	var err error
	var page int
	strLimit := r.URL.Query().Get("limit")
	if strLimit == "" {
		limit = 10
	} else {
		limit, err = strconv.Atoi(strLimit)
		if err != nil {
			logrus.Printf("Limit: cannot get limit:%v", err)
			return filtersCheck, err
		}
	}

	strPage := r.URL.Query().Get("page")
	if strPage == "" {
		page = 0
	} else {
		page, err = strconv.Atoi(strPage)
		if err != nil {
			logrus.Printf("Page: cannot get page:%v", err)
			return filtersCheck, err
		}
	}

	filtersCheck = models.FiltersCheck{
		IsSearched:   isSearched,
		SearchedName: searchedName,
		Page:         page,
		Limit:        limit}
	return filtersCheck, nil
}

func GetAllRestaurants(w http.ResponseWriter, r *http.Request) {
	filtersCheck, err := filters(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("GetAllRestaurants: filterCheck error: %v", err)
		return
	}

	restaurantsDetails, restaurantsDetailsErr := helper.GetAllRestaurants(filtersCheck)
	if restaurantsDetailsErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("GetRestaurants: not able to get all restaurant :%v ", restaurantsDetailsErr)
		return
	}

	err = utilities.Encoder(w, restaurantsDetails)
	if err != nil {
		logrus.Printf("GetAllRestaurants:%v", err)
		return
	}
}

func GetRestaurants(w http.ResponseWriter, r *http.Request) {
	filtersCheck, err := filters(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("GetRestaurants: filterCheck error: %v", err)
		return
	}

	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateRestaurant:QueryParam for ID:%v", ok)
		return
	}

	restaurantsDetails, restaurantsDetailsErr := helper.GetRestaurants(contextValues.ID, filtersCheck)
	if restaurantsDetailsErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("SubAdminGetRestaurants: not able to get restaurant:%v ", restaurantsDetailsErr)
		return
	}

	err = utilities.Encoder(w, restaurantsDetails)
	if err != nil {
		logrus.Printf("GetRestaurants:%v", err)
		return
	}
}

func GetDishes(w http.ResponseWriter, r *http.Request) {
	filtersCheck, err := filters(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("GetDishes: filterCheck error: %v", err)
		return
	}

	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateRestaurant:QueryParam for ID:%v", ok)
		return
	}

	restaurantID, err := strconv.Atoi(chi.URLParam(r, "restaurantID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("parsing error :%v", err)
		return
	}

	dishDetails, dishDetailsErr := helper.GetDishes(restaurantID, contextValues.ID, filtersCheck)
	if dishDetailsErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("GetDishes: not able to get dishes :%v ", dishDetailsErr)
		return
	}

	err = utilities.Encoder(w, dishDetails)
	if err != nil {
		logrus.Printf("GetDishes:%v", err)
		return
	}
}

func GetDistance(w http.ResponseWriter, r *http.Request) {
	addressID, err := strconv.Atoi(chi.URLParam(r, "addressID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("parsing error :%v", err)
		return
	}

	userLocation, userLocationErr := helper.FetchUserLocation(addressID)
	if userLocationErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("FetchUserLocation: error is:%v ", userLocationErr)
		return
	}

	restaurantID, err := strconv.Atoi(chi.URLParam(r, "restaurantID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("parsing error :%v", err)
		return
	}

	restaurantLocation, restaurantLocationErr := helper.FetchRestaurantLocation(restaurantID)
	if restaurantLocationErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("FetchRestaurantLocation: error is:%v ", restaurantLocationErr)
		return
	}

	restaurantDistance := helper.GetDistance(userLocation.Lat, userLocation.Lon, restaurantLocation.Lat, restaurantLocation.Lon)
	err = utilities.Encoder(w, restaurantDistance)
	if err != nil {
		logrus.Printf("GetDistance: not able to fetch distance between user and restaurant :%v", err)
		return
	}
}

// func UpdateAdmin(w http.ResponseWriter, r *http.Request) {
//	id, ok := r.Context().Value("ID").(int)
//	if !ok {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("UpdateAdmin:QueryParam for ID:%v", ok)
//		return
//	}
//
//	var adminDetails models.UserDetails
//	err := json.NewDecoder(r.Body).Decode(&adminDetails)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		log.Printf("decoder error:%v", err)
//		return
//	}
//
//	updateAdminErr := helper.UpdateAdmin(id, adminDetails)
//	if updateAdminErr != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("UpdateAdmin: not able  to update admin:%v", err)
//		return
//	}
// }

// func UpdateSubAdmin(w http.ResponseWriter, r *http.Request) {
//	id, ok := r.Context().Value("ID").(int)
//	if !ok {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("UpdateSubAdmin:QueryParam for ID:%v", ok)
//		return
//	}
//
//	var subAdminDetails models.UserDetails
//	err := json.NewDecoder(r.Body).Decode(&subAdminDetails)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		log.Printf("decoder error:%v", err)
//		return
//	}
//
//	updateSubAdminErr := helper.UpdateSubAdmin(id, subAdminDetails)
//	if updateSubAdminErr != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("UpdateSubAdmin: not able  to update subAdmin:%v", err)
//		return
//	}
// }

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateRestaurant:QueryParam for ID:%v", ok)
		return
	}

	var userDetails models.UserDetails
	decoderErr := utilities.Decoder(r, &userDetails)
	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	updateUserErr := helper.UpdateUser(contextValues.ID, &userDetails)
	if updateUserErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("UpdateUser: not able  to update user:%v", updateUserErr)
		return
	}

	message := "updated User Successfully"
	err := utilities.Encoder(w, message)
	if err != nil {
		logrus.Printf("UpdateUser:%v", err)
		return
	}
}

// func UpdateRestaurants(w http.ResponseWriter, r *http.Request) {
//	id, ok := r.Context().Value("ID").(int)
//	if !ok {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("UpdateRestaurants:QueryParam for ID:%v", ok)
//		return
//	}
//
//	var restaurantDetails models.Restaurants
//	err := json.NewDecoder(r.Body).Decode(&restaurantDetails)
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		log.Printf("decoder error:%v", err)
//		return
//	}
//
//	updateRestaurantErr := helper.UpdateRestaurants(id, restaurantDetails)
//	if updateRestaurantErr != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("UpdateRestaurants: not able  to update restaurant:%v", err)
//		return
//	}
// }

func UpdateDish(w http.ResponseWriter, r *http.Request) {
	dishID, err := strconv.Atoi(chi.URLParam(r, "dishID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("parsing error :%v", err)
		return
	}

	var dishesDetails models.Dishes
	decoderErr := utilities.Decoder(r, &dishesDetails)
	if decoderErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("Decoder error:%v", decoderErr)
		return
	}

	updateDishesErr := helper.UpdateDish(dishID, dishesDetails)
	if updateDishesErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("UpdateDishes: not able  to update dishes:%v", err)
		return
	}

	message := "updated dish successfully"
	err = utilities.Encoder(w, message)
	if err != nil {
		logrus.Printf("UpdateDish:%v", err)
		return
	}
}

// func DeleteAdmin(w http.ResponseWriter, r *http.Request) {
//	id, ok := r.Context().Value("ID").(int)
//	if !ok {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("DeleteAdmin:QueryParam for ID:%v", ok)
//		return
//	}
//	err := helper.DeleteAdmin(id)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("DeleteAdmin:not able to delete admin:%v", err)
//		return
//	}
// }
//
// func DeleteSubAdmin(w http.ResponseWriter, r *http.Request) {
//	id, ok := r.Context().Value("ID").(int)
//	if !ok {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("DeleteSubAdmin:QueryParam for ID:%v", ok)
//		return
//	}
//	err := helper.DeleteSubAdmin(id)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("DeleteSubAdmin:not able to delete subAdmin:%v", err)
//		return
//	}
// }

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	contextValues, ok := r.Context().Value(utilities.UserContextKey).(models.ContextValues)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("CreateRestaurant:QueryParam for ID:%v", ok)
		return
	}

	err := helper.DeleteUser(contextValues.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("DeleteUser:not able to delete user:%v", err)
		return
	}

	message := "deleted user successfully"
	err = utilities.Encoder(w, message)
	if err != nil {
		logrus.Printf("DeleteUser:%v", err)
		return
	}
}

// func DeleteRestaurants(w http.ResponseWriter, r *http.Request) {
//	id, ok := r.Context().Value("ID").(int)
//	if !ok {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("DeleteRestaurants:QueryParam for ID:%v", ok)
//		return
//	}
//	err := helper.DeleteRestaurants(id)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		log.Printf("DeleteRestaurants:not able to delete restaurants:%v", err)
//		return
//	}
// }

func DeleteDish(w http.ResponseWriter, r *http.Request) {
	dishID, err := strconv.Atoi(chi.URLParam(r, "dishID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Printf("parsing error :%v", err)
		return
	}

	err = helper.DeleteDish(dishID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.Printf("DeleteDishes:not able to delete dishes:%v", err)
		return
	}

	message := "deleted dish successfully"
	err = utilities.Encoder(w, message)
	if err != nil {
		logrus.Printf("DeleteDish")
		return
	}
}
