package mock_data

import (
	"fmt"
	"math/rand"

	"github.com/ceciivanov/platform-go-challenge/internal/models"
)

// Age groups
const (
	AgeGroupTeen       string = "0-17"
	AgeGroupYoungAdult string = "18-25"
	AgeGroupAdult      string = "26-40"
	AgeGroupMiddleAged string = "41-65"
	AgeGroupSenior     string = "66+"
)

// Genders
const (
	Male      string = "Male"
	Female    string = "Female"
	NonBinary string = "Non-Binary"
)

// GetRandomNumber returns a random number in the [0, n] range
func GetRandomNumber(max int) uint {
	return uint(rand.Intn(max))
}

// GetRandomPoints generates a random number of random points
func GetRandomPoints(MinPoints, MaxPoints int) []models.Point {
	numPoints := rand.Intn(MaxPoints-MinPoints+1) + MinPoints
	points := make([]models.Point, numPoints)
	for i := 0; i < numPoints; i++ {
		points[i] = models.Point{
			X: rand.Float32(), // Adjust range as needed
			Y: rand.Float32(), // Adjust range as needed
		}
	}
	return points
}

// GetRandomAgeGroup returns a random age group
func GetRandomAgeGroup() string {
	ageGroups := []string{AgeGroupTeen, AgeGroupAdult, AgeGroupSenior}
	return ageGroups[rand.Intn(len(ageGroups))]
}

// GetRandomGender returns a random gender
func GetRandomGender() string {
	genders := []string{Male, Female, NonBinary}
	return genders[rand.Intn(len(genders))]
}

// GetRandomCountry returns a random country from the predefined list
func GetRandomCountry() string {
	countries := []string{
		"Afghanistan", "Albania", "Algeria", "Andorra", "Angola", "Argentina", "Armenia", "Australia", "Austria", "Azerbaijan",
		"Bahamas", "Bahrain", "Bangladesh", "Barbados", "Belarus", "Belgium", "Belize", "Benin", "Bhutan", "Bolivia",
		"Bosnia and Herzegovina", "Botswana", "Brazil", "Brunei", "Bulgaria", "Burkina Faso", "Burundi", "Cambodia", "Cameroon", "Canada",
		"Cape Verde", "Central African Republic", "Chad", "Chile", "China", "Colombia", "Comoros", "Congo", "Costa Rica", "Croatia",
		"Cuba", "Cyprus", "Czech Republic", "Denmark", "Djibouti", "Dominica", "Dominican Republic", "Ecuador", "Egypt", "El Salvador",
		"Equatorial Guinea", "Eritrea", "Estonia", "Eswatini", "Ethiopia", "Fiji", "Finland", "France", "Gabon", "Gambia",
		"Georgia", "Germany", "Ghana", "Greece", "Grenada", "Guatemala", "Guinea", "Guinea-Bissau", "Guyana", "Haiti",
		"Honduras", "Hungary", "Iceland", "India", "Indonesia", "Iran", "Iraq", "Ireland", "Israel", "Italy",
		"Jamaica", "Japan", "Jordan", "Kazakhstan", "Kenya", "Kiribati", "Kuwait", "Kyrgyzstan", "Laos", "Latvia",
		"Lebanon", "Lesotho", "Liberia", "Libya", "Liechtenstein", "Lithuania", "Luxembourg", "Madagascar", "Malawi", "Malaysia",
		"Maldives", "Mali", "Malta", "Marshall Islands", "Mauritania", "Mauritius", "Mexico", "Micronesia", "Moldova", "Monaco",
		"Mongolia", "Montenegro", "Morocco", "Mozambique", "Myanmar", "Namibia", "Nauru", "Nepal", "Netherlands", "New Zealand",
		"Nicaragua", "Niger", "Nigeria", "North Macedonia", "Norway", "Oman", "Pakistan", "Palau", "Panama", "Papua New Guinea",
		"Paraguay", "Peru", "Philippines", "Poland", "Portugal", "Qatar", "Romania", "Russia", "Rwanda", "Saint Kitts and Nevis",
		"Saint Lucia", "Saint Vincent and the Grenadines", "Samoa", "San Marino", "Sao Tome and Principe", "Saudi Arabia", "Senegal", "Serbia", "Seychelles", "Sierra Leone",
		"Singapore", "Slovakia", "Slovenia", "Solomon Islands", "Somalia", "South Africa", "South Sudan", "Spain", "Sri Lanka", "Sudan",
		"Suriname", "Sweden", "Switzerland", "Syria", "Taiwan", "Tajikistan", "Tanzania", "Thailand", "Timor-Leste", "Togo",
		"Tonga", "Trinidad and Tobago", "Tunisia", "Turkey", "Turkmenistan", "Tuvalu", "Uganda", "Ukraine", "United Arab Emirates", "United Kingdom",
		"United States", "Uruguay", "Uzbekistan", "Vanuatu", "Vatican City", "Venezuela", "Vietnam", "Yemen", "Zambia", "Zimbabwe",
	}
	return countries[rand.Intn(len(countries))]
}

// GenerateMockData generates mock data for users and assets and returns a map of users
func GenerateMockData(NumberOfUsers, NumberOfAssets int) map[int]models.User {
	Users := make(map[int]models.User)

	for i := 1; i <= NumberOfUsers; i++ {
		userID := i
		user := models.User{
			ID:         userID,
			Favourites: make(map[int]models.Asset),
		}

		for j := 1; j <= NumberOfAssets; j++ {
			assetID := j

			// Randomly choose an asset type
			assetType := j % 3

			var asset models.Asset
			switch assetType {
			case 0:
				asset = &models.Chart{
					ID:          assetID,
					Type:        models.ChartType,
					Description: "Sample Chart for GWI",
					Title:       fmt.Sprintf("GWI Chart %d", j),
					XAxesTitle:  "X-Axis",
					YAxesTitle:  "Y-Axis",
					DataPoints:  GetRandomPoints(1, 5),
				}
			case 1:
				asset = &models.Insight{
					ID:          assetID,
					Type:        models.InsightType,
					Description: "Sample Insight for GWI",
					Text:        fmt.Sprintf("GWI Insight %d", j),
				}
			case 2:
				asset = &models.Audience{
					ID:                assetID,
					Type:              models.AudienceType,
					Description:       "Sample Audience for GWI",
					Age:               GetRandomNumber(100),
					AgeGroup:          GetRandomAgeGroup(),
					Gender:            GetRandomGender(),
					BirthCountry:      GetRandomCountry(),
					HoursSpentOnMedia: GetRandomNumber(100),
					NumberOfPurchases: GetRandomNumber(100),
				}
			}
			// Add the asset to the user's favourites
			user.Favourites[assetID] = asset
		}
		// Update the user in the Users map
		Users[userID] = user
	}

	return Users
}
