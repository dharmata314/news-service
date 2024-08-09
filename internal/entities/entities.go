package entities

type News struct {
	ID      int    `json:"news_id"`
	Title   string `json:"news_title"`
	Content string `json:"news_content"`
}

type Categorie struct {
	ID   int `json:"categorie_id"`
	Name int `json:"categorie_name"`
}

type NewsCategories struct {
	ID         int
	CategoryID int
	NewsID     int
}

type User struct {
	ID       int    `json:"user_id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
