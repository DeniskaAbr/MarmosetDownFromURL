package main

import "time"

// структура ответа сервера, JSON файла
type responce struct {
	Data []Project `json;"data"`
}

// структура данных проекта части JSON
type Project struct {
	Id                  int       `json:"id"`
	User_id             int       `json:"user_id"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	Created_at          time.Time `json:"created_at"` /* 2019-09-19T12:25:51.895-05:00 */
	Updated_at          time.Time `json:"updated_at"` /* 2019-09-24T04:56:38.873-05:00 */
	Likes_count         int       `json:"likes_count"`
	Slug                string    `json:"slug"`
	Published_at        time.Time `json:"published_at"` /* 2019-09-19T12:28:06.632-05:00 */
	Adult_content       bool      `json:"adult_content"`
	Cover_asset_id      int       `json:"cover_asset_id"`
	Admin_adult_content bool      `json:"admin_adult_content"`
	Views_count         int       `json:"views_count"`
	Hash_id             string    `json:"hash_id"`
	Permalink           string    `json:"permalink"`
	Hide_as_adult       bool      `json:"hide_as_adult"`
	User                *User     `json:"user"`
	Cover               *Cover    `json:"cover"`
	Icons               *Icons    `json:"icons"`
	URLs				[]*URLs
	Checked				bool
	Assets_count		int			`json:"assets_count"`
}

// структура данных пользователя части JSON
type User struct {
	Id                     int    `json:""` // 45807
	Username               string `json:""` // alishermirzoev
	First_name             string `json:""` // Alisher,
	Last_name              string `json:""` // Mirzoev,
	Avatar_file_name       string `json:""` // bdfe79dac80113103fbe0c103ce85343.jpg,
	Country                string `json:""` // Russia,
	City                   string `json:""` // Moscow,
	Subdomain              string `json:""` // alishermirzoev,
	Headline               string `json:""` // Concept artist | Matte painter | yellomice@gmail.com,
	Pro_member             bool   `json:""` // true,
	Is_staff               bool   `json:""` // false,
	Is_plus_member         bool   `json:""` // false,
	Medium_avatar_url      string `json:""` // https//cdnb.artstation.com/p/users/avatars/000/045/807/medium/bdfe79dac80113103fbe0c103ce85343.jpg,
	Large_avatar_url       string `json:""` // https//cdnb.artstation.com/p/users/avatars/000/045/807/large/bdfe79dac80113103fbe0c103ce85343.jpg,
	Full_name              string `json:""` // Alisher Mirzoev,
	Permalink              string `json:""` // http//www.artstation.com/alishermirzoev,
	Artstation_profile_url string `json:""` // https//www.artstation.com/alishermirzoev,
	Location               string `json:""` // Moscow, Russia
}

// структура данных обложки части JSON
type Cover struct {
	Id                     int    `json:"id"`                     // 20724456,
	Small_square_url       string `json:"small_square_url"`       // https://cdna.artstation.com/p/assets/images/images/020/724/456/20190919124105/small_square/alisher-mirzoev-alishermirzoev-thumb-03.jpg?1568914866,
	Micro_square_image_url string `json:"micro_square_image_url"` // https://cdna.artstation.com/p/assets/images/images/020/724/456/20190919124105/micro_square/alisher-mirzoev-alishermirzoev-thumb-03.jpg?1568914866,
	Thumb_url              string `json:"thumb_url"`              // https://cdna.artstation.com/p/assets/images/images/020/724/456/20190919124105/smaller_square/alisher-mirzoev-alishermirzoev-thumb-03.jpg?1568914866
}

// структура данных иконки части JSON
type Icons struct {
	Image    bool `json:"image"`
	Video    bool `json:"video"`
	Video_clip bool `json:"video_clip"`
	Model3d  bool `json:"model3d"`
	Marmoset bool `json:"marmoset"`
	Pano     bool `json:"pano"`
}

type URLs struct {
	PageURL string
	FileURL string
}

type urlarr struct {
	Data []*URLs `json;"data"`
}

