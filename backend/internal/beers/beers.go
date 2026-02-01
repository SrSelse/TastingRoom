package beers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type BeerRepo struct {
	db *sql.DB
}

type Beer struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Style      *string `json:"style"`
	PictureUrl *string `json:"pictureUrl"`
	RoomId     int     `json:"roomId"`
	Published  bool    `json:"published"`
}

type Vote struct {
	Id       int     `json:"id"`
	UserId   int     `json:"userId"`
	UserName string  `json:"name"`
	Value    int     `json:"rating"`
	BeerId   int     `json:"beerId"`
	Note     *string `json:"note"`
}

func NewBeerRepo(db *sql.DB) *BeerRepo {
	return &BeerRepo{db}
}

func (br *BeerRepo) getVotesByBeerId(ctx context.Context, beerId int, roomId int) ([]Vote, error) {

	rows, err := br.db.QueryContext(ctx,
		`
      SELECT
        votes.id,
        user_room.user_id,
        IF(users.name != '', users.name, users.username) as userName,
        votes.points,
        votes.note
      FROM user_room
      JOIN users ON users.id = user_room.user_id
      LEFT OUTER JOIN votes ON votes.user_id = users.id
      WHERE votes.beer_id = ?
    `,
		beerId,
	)
	if err != nil {
		return nil, err
	}
	voteMap := map[int]Vote{}
	for rows.Next() {
		var v Vote
		err := rows.Scan(
			&v.Id,
			&v.UserId,
			&v.UserName,
			&v.Value,
			&v.Note,
		)
		if err != nil {
			return nil, err
		}
		voteMap[v.UserId] = v
	}
	rows, err = br.db.QueryContext(ctx,
		`
      SELECT
        user_room.user_id,
        IF(users.name != '', users.name, users.username) as userName
      FROM user_room
      JOIN users ON users.id = user_room.user_id
      WHERE user_room.room_id = ?
    `,
		roomId,
	)
	if err != nil {
		return nil, err
	}

	votes := []Vote{}
	for rows.Next() {
		var u struct {
			Id   int
			Name string
		}
		err := rows.Scan(&u.Id, &u.Name)
		if err != nil {
			return nil, err
		}
		if vote, ok := voteMap[u.Id]; ok {
			votes = append(votes, vote)
		} else {
			votes = append(votes, Vote{UserId: u.Id, UserName: u.Name})
		}
	}
	return votes, nil
}

func (br *BeerRepo) getBeerById(ctx context.Context, beerId int) (*Beer, error) {
	row := br.db.QueryRowContext(ctx,
		`
      SELECT id, name, style, published, pictureurl
      FROM beers
      WHERE id = ?
    `,
		beerId,
	)
	var beer Beer
	err := row.Scan(
		&beer.Id,
		&beer.Name,
		&beer.Style,
		&beer.Published,
		&beer.PictureUrl,
	)
	return &beer, err
}

func (br *BeerRepo) addVoteOnBeerId(ctx context.Context, vote Vote) error {
	note := ""
	if vote.Note != nil {
		note = *vote.Note
	}
	fmt.Printf("%v %s\n", vote, note)
	_, err := br.db.ExecContext(ctx, `
      INSERT INTO votes (beer_id, user_id, points, note)
      VALUES (?, ?, ?, ?)
    `,
		vote.BeerId,
		vote.UserId,
		vote.Value,
		note,
	)
	return err
}

func (br *BeerRepo) updateVoteOnBeerId(ctx context.Context, vote Vote) error {
	note := ""
	if vote.Note != nil {
		note = *vote.Note
	}
	_, err := br.db.ExecContext(ctx, `
      UPDATE votes SET points = ?, note = ?
      WHERE id = ?

    `,
		vote.Value,
		note,
		vote.Id,
	)
	return err
}

func (br *BeerRepo) addNewBeer(ctx context.Context, beer Beer) error {
	_, err := br.db.ExecContext(ctx, `
    INSERT INTO beers (name, style, pictureurl, room_id)
    VALUES(?, ?, ?, ?)
    `,
		beer.Name,
		beer.Style,
		beer.PictureUrl,
		beer.RoomId,
	)
	return err
}

func (br *BeerRepo) getBeersByUserVotes(ctx context.Context, userId int) ([]Beer, error) {
	return []Beer{}, nil
}

func (br *BeerRepo) getRandomBeerInRoom(ctx context.Context, roomId int) (*Beer, error) {
	row := br.db.QueryRowContext(ctx,
		`
      SELECT id, name, style, pictureurl
      FROM beers
      WHERE room_id = ?
      AND published = 0
      ORDER BY RAND()
      LIMIT 1
    `,
		roomId,
	)

	var beer Beer
	err := row.Scan(
		&beer.Id,
		&beer.Name,
		&beer.Style,
		&beer.PictureUrl,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &beer, nil
}

func (br *BeerRepo) getFirstBeerInRoom(ctx context.Context, roomId int) (*Beer, error) {
	row := br.db.QueryRowContext(ctx,
		`
      SELECT id, name, style, pictureurl
      FROM beers
      WHERE room_id = ?
	  ORDER BY id ASC
      LIMIT 1
    `,
		roomId,
	)

	var beer Beer
	err := row.Scan(
		&beer.Id,
		&beer.Name,
		&beer.Style,
		&beer.PictureUrl,
	)
	fmt.Printf("%v\n", beer)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &beer, nil
}

func (br *BeerRepo) getNextBeerInRoom(ctx context.Context, roomId int, oldBeerId int) (*Beer, error) {
	row := br.db.QueryRowContext(ctx,
		`
      SELECT id, name, style, pictureurl
      FROM beers
      WHERE room_id = ?
	  AND id > ?
	  ORDER BY id ASC
      LIMIT 1
    `,
		roomId,
		oldBeerId,
	)

	var beer Beer
	err := row.Scan(
		&beer.Id,
		&beer.Name,
		&beer.Style,
		&beer.PictureUrl,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			b, err := br.getFirstBeerInRoom(ctx, roomId)
			if err != nil {
				return nil, err
			}
			beer = *b
		} else {
			return nil, err
		}
	}
	return &beer, nil
}

func (br *BeerRepo) publishRatingsForBeer(ctx context.Context, beerId int, roomId int) error {
	_, err := br.db.ExecContext(ctx,
		`
      UPDATE beers SET published = 1
      WHERE beers.id = ? AND room_id = ?
    `,
		beerId,
		roomId,
	)
	return err
}

func (br *BeerRepo) unpublishRatingsForBeer(ctx context.Context, beerId int, roomId int) error {
	_, err := br.db.ExecContext(ctx,
		`
      UPDATE beers SET published = 0
      WHERE beers.id = ? AND room_id = ?
    `,
		beerId,
		roomId,
	)
	return err
}

func (br *BeerRepo) getMyRatingOnBeer(ctx context.Context, beerId int, userId int) (*Vote, error) {
	row := br.db.QueryRowContext(ctx,
		`
      SELECT id, points, note
      FROM votes
      WHERE beer_id = ?
      AND user_id = ?
    `,
		beerId,
		userId,
	)
	var vote Vote
	var note *string
	err := row.Scan(
		&vote.Id,
		&vote.Value,
		&note,
	)
	if note != nil {
		vote.Note = note
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &vote, err
}

func (br *BeerRepo) updateBeer(ctx context.Context, beer Beer, roomId int) error {
	_, err := br.db.ExecContext(ctx, `
      UPDATE beers
      SET name = ?, style = ?, pictureurl = ?
      WHERE id = ? AND room_id = ?
    `,
		beer.Name,
		beer.Style,
		beer.PictureUrl,
		beer.Id,
		roomId,
	)
	return err
}
