package program

import (
	"context"
	"log"

	"github.com/JUNGHUNKIM-7/cli_gpt/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	Coll   *mongo.Collection
)

func Initialize(mongoUri string) {
	if mongoUri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	coll := client.Database("gpt_personal").Collection("qna")
	Client = client
	Coll = coll
}

func SetData(body *model.QnaBody) {
	_, err := Coll.InsertOne(context.TODO(), *body)
	if err != nil {
		log.Fatal(err)
	}
}

func SetAll(body []model.QnaBody) {
	bodies := make([]interface{}, len(body))
	for i, v := range body {
		bodies[i] = v
	}
	_, err := Coll.InsertMany(context.TODO(), bodies)
	if err != nil {
		log.Fatal(err)
	}
}

func GetData(q string) []*model.QnaBody {
	historiesFrom := make([]*model.QnaBody, 0)
	ctx := context.Background()
	filter := bson.M{"q": bson.M{"$regex": q, "$options": "i"}}
	cursor, err := Coll.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var body model.QnaBody
		if err := cursor.Decode(&body); err != nil {
			log.Fatal(err)
		}
		historiesFrom = append(historiesFrom, &body)
	}

	return historiesFrom
}
