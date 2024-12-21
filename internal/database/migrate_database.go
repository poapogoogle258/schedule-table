package database

import (
	"schedule_table/internal/model/dao"
	"time"

	"github.com/google/uuid"

	"schedule_table/util"

	rrule "github.com/teambition/rrule-go"

	"gorm.io/gorm"
)

func MigrateSetUpAndInitData(db *gorm.DB) error {

	dropTable_error := db.Migrator().DropTable(&dao.Tasks{}, &dao.Responsible{}, &dao.Schedules{}, &dao.Leaves{}, &dao.Members{}, &dao.Calendars{}, &dao.Users{})

	if dropTable_error != nil {
		return dropTable_error
	}

	migrate_err := db.AutoMigrate(&dao.Users{}, &dao.Leaves{}, &dao.Members{}, &dao.Schedules{}, &dao.Responsible{}, &dao.Tasks{}, &dao.Calendars{})

	if migrate_err != nil {
		return migrate_err
	}

	initUser := &dao.Users{
		Id:       uuid.New(),
		Name:     "user_testing",
		Email:    "user_testing@gmail.com",
		Password: util.HashPassword("password123"),
	}

	initCalendar := &dao.Calendars{
		Id:     uuid.New(),
		UserId: initUser.Id,
		Name:   "default",
	}

	initMembers := []*dao.Members{
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Danny Ownsworth",
			Nickname:    "Danny",
			ImageURL:    "https://robohash.org/etquiquia.png?size=200x200&set=set1",
			Color:       "#9612ff",
			Email:       "downsworth0@nytimes.com",
			Description: "Superficial foreign body of penis",
			Position:    "Barr Laboratories Inc.",
			Telephone:   "633 482 6982",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Berkly Maun",
			Nickname:    "Berkly",
			ImageURL:    "https://robohash.org/quiablanditiislaborum.png?size=200x200&set=set1",
			Color:       "#e5559a",
			Email:       "bmaun1@gmpg.org",
			Description: "Unsp injury at C4 level of cervical spinal cord, init encntr",
			Position:    "Infinite Therapies of Sarasota, Inc.",
			Telephone:   "410 886 3030",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Web Stepto",
			Nickname:    "Web",
			ImageURL:    "https://robohash.org/estaliquidporro.png?size=200x200&set=set1",
			Color:       "#a5c6d4",
			Email:       "wstepto2@tripadvisor.com",
			Description: "Quadriplegia, C1-C4 complete",
			Position:    "McKesson Contract Packaging",
			Telephone:   "560 857 3267",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Haywood Hordell",
			Nickname:    "Haywood",
			ImageURL:    "https://robohash.org/sedpraesentiumquas.png?size=200x200&set=set1",
			Color:       "#64fabb",
			Email:       "hhordell3@xinhuanet.com",
			Description: "Greenstick fx shaft of humer, r arm, 7thG",
			Position:    "Juice Beauty",
			Telephone:   "901 410 4362",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Phillipp Devenny",
			Nickname:    "Phillipp",
			ImageURL:    "https://robohash.org/suscipiteligendimaiores.png?size=200x200&set=set1",
			Color:       "#9952b1",
			Email:       "pdevenny4@devhub.com",
			Description: "Complete traumatic amputation at elbow level, right arm",
			Position:    "Unifirst First Aid Corporation",
			Telephone:   "196 128 7041",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Kattie Wixey",
			Nickname:    "Kattie",
			ImageURL:    "https://robohash.org/nisiquiaut.png?size=200x200&set=set1",
			Color:       "#c30e64",
			Email:       "kwixey5@privacy.gov.au",
			Description: "Disp fx of dist phalanx of unsp great toe, 7thK",
			Position:    "Seton Pharmaceuticals",
			Telephone:   "951 147 4079",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Justin Lyons",
			Nickname:    "Justin",
			ImageURL:    "https://robohash.org/odiopossimusquis.png?size=200x200&set=set1",
			Color:       "#f74514",
			Email:       "jlyons6@dailymotion.com",
			Description: "Insect bite (nonvenomous), unspecified hip, subs encntr",
			Position:    "H E B",
			Telephone:   "118 601 3195",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Amelia Hellin",
			Nickname:    "Amelia",
			ImageURL:    "https://robohash.org/impeditlaudantiumunde.png?size=200x200&set=set1",
			Color:       "#27e206",
			Email:       "ahellin7@unicef.org",
			Description: "Nondisp posterior arch fx first cervcal vertebra, sequela",
			Position:    "Boehringer Ingelheim Pharmaceuticals, Inc.",
			Telephone:   "549 402 8354",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Tammy Focke",
			Nickname:    "Tammy",
			ImageURL:    "https://robohash.org/facereeteos.png?size=200x200&set=set1",
			Color:       "#beabbf",
			Email:       "tfocke8@about.me",
			Description: "Bipolar disorder, currently in remission",
			Position:    "REMEDYREPACK INC.",
			Telephone:   "878 593 5108",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Riva Klees",
			Nickname:    "Riva",
			ImageURL:    "https://robohash.org/utsitharum.png?size=200x200&set=set1",
			Color:       "#7cfdfa",
			Email:       "rklees9@oracle.com",
			Description: "Inj blood vessel of left little finger, init encntr",
			Position:    "Sandoz Inc",
			Telephone:   "208 742 6525",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Kizzie Reilingen",
			Nickname:    "Kizzie",
			ImageURL:    "https://robohash.org/dolorumesseinventore.png?size=200x200&set=set1",
			Color:       "#e8ee79",
			Email:       "kreilingena@go.com",
			Description: "Monoplg low lmb fol ntrm subarach hemor aff r nondom side",
			Position:    "Bryant Ranch Prepack",
			Telephone:   "320 532 3284",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Stevie Ferebee",
			Nickname:    "Stevie",
			ImageURL:    "https://robohash.org/fuganemovoluptatem.png?size=200x200&set=set1",
			Color:       "#7e0154",
			Email:       "sferebeeb@hibu.com",
			Description: "Periodic fever syndromes",
			Position:    "Premium Formulations LLC",
			Telephone:   "908 349 2975",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Nicolas Leyes",
			Nickname:    "Nicolas",
			ImageURL:    "https://robohash.org/quosutet.png?size=200x200&set=set1",
			Color:       "#bbeebe",
			Email:       "nleyesc@lulu.com",
			Description: "Oth personality & behavrl disord due to known physiol cond",
			Position:    "EQUALINE (SuperValu)",
			Telephone:   "973 205 9986",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Faunie Zisneros",
			Nickname:    "Faunie",
			ImageURL:    "https://robohash.org/autdoloremaccusamus.png?size=200x200&set=set1",
			Color:       "#c58cf1",
			Email:       "fzisnerosd@smh.com.au",
			Description: "Sprain of unspecified site of right knee",
			Position:    "Bryant Ranch Prepack",
			Telephone:   "531 937 5218",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Shay Chung",
			Nickname:    "Shay",
			ImageURL:    "https://robohash.org/illovelitut.png?size=200x200&set=set1",
			Color:       "#bc72eb",
			Email:       "schunge@linkedin.com",
			Description: "Acute ischemia of small intestine, extent unspecified",
			Position:    "Major Pharmaceuticals",
			Telephone:   "287 107 7254",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Marietta Saladin",
			Nickname:    "Marietta",
			ImageURL:    "https://robohash.org/sedasperioreset.png?size=200x200&set=set1",
			Color:       "#0acd85",
			Email:       "msaladinf@ed.gov",
			Description: "Complex tear of medial meniscus, current injury",
			Position:    "Rx Pak Division of McKesson Corporation",
			Telephone:   "770 604 7071",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Sophie Stoffel",
			Nickname:    "Sophie",
			ImageURL:    "https://robohash.org/magnamipsumsed.png?size=200x200&set=set1",
			Color:       "#8a8d6f",
			Email:       "sstoffelg@admin.ch",
			Description: "Other superficial bite of left wrist, subsequent encounter",
			Position:    "HOMEOLAB USA INC.",
			Telephone:   "393 812 1026",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Orville Bunten",
			Nickname:    "Orville",
			ImageURL:    "https://robohash.org/facilisvelatque.png?size=200x200&set=set1",
			Color:       "#b4f05a",
			Email:       "obuntenh@howstuffworks.com",
			Description: "Sprain of unspecified ligament of right ankle",
			Position:    "Uriel Pharmacy Inc.",
			Telephone:   "210 997 0229",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Issie Frier",
			Nickname:    "Issie",
			ImageURL:    "https://robohash.org/fugiatquiacorrupti.png?size=200x200&set=set1",
			Color:       "#5256c6",
			Email:       "ifrieri@virginia.edu",
			Description: "Other non-in-line roller-skating accident, sequela",
			Position:    "GlaxoSmithKline LLC",
			Telephone:   "612 593 2537",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Jonah Corkel",
			Nickname:    "Jonah",
			ImageURL:    "https://robohash.org/quiaetet.png?size=200x200&set=set1",
			Color:       "#6ae5ef",
			Email:       "jcorkelj@weibo.com",
			Description: "Bitten by duck",
			Position:    "PD-Rx Pharmaceuticals, Inc.",
			Telephone:   "615 512 4885",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Elise Capinetti",
			Nickname:    "Elise",
			ImageURL:    "https://robohash.org/quiaidadipisci.png?size=200x200&set=set1",
			Color:       "#f74128",
			Email:       "ecapinettik@myspace.com",
			Description: "Person outside car inj in clsn w rail trn/veh in traf, subs",
			Position:    "ALK-Abello, Inc.",
			Telephone:   "399 628 0632",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Tracee Cullimore",
			Nickname:    "Tracee",
			ImageURL:    "https://robohash.org/adtemporeconsectetur.png?size=200x200&set=set1",
			Color:       "#a5c262",
			Email:       "tcullimorel@seattletimes.com",
			Description: "Underdosing of inhaled anesthetics, sequela",
			Position:    "Menper Distributors, Inc.",
			Telephone:   "811 653 6487",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Caldwell Streight",
			Nickname:    "Caldwell",
			ImageURL:    "https://robohash.org/necessitatibusducimusqui.png?size=200x200&set=set1",
			Color:       "#16b8eb",
			Email:       "cstreightm@pinterest.com",
			Description: "Oth comp specific to multiple gest, first trimester, fetus 2",
			Position:    "Similasan Corporation",
			Telephone:   "989 991 9056",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Asa Clifft",
			Nickname:    "Asa",
			ImageURL:    "https://robohash.org/dignissimosutdolorem.png?size=200x200&set=set1",
			Color:       "#14b3d8",
			Email:       "aclifftn@vkontakte.ru",
			Description: "Sprain of left sternoclavicular joint, initial encounter",
			Position:    "Golden State Medical Supply, Inc.",
			Telephone:   "102 625 2471",
		},
		{
			Id:          uuid.New(),
			CalendarId:  initCalendar.Id,
			Name:        "Sallee Bottrell",
			Nickname:    "Sallee",
			ImageURL:    "https://robohash.org/illoinerror.png?size=200x200&set=set1",
			Color:       "#89d395",
			Email:       "sbottrello@engadget.com",
			Description: "Poisn by unsp prim sys and hematolog agent, slf-hrm, sequela",
			Position:    "Advanced Beauty Systems, Inc.",
			Telephone:   "130 615 2527",
		},
	}

	initSchedule := []*dao.Schedules{
		{
			Id:                   uuid.New(),
			CalendarId:           initCalendar.Id,
			Name:                 "setting_first_schedule 0700 - 1000 ",
			Description:          "test setting",
			Priority:             1,
			Start:                time.Time{},
			End:                  time.Time{},
			Hr_start:             "07:00",
			Hr_end:               "10:00",
			BreakTime:            2 * 60 * 60,
			Tzid:                 "Asia/Bangkok",
			Recurrence_freq:      int8(rrule.DAILY),
			Recurrence_interval:  1,
			Recurrence_wkst:      "",
			Recurrence_bymonth:   "",
			Recurrence_byweekday: "0,1,2,3,4",
		},
		{
			Id:                   uuid.New(),
			CalendarId:           initCalendar.Id,
			Name:                 "setting_first_schedule 1000 - 1400 ",
			Description:          "test setting",
			Priority:             1,
			Start:                time.Time{},
			End:                  time.Time{},
			Hr_start:             "11:00",
			Hr_end:               "14:00",
			BreakTime:            2 * 60 * 60,
			Tzid:                 "Asia/Bangkok",
			Recurrence_freq:      int8(rrule.DAILY),
			Recurrence_interval:  1,
			Recurrence_wkst:      "",
			Recurrence_bymonth:   "",
			Recurrence_byweekday: "0,1,2,3,4",
		},
		{
			Id:                   uuid.New(),
			CalendarId:           initCalendar.Id,
			Name:                 "setting_first_schedule 1400 - 1800 ",
			Description:          "test setting",
			Priority:             1,
			Start:                time.Time{},
			End:                  time.Time{},
			Hr_start:             "14:00",
			Hr_end:               "18:00",
			BreakTime:            2 * 60 * 60,
			Tzid:                 "Asia/Bangkok",
			Recurrence_freq:      int8(rrule.DAILY),
			Recurrence_interval:  1,
			Recurrence_wkst:      "",
			Recurrence_bymonth:   "",
			Recurrence_byweekday: "0,1,2,3,4",
		},
	}

	initResponsible := make([]*dao.Responsible, 0)
	for i := 0; i < len(initSchedule); i++ {
		for j := 0; j < len(initMembers); j++ {
			initResponsible = append(initResponsible, &dao.Responsible{ScheduleId: initSchedule[i].Id, MemberId: initMembers[j].Id, Queue: int8(j + 1)})
		}
	}

	result_created_user := db.Save(initUser)

	if result_created_user.Error != nil {
		return result_created_user.Error
	}

	result_created_calendar := db.Save(initCalendar)

	if result_created_calendar.Error != nil {
		return result_created_user.Error
	}

	result_created_members := db.Save(initMembers)

	if result_created_members.Error != nil {
		return result_created_members.Error
	}

	result_created_schedule := db.Save(initSchedule)

	if result_created_schedule.Error != nil {
		return result_created_schedule.Error
	}

	result_created_responsible := db.Save(initResponsible)

	if result_created_responsible.Error != nil {
		return result_created_responsible.Error
	}

	return nil
}
