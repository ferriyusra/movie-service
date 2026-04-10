package platform

import (
	"fmt"
	"log"
	"time"

	"github.com/ferriyusra/movie-service/internal/model/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedMoviesAndTheaters seeds sample movies, genres, movie_genres, theaters with seats, and showtimes.
func SeedMoviesAndTheaters(db *gorm.DB) error {
	if err := seedGenresAndMovies(db); err != nil {
		return err
	}
	if err := seedTheaters(db); err != nil {
		return err
	}
	if err := seedShowtimes(db); err != nil {
		return err
	}
	return nil
}

func seedGenresAndMovies(db *gorm.DB) error {
	// Check if movies already seeded
	var movieCount int64
	if err := db.Model(&entity.MovieEntity{}).Count(&movieCount).Error; err != nil {
		return err
	}
	if movieCount > 0 {
		return nil
	}

	// Create genres
	genreAction := entity.GenreEntity{ID: uuid.New(), Name: "Action", CreatedAt: time.Now()}
	genreDrama := entity.GenreEntity{ID: uuid.New(), Name: "Drama", CreatedAt: time.Now()}
	genreSciFi := entity.GenreEntity{ID: uuid.New(), Name: "Sci-Fi", CreatedAt: time.Now()}
	genreThriller := entity.GenreEntity{ID: uuid.New(), Name: "Thriller", CreatedAt: time.Now()}
	genreComedy := entity.GenreEntity{ID: uuid.New(), Name: "Comedy", CreatedAt: time.Now()}
	genreHorror := entity.GenreEntity{ID: uuid.New(), Name: "Horror", CreatedAt: time.Now()}
	genreRomance := entity.GenreEntity{ID: uuid.New(), Name: "Romance", CreatedAt: time.Now()}
	genreAnimation := entity.GenreEntity{ID: uuid.New(), Name: "Animation", CreatedAt: time.Now()}

	genres := []entity.GenreEntity{
		genreAction, genreDrama, genreSciFi, genreThriller,
		genreComedy, genreHorror, genreRomance, genreAnimation,
	}

	// Upsert genres (skip if name already exists)
	for i := range genres {
		var existing entity.GenreEntity
		result := db.Where("name = ?", genres[i].Name).First(&existing)
		if result.Error == nil {
			genres[i].ID = existing.ID
			continue
		}
		if err := db.Create(&genres[i]).Error; err != nil {
			return err
		}
	}
	// Re-assign after potential ID updates
	genreAction = genres[0]
	genreDrama = genres[1]
	genreSciFi = genres[2]
	genreThriller = genres[3]
	genreComedy = genres[4]
	genreHorror = genres[5]
	genreRomance = genres[6]
	genreAnimation = genres[7]

	log.Println("Seeded genres")

	// Create movies
	releaseDate := func(y, m, d int) *time.Time {
		t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
		return &t
	}

	movies := []struct {
		entity entity.MovieEntity
		genres []uuid.UUID
	}{
		{
			entity: entity.MovieEntity{
				ID:          uuid.New(),
				Title:       "The Dark Knight",
				Description: "When the menace known as the Joker wreaks havoc on Gotham, Batman must accept one of the greatest tests.",
				PosterURL:   "https://image.tmdb.org/t/p/w500/qJ2tW6WMUDux911BytUy5a3dJB5.jpg",
				DurationMin: 152,
				Language:    "English",
				ReleaseDate: releaseDate(2008, 7, 18),
				Rating:      9.0,
			},
			genres: []uuid.UUID{genreAction.ID, genreDrama.ID, genreThriller.ID},
		},
		{
			entity: entity.MovieEntity{
				ID:          uuid.New(),
				Title:       "Inception",
				Description: "A thief who steals corporate secrets through dream-sharing technology is given the task of planting an idea.",
				PosterURL:   "https://image.tmdb.org/t/p/w500/ljsZTbVsrQSqZgWeep2B1QiDKuh.jpg",
				DurationMin: 148,
				Language:    "English",
				ReleaseDate: releaseDate(2010, 7, 16),
				Rating:      8.8,
			},
			genres: []uuid.UUID{genreAction.ID, genreSciFi.ID, genreThriller.ID},
		},
		{
			entity: entity.MovieEntity{
				ID:          uuid.New(),
				Title:       "Interstellar",
				Description: "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival.",
				PosterURL:   "https://image.tmdb.org/t/p/w500/gEU2QniE6E77NI6lCU6MxlNBvIx.jpg",
				DurationMin: 169,
				Language:    "English",
				ReleaseDate: releaseDate(2014, 11, 7),
				Rating:      8.7,
			},
			genres: []uuid.UUID{genreSciFi.ID, genreDrama.ID},
		},
		{
			entity: entity.MovieEntity{
				ID:          uuid.New(),
				Title:       "The Shawshank Redemption",
				Description: "Two imprisoned men bond over a number of years, finding solace and eventual redemption through acts of common decency.",
				PosterURL:   "https://image.tmdb.org/t/p/w500/9cjIGRRa4bOoEkBfUsqbCv7shwV.jpg",
				DurationMin: 142,
				Language:    "English",
				ReleaseDate: releaseDate(1994, 9, 23),
				Rating:      9.3,
			},
			genres: []uuid.UUID{genreDrama.ID},
		},
		{
			entity: entity.MovieEntity{
				ID:          uuid.New(),
				Title:       "Spirited Away",
				Description: "During her family's move to the suburbs, a sulky 10-year-old girl wanders into a world ruled by gods and witches.",
				PosterURL:   "https://image.tmdb.org/t/p/w500/39wmItIWsg5sZMyRUHLkWBcuVCM.jpg",
				DurationMin: 125,
				Language:    "Japanese",
				ReleaseDate: releaseDate(2001, 7, 20),
				Rating:      8.6,
			},
			genres: []uuid.UUID{genreAnimation.ID, genreDrama.ID},
		},
		{
			entity: entity.MovieEntity{
				ID:          uuid.New(),
				Title:       "Get Out",
				Description: "A young African-American visits his white girlfriend's parents for the weekend, where his simmering uneasiness about their reception of him eventually reaches a boiling point.",
				PosterURL:   "https://image.tmdb.org/t/p/w500/tFXcEccSQMf3lfhfXKSU9iRBpa3.jpg",
				DurationMin: 104,
				Language:    "English",
				ReleaseDate: releaseDate(2017, 2, 24),
				Rating:      7.7,
			},
			genres: []uuid.UUID{genreHorror.ID, genreThriller.ID},
		},
		{
			entity: entity.MovieEntity{
				ID:          uuid.New(),
				Title:       "La La Land",
				Description: "While navigating their careers in Los Angeles, a pianist and an actress fall in love while attempting to reconcile their aspirations for the future.",
				PosterURL:   "https://image.tmdb.org/t/p/w500/uDO8zWDhfWwoFdKS4fzkUJt0Rf0.jpg",
				DurationMin: 128,
				Language:    "English",
				ReleaseDate: releaseDate(2016, 12, 9),
				Rating:      8.0,
			},
			genres: []uuid.UUID{genreComedy.ID, genreDrama.ID, genreRomance.ID},
		},
		{
			entity: entity.MovieEntity{
				ID:          uuid.New(),
				Title:       "Parasite",
				Description: "Greed and class discrimination threaten the newly formed symbiotic relationship between the wealthy Park family and the destitute Kim clan.",
				PosterURL:   "https://image.tmdb.org/t/p/w500/7IiTTgloJzvGI1TAYymCfbfl3vT.jpg",
				DurationMin: 132,
				Language:    "Korean",
				ReleaseDate: releaseDate(2019, 5, 30),
				Rating:      8.5,
			},
			genres: []uuid.UUID{genreDrama.ID, genreThriller.ID, genreComedy.ID},
		},
	}

	for _, m := range movies {
		if err := db.Create(&m.entity).Error; err != nil {
			return err
		}
		for _, genreID := range m.genres {
			mg := entity.MovieGenreEntity{
				MovieEntityID: m.entity.ID,
				GenreEntityID: genreID,
			}
			if err := db.Create(&mg).Error; err != nil {
				return err
			}
		}
	}

	log.Printf("Seeded %d movies with genre associations", len(movies))
	return nil
}

func seedTheaters(db *gorm.DB) error {
	var theaterCount int64
	if err := db.Model(&entity.TheaterEntity{}).Count(&theaterCount).Error; err != nil {
		return err
	}
	if theaterCount > 0 {
		return nil
	}

	theaters := []entity.TheaterEntity{
		{
			ID:          uuid.New(),
			Name:        "Theater 1 - Standard",
			Location:    "Ground Floor, Hall A",
			TotalSeats:  80,
			Rows:        8,
			SeatsPerRow: 10,
		},
		{
			ID:          uuid.New(),
			Name:        "Theater 2 - Premium",
			Location:    "First Floor, Hall B",
			TotalSeats:  48,
			Rows:        6,
			SeatsPerRow: 8,
		},
		{
			ID:          uuid.New(),
			Name:        "Theater 3 - IMAX",
			Location:    "Second Floor, Hall C",
			TotalSeats:  120,
			Rows:        10,
			SeatsPerRow: 12,
		},
	}

	for _, theater := range theaters {
		if err := db.Create(&theater).Error; err != nil {
			return err
		}

		// Generate seats for this theater
		seats := generateSeats(theater.ID, theater.Rows, theater.SeatsPerRow)
		for _, seat := range seats {
			if err := db.Create(&seat).Error; err != nil {
				return err
			}
		}
	}

	log.Printf("Seeded %d theaters with seats", len(theaters))
	return nil
}

func seedShowtimes(db *gorm.DB) error {
	var showtimeCount int64
	if err := db.Model(&entity.ShowtimeEntity{}).Count(&showtimeCount).Error; err != nil {
		return err
	}
	if showtimeCount > 0 {
		return nil
	}

	// Load all movies and theaters from DB
	var movies []entity.MovieEntity
	if err := db.Find(&movies).Error; err != nil {
		return err
	}
	var theaters []entity.TheaterEntity
	if err := db.Find(&theaters).Error; err != nil {
		return err
	}
	if len(movies) == 0 || len(theaters) == 0 {
		return nil
	}

	// Base date: tomorrow at 10:00 UTC
	tomorrow := time.Now().AddDate(0, 0, 1)
	baseDate := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, time.UTC)

	// Schedule showtimes across 3 days, rotating movies through theaters
	// Each theater gets multiple showtimes per day with 30-min gaps between shows
	type showtimeDef struct {
		dayOffset     int
		theaterIdx    int
		movieIdx      int
		startHour     int
		startMinute   int
		price         float64
	}

	defs := []showtimeDef{
		// Day 1
		{0, 0, 0, 10, 0, 45000},   // Theater 1: The Dark Knight - morning
		{0, 0, 1, 13, 30, 45000},   // Theater 1: Inception - afternoon
		{0, 0, 4, 17, 0, 45000},    // Theater 1: Spirited Away - evening
		{0, 1, 2, 10, 0, 65000},    // Theater 2 (Premium): Interstellar - morning
		{0, 1, 6, 14, 0, 65000},    // Theater 2 (Premium): La La Land - afternoon
		{0, 1, 7, 17, 30, 65000},   // Theater 2 (Premium): Parasite - evening
		{0, 2, 0, 11, 0, 85000},    // Theater 3 (IMAX): The Dark Knight - morning
		{0, 2, 2, 14, 30, 85000},   // Theater 3 (IMAX): Interstellar - afternoon
		{0, 2, 1, 18, 0, 85000},    // Theater 3 (IMAX): Inception - evening

		// Day 2
		{1, 0, 3, 10, 0, 45000},    // Theater 1: The Shawshank Redemption
		{1, 0, 5, 13, 0, 45000},    // Theater 1: Get Out
		{1, 0, 7, 16, 0, 45000},    // Theater 1: Parasite
		{1, 1, 1, 10, 30, 65000},   // Theater 2: Inception
		{1, 1, 4, 14, 0, 65000},    // Theater 2: Spirited Away
		{1, 2, 2, 11, 0, 85000},    // Theater 3: Interstellar
		{1, 2, 0, 15, 0, 85000},    // Theater 3: The Dark Knight
		{1, 2, 7, 18, 30, 85000},   // Theater 3: Parasite

		// Day 3
		{2, 0, 6, 10, 0, 45000},    // Theater 1: La La Land
		{2, 0, 2, 13, 30, 45000},   // Theater 1: Interstellar
		{2, 0, 5, 17, 0, 45000},    // Theater 1: Get Out
		{2, 1, 3, 10, 0, 65000},    // Theater 2: The Shawshank Redemption
		{2, 1, 0, 13, 30, 65000},   // Theater 2: The Dark Knight
		{2, 2, 1, 11, 0, 85000},    // Theater 3: Inception
		{2, 2, 4, 14, 30, 85000},   // Theater 3: Spirited Away
		{2, 2, 6, 17, 30, 85000},   // Theater 3: La La Land
	}

	var count int
	for _, d := range defs {
		if d.movieIdx >= len(movies) || d.theaterIdx >= len(theaters) {
			continue
		}
		movie := movies[d.movieIdx]
		theater := theaters[d.theaterIdx]

		startTime := baseDate.AddDate(0, 0, d.dayOffset).
			Add(time.Duration(d.startHour-10)*time.Hour + time.Duration(d.startMinute)*time.Minute)
		// end_time = start_time + duration + 15min cleaning buffer (matches CreateShowtime logic)
		endTime := startTime.Add(time.Duration(movie.DurationMin)*time.Minute + 15*time.Minute)

		showtime := entity.ShowtimeEntity{
			ID:        uuid.New(),
			MovieID:   movie.ID,
			TheaterID: theater.ID,
			StartTime: startTime,
			EndTime:   endTime,
			Price:     d.price,
		}
		if err := db.Create(&showtime).Error; err != nil {
			return err
		}
		count++
	}

	log.Printf("Seeded %d showtimes across 3 days", count)
	return nil
}

func generateSeats(theaterID uuid.UUID, rows, seatsPerRow int) []entity.SeatEntity {
	var seats []entity.SeatEntity
	for r := 0; r < rows; r++ {
		rowLetter := string(rune('A' + r))
		for s := 1; s <= seatsPerRow; s++ {
			seatType := "standard"
			// Last 2 rows are premium
			if r >= rows-2 {
				seatType = "premium"
			}
			seats = append(seats, entity.SeatEntity{
				ID:        uuid.New(),
				TheaterID: theaterID,
				Row:       rowLetter,
				Number:    s,
				Label:     fmt.Sprintf("%s%d", rowLetter, s),
				Type:      seatType,
			})
		}
	}
	return seats
}
