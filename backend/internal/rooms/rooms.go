package rooms

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/google/uuid"
)

type RoomRepo struct {
	db *sql.DB
}

type Room struct {
	Id          int     `db:"id" json:"id"`
	Name        string  `db:"name" json:"name"`
	Code        *string `db:"code" json:"code"`
	CreatedAt   string  `db:"created_at" json:"createdAt"`
	Description string  `db:"description" json:"description"`
	PlannedDate string  `db:"planned_date" json:"plannedDate"`
	Members     int     `db:"members" json:"members"`
}

type RelatedBeer struct {
	Id         int      `json:"id"`
	Name       string   `json:"name"`
	Style      *string  `json:"style"`
	PictureUrl *string  `json:"pictureUrl"`
	Average    *float64 `db:"average" json:"average"`
	Published  bool     `db:"published" json:"published"`
}

type RelatedUser struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"isAdmin"`
}

func NewRoomRepo(db *sql.DB) *RoomRepo {
	return &RoomRepo{db}
}

func (rr *RoomRepo) getRoomsByUserId(ctx context.Context, userId int) ([]Room, error) {
	rows, err := rr.db.QueryContext(ctx, `
    SELECT
      rooms.id,
      rooms.name,
      rooms.created_at,
      rooms.description,
      rooms.planned_date,
      (
      	SELECT count(*)
      	FROM user_room
      	WHERE user_room.room_id = rooms.id
      ) as members
    FROM rooms
    JOIN user_room ON user_room.room_id = rooms.id
    WHERE user_room.user_id = ?
    GROUP BY rooms.id
  `, userId)
	if err != nil {
		return []Room{}, err
	}
	rooms := []Room{}

	for rows.Next() {
		var room Room
		err = rows.Scan(
			&room.Id,
			&room.Name,
			&room.CreatedAt,
			&room.Description,
			&room.PlannedDate,
			&room.Members,
		)
		if err != nil {
			return []Room{}, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (rr *RoomRepo) getRoomByCode(ctx context.Context, code string) (*Room, error) {
	row := rr.db.QueryRowContext(ctx, `
    SELECT
      rooms.id
    FROM rooms
    WHERE code = ?
  `,
		code,
	)
	var room Room
	err := row.Scan(
		&room.Id,
	)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (rr *RoomRepo) getRoomById(ctx context.Context, roomId int) (*Room, error) {
	row := rr.db.QueryRowContext(ctx, `
    SELECT
      rooms.id,
      rooms.name,
      rooms.code,
      rooms.description,
      rooms.planned_date,
      (
      	SELECT count(*)
      	FROM user_room
      	WHERE user_room.room_id = rooms.id
      ) as members
    FROM rooms
    WHERE id = ?
`, roomId)
	var room Room
	err := row.Scan(
		&room.Id,
		&room.Name,
		&room.Code,
		&room.Description,
		&room.PlannedDate,
		&room.Members,
	)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (rr *RoomRepo) getUsersInRoom(ctx context.Context, roomId int) ([]RelatedUser, error) {
	rows, err := rr.db.QueryContext(ctx, `
    SELECT users.id, users.username, user_room.is_admin
    FROM users
    JOIN user_room ON user_room.user_id = users.id
    WHERE user_room.room_id = ?
  `, roomId)
	if err != nil {
		return []RelatedUser{}, err
	}

	relatedUsers := []RelatedUser{}
	for rows.Next() {
		var u RelatedUser
		err := rows.Scan(&u.Id, &u.Name, &u.IsAdmin)
		if err != nil {
			return []RelatedUser{}, err
		}
		relatedUsers = append(relatedUsers, u)
	}
	return relatedUsers, nil
}

func (rr *RoomRepo) getBeersInRoom(ctx context.Context, roomId int) ([]RelatedBeer, error) {
	rooms, _ := rr.getRoomsByUserId(ctx, ctx.Value("userID").(int))
	if !slices.ContainsFunc(rooms, func(r Room) bool {
		return r.Id == roomId
	}) {
		return []RelatedBeer{}, fmt.Errorf("Unauthorized")
	}
	rows, err := rr.db.QueryContext(ctx, `
    SELECT
      beers_votes.id,
      beers_votes.name,
      beers_votes.style,
      beers_votes.pictureurl,
      beers_votes.average,
      beers_votes.published
    FROM beers_votes
    WHERE beers_votes.room_id = ?
    ORDER BY COALESCE(beers_votes.average, beers_votes.id) DESC
  `, roomId)
	if err != nil {
		return []RelatedBeer{}, err
	}

	beers := []RelatedBeer{}
	for rows.Next() {
		var beer RelatedBeer
		err := rows.Scan(
			&beer.Id,
			&beer.Name,
			&beer.Style,
			&beer.PictureUrl,
			&beer.Average,
			&beer.Published,
		)
		if err != nil {
			return []RelatedBeer{}, err
		}
		beers = append(beers, beer)
	}

	return beers, nil
}

func (rr *RoomRepo) createNewRoom(ctx context.Context, room Room) (int, error) {
	code := uuid.NewString()
	res, err := rr.db.ExecContext(ctx, `
    INSERT INTO rooms (name, code, planned_date, description)
    VALUES (?, ?, ?, ?)
  `,
		room.Name,
		code,
		room.PlannedDate,
		room.Description,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (rr *RoomRepo) addUserToRoom(ctx context.Context, userId int, roomId int, admin bool) error {
	_, err := rr.db.ExecContext(ctx, `
    INSERT INTO user_room (room_id, user_id, is_admin)
    VALUES (?, ?, ?)
    `,
		roomId,
		userId,
		admin,
	)
	return err
}

func (rr *RoomRepo) removeUserFromRoom(ctx context.Context, userId int, roomId int) error {
	_, err := rr.db.ExecContext(ctx, `
      DELETE FROM user_room
      WHERE room_id = ?
      AND user_id = ?
    `,
		roomId,
		userId,
	)
	return err
}

func (rr *RoomRepo) updateIsAdmin(ctx context.Context, roomId int, userId int, admin bool) error {
	_, err := rr.db.ExecContext(ctx, `
      UPDATE user_room SET is_admin = ?
      WHERE user_id = ? AND room_id = ?
    `,
		admin,
		userId,
		roomId,
	)
	return err
}

func (rr *RoomRepo) checkIfUserInRoom(ctx context.Context, roomId int, userId int) (bool, error) {
	row := rr.db.QueryRowContext(ctx, `
    SELECT EXISTS (
      SELECT room_id
      FROM user_room
      WHERE room_id = ?
      AND user_id = ?
    )
    `,
		roomId,
		userId,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

func (rr *RoomRepo) checkIfUserIsAdminInRoom(ctx context.Context, roomId int, userId int) (bool, error) {
	row := rr.db.QueryRowContext(ctx, `
      SELECT is_admin
      FROM user_room
      WHERE room_id = ?
      AND user_id = ?
    `,
		roomId,
		userId,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

func (rr *RoomRepo) checkIfOtherAdminInRoom(ctx context.Context, roomId int, userId int) (bool, error) {
	row := rr.db.QueryRowContext(ctx, `
      SELECT count(*)
      FROM user_room
      WHERE room_id = ?
      AND user_id != ?
      AND is_admin = 1
    `,
		roomId,
		userId,
	)
	var admins int
	err := row.Scan(&admins)
	return admins > 0, err
}

func (rr *RoomRepo) checkIfBeerInRoom(ctx context.Context, roomId int, beerId int) (bool, error) {
	row := rr.db.QueryRowContext(ctx, `
    SELECT EXISTS (
      SELECT room_id
      FROM beers
      WHERE room_id = ?
      AND id = ?
    )
    `,
		roomId,
		beerId,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

func (rr *RoomRepo) updateRoom(ctx context.Context, room Room) error {
	_, err := rr.db.ExecContext(ctx, `
    UPDATE rooms
    SET name = ?, description = ?, planned_date = ?
    WHERE id = ?
    `,
		room.Name,
		room.Description,
		room.PlannedDate,
		room.Id,
	)
	return err
}
