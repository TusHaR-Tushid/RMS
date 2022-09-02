package helper

import (
	"RMS/database"
	"RMS/models"
	"log"
	"math"
	"strings"

	// using context
	_ "context"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func FetchRole(userID int) (string, error) {
	SQL := `SELECT roles 
            FROM   roles
            JOIN   users ON roles.user_id=users.id
            WHERE  users.id=$1
            AND    users.archived_at IS NULL `
	var userRole string
	err := database.RmsDB.Get(userRole, SQL, userID)
	if err != nil {
		log.Printf("FetchRole:%v", err)
		return "", err
	}
	return userRole, nil
}

func FetchPasswordAndIDANDRole(userMail, userRole string) (models.UserCredentials, error) {
	SQL := `SELECT users.id,password,role
            FROM   users
            JOIN   roles ON users.id=roles.user_id
            WHERE  email=$1
            AND    roles.role=$2 
            AND    archived_at IS NULL `

	var userCredentials models.UserCredentials

	err := database.RmsDB.Get(&userCredentials, SQL, userMail, userRole)
	if err != nil {
		logrus.Printf("FetchPassword: Not able to fetch password, ID or role: %v", err)
		return userCredentials, err
	}
	return userCredentials, nil
}

func CreateUser(userDetails *models.Users, tx *sqlx.Tx) (int, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Printf("CreateUser: Not able to hash password:%v", err)
		return userDetails.ID, err
	}

	// language=SQL
	SQL := `INSERT  INTO  users(name, email, password, phone_no, age, gender)
           VALUES  ($1, $2, $3, $4, $5, $6)
           RETURNING id`

	var id int
	userDetails.Email = strings.ToLower(userDetails.Email)

	err = tx.Get(&id, SQL, userDetails.Name, userDetails.Email, hashPassword, userDetails.PhoneNo, userDetails.Age, userDetails.Gender)

	if err != nil {
		logrus.Printf("CreateUser: Not able to Create User :%v", err)
		return id, err
	}

	return id, nil
}

func CreateAddress(id int, userDetails *models.Users, tx *sqlx.Tx) error {
	SQL := `INSERT INTO user_address(user_id, address, lat, lon)
          VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(SQL, id, userDetails.Address, userDetails.Lat, userDetails.Lon)

	if err != nil {
		logrus.Printf("CreateAddress: unable to create address:%v", err)
		return err
	}

	return nil
}

func AddAddress(userID int, userDetails *models.Users) error {
	SQL := `INSERT INTO user_address(user_id, address, lat, lon)
          VALUES ($1, $2, $3, $4)`
	_, err := database.RmsDB.Exec(SQL, userID, userDetails.Address, userDetails.Lat, userDetails.Lon)

	if err != nil {
		logrus.Printf("CreateAddress: unable to create address:%v", err)
		return err
	}

	return nil
}

func CreateRole(id int, role models.UserProfile, tx *sqlx.Tx) error {
	SQL := `INSERT INTO roles(user_id, role)
           VALUES ($1, $2)`

	_, err := tx.Exec(SQL, id, role)

	if err != nil {
		logrus.Printf("CreateRoleUser: Not able to set user role: %v", err)
		return err
	}

	return nil
}

// func CreateRoleAdmin(id int, tx *sqlx.Tx) error {
//	SQL := `INSERT INTO roles(user_id, role)
//           VALUES ($1, $2)`
//
//	_, err := tx.Exec(SQL, id, models.UserRoleAdmin)
//
//	if err != nil {
//		logrus.Printf("CreateRoleAdmin:%v", err)
//		return err
//	}
//
//	return nil
// }
//
// func CreateRoleSubAdmin(id int, tx *sqlx.Tx) error {
//	SQL := `INSERT INTO roles(user_id, role)
//           VALUES ($1, $2)`
//
//	_, err := tx.Exec(SQL, id, models.UserRoleSubAdmin)
//
//	if err != nil {
//		logrus.Printf("CreateRoleSubAdmin: Not able to set SubAdmin role of user: %v", err)
//		return err
//	}
//
//	return nil
// }

// func CreateAdmin(userDetails models.Users) (int, error) {
//	hash, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
//	if err != nil {
//		log.Printf("CreateUser: error is:%v", err)
//
//	}
//
//	language=SQL
//	SQL := `INSERT  INTO  users(name,email,password,phone_no,age,gender)
//           VALUES  ($1,$2,$3,$4,$5,$6)
//           RETURNING id`
//
//	var id int
//	userDetails.Email = strings.ToLower(userDetails.Email)
//
//	err = database.RmsDb.Get(&id, SQL, userDetails.Name, userDetails.Email, hash, userDetails.PhoneNo, userDetails.Age, userDetails.Gender)
//	if err != nil {
//		log.Printf("CreateUser:error is :%v", err)
//		return id, err
//	}
//
//	userDetails.UserId = id
//	language=SQL
//	SQL = `INSERT INTO user_address(user_id,address,lat,lon)
//          VALUES ($1,$2,$3,$4)`
//	_, err = database.RmsDb.Exec(SQL, userDetails.UserId, userDetails.Address, userDetails.Lat, userDetails.Lon)
//	if err != nil {
//		log.Printf("CreateUser: unable to create address:%v", err)
//		return id, err
//	}
//
//  language=SQL
//	SQL = `INSERT INTO roles(user_id,roles)
//           VALUES ($1,$2)`
//
//	_, err = database.RmsDb.Exec(SQL, userDetails.UserId, "admin")
//
//	return id, nil
// }
//
//  func CreateSubAdmin(userDetails models.Users) (int, error) {
//	hash, err := bcrypt.GenerateFromPassword([]byte(userDetails.Password), bcrypt.DefaultCost)
//	if err != nil {
//		log.Printf("CreateUser: error is:%v", err)
//
//	}
//
//  language=SQL
//	SQL := `INSERT  INTO  users(name,email,password,phone_no,age,gender)
//           VALUES  ($1,$2,$3,$4,$5,$6)
//           RETURNING id`
//
//	var id int
//	userDetails.Email = strings.ToLower(userDetails.Email)
//
//	err = database.RmsDb.Get(&id, SQL, userDetails.Name, userDetails.Email, hash, userDetails.PhoneNo, userDetails.Age, userDetails.Gender)
//	if err != nil {
//		log.Printf("CreateUser:error is :%v", err)
//		return id, err
//	}
//
//	userDetails.UserId = id
//	//language=SQL
//	SQL = `INSERT INTO user_address(user_id,address,lat, lon)
//          VALUES ($1,$2,$3,$4)`
//	_, err = database.RmsDb.Exec(SQL, userDetails.UserId, userDetails.Address, userDetails.Lat, userDetails.Lon)
//	if err != nil {
//		log.Printf("CreateUser: unable to create address:%v", err)
//		return id, err
//	}
//
//	language=SQL
//	SQL = `INSERT INTO roles(user_id,roles)
//           VALUES ($1,$2)`
//
//	_, err = database.RmsDb.Exec(SQL, userDetails.UserId, "sub_admin")
//
//	return id, nil
// }

func CreateRestaurants(restaurantDetails *models.Restaurants, id int, tx *sqlx.Tx) (int, error) {
	// language=SQL
	SQL := `INSERT INTO restaurants(name, created_by)
            VALUES ($1, $2)
            RETURNING id`

	var restaurantID int
	err := tx.Get(&restaurantID, SQL, restaurantDetails.Name, id)

	if err != nil {
		logrus.Printf("CreateRestaurants: Cannot create restaurant :%v", err)
		return restaurantID, err
	}

	return restaurantID, nil
}

func CreateRestaurantAddress(restaurantDetails *models.Restaurants, restaurantID int, tx *sqlx.Tx) error {
	// language=SQL
	SQL := `INSERT INTO restaurant_address(restaurant_id, address, lat, lon)
          VALUES ($1, $2, $3, $4)`
	_, err := tx.Exec(SQL, restaurantID, restaurantDetails.Address, restaurantDetails.Lat, restaurantDetails.Lon)

	if err != nil {
		logrus.Printf("CreateRestaurantAddress: unable to create address:%v", err)
		return err
	}

	return nil
}

// func CreateDishes(dishDetails models.DishDetails, createdBy, restaurantID int) (int, error) {
//	// language=SQL
//	SQL := `INSERT INTO dishes(name, price, created_by, restaurant_id)
//            VALUES ($1, $2, $3, $4)
//            ON  CONFLICT (name, restaurant_id)
//            WHERE archived_at IS NULL
//            DO NOTHING
//            RETURNING id`
//
//	var dishID int
//
//	err := database.RmsDB.Get(&dishID, SQL, dishDetails.Name, dishDetails.Price, createdBy, restaurantID)
//	if err != nil {
//		logrus.Printf("CreateDishes: unable to create dish :%v", err)
//		return dishID, err
//	}
//
//	return dishID, nil
// }

func CreateDishes(dishDetails models.DishDetails, createdBy, restaurantID int, tx *sqlx.Tx) (int, error) {
	// language=SQL
	SQL := `INSERT INTO dishes(name, price, created_by, restaurant_id)
            VALUES ($1, $2, $3, $4)
            ON  CONFLICT (name, restaurant_id)
            WHERE archived_at IS NULL 
            DO NOTHING
            RETURNING id`

	var dishID int

	err := tx.Get(&dishID, SQL, dishDetails.Name, dishDetails.Price, createdBy, restaurantID)
	if err != nil {
		logrus.Printf("CreateDishes: unable to create dish :%v", err)
		return dishID, err
	}

	return dishID, nil
}

func InsertDishImage(imageID, dishID int, tx *sqlx.Tx) error {
	if imageID == 0 {
		return nil
	}
	SQL := `INSERT INTO dish_per_image(dish_id, image_id)
            VALUES ($1, $2)`

	_, err := tx.Exec(SQL, dishID, imageID)
	if err != nil {
		logrus.Printf("InsertDishImage: unable to insert image for ths dish:%v", err)
		return err
	}
	return nil
}

// func BulkDishes(bulkDishes []models.BulkDishes) error {
//	valueStrings := make([]string, 0, len(bulkDishes))
//	valueArgs := make([]interface{}, 0, len(bulkDishes)*3)
//	i := 0
//	for _, post := range bulkDishes {
//		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
//		valueArgs = append(valueArgs, post.Name)
//		valueArgs = append(valueArgs, post.Price)
//		valueArgs = append(valueArgs, post.CreatedBy)
//		valueArgs = append(valueArgs, post.RestaurantID)
//		valueArgs = append(valueArgs, post.ImageID)
//		i++
//	}
//
//	stmt := fmt.Sprintf("INSERT INTO dishes (name, price, created_by, restaurant_id, image_id) VALUES %s",
//		strings.Join(valueStrings, ","))
//
//	// stmt := ("INSERT INTO images(url, image_type, image_id)
//	//        VALUES  %s", strings.Join(valueStrings, ","))
//
//	_, err := database.RmsDB.Exec(stmt, valueArgs...)
//	if err != nil {
//		logrus.Printf("BulkDishes: error is: %v", err)
//		return err
//	}
//	return nil
// }

// func BulkInsert(bulkImages []*models.BulkImages) error {
//	valueStrings := make([]string, 0, len(bulkImages))
//	valueArgs := make([]interface{}, 0, len(bulkImages)*3)
//	for _, post := range bulkImages {
//		valueStrings = append(valueStrings, "($1, $2, $3)")
//		valueArgs = append(valueArgs, post.URL)
//		valueArgs = append(valueArgs, post.ImageType)
//		valueArgs = append(valueArgs, post.ImageID)
//	}
//
//	stmt := fmt.Sprintf("INSERT INTO images (url, image_type, image_id) VALUES %s",
//		strings.Join(valueStrings, ","))
//
//	// stmt := ("INSERT INTO images(url, image_type, image_id)
//	//        VALUES  %s", strings.Join(valueStrings, ","))
//
//	_, err := database.RmsDB.Exec(stmt, valueArgs...)
//	if err != nil {
//		logrus.Printf("BulkInsert: error is: %v", err)
//		return err
//	}
//	return nil
// }

func UploadImage(imageURL string, userID int, imageType string) (int, error) {
	// language=SQL
	var imageID int
	SQL := `INSERT INTO images(url, image_type, uploaded_by)
          VALUES  ($1, $2, $3)`
	err := database.RmsDB.Get(&imageID, SQL, imageURL, imageType, userID)
	if err != nil {
		logrus.Printf("StoreImage: not able to store image in db:%v", err)
		return imageID, err
	}
	return imageID, nil
}

func GetAdminDetails(id int) ([]models.UserDetails, error) {
	// language=SQL
	SQL := `SELECT  
                     users.id,
                     users.name,
                     users.email,
                     users.phone_no,
                     users.password,
                     users.age,
                     users.gender,
                     users.created_at,
                     users.updated_at,
                     roles.role,
                     user_address.address,
                     user_address.lat,
                     user_address.lon
            FROM  users 
            JOIN roles ON users.id=roles.user_id
            JOIN  user_address
            ON   roles.user_id=user_address.user_id
            WHERE users.archived_at IS NULL
            AND   user_address.archived_at IS NULL 
            AND   role=$1
            AND   users.id=$2`

	adminDetails := make([]models.UserDetails, 0)

	err := database.RmsDB.Select(&adminDetails, SQL, models.UserRoleAdmin, id)
	if err != nil {
		logrus.Printf("GetAdmin: unable to fetch admin details:%v", err)
		return nil, err
	}
	return adminDetails, nil
}

// func AddDishes(id int) error {
//	var restaurantID int
//	var dishID int
//
//	SQL := `SELECT id
//            FROM   restaurants
//            WHERE  created_by=$1`
//	err := database.RmsDb.Get(&restaurantID, SQL, id)
//	if err != nil {
//		log.Printf("AddDish: unable to get restaurant id:%v", err)
//		return err
//	}
//
//	SQL = `SELECT id
//            FROM   dishes
//            WHERE  created_by=$1`
//	err = database.RmsDb.Get(&dishID, SQL, id)
//	if err != nil {
//		log.Printf("AddDish: unable to get dish id:%v", err)
//		return err
//	}
//
//	SQL = `INSERT INTO dish_per_restaurant( restaurant_id, dish_id)
//           VALUES ($1,$2)`
//	_, err = database.RmsDb.Exec(SQL, restaurantID, dishID)
//	if err != nil {
//		log.Printf("AddDish: unable to add dish:%v", err)
//		return err
//	}
//	return nil
// }

func SetRole(id int, role string) error {
	// language=SQL
	SQL := `INSERT INTO roles(user_id, role)
            VALUES  ($1, $2)`

	_, err := database.RmsDB.Exec(SQL, id, role)
	if err != nil {
		logrus.Printf("SetRole:error is:%v", err)
		return err
	}
	return nil
}

// func AdminGetSubAdmin() ([]models.UserDetails, error) {
//	language=SQL
//	SQL := `SELECT
//                     id,
//                     name,
//                     email,
//                     phone_no,
//                     password,
//                     age,
//                     gender,
//                     created_at,
//                     updated_at,
//                     roles,
//                     address,
//                     lat,
//                     lon,
//                     archived_at
//            FROM  users JOIN roles
//            ON    users.id=roles.user_id
//            JOIN   user_address
//            ON     roles.user_id=user_address.user_id
//            WHERE users.archived_at IS NULL
//            AND   user_address.archived_at IS NULL
//            AND   roles=$1`
//
//	subAdminDetails := make([]models.UserDetails, 0)
//
//	err := database.RmsDb.Select(&subAdminDetails, SQL, models.UserRoleSubAdmin)
//	if err != nil {
//		logrus.Printf("AdminGetSubAdmin: unable to fetch sub-admin details:%v", err)
//		return nil, err
//	}
//	return subAdminDetails, nil
// }

func GetSubAdmin(filterCheck models.FiltersCheck) (models.TotalUser, error) {
	// language=SQL
	var totalUser models.TotalUser
	SQL := `WITH  cte_sub_admin AS(
            SELECT  
                     count(*) over () as total_count,
                     users.id as id,
                     name ,
                     email ,
                     phone_no ,
                     password ,
                     age,
                     gender,
                     users.created_at as created_at,
                     users.updated_at as updated_at,
                     role,
                     address,
                     lat,
                     lon
            FROM  users 
            JOIN roles ON users.id=roles.user_id
            JOIN  user_address ON  roles.user_id=user_address.user_id
            WHERE users.archived_at IS NULL
            AND   user_address.archived_at IS NULL
            AND   role=$1
            AND   ($2 or name ilike '%' || $3 || '%')
            LIMIT $4 OFFSET $5
            )
            SELECT 
                     total_count,
                      id,
                     name,
                     email,
                     phone_no,
                     password,
                     age,
                     gender,
                     created_at,
                     updated_at,
                     role,
                     address,
                     lat,
                     lon
            FROM     cte_sub_admin`

	subAdminDetails := make([]models.UserDetails, 0)

	err := database.RmsDB.Select(&subAdminDetails, SQL, models.UserRoleSubAdmin, !filterCheck.IsSearched, filterCheck.SearchedName, filterCheck.Limit, filterCheck.Limit*filterCheck.Page)

	if err != nil {
		logrus.Printf("GetSubAdmin: unable to fetch subadmin details:%v", err)
		return totalUser, err
	}

	if len(subAdminDetails) == 0 {
		logrus.Printf("GetAllRestaurants:%v", err)
		return totalUser, err
	}

	totalUser.UserDetails = subAdminDetails
	totalUser.TotalCount = subAdminDetails[0].TotalCount
	return totalUser, nil
}

func GetUsers(filterCheck models.FiltersCheck) (models.TotalUser, error) {
	var totalUser models.TotalUser
	// language=SQL
	SQL := `WITH  cte_User AS(

            SELECT  
                     count(*) over () as total_count,
                     users.id as id,
                     name,
                     email,
                     phone_no,
                     password,
                     age,
                     gender,
                     users.created_at as created_at,
                     users.updated_at as updated_at,
                     role,
                     address,
                     lat,
                     lon
            FROM  users 
            JOIN roles ON users.id=roles.user_id
            JOIN   user_address ON roles.user_id=user_address.user_id
            WHERE users.archived_at IS NULL
            AND   user_address.archived_at IS NULL
            AND   role=$1
            AND    ($2 or name ilike '%'||$3||'%')
            ORDER BY name
            LIMIT $4 OFFSET $5
            )
            SELECT total_count,
                      id,
                     name,
                     email,
                     phone_no,
                     password,
                     age,
                     gender,
                     created_at,
                     updated_at,
                     role,
                     address,
                     lat,
                     lon
            FROM     cte_User`

	userDetails := make([]models.UserDetails, 0)

	err := database.RmsDB.Select(&userDetails, SQL, models.UserRoleUser, !filterCheck.IsSearched, filterCheck.SearchedName, filterCheck.Limit, filterCheck.Limit*filterCheck.Page)

	if err != nil {
		logrus.Printf("GetUsers: unable to fetch user details:%v", err)
		return totalUser, err
	}

	if len(userDetails) == 0 {
		logrus.Printf("GetUsers:%v", err)
		return totalUser, err
	}

	totalUser.UserDetails = userDetails
	totalUser.TotalCount = userDetails[0].TotalCount
	return totalUser, nil
}

// func GetUser(id int) ([]models.UserDetails, error) {
//	language=SQL
//	SQL := `SELECT
//                     id,
//                     name,
//                     email,
//                     phone_no,
//                     password,
//                     age,
//                     gender,
//                     created_at,
//                     updated_at,
//                     roles,
//                     address,
//                     lat,
//                     lon,
//                     archived_at
//            FROM  users JOIN roles
//            ON    users.id=roles.user_id
//            JOIN   user_address
//            ON     roles.user_id=user_address.user_id
//            WHERE users.archived_at IS NULL
//            AND   user_address.archived_at IS NULL
//            AND   roles=$1
//            AND   id=$2`
//
//	userDetails := make([]models.UserDetails, 0)
//
//	err := database.RmsDb.Select(&userDetails, SQL, "user", id)
//	if err != nil {
//		logrus.Printf("GetUsers: unable to fetch user details:%v", err)
//		return nil, err
//	}
//	return userDetails, nil
// }

func GetAllRestaurants(filtersCheck models.FiltersCheck) (models.TotalRestaurant, error) {
	// language=SQL
	var totalRestaurant models.TotalRestaurant
	SQL := `WITH cte_restaurant AS(
                    SELECT      count(*) over () as total_count,
                     restaurants.id as id,
                     name,
                     restaurants.created_at as created_at,
                     restaurants.updated_at as updated_at,
                     address,
                     lat,
                     lon
                    FROM  restaurants 
                    JOIN restaurant_address ON restaurants.id=restaurant_id
                    WHERE restaurants.archived_at IS NULL
                    AND   restaurant_address.archived_at IS NULL
                    AND   ($1 or name ilike '%' || $2 || '%')
                    ORDER BY name
                    LIMIT $3 OFFSET $4
            )
            
            SELECT  
                     total_count,
                     id,
                     name,
                     created_at,
                     updated_at,
                     address,
                     lat,
                     lon
            FROM     cte_restaurant`

	restaurantDetails := make([]models.Restaurants, 0)

	err := database.RmsDB.Select(&restaurantDetails, SQL, !filtersCheck.IsSearched, filtersCheck.SearchedName, filtersCheck.Limit, filtersCheck.Limit*filtersCheck.Page)

	if err != nil {
		logrus.Printf("GetRestaurants: unable to fetch all restaurants :%v", err)
		return totalRestaurant, err
	}

	if len(restaurantDetails) == 0 {
		logrus.Printf("GetAllRestaurants:%v", err)
		return totalRestaurant, err
	}

	totalRestaurant.Restaurants = restaurantDetails
	totalRestaurant.TotalCount = restaurantDetails[0].TotalCount
	return totalRestaurant, nil
}

func GetRestaurants(subAdminID int, filtersCheck models.FiltersCheck) (models.TotalRestaurant, error) {
	// language=SQL
	var totalRestaurant models.TotalRestaurant
	SQL := `WITH cte_restaurant AS(
                    SELECT      count(*) over () as total_count,
                     restaurants.id as id,
                     name,
                     restaurants.created_at as created_at,
                     restaurants.updated_at as updated_at,
                     address,
                     lat,
                     lon
                    FROM  restaurants 
                    JOIN restaurant_address ON  restaurants.id=restaurant_id
                    WHERE restaurants.archived_at IS NULL
                    AND   created_by=$1
                    AND   restaurant_address.archived_at IS NULL
                    AND   ($1 or name ilike '%' || $2 || '%')
                    ORDER BY name
                    LIMIT $3 OFFSET $4
            )
            
            SELECT  
                     total_count,
                     id,
                     name,
                     created_at,
                     updated_at,
                     address,
                     lat,
                     lon
            FROM     cte_restaurant`

	restaurantDetails := make([]models.Restaurants, 0)

	err := database.RmsDB.Select(&restaurantDetails, SQL, subAdminID, !filtersCheck.IsSearched, filtersCheck.SearchedName, filtersCheck.Limit, filtersCheck.Limit*filtersCheck.Page)

	if err != nil {
		logrus.Printf("GetRestaurants: unable to fetch restaurants details:%v", err)
		return totalRestaurant, err
	}

	if len(restaurantDetails) == 0 {
		logrus.Printf("GetAllRestaurants:%v", err)
		return totalRestaurant, err
	}

	totalRestaurant.Restaurants = restaurantDetails
	totalRestaurant.TotalCount = restaurantDetails[0].TotalCount
	return totalRestaurant, nil
}

func FetchUserLocation(userID int) (models.UserLocation, error) {
	// language=SQL
	SQL := `SELECT lat,
                   lon
            FROM   user_address 
            WHERE  user_address.id=$1`

	var userLocation models.UserLocation

	err := database.RmsDB.Get(&userLocation, SQL, userID)

	if err != nil {
		logrus.Printf("FetchUserLocation: unable to fetch user location:%v", err)
		return userLocation, err
	}

	return userLocation, nil
}

func FetchRestaurantLocation(restaurantID int) (models.UserLocation, error) {
	// language=SQL
	SQL := `SELECT lat,
                   lon
            FROM   restaurant_address 
            WHERE  restaurant_id=$1
            `

	var restaurantLocation models.UserLocation

	err := database.RmsDB.Get(&restaurantLocation, SQL, restaurantID)

	if err != nil {
		logrus.Printf("FetchRestaurantLocation: unable to fetch restaurant location:%v", err)
		return restaurantLocation, err
	}

	return restaurantLocation, nil
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func GetDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return (2 * r * math.Asin(math.Sqrt(h))) / 1000
}

func GetDishes(restaurantID, userID int, filtersCheck models.FiltersCheck) (models.CountDishes, error) {
	// language=SQL
	var totalDishes models.CountDishes
	SQL := ` WITH cte_count AS (
                                  SELECT 
         
                                          COUNT(*) over() as total_count,
                                          dishes.id as id,
                                          name,
                                          price,
                                          created_by,
                                          restaurant_id,
                                          url
                                  FROM
                                          dishes LEFT JOIN dish_per_image ON dishes.id = dish_per_image.dish_id
                                                 LEFT JOIN images         ON dish_per_image.image_id = images.id
                                  WHERE              dishes.archived_at IS NULL 
                                  AND        restaurant_id = $1
                                  AND        created_by = $2
                                  AND        ($3 or name ilike '%'|| $4 ||'%')
                                  ORDER BY    name
                                  LIMIT      $5
                                  OFFSET     $6 

                                 )
            SELECT 
                       total_count, 
                       id,
                       name,
                       price,
                       created_by,
                       restaurant_id,
                       COALESCE (url, '') as url
            FROM
                       cte_count`

	dishDetails := make([]models.Dishes, 0)

	err := database.RmsDB.Select(&dishDetails, SQL, restaurantID, userID, !filtersCheck.IsSearched, filtersCheck.SearchedName, filtersCheck.Limit, filtersCheck.Limit*filtersCheck.Page)

	if err != nil {
		logrus.Printf("GetDishes: unable to fetch dish details: %v", err)
		return totalDishes, err
	}

	totalDishes.Dishes = dishDetails
	if len(dishDetails) == 0 {
		return totalDishes, nil
	}

	totalDishes.TotalCount = dishDetails[0].TotalCount
	return totalDishes, nil
}

// func SubAdminGetDishes(id int) ([]models.Dishes, error) {
//	language=SQL
//	SQL := `SELECT
//                       id,
//                       name,
//                       price,
//                       created_by,
//                       restaurant_id
//            FROM
//                       dishes
//            JOIN       dish_per_restaurant
//            ON         dishes.id=dish_id
//            WHERE      archived_at IS NULL
//            AND        created_by=$1`
//
//	dishDetails := make([]models.Dishes, 0)
//
//	err := database.RmsDb.Select(&dishDetails, SQL, id)
//	if err != nil {
//		log.Printf("SubAdminGetDishes: unable to fetch dish details:%v", err)
//		return nil, err
//	}
//	return dishDetails, nil
// }

// func UserGetDishes() ([]models.Dishes, error) {
//	language=SQL
//	SQL := `SELECT
//                       id,
//                       name,
//                       price,
//                       created_by,
//                       restaurant_id
//            FROM
//                       dishes
//            JOIN       dish_per_restaurant
//            ON         dishes.id=dish_id
//            WHERE      archived_at IS NULL `
//
//	dishDetails := make([]models.Dishes, 0)
//
//	err := database.RmsDb.Select(&dishDetails, SQL)
//	if err != nil {
//		log.Printf("UserGetDishes: unable to fetch dish details:%v", err)
//		return nil, err
//	}
//	return dishDetails, nil
// }

// func UpdateAdmin(id int, adminDetails *models.UserDetails) error {
//	language=SQL
//	SQL := `UPDATE users
//           SET
//                  name=$1,
//                  email=$2,
//                  phone_no=$3,
//                  password=$4,
//                  age=$5,
//                  gender=$6
//           WHERE
//                  users.id=$7`
//
//	_, err := database.RmsDb.Exec(SQL, adminDetails.Name, adminDetails.Email, adminDetails.PhoneNo, adminDetails.PhoneNo, adminDetails.Password, adminDetails.Age, adminDetails.Gender, id)
//	if err != nil {
//		logrus.Printf("UpdateAdmin: cannot update admin:%v", err)
//		return err
//	}
//	SQL = `UPDATE user_address
//           SET
//                  address=$1,
//                  lat=$2,
//                  lon=$3
//           WHERE
//                  user_id=$4`
//
//	_, err = database.RmsDb.Exec(SQL, adminDetails.Address, adminDetails.Lat, adminDetails.Lon, id)
//	if err != nil {
//		logrus.Printf("UpdateAdmin: cannot update admin:%v", err)
//		return err
//	}
//	return nil
// }

func UpdateSubAdmin(id int, subAdminDetails *models.UserDetails) error {
	//  language=SQL
	SQL := `UPDATE users  
           SET
                  name=$1,
                  email=$2,
                  phone_no=$3,
                  password=$4,
                  age=$5,
                  gender=$6
           WHERE  
                  users.id=$7`

	_, err := database.RmsDB.Exec(SQL, subAdminDetails.Name, subAdminDetails.Email, subAdminDetails.PhoneNo, subAdminDetails.PhoneNo, subAdminDetails.Password, subAdminDetails.Age, subAdminDetails.Gender, id)
	if err != nil {
		logrus.Printf("UpdateSubAdmin: cannot update subAdmin:%v", err)
		return err
	}

	SQL = `UPDATE user_address 
           SET
                  address=$1,
                  lat=$2,
                  lon=$3
           WHERE  
                  user_id=$4`

	_, err = database.RmsDB.Exec(SQL, subAdminDetails.Address, subAdminDetails.Lat, subAdminDetails.Lon, id)
	if err != nil {
		logrus.Printf("UpdateSubAdmin: cannot update subAdmin:%v", err)
		return err
	}
	return nil
}

func UpdateUser(id int, userDetails *models.UserDetails) error {
	// language=SQL
	SQL := `UPDATE users 
           SET
                  name=$1,
                  email=$2,
                  phone_no=$3,
                  password=$4,
                  age=$5,
                  gender=$6
           WHERE  
                  users.id=$7`

	_, err := database.RmsDB.Exec(SQL, userDetails.Name, userDetails.Email, userDetails.PhoneNo, userDetails.PhoneNo, userDetails.Password, userDetails.Age, userDetails.Gender, id)

	if err != nil {
		logrus.Printf("UpdateUser: cannot update user details:%v", err)
		return err
	}

	SQL = `UPDATE user_address 
           SET
                  address=$1,
                  lat=$2,
                  lon=$3
           WHERE  
                  user_id=$4`

	_, err = database.RmsDB.Exec(SQL, userDetails.Address, userDetails.Lat, userDetails.Lon, id)

	if err != nil {
		logrus.Printf("UpdateUser: cannot update user address:%v", err)
		return err
	}

	return nil
}

// func UpdateRestaurants(id int, restaurantDetails models.Restaurants) error {
//	var restaurantID int
//	language=SQL
//	SQL := `UPDATE  restaurants
//            SET
//                    name=$1
//
//            WHERE
//                    created_by=$2`
//
//	_, err := database.RmsDb.Exec(SQL, restaurantDetails.Name, id)
//	if err != nil {
//		logrus.Printf("UpdateRestaurants: cannot update restaurant:%v", err)
//		return err
//	}
//
//	SQL = `SELECT id
//            FROM   restaurants
//            WHERE  created_by=$1`
//
//	err = database.RmsDb.Get(&restaurantID, SQL, id)
//	if err != nil {
//		logrus.Printf("UpdateRestaurants: cannot get restaurantID:%v", err)
//		return err
//	}
//
//	SQL = `UPDATE restaurant_address
//           SET
//                  address=$1,
//                  lat=$2,
//                  lon=$3
//           WHERE
//                  restaurant_id=$4`
//
//	_, err = database.RmsDb.Exec(SQL, restaurantDetails.Address, restaurantDetails.Lat, restaurantDetails.Lon, restaurantID)
//	if err != nil {
//		logrus.Printf("UpdateSubAdmin: cannot update subAdmin:%v", err)
//		return err
//	}
//	return nil
// }

func UpdateDish(id int, dishesDetails models.Dishes) error {
	//  language=SQL
	SQL := `UPDATE  dishes   
            SET    
                    name=$1,
                    price=$2
            WHERE   dishes.id=$3`

	_, err := database.RmsDB.Exec(SQL, dishesDetails.Name, dishesDetails.Price, id)

	if err != nil {
		logrus.Printf("UpdateDishes: cannot update dish:%v", err)
		return err
	}

	return nil
}

func DeleteAdmin(id int) error {
	// language=SQL
	SQL := `UPDATE users
            SET archived_at=now()
            WHERE id=$1`

	_, err := database.RmsDB.Exec(SQL, id)
	if err != nil {
		logrus.Printf("DeleteAdmin: cannot delete admin:%v", err)
		return err
	}
	return nil
}

func DeleteSubAdmin(id int) error {
	// language=SQL
	SQL := `UPDATE users
            SET    archived_at=now()
            WHERE id=$1`

	_, err := database.RmsDB.Exec(SQL, id)
	if err != nil {
		logrus.Printf("DeleteSubAdmin: cannot delete subAdmin:%v", err)
		return err
	}
	return nil
}

func DeleteUser(id int) error {
	// language=SQL
	SQL := `UPDATE users
            SET    archived_at=now()
            WHERE id=$1`

	_, err := database.RmsDB.Exec(SQL, id)

	if err != nil {
		logrus.Printf("DeleteUser: cannot delete user:%v", err)
		return err
	}

	return nil
}

func DeleteRestaurants(id int) error {
	//  language=SQL
	SQL := `UPDATE restaurants
            SET    archived_at=now()
            WHERE created_by=$1`
	_, err := database.RmsDB.Exec(SQL, id)
	if err != nil {
		logrus.Printf("DeleteRestaurants: cannot delete restaurant:%v", err)
		return err
	}
	return nil
}

func DeleteDish(id int) error {
	// language=SQL
	SQL := `UPDATE dishes
            SET    archived_at=now()
            WHERE dishes.id=$1`
	_, err := database.RmsDB.Exec(SQL, id)

	if err != nil {
		logrus.Printf("DeleteDishes: cannot delete dishes:%v", err)
		return err
	}

	return nil
}
