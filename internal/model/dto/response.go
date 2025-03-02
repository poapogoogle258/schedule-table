package dto

// type ResponseTaskDescription struct {
// 	Id          string `json:"id"`
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	ImageURL    string `json:"imageURL"`
// 	Priority    int8   `json:"priority"`
// }

// type ResponseTask struct {
// 	Id          string                  `json:"id"`
// 	CalendarId  string                  `json:"calendar_id"`
// 	ScheduleId  string                  `json:"schedule_id"`
// 	Start       time.Time               `json:"start"`
// 	End         time.Time               `json:"end"`
// 	Status      int8                    `json:"status"`
// 	Person      RequestMember           `json:"person"`
// 	Description ResponseTaskDescription `json:"description"`
// }

type Pagination struct {
	Total       int64 `json:"total_records"`
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
}
