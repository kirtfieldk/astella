package userservice

import (
	"database/sql"
	"log"

	"github.com/kirtfieldk/astella/src/structures"
)

func MapUserRows(rows *sql.Rows) []structures.User {
	var users []structures.User
	for rows.Next() {
		var usr structures.User
		if err := rows.Scan(&usr.UUID, &usr.Username, &usr.Created, &usr.Ig, &usr.Twitter, &usr.Toktok, &usr.AvatarUrl,
			&usr.ImgOne, &usr.ImgTwo, &usr.ImgThree); err != nil {
			log.Println("Issue mapping DB row for user")
		}

		users = append(users, usr)
	}
	return users
}
