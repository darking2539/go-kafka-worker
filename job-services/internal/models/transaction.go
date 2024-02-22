package models

type TransactionModel struct {
    Title       string `bson:"title"`
    Genre       string `bson:"genre"`
    ReleaseYear string `bson:"release_year"`
    Platform    string `bson:"platform"`
    Developer   string `bson:"developer"`
    Publisher   string `bson:"publisher"`
    Rating      string `bson:"rating"`
}
