package emails

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
)

type Spec struct {
	Processor   string `json:"processor"`
	Motherboard string `json:"motherBoard"`
	Ram         string `json:"ram"`
	GPU         string `json:"gpu"`
	SMPS        string `json:"smps"`
	Case        string `json:"case"`
	Cooler      string `json:"cooler"`
	Storage     string `json:"storage"`
}
type EmailDetails struct {
	ID             int    `json:"id"`
	Mail           string `json:"email"`
	Name           string `json:"name"`
	PhoneNumber    string `json:"phone"`
	Specifications Spec   `json:"specifications"`
}

func EmailSending(data EmailDetails) (string, error) {
	fmt.Printf("mail %v", data.Mail)
	fmt.Printf("name %v", data.Name)
	fmt.Printf("phone %v", data.PhoneNumber)
	fmt.Printf("spec %v", data.Specifications)
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	supportEmail := "sheik.mohammed@holoware.co"

	if smtpUsername == "" || smtpPassword == "" || smtpHost == "" || smtpPort == "" {
		return "SMTP credentials not set", fmt.Errorf("SMTP credentials missing in environment variables")
	}
	fmt.Printf("Email: %v\nPhone: %v\nName: %v\n", data.Mail, data.PhoneNumber, data.Name)

	apiKey := os.Getenv("HUNTER_API_KEY")
	isPhoneValid := validateNumber(data.PhoneNumber)
	isValidMail, err := validateEmail(data.Mail, apiKey)
	if err != nil {
		return "", fmt.Errorf("error validating email: %v", err)
	}
	if !isValidMail || !isPhoneValid {
		return "", fmt.Errorf("invalid email address or phone number, please check and try again")
	}
	//Initialize DB Connection
	specJson, err := json.Marshal(data.Specifications)
	if err != nil {
		return "", fmt.Errorf("error marshaling specifications: %v", err)
	}
	res1, err := StoringDetailsDB(data.Mail, data.Name, data.PhoneNumber, string(specJson))
	if err != nil {
		return "", fmt.Errorf("error storing order details: %v", err)
	}
	fmt.Print("Order Details Stored Successfully", res1)
	// Generate the email template
	body, err := generateEmailTemplate(data)
	if err != nil {
		return "", fmt.Errorf("error generating email template: %v", err)
	}
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Order Details\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n", smtpUsername, supportEmail)
	message := []byte(headers + "\r\n" + body)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	recipients := []string{supportEmail}

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUsername, recipients, message)
	if err != nil {
		return "", fmt.Errorf("error sending email: %v", err)
	}

	res, err := ConformationMail(data)
	if err != nil {
		return "", fmt.Errorf("error sending confirmation email: %v", err)
	}
	fmt.Println("Email sent successfully!")
	return res, nil
} // Generating Email Template
func generateEmailTemplate(data EmailDetails) (string, error) {
	const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
        }
        .container {
            width: 100%;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            border: 1px solid #ddd;
            border-radius: 8px;
            background-color: #f9f9f9;
        }
        .header {
            text-align: center;
            margin-bottom: 20px;
        }
        .header h2 {
            color: #333;
        }
        .section {
            margin-bottom: 20px;
        }
        .section h3 {
            color: #555;
            margin-bottom: 10px;
            border-bottom: 1px solid #ddd;
            padding-bottom: 5px;
        }
        .section p {
            margin: 5px 0;
        }
        .spec-table {
            width: 100%;
            border-collapse: collapse;
        }
        .spec-table th, .spec-table td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        .spec-table th {
            background-color: #f4f4f4;
            color: #333;
        }
        .footer {
            text-align: center;
            margin-top: 20px;
            font-size: 0.9em;
            color: #777;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>New Order Placement Notification</h2>
        </div>

        <div class="section">
            <h3>Customer Details</h3>
            <p><strong>Name:</strong> {{.Name}}</p>
            <p><strong>Email:</strong> {{.Mail}}</p>
            <p><strong>Phone Number:</strong> {{.PhoneNumber}}</p>
        </div>

        <div class="section">
            <h3>Order Specifications</h3>
            <table class="spec-table">
                <tr>
                    <th>Component</th>
                    <th>Details</th>
                </tr>
                <tr>
                    <td>Processor</td>
                    <td>{{.Specifications.Processor}}</td>
                </tr>
                <tr>
                    <td>Motherboard</td>
                    <td>{{.Specifications.Motherboard}}</td>
                </tr>
                <tr>
                    <td>RAM</td>
                    <td>{{.Specifications.Ram}}</td>
                </tr>
                <tr>
                    <td>GPU</td>
                    <td>{{.Specifications.GPU}}</td>
                </tr>
                <tr>
                    <td>SMPS</td>
                    <td>{{.Specifications.SMPS}}</td>
                </tr>
                <tr>
                    <td>Case</td>
                    <td>{{.Specifications.Case}}</td>
                </tr>
                <tr>
                    <td>Cooler</td>
                    <td>{{.Specifications.Cooler}}</td>
                </tr>
                <tr>
                    <td>Storage</td>
                    <td>{{.Specifications.Storage}}</td>
                </tr>
            </table>
        </div>

        <div class="center-text">
            <p>Please process this order promptly and contact the customer if additional information is required.</p>
            <p>Best regards,<br>Holoware CustomizedCart System Notification</p>
        </div>
    </div>
</body>
</html>`

	tmpl, err := template.New("emailTemplate").Parse(emailTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var content bytes.Buffer
	err = tmpl.Execute(&content, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return content.String(), nil
}

// Email Validation
type EmailValidationResponse struct {
	Data struct {
		Status string `json:"status"`
		Score  int    `json:"score"`
	} `json:"data"`
}

func validateEmail(email, apiKey string) (bool, error) {
	url := fmt.Sprintf("https://api.hunter.io/v2/email-verifier?email=%s&api_key=%s", email, apiKey)
	fmt.Printf("Email%v", email)
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to validate email: HTTP %d", resp.StatusCode)
	}

	var result EmailValidationResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return false, fmt.Errorf("failed to decode response: %v", err)
	}

	fmt.Printf("Email Status: %s, Score: %d\n", result.Data.Status, result.Data.Score)
	if result.Data.Status == "valid" || result.Data.Status == "accept_all" {
		return true, nil
	}

	return false, err
}

func validateNumber(PhoneNumber string) bool {
	// Define the regex pattern
	pattern := `^(\+91[\-\s]?|91[\-\s]?|0)?[6-9]\d{9}$`
	re := regexp.MustCompile(pattern)
	ra := re.MatchString(PhoneNumber)
	fmt.Printf("Phone Number %v", ra)
	return ra
}

// Email function
func EmailUserDetails() (EmailDetails, error) {
	var data EmailDetails
	var spec string
	query := `SELECT * FROM orderdetails`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("executing an error %v", err)
		return EmailDetails{}, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&data.Name, &data.Mail, &data.PhoneNumber, &spec)
		if err != nil {
			log.Printf("executing an error %v", err)
			return EmailDetails{}, err
		}
		err = json.Unmarshal([]byte(spec), &data.Specifications)
		if err != nil {
			log.Printf("executing an error %v", err)
			return EmailDetails{}, err
		}
	}
	return data, nil
}

// DB Connection Storing and retrying details
var db *sql.DB

func StoringDetailsDB(Mail, Name, PhoneNumber, Specification string) (string, error) {
	var data EmailDetails
	//DB initialization
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Initialize the database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", fmt.Errorf("error connecting to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return "", fmt.Errorf("error connecting to database: %v", err)
	}
	//Storing the Order Details in database
	query := `INSERT INTO orderdetails (UserName, UserMail, UserNumber, OrderSpecification) VALUES (?, ?, ?, ?)`
	if err != nil {
		return "", fmt.Errorf("error marshaling specifications: %v", err)
	}
	res, err := db.Exec(query, data.Name, data.Mail, data.PhoneNumber, Specification)
	if err != nil {
		return "", fmt.Errorf("error storing order details: %v", err)
	}
	fmt.Printf("Order Details Stored Successfully%v", res)
	return "", nil
}

// Function getting user details
func GetUserDetails() (EmailDetails, error) {
	var data EmailDetails
	var spec string
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return data, fmt.Errorf("error connecting to database: %v", err)
	}
	defer db.Close()
	query := `SELECT id, UserName, UserMail, UserNumber, OrderSpecification FROM orderdetails`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("executing an error %v", err)
		return EmailDetails{}, err
	}
	for rows.Next() {
		err := rows.Scan(&data.ID, &data.Name, &data.Mail, &data.PhoneNumber, &spec)
		if err != nil {
			log.Printf("executing an error %v", err)
			return EmailDetails{}, err
		}
		err = json.Unmarshal([]byte(spec), &data.Specifications)
		if err != nil {
			log.Printf("executing an error n unMarshall function %v", err)
			return data, err
		}
	}
	return data, nil
}
