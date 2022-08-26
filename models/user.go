package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserProfile string

const (
	UserRoleAdmin    UserProfile = "admin"
	UserRoleSubAdmin UserProfile = "sub_admin"
	UserRoleUser     UserProfile = "user"
)

type FiltersCheck struct {
	IsSearched   bool
	SearchedName string
	Limit        int
	Page         int
}

type UserLocation struct {
	Lat float64 `json:"lat" db:"lat"`
	Lon float64 `json:"lon" db:"lon"`
}

type UserDetails struct {
	TotalCount int       `json:"-" db:"total_count"`
	ID         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Email      string    `json:"email" db:"email"`
	Password   string    `json:"password" db:"password"`
	PhoneNo    int       `json:"phoneNo" db:"phone_no"`
	Age        int       `json:"age" db:"age"`
	Gender     string    `json:"gender" db:"gender"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
	Address    string    `json:"address" db:"address"`
	Lat        float64   `json:"lat" db:"lat"`
	Lon        float64   `json:"lon" db:"lon"`
	Role       string    `json:"role"  db:"role"`
}
type TotalUser struct {
	UserDetails []UserDetails
	TotalCount  int `json:"totalCount" db:"total_count"`
}

type Users struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	PhoneNo  int     `json:"phoneNo"`
	Age      int     `json:"age"`
	Gender   string  `json:"gender"`
	UserID   int     `json:"userId"`
	Address  string  `json:"address"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

type Restaurants struct {
	TotalCount   int       `json:"-" db:"total_count"`
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Address      string    `json:"address" db:"address"`
	Lat          float64   `json:"lat" db:"lat"`
	Lon          float64   `json:"lon" db:"lon"`
	RestaurantID int       `json:"restaurantId" db:"restaurant_id"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}
type TotalRestaurant struct {
	Restaurants []Restaurants
	TotalCount  int `json:"totalCount" db:"total_count"`
}

type Dishes struct {
	TotalCount   int    `json:"-" db:"total_count"`
	ID           int    `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Price        int    `json:"price" db:"price"`
	DishID       int    `json:"dishId" db:"dish_id"`
	RestaurantID int    `json:"restaurantId" db:"restaurant_id"`
	CreatedBy    int    `json:"createdBy" db:"created_by"`
}

type CountDishes struct {
	Dishes     []Dishes `json:"dishes"`
	TotalCount int      `json:"totalCount"`
}

type UserCredentials struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UsersLoginDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Claims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}
