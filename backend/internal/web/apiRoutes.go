package web

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"skafteresort.se/beers/internal/auth"
	"skafteresort.se/beers/internal/beers"
	"skafteresort.se/beers/internal/rooms"
)

func addApiRoutes(
	logger *slog.Logger,
	userService *auth.UserService,
	roomService *rooms.RoomService,
	beerService *beers.BeerService,
) *http.ServeMux {

	mux := http.NewServeMux()

	mux.Handle(
		"/api/rooms",
		handleRooms(roomService, logger),
	)

	mux.Handle(
		"/api/room/join",
		handleJoinRoom(roomService, logger),
	)

	mux.Handle(
		"/api/room/create",
		handleCreateRoom(roomService, logger),
	)

	mux.Handle(
		"/api/room/{room}",
		handleRoom(roomService, logger),
	)

	mux.Handle(
		"/api/room/{room}/edit",
		handleEditRoom(roomService, logger),
	)

	mux.Handle(
		"/api/room/{room}/leave",
		handleLeaveRoom(roomService, logger),
	)

	mux.Handle(
		"/api/room/{room}/is-admin",
		handleCheckIfUserIsAdminInRoom(roomService, logger),
	)

	mux.Handle(
		"/api/room/{room}/beers/{beer}/my-rating",
		handleGetMyRatingForBeer(roomService, beerService, logger),
	)

	mux.Handle(
		"/api/room/{room}/beers/{beer}/ratings",
		handleGetRatingsForBeer(roomService, beerService, logger),
	)

	mux.Handle(
		"/api/room/{room}/beers/{beer}/edit",
		handleEditBeer(roomService, beerService, logger),
	)

	mux.Handle(
		"/api/room/{room}/users",
		handleUsersInRoom(roomService, logger),
	)

	mux.Handle(
		"/api/room/{room}/users/{user}/admin",
		handleUpdateIsAdminForUser(roomService, logger),
	)

	mux.Handle(
		"/api/room/{room}/users/{user}/remove",
		handleRemoveUserFromRoom(roomService, logger),
	)

	mux.Handle(
		"/api/room/{room}/beers",
		handleBeersInRoom(roomService, logger),
	)

	mux.Handle(
		"/api/room/{room}/beers/new",
		handleAddBeer(beerService, logger),
	)

	mux.Handle(
		"/api/room/{room}/beers/random",
		handleGetRandomBeer(beerService, roomService, logger),
	)

	mux.Handle(
		"/api/room/{room}/beers/next",
		handleGetNextBeer(beerService, roomService, logger),
	)

	mux.Handle(
		"/api/room/{room}/beers/{beer}",
		handleGetSingleBeer(beerService, roomService, logger),
	)

	mux.Handle(
		"/api/room/{room}/beers/{beer}/publish",
		handlePublishRatingsForBeer(roomService, beerService, logger),
	)

	mux.Handle(
		"/api/room/{room}/beers/{beer}/unpublish",
		handleUnpublishRatingsForBeer(roomService, beerService, logger),
	)

	mux.Handle(
		"/api/room/{room}/beers/{beer}/rate",
		handleVoteOnBeer(beerService, roomService, logger),
	)

	mux.Handle(
		"/api/user/profile",
		handleGetUserProfile(userService, logger),
	)

	mux.Handle(
		"/api/user/updateProfile",
		handleUpdateUserProfile(userService, logger),
	)

	mux.Handle(
		"/api/verifyToken",
		handleTestToken(logger),
	)

	return mux
}

func handleRooms(
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(ContextUserKey)
			rooms, err := rs.GetRoomsByUserId(r.Context(), userId.(int))
			if err != nil {
				if err.Error() == "Unauthorized" {
					logger.Error("handleRooms", "err", err)
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				logger.Error("handleRooms", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(rooms)
		},
	)
}

func handleEditRoom(
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(ContextUserKey)
			roomId, err := strconv.Atoi(r.PathValue("room"))
			if err != nil {
				logger.Error("handleEditRoom/strconv", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if ok, err := rs.CheckIfUserIsAdminInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handleEditRoom", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			var data rooms.Room

			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				logger.Error("handleEditRoom/strconv", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			data.Id = roomId
			err = rs.UpdateRoom(r.Context(), data)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					logger.Error("handleEditRoom/room", "err", "Room not found")
					http.Error(w, "Room not found", http.StatusNotFound)
					return
				}
				logger.Error("handleEditRoom/db", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Success")
		},
	)
}

func handleRoom(
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(ContextUserKey)
			roomId, _ := strconv.Atoi(r.PathValue("room"))
			room, err := rs.GetRoomById(r.Context(), roomId)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					logger.Error("handleRoom/room", "err", "Room not found")
					http.Error(w, "Room not found", http.StatusNotFound)
					return
				}
				logger.Error("handleRoom/db", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if ok, err := rs.CheckIfUserInRoom(r.Context(), room.Id, userId.(int)); !ok || err != nil {
				logger.Error("hanldeGetRoom", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(&room)
			return
		},
	)
}

func handleUsersInRoom(
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			roomId, _ := strconv.Atoi(r.PathValue("room"))
			users, err := rs.GetUsersInRoom(r.Context(), roomId)

			if err != nil {
				logger.Error("handleUsersInRooms", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			userId := r.Context().Value(ContextUserKey)

			if ok, err := rs.CheckIfUserInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handleUsersInRoom", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users)
		},
	)
}

func handleUpdateIsAdminForUser(
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(ContextUserKey)
			roomId, err := strconv.Atoi(r.PathValue("room"))
			if err != nil {
				logger.Error("handleIsAdminForUser", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			targetUserId, err := strconv.Atoi(r.PathValue("user"))
			if err != nil {
				logger.Error("handleIsAdminForUser", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if ok, err := rs.CheckIfUserIsAdminInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handlePublishRatingsForBeer", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var data struct {
				IsAdmin bool `json:"isAdmin"`
			}

			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				logger.Error("handleIsAdminForUser", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !data.IsAdmin {
				if ok, err := rs.CheckIfOtherAdminInRoom(r.Context(), roomId, targetUserId); !ok || err != nil {
					logger.Error("handleGetSingleBeer", "err", err)
					http.Error(w, "Only admin left", http.StatusUnprocessableEntity)
					return
				}
			}

			err = rs.UpdateIsAdmin(r.Context(), roomId, targetUserId, data.IsAdmin)
			if err != nil {
				logger.Error("handleIsAdminForUser", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Success")
		},
	)
}

func handleCheckIfUserIsAdminInRoom(
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			roomId, err := strconv.Atoi(r.PathValue("room"))
			if err != nil {
				logger.Error("handleCheckIfUserIsAdminInRoom", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			userId := r.Context().Value(ContextUserKey)
			isAdmin, err := rs.CheckIfUserIsAdminInRoom(r.Context(), roomId, userId.(int))
			if err != nil {
				logger.Error("handleCheckIfUserIsAdminInRoom", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(isAdmin)
		},
	)
}

func handleBeersInRoom(
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(ContextUserKey)
			roomId, _ := strconv.Atoi(r.PathValue("room"))
			if ok, err := rs.CheckIfUserInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("hanldeGetRoom", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			beers, err := rs.GetBeersInRoom(r.Context(), roomId)
			if err != nil {
				if err.Error() == "Unauthorized" {
					logger.Error("handleRooms", "err", err)
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				logger.Error("handleBeersInRoom", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(beers)
		},
	)
}

func handleAddBeer(
	bs *beers.BeerService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			roomId, _ := strconv.Atoi(r.PathValue("room"))
			var beer beers.Beer
			json.NewDecoder(r.Body).Decode(&beer)
			logger.Info("HandleAddBeer", "beer", beer)
			err := bs.AddNewBeer(r.Context(), beer.Name, beer.Style, beer.PictureUrl, roomId)

			if err != nil {
				logger.Error("handleBeersInRoom", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(beer)
		},
	)
}

func handleRemoveUserFromRoom(
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			userId := r.Context().Value(ContextUserKey)

			roomId, err := strconv.Atoi(r.PathValue("room"))

			if err != nil {
				logger.Error("handleGetSingleBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if ok, err := rs.CheckIfUserIsAdminInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handleGetSingleBeer", "err", err)
				http.Error(w, "Not in room", http.StatusUnprocessableEntity)
				return
			}
			targetUserId, err := strconv.Atoi(r.PathValue("user"))
			if err != nil {
				logger.Error("handleJoinRoom", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			err = rs.RemoveUserFromRoom(r.Context(), roomId, targetUserId)
			if err != nil {
				logger.Error("handleJoinRoom", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Success")
		},
	)
}

func handleLeaveRoom(
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			userId := r.Context().Value(ContextUserKey)

			roomId, err := strconv.Atoi(r.PathValue("room"))

			if err != nil {
				logger.Error("handleGetSingleBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if ok, err := rs.CheckIfUserInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handleGetSingleBeer", "err", err)
				http.Error(w, "Not in room", http.StatusUnprocessableEntity)
				return
			}

			logger.Info("LeaveRoom", "userId", userId.(int), "roomId", roomId)
			if ok, err := rs.CheckIfOtherAdminInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handleGetSingleBeer", "err", err)
				http.Error(w, "Only admin left", http.StatusUnprocessableEntity)
				return
			}

			err = rs.RemoveUserFromRoom(r.Context(), roomId, userId.(int))
			if err != nil {
				logger.Error("handleJoinRoom", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Success")
		},
	)
}

func handleJoinRoom(
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			userId := r.Context().Value(ContextUserKey)
			var data struct {
				Code string `json:"code"`
			}
			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				logger.Error("handleJoinRoom", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			room, err := rs.GetRoomByCode(r.Context(), data.Code)

			if err != nil {
				logger.Error("handleJoinRoom", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if ok, err := rs.CheckIfUserInRoom(r.Context(), room.Id, userId.(int)); ok || err != nil {
				logger.Error("handleGetSingleBeer", "err", err)
				http.Error(w, "Already in room", http.StatusUnprocessableEntity)
				return
			}
			logger.Info("handleJoinRoom", "user", userId.(int), "room", room)

			err = rs.AddUserToRoom(r.Context(), room.Id, userId.(int), false)
			if err != nil {
				logger.Error("handleJoinRoom", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(room)
		},
	)
}

func handleCreateRoom(
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(ContextUserKey)
			var room rooms.Room
			err := json.NewDecoder(r.Body).Decode(&room)
			if err != nil {
				logger.Error("handleCreateRoom", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			id, err := rs.CreateNewRoom(r.Context(), userId.(int), room)
			if err != nil {
				logger.Error("handleCreateRoom", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			room.Id = id
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(room)
		},
	)
}

func handleGetSingleBeer(
	bs *beers.BeerService,
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			roomId, err := strconv.Atoi(r.PathValue("room"))

			if err != nil {
				logger.Error("handleGetSingleBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			userId := r.Context().Value(ContextUserKey)

			if ok, err := rs.CheckIfUserInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handleGetSingleBeer", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			beerId, err := strconv.Atoi(r.PathValue("beer"))

			if err != nil {
				logger.Error("handleGetSingleBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			isAdmin := false
			if ok, err := rs.CheckIfUserIsAdminInRoom(r.Context(), roomId, userId.(int)); ok && err == nil {
				isAdmin = true
			}

			beer, err := bs.GetBeerById(r.Context(), beerId, roomId, isAdmin)
			if err != nil {
				logger.Error("handleGetSingleBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(beer)

		},
	)
}

func handleGetMyRatingForBeer(
	rs *rooms.RoomService,
	bs *beers.BeerService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			roomId, err := strconv.Atoi(r.PathValue("room"))

			if err != nil {
				logger.Error("handleGetMyRatingForBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			userId := r.Context().Value(ContextUserKey)

			if ok, err := rs.CheckIfUserInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handleGetMyRatingForBeer", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			beerId, err := strconv.Atoi(r.PathValue("beer"))

			if err != nil {
				logger.Error("handleGetMyRatingForBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			vote, err := bs.GetMyRatingOnBeer(r.Context(), beerId, userId.(int))
			if err != nil {
				logger.Error("handleGetMyRatingForBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if vote != nil {
				vote.BeerId = beerId
				vote.UserId = userId.(int)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(vote)

		},
	)
}

func handleGetRatingsForBeer(
	rs *rooms.RoomService,
	bs *beers.BeerService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			roomId, err := strconv.Atoi(r.PathValue("room"))

			if err != nil {
				logger.Error("handleGetRatingsForBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			userId := r.Context().Value(ContextUserKey)

			if ok, err := rs.CheckIfUserInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handleGetRatingsForBeer", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			beerId, err := strconv.Atoi(r.PathValue("beer"))

			if err != nil {
				logger.Error("handleGetRatingsForBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			ratings, err := bs.GetVotesByBeerId(r.Context(), beerId, roomId)
			if err != nil {
				logger.Error("handleGetRatingsForBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ratings)
		},
	)
}

func handleUnpublishRatingsForBeer(
	rs *rooms.RoomService,
	bs *beers.BeerService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			roomId, err := strconv.Atoi(r.PathValue("room"))

			if err != nil {
				logger.Error("handlePublishRatingsForBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			userId := r.Context().Value(ContextUserKey)

			if ok, err := rs.CheckIfUserIsAdminInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handlePublishRatingsForBeer", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			beerId, err := strconv.Atoi(r.PathValue("beer"))

			if err != nil {
				logger.Error("handlePublishRatingsForBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			err = bs.UnpublishRatingsForBeer(r.Context(), beerId, roomId)
			if err != nil {
				logger.Error("handlePublishRatingsForBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Published")
		},
	)
}

func handlePublishRatingsForBeer(
	rs *rooms.RoomService,
	bs *beers.BeerService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			roomId, err := strconv.Atoi(r.PathValue("room"))

			if err != nil {
				logger.Error("handlePublishRatingsForBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			userId := r.Context().Value(ContextUserKey)

			if ok, err := rs.CheckIfUserIsAdminInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handlePublishRatingsForBeer", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			beerId, err := strconv.Atoi(r.PathValue("beer"))

			if err != nil {
				logger.Error("handlePublishRatingsForBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			err = bs.PublishRatingsForBeer(r.Context(), beerId, roomId)
			if err != nil {
				logger.Error("handlePublishRatingsForBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Published")
		},
	)
}

func handleVoteOnBeer(
	bs *beers.BeerService,
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			roomId, err := strconv.Atoi(r.PathValue("room"))

			if err != nil {
				logger.Error("handleVoteOnBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			userId := r.Context().Value(ContextUserKey)

			if ok, err := rs.CheckIfUserInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handleVoteOnBeer", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			beerId, err := strconv.Atoi(r.PathValue("beer"))

			if err != nil {
				logger.Error("handleVoteOnBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if ok, err := rs.CheckIfBeerInRoom(r.Context(), roomId, beerId); !ok || err != nil {
				logger.Error("handleVoteOnBeer", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var vote beers.Vote
			err = json.NewDecoder(r.Body).Decode(&vote)
			if err != nil {
				logger.Error("handleVoteOnBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			logger.Info("handleVoteOnBeer", "vote", vote)
			vote.UserId = userId.(int)
			vote.BeerId = beerId

			err = bs.UpdateVoteOnBeerId(r.Context(), vote)
			if err != nil {
				logger.Error("handleVoteOnBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Success")
		},
	)
}
func handleGetRandomBeer(
	bs *beers.BeerService,
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			roomId, err := strconv.Atoi(r.PathValue("room"))

			if err != nil {
				logger.Error("handleGetRandomBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			userId := r.Context().Value(ContextUserKey)

			if ok, err := rs.CheckIfUserIsAdminInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handleGetRandomBeer", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			beer, err := bs.GetRandomBeer(r.Context(), roomId)
			if err != nil {
				logger.Error("handleGetRandomBeer", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(beer)
		},
	)
}

func handleGetNextBeer(
	bs *beers.BeerService,
	rs *rooms.RoomService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			roomId, err := strconv.Atoi(r.PathValue("room"))

			if err != nil {
				logger.Error("handleGetNextBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			var data struct {
				OldBeerId int `json:"oldBeerId"`
			}

			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				logger.Error("handleGetNextBeer", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			userId := r.Context().Value(ContextUserKey)

			if ok, err := rs.CheckIfUserIsAdminInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handleGetNextBeer", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			logger.Info("OldBeerId check", "OldBeerId", data.OldBeerId)
			beer, err := bs.GetNextBeer(r.Context(), roomId, data.OldBeerId)
			logger.Error("handleGetNextBeer", "beer", beer)
			if err != nil {
				logger.Error("handleGetNextBeer", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(beer)
		},
	)
}

func handleEditBeer(
	rs *rooms.RoomService,
	bs *beers.BeerService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(ContextUserKey)
			roomId, err := strconv.Atoi(r.PathValue("room"))
			if err != nil {
				logger.Error("handleEditRoom/strconv", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			beerId, err := strconv.Atoi(r.PathValue("beer"))
			if err != nil {
				logger.Error("handleEditRoom/strconv", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if ok, err := rs.CheckIfUserInRoom(r.Context(), roomId, userId.(int)); !ok || err != nil {
				logger.Error("handleEditRoom", "err", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			var data beers.Beer

			err = json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				logger.Error("handleEditRoom/strconv", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			data.Id = beerId
			err = bs.UpdateBeer(r.Context(), data, roomId)
			if err != nil {
				logger.Error("handleEditRoom/db", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Success")
		},
	)
}

func handleGetUserProfile(
	us *auth.UserService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(ContextUserKey)
			u, err := us.GetUserById(r.Context(), userId.(int))
			if err != nil {
				logger.Error("handleGetUserProfile/GetUserById", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(u)

		},
	)
}

func handleUpdateUserProfile(
	us *auth.UserService,
	logger *slog.Logger,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(ContextUserKey)
			var data auth.UpdateProfile

			err := json.NewDecoder(r.Body).Decode(&data)
			if err != nil {
				logger.Error("handleUpdateUserProfile/bodyDecode", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			err = us.UpdateUserProfile(r.Context(), userId.(int), data)
			if err != nil {
				logger.Error("handleUpdateUserProfile/dbUpdate", "err", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("Success")
		},
	)
}
