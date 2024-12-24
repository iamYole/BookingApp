package dbrepo

import (
	"context"
	"time"

	"github.com/iamYole/BookingApp/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newId int
	stmt := `INSERT INTO reservations
					 (first_name, last_name, email, phone, start_date, end_date,
					 room_id, created_at, updated_at)
			VALUES
					($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestrictions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO room_restrictions
						(start_date, end_date, room_id, reservation_id, created_at, updated_at,
						restriction_id)
					VALUES ($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		time.Now(),
		time.Now(),
		res.RestrictionID,
	)

	if err != nil {
		return err
	}

	return nil
}
func (m *postgresDBRepo) SearchAvailablityByDateByRoomID(start_date, end_date time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numrows int

	stmt := `SELECT count(id) FROM
				room_restrictions 
				WHERE roomID = $1 AND 
				($2 < end_date and $3 > start_date);`

	row := m.DB.QueryRowContext(ctx, stmt, roomId, start_date, end_date)
	err := row.Scan(&numrows)
	if err != nil {
		return false, err
	}

	if numrows == 0 {
		return true, nil
	}

	return false, nil
}

func (m *postgresDBRepo) SearchAvailablityForAllRooms(start_date, end_date time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var availableRooms []models.Room

	stmt := `select r.id ,r.room_name 
			from rooms r 
			where r.id not in 
					(select rr.room_id from room_restrictions rr
							where $1 < rr.end_date and $2 > rr.start_date)`
	rows, err := m.DB.QueryContext(ctx, stmt, start_date, end_date)
	if err != nil {
		return availableRooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID, &room.RoomName,
		)
		if err != nil {
			return availableRooms, err
		}

		availableRooms = append(availableRooms, room)
	}
	if err = rows.Err(); err != nil {
		return availableRooms, err
	}

	return availableRooms, nil

}
