package main

import (
	"encoding/json"
	"fmt"
	"gamegolang/httpserver"
	"gamegolang/repository/mysql"
	categoryservice "gamegolang/service/category_service"
	userservice "gamegolang/service/user_service"
	"io"
	"net/http"
)

func main() {

	// mux := http.NewServeMux()

	// mux.HandleFunc("/sign-up", RegisterUser)
	// mux.HandleFunc("/create-category", CreateCategory)
	// mux.HandleFunc("/login", Login)
	// fmt.Println("Server starting on port 5000...")
	// if err := http.ListenAndServe(":5000", mux); err != nil {
	// 	fmt.Printf("Server error: %v\n", err)
	// }

	httpserver.Server()
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

func Login(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{
		"message": "Invalid HTTP method. Use POST for this API",
	}
	if req.Method != http.MethodPost {
		jsonData, _ := json.Marshal(errorResponse)
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

	var bodyData userservice.LoginCredentials
	UnmarshalError := json.Unmarshal(bData, &bodyData)
	if UnmarshalError != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, UnmarshalError.Error())
		return
	}

	mysqlRepo := mysql.NewDB()
	LoginRepo := userservice.LoginService{
		Repo: mysqlRepo,
	}

	token, LoginError := LoginRepo.Login(bodyData)
	if LoginError != nil {

		switch LoginError.Error() {
		case "user not found":
			res.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(res, "user not found")
			return

		case "password or phone number does not match":
			res.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(res, "password or phone number does not match")
			return

		case "server Error":
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(res, "server Error")
			return

		default:
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(res, LoginError.Error())
			return
		}

	}

	response := fmt.Sprintf(`{"token":"%s"}`, token)

	res.WriteHeader(http.StatusOK)
	fmt.Fprint(res, response)
	return

}
