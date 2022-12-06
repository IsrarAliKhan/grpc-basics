package services

import (
	"context"
	"errors"
	"item/pb"
	"item/server/db"
	"item/server/db/models"
	"item/server/log"
	"time"
)

type ItemServer struct {
	pb.UnimplementedItemsServer
}

func (i *ItemServer) GetItems(req *pb.Empty, stream pb.Items_GetItemsServer) error {
	// get items form db
	var items []models.Item
	err := db.Conn().
		Find(&items).
		Error
	if err != nil {
		log.Error.Println(err)
		return err
	}

	// send items in stream
	for _, i := range items {
		item := &pb.Item{
			Id:       i.Id,
			Name:     i.Name,
			Price:    i.Price,
			Quantity: i.Quantity,
		}

		if err := stream.Send(&pb.ItemInfo{Item: item}); err != nil {
			return err
		}
	}

	// respond
	return nil
}

func (i *ItemServer) GetItem(ctx context.Context, req *pb.Id) (*pb.ItemInfo, error) {
	log.Info.Printf("Recieved: %v\n", req)

	// validate request
	if req.GetId() == 0 {
		log.Error.Println("invalid request id")
		return nil, errors.New("invalid request id")
	}

	// get item form db
	var item models.Item
	err := db.Conn().
		Where("id = ?", req.GetId()).
		First(&item).
		Error
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	// build response
	res := &pb.ItemInfo{
		Item: &pb.Item{
			Id:       item.Id,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
		},
	}

	// respond
	return res, nil
}

func (i *ItemServer) CreateItem(ctx context.Context, req *pb.ItemInfo) (*pb.Id, error) {
	log.Info.Printf("Recieved: %v\n", req)

	// build db object
	var item = models.Item{
		Name:     req.Item.Name,
		Price:    req.Item.Price,
		Quantity: req.Item.Quantity,
	}

	// save item in db
	err := db.Conn().Create(&item).Error
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	// build response
	res := &pb.Id{
		Id: item.Id,
	}

	// respond
	return res, nil
}

func (i *ItemServer) UpdateItem(ctx context.Context, req *pb.ItemInfo) (*pb.Status, error) {
	log.Info.Printf("Recieved: %v\n", req)

	// validate request
	if req.Item.Id == 0 {
		log.Error.Println("invalid request id")
		return nil, errors.New("invalid request id")
	}

	// get item form db
	var item models.Item
	err := db.Conn().
		Where("id = ?", req.Item.Id).
		First(&item).
		Error
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	// update item
	item.Name = req.Item.Name
	item.Price = req.Item.Price
	item.Quantity = req.Item.Quantity

	// update item in db
	err = db.Conn().Save(&item).Error
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	// build response
	res := &pb.Status{
		Status: "item updated successfully",
	}

	// respond
	return res, nil
}

func (i *ItemServer) DeleteItem(ctx context.Context, req *pb.Id) (*pb.Status, error) {
	log.Info.Printf("Recieved: %v\n", req)

	// validate request
	if req.Id == 0 {
		log.Error.Println("invalid request id")
		return nil, errors.New("invalid request id")
	}

	// check item in db
	err := db.Conn().
		Where("id = ?", req.Id).
		First(&models.Item{}).
		Error
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	// delete item in db
	err = db.Conn().
		Where("id = ?", req.Id).
		Model(&models.Item{}).
		Update("deleted_at", time.Now()).
		Error
	if err != nil {
		log.Error.Println(err)
		return nil, err
	}

	// build response
	res := &pb.Status{
		Status: "item deleled successfully",
	}

	// respond
	return res, nil
}
