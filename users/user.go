package users

import (
	"context"
	"fmt"
	"time"

	"github.com/danmuck/dps_lib/mongo_client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var TMP_STORAGE *mongo_client.MongoClient
var TMP_USER_DB = "users"
var SECRET = "temp-token-signature"

func init() {
	// Initialize the temporary storage client
	var err error
	TMP_STORAGE, err = mongo_client.NewMongoStore(
		"mongodb://localhost:27017", // Replace with your MongoDB URI
		TMP_USER_DB,                 // Database name for temporary user storage
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}

}

type ContactInfo struct {
	Email    string `json:"email" bson:"email"`
	Phone    string `json:"phone,omitempty" bson:"phone,omitempty"`
	Location string `json:"location,omitempty" bson:"location,omitempty"`
	Website  string `json:"website,omitempty" bson:"website,omitempty"`
	LinkedIn string `json:"linkedin,omitempty" bson:"linkedin,omitempty"`
	GitHub   string `json:"github,omitempty" bson:"github,omitempty"`
}

type ProjectInfo struct {
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description,omitempty" bson:"description,omitempty"`
	Highlights  string   `json:"highlights,omitempty" bson:"highlights,omitempty"`
	Role        string   `json:"role,omitempty" bson:"role,omitempty"`
	URL         string   `json:"url,omitempty" bson:"url,omitempty"`
	TechStack   []string `json:"tech_stack,omitempty" bson:"tech_stack,omitempty"`
}

type ProfileInfo struct {
	Bio       string             `json:"bio,omitempty" bson:"bio,omitempty"`
	AvatarURL string             `json:"avatar_url,omitempty" bson:"avatar_url,omitempty"`
	CreatedAt primitive.DateTime `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type CareerInfo struct {
	CurrentRole string   `json:"current_role,omitempty" bson:"current_role,omitempty"`
	Skills      []string `json:"skills,omitempty" bson:"skills,omitempty"`
	Experience  string   `json:"experience,omitempty" bson:"experience,omitempty"`
	Education   string   `json:"education,omitempty" bson:"education,omitempty"`
	Interests   []string `json:"interests,omitempty" bson:"interests,omitempty"`
}

func (c *ContactInfo) String() string {
	return fmt.Sprintf(`
	Email: %s
	Phone: %s
	Location: %s
	Website: %s
	LinkedIn: %s
	GitHub: %s`,
		c.Email, c.Phone, c.Location, c.Website, c.LinkedIn, c.GitHub)
}
func (p *ProjectInfo) String() string {
	return fmt.Sprintf(`
	Name: %s
	Description: %s
	Highlights: %s
	Role: %s
	URL: %s
	TechStack: %v`,
		p.Name, p.Description, p.Highlights, p.Role, p.URL, p.TechStack)
}
func (p *ProfileInfo) String() string {
	return fmt.Sprintf(`
	Bio: %s
	AvatarURL: %s
	CreatedAt: %s
	UpdatedAt: %s`,
		p.Bio, p.AvatarURL, p.CreatedAt.Time(), p.UpdatedAt.Time())
}
func (c *CareerInfo) String() string {
	return fmt.Sprintf(`
	CurrentRole: %s
	Skills: %v
	Experience: %s
	Education: %s
	Interests: %v`,
		c.CurrentRole, c.Skills, c.Experience, c.Education, c.Interests)
}

// User represents an authenticated user
type User struct {
	ID           primitive.ObjectID `bson:"id,omitempty" json:"id"`
	Username     string             `bson:"username" json:"username"`
	Email        string             `bson:"email" json:"email"`
	PasswordHash string             `bson:"password_hash" json:"-"`
	Token        string             `bson:"token" json:"-"`
	Roles        []string           `bson:"roles" json:"roles"`

	Bio       string             `bson:"bio,omitempty" json:"bio,omitempty"`
	AvatarURL string             `bson:"avatar_url,omitempty" json:"avatar_url,omitempty"`
	CreatedAt primitive.DateTime `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at,omitempty"`

	ContactInfo ContactInfo   `bson:"contact_info,omitempty" json:"contact_info,omitempty"`
	Projects    []ProjectInfo `bson:"projects,omitempty" json:"projects,omitempty"`
	ProfileInfo ProfileInfo   `bson:"profile_info,omitempty" json:"profile_info,omitempty"`
	CareerInfo  CareerInfo    `bson:"career_info,omitempty" json:"career_info"`
}

// var db = "usersv1"
// var bucket storage.Bucket

func (u *User) String() string {
	var token string = "[no token]"
	if len(u.Token) > 20 {
		token = u.Token[:10] + " ... " + u.Token[len(u.Token)-10:]
	}
	return fmt.Sprintf(`
	User: %s
	Password: %s
	Email: %s
	Roles: %v
	Token: [%s]

	CreatedAt: %s
	UpdatedAt: %s
	Bio: %s
	AvatarURL: %s

	ContactInfo: %s
	Projects: %+v
	ProfileInfo: %s
	CareerInfo: %s
`,
		u.Username, u.PasswordHash, u.Email, u.Roles, token,
		u.CreatedAt.Time(),
		u.UpdatedAt.Time(), u.Bio, u.AvatarURL,
		u.ContactInfo.String(),
		u.Projects,
		u.ProfileInfo.String(),
		u.CareerInfo.String())

}

func (u *User) UpdateJWTToken(token string) {
	u.Token = token
	u.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())
	// Update the user in the database
	_, err := TMP_STORAGE.Collection(TMP_USER_DB).UpdateOne(
		context.Background(),
		bson.M{"id": u.ID},
		bson.M{"$set": bson.M{"token": u.Token, "updated_at": u.UpdatedAt}},
	)
	if err != nil {
		fmt.Printf("Failed to update user token: %v\n", err)
	}
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashed, password string) bool {
	check := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return check == nil
}
