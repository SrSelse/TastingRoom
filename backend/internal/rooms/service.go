package rooms

import (
	"context"
	"log/slog"

	"skafteresort.se/beers/internal/providers"
)

type RoomService struct {
	roomRepo   *RoomRepo
	logger     *slog.Logger
	centrifugo *providers.CentrifugoProvider
}

func NewRoomService(rr *RoomRepo, logger *slog.Logger, gp *providers.CentrifugoProvider) *RoomService {
	ts := RoomService{
		roomRepo:   rr,
		logger:     logger,
		centrifugo: gp,
	}
	return &ts
}

func (s *RoomService) GetUsersInRoom(ctx context.Context, roomId int) ([]RelatedUser, error) {
	return s.roomRepo.getUsersInRoom(ctx, roomId)
}

func (s *RoomService) GetBeersInRoom(ctx context.Context, roomId int) ([]RelatedBeer, error) {
	return s.roomRepo.getBeersInRoom(ctx, roomId)
}

func (s *RoomService) GetRoomById(ctx context.Context, roomId int) (*Room, error) {
	return s.roomRepo.getRoomById(ctx, roomId)
}

func (s *RoomService) GetRoomByCode(ctx context.Context, code string) (*Room, error) {
	return s.roomRepo.getRoomByCode(ctx, code)
}

func (s *RoomService) GetRoomsByUserId(ctx context.Context, userId int) ([]Room, error) {
	return s.roomRepo.getRoomsByUserId(ctx, userId)
}

func (s *RoomService) AddUserToRoom(ctx context.Context, roomId int, userId int, admin bool) error {
	return s.roomRepo.addUserToRoom(ctx, userId, roomId, admin)
}

func (s *RoomService) RemoveUserFromRoom(ctx context.Context, roomId int, userId int) error {
	return s.roomRepo.removeUserFromRoom(ctx, userId, roomId)
}

func (s *RoomService) UpdateIsAdmin(ctx context.Context, roomId int, targetUserId int, isAdmin bool) error {
	return s.roomRepo.updateIsAdmin(ctx, roomId, targetUserId, isAdmin)
}

func (s *RoomService) CreateNewRoom(ctx context.Context, userId int, room Room) (int, error) {
	id, err := s.roomRepo.createNewRoom(ctx, room)
	if err != nil {
		return id, err
	}
	err = s.roomRepo.addUserToRoom(ctx, userId, id, true)
	return id, err
}

func (s *RoomService) CheckIfUserInRoom(ctx context.Context, roomId int, userId int) (bool, error) {
	return s.roomRepo.checkIfUserInRoom(ctx, roomId, userId)
}

func (s *RoomService) CheckIfUserIsAdminInRoom(ctx context.Context, roomId int, userId int) (bool, error) {
	return s.roomRepo.checkIfUserIsAdminInRoom(ctx, roomId, userId)
}

func (s *RoomService) CheckIfOtherAdminInRoom(ctx context.Context, roomId int, userId int) (bool, error) {
	return s.roomRepo.checkIfOtherAdminInRoom(ctx, roomId, userId)
}

func (s *RoomService) CheckIfBeerInRoom(ctx context.Context, roomId int, beerId int) (bool, error) {
	return s.roomRepo.checkIfBeerInRoom(ctx, roomId, beerId)
}

func (s *RoomService) UpdateRoom(ctx context.Context, room Room) error {
	return s.roomRepo.updateRoom(ctx, room)
}
