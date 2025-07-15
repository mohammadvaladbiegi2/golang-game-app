package main

import (
	"encoding/json"
	"fmt"
	"gamegolang/repository/mysql"
	"gamegolang/service/category_service"
	userservice "gamegolang/service/user_service"
	"io"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/sign-up", RegisterUser)
	mux.HandleFunc("/create-category", CreateCategory)
	fmt.Println("Server starting on port 3000...")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func RegisterUser(res http.ResponseWriter, req *http.Request) {
	errorResponse := map[string]string{
		"message": "Invalid HTTP method. Use POST for this API",
	}

	jsonData, _ := json.Marshal(errorResponse)
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(res, string(jsonData))
		return
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(
			fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		))
		return
	}
	defer req.Body.Close()
	var UserRequestStruct userservice.RegisterRequest
	MarshalError := json.Unmarshal(data, &UserRequestStruct)
	if MarshalError != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, "server cant read the data")
		return
	}
	mysqlRepo := mysql.NewDB()
	userRepo := userservice.Service{Repo: mysqlRepo}

	respons, RegisterUserError := userRepo.Register(UserRequestStruct)
	if RegisterUserError != nil {

		switch RegisterUserError.Error() {
		case "phone number is not valid":
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(res, "phone number is not valid")
			return

		case "phone number is not unique":
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(res, "phone number is not unique")
			return

		case "password length should be greater than 8":
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(res, "password length should be greater than 8")
			return

		case "name length should be greater than 3":
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(res, "name length should be greater than 3")
			return

		default:
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(res, RegisterUserError.Error())
			return
		}

	}

	response := struct {
		ID          uint
		Name        string
		PhoneNumber string
		Password    string
	}{
		ID:          respons.ID,
		Name:        respons.Name,
		PhoneNumber: respons.PhoneNumber,
		Password:    "*********",
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(response)
	return

}

func CreateCategory(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{
		"message": "Invalid HTTP method. Use POST for this API",
	}
	jsonData, _ := json.Marshal(errorResponse)
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(res, string(jsonData))
		return
	}

	bData, Rerror := io.ReadAll(req.Body)

	if Rerror != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, Rerror)
		return
	}
	defer req.Body.Close()
	var bodyData categoryservice.CreateRequest
	MarshalError := json.Unmarshal(bData, &bodyData)
	if MarshalError != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, MarshalError.Error())
		return
	}

	mysqlRepo := mysql.NewDB()
	categoryRepo := categoryservice.Service{
		Repo: mysqlRepo,
	}

	resultCreateCategory, ErrorCreateCategory := categoryRepo.Create(bodyData)
	if ErrorCreateCategory != nil {
		fmt.Println(ErrorCreateCategory)
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, ErrorCreateCategory.Error())
		return
	}

	response := struct {
		ID          uint
		Title       string
		Description string
	}{
		ID:          resultCreateCategory.ID,
		Title:       resultCreateCategory.Title,
		Description: resultCreateCategory.Description,
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(response)
	return
}
