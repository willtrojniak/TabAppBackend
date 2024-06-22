package types

type UserStore interface {

}

type UserCreate struct {
  Name string `json:"name" db:"name"`
  Email string `json:"email" db:"email"`
}

type User struct {
  Id uint32 `json:"id" db:"id"`
  *UserCreate
}
