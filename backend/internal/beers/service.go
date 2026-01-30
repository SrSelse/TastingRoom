package beers

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"skafteresort.se/beers/internal/providers"
)

type BeerService struct {
	beerRepo   *BeerRepo
	logger     *slog.Logger
	centrifugo *providers.CentrifugoProvider
}

func NewBeerService(br *BeerRepo, logger *slog.Logger, gp *providers.CentrifugoProvider) *BeerService {
	ts := BeerService{
		beerRepo:   br,
		logger:     logger,
		centrifugo: gp,
	}
	return &ts
}

func (s *BeerService) GetVotesByBeerId(ctx context.Context, beerId int, roomId int) ([]Vote, error) {
	return s.beerRepo.getVotesByBeerId(ctx, beerId, roomId)
}

func (s *BeerService) GetBeerById(ctx context.Context, beerId int, roomId int, isAdmin bool) (*Beer, error) {
	beer, err := s.beerRepo.getBeerById(ctx, beerId)
	if err != nil || beer == nil {
		return nil, err
	}
	if isAdmin {
		cMessage := s.centrifugo.CreateMessageWithChannel(
			fmt.Sprintf("rooms:%d-next-beer", roomId),
			map[string]string{
				"beerId": strconv.Itoa(beer.Id),
			},
		)
		s.centrifugo.HandleMessage(ctx, cMessage)
	}

	return beer, nil
}

func (s *BeerService) UpdateVoteOnBeerId(
	ctx context.Context,
	vote Vote,
) error {
	if vote.Id != 0 {
		err := s.beerRepo.updateVoteOnBeerId(ctx, vote)
		if err != nil {
			return err
		}
	} else {
		err := s.beerRepo.addVoteOnBeerId(ctx, vote)
		if err != nil {
			return err
		}
	}

	note := ""
	if vote.Note != nil {
		note = *vote.Note
	}

	cMessage := s.centrifugo.CreateVoteMessage(
		vote.BeerId,
		vote.UserId,
		vote.Value,
		note,
		vote.UserName,
		"vote-updated",
	)
	s.centrifugo.HandleMessage(ctx, cMessage)
	return nil

}

func (s *BeerService) AddNewBeer(
	ctx context.Context,
	name string,
	beerType *string,
	pictureUrl *string,
	roomId int,
) error {
	b := Beer{
		Name:       name,
		Style:      beerType,
		PictureUrl: pictureUrl,
		RoomId:     roomId,
	}
	err := s.beerRepo.addNewBeer(ctx, b)
	if err != nil {
		return err
	}
	cMessage := s.centrifugo.CreateBeerMessage(
		b.RoomId,
		"beer-added",
	)
	s.centrifugo.HandleMessage(ctx, cMessage)
	return nil
}

func (s *BeerService) GetBeersByUserVotes(ctx context.Context, userId int) ([]Beer, error) {
	return s.beerRepo.getBeersByUserVotes(ctx, userId)
}

func (s *BeerService) PublishRatingsForBeer(ctx context.Context, beerId int, roomId int) error {
	return s.beerRepo.publishRatingsForBeer(ctx, beerId, roomId)
}

func (s *BeerService) UnpublishRatingsForBeer(ctx context.Context, beerId int, roomId int) error {
	return s.beerRepo.unpublishRatingsForBeer(ctx, beerId, roomId)
}
func (s *BeerService) GetMyRatingOnBeer(ctx context.Context, beerId int, userId int) (*Vote, error) {
	return s.beerRepo.getMyRatingOnBeer(ctx, beerId, userId)
}

func (s *BeerService) GetRandomBeer(ctx context.Context, roomId int) (*Beer, error) {
	beer, err := s.beerRepo.getRandomBeerInRoom(ctx, roomId)
	if err != nil || beer == nil {
		return nil, err
	}

	cMessage := s.centrifugo.CreateMessageWithChannel(
		fmt.Sprintf("rooms:%d-next-beer", roomId),
		map[string]string{
			"beerId": strconv.Itoa(beer.Id),
		},
	)
	s.centrifugo.HandleMessage(ctx, cMessage)

	return beer, nil
}

func (s *BeerService) GetNextBeer(ctx context.Context, roomId int, oldBeerId int) (*Beer, error) {
	beer, err := s.beerRepo.getNextBeerInRoom(ctx, roomId, oldBeerId)
	if err != nil || beer == nil {
		return nil, err
	}

	cMessage := s.centrifugo.CreateMessageWithChannel(
		fmt.Sprintf("rooms:%d-next-beer", roomId),
		map[string]string{
			"beerId": strconv.Itoa(beer.Id),
		},
	)
	s.centrifugo.HandleMessage(ctx, cMessage)

	return beer, nil
}

func (s *BeerService) UpdateBeer(ctx context.Context, beer Beer, roomId int) error {
	return s.beerRepo.updateBeer(ctx, beer, roomId)
}
