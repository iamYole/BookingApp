package repository

import (
	"time"

	"github.com/iamYole/BookingApp/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestrictions) error
	SearchAvailablityByDateByRoomID(start_date, end_date time.Time, roomId int) (bool, error)
	SearchAvailablityForAllRooms(start_date, end_date time.Time) ([]models.Room, error)
}
