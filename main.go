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
	//mysqlRepo := mysql.NewDB()
	//
	//if status, err := mysqlRepo.IsPhoneNumberUnique("0918"); err != nil {
	//	fmt.Println(err)
	//} else {
	//	if status {
	//		fmt.Println("user is have uniq phone")
	//	} else {
	//		fmt.Println("user is not have uniq phone")
	//	}
	//}
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

		case "phone number is not unique":
			res.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(res, "phone number is not unique")

		default:
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(res, RegisterUserError.Error())
		}

	}

	res.WriteHeader(http.StatusCreated)
	fmt.Fprint(res, respons.User)
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
	var bodyData category_service.CreateCategoryRequestStruct
	MarshalError := json.Unmarshal(bData, &bodyData)
	if MarshalError != nil {
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, MarshalError.Error())
		return
	}

	mysqlRepo := mysql.NewDB()
	categoryRepo := category_service.CategoryService{
		Repo: mysqlRepo,
	}

	resultCreateCategory, ErrorCreateCategory := categoryRepo.CreateCategory(bodyData)
	if ErrorCreateCategory != nil {
		fmt.Println(ErrorCreateCategory)
		res.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(res, ErrorCreateCategory.Error())
	}

	response := struct {
		ID          uint
		Title       string
		Description string
	}{
		ID:          resultCreateCategory.Category.ID,
		Title:       resultCreateCategory.Category.Title,
		Description: resultCreateCategory.Category.Description,
	}

	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(response)
	return
}
