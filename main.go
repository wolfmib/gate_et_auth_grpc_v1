package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/sirupsen/logrus"

	pb "github.com/wolfmib/user_grpc_v1/user_proto"
	"io/ioutil"
	"encoding/json"
	"google.golang.org/grpc"
	"context"
	"reflect"

	



)

const (
	SecretKey = "ja_jwt_gate"
	address = "localhost:5001"
	defaultFilename = "user.json"

)



func chk_json_req_match_grpc_req( pb_request interface{}) error {
	
	logrus.Warn("[checking_json_req_match_grpc_req]:  Put me in the ja_golang_grpc package")
	data, err := ioutil.ReadFile("user.json")
	
	err = json.Unmarshal(data,&pb_request)
	if err != nil{
		logrus.Error("Data is not matched:   data:  ",data, "pb_request:   ",pb_request, "\n")
	} else{
		logrus.Info("Mached Data: ", data)
		logrus.Info("  pb_request:", pb_request)
	}	
	return nil
}




func register_parseFile(file string) (*pb.RegisterRequest, error){
	var register_request *pb.RegisterRequest
	data, err := ioutil.ReadFile(file)

	if err != nil{
		return nil, err
	}

	json.Unmarshal(data, &register_request)
	return register_request, err
}


// Save Time to check err all the time
func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func StartServer() {

	http.HandleFunc("/login", LoginHandler)

	// Make Register endpoit to access user_server via grpc
	http.HandleFunc("/register",RegisterHandler)

	http.Handle("/resource", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))

	logrus.Info("Now listening...")
	logrus.Warn("Now listening...")
	


	logrus.Info("\nPOST:  http://localhost:8080/login  with body")
	logrus.Info("       { \"username\": \"johnny\" " )
	logrus.Info("        \"password\": \"P@sswordXXXX\" } ")
	logrus.Info("-------------------------------------\n")
	logrus.Info(" GET    http://localhost:8080/resource with head")
	logrus.Info("       authorization:     YOUR_TOKEN_STRING ")
	logrus.Info("--------------------------------------------\n")
	logrus.Info(" POST  http://localhost:8080/register with body")
	logrus.Info("       { \"first_name\": \"Mary\" " )
	logrus.Info("        \"family_name\": \"Jean\" ")
	logrus.Info("        \"email\": \"MaryJean@ja.com\" } ")

	http.ListenAndServe(":8080", nil)
	

}

func main() {


	StartServer()


}


func RegisterHandler(w http.ResponseWriter, r *http.Request){
	logrus.Info("Get the Register Enpoint")

	var register_request pb.RegisterRequest


	logrus.Info(r)

	// Checking the Requst data match the protobuffer format or not
	err := json.NewDecoder(r.Body).Decode(&register_request)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		logrus.Error(err)
		return 
	}


	// Dial
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil{
		logrus.Error("Didn't oconnect with err: ",err)
	}

	defer conn.Close()
	client := pb.NewUserServiceClient(conn)


	logrus.Info("..................Sending the grpc_api_request...............")
	logrus.Info(register_request)
	logrus.Warn("Type: ",reflect.TypeOf(register_request))
	logrus.Info(".................................................................\n")

	// Call API
	res, err := client.RegisterApi(context.Background(), &register_request)
	if err != nil{
		logrus.Error("Could not register ",err)
	}

	logrus.Info("----------Finish Register Handler------------")
	logrus.Info(res)
	logrus.Info("---------------------------------------------")

}



func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var user UserCredentials

	//Checking the post data match or not
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	// Implement MongoDB Here
	if strings.ToLower(user.Username) != "johnny" {
		if user.Password != "P@sswordXXXX" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Invalid credentials")
			return
		}
	}

	// Creer TOKEN
	// Declare Sign Method avec HS256
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)

	// Expire Time Setting
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()

	// Issue At time Setting
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error extracting the key")
		fatal(err)
	}

	// Signed by the ScretKey
		// var jwtKey = []byte("my_secret_key")
	    // tokenString, err := token.SignedString(jwtKey)
	// Johnny: Try code , b/a, signing process of tokenString and token
	logrus.Info("Before Signing...")
	fmt.Printf("%+v",token)

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}

	logrus.Info("\nAfter Signing...\n")
	fmt.Printf("%+v",token)




	response := Token{tokenString}
	JsonResponse(response, w)

}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//Check the token is valid or not
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			fmt.Println("Inside key func:")
			fmt.Println("--------------")
			claims := token.Claims.(jwt.MapClaims)
		    fmt.Printf("Token for user %v expires %v", claims["user"], claims["exp"])
			fmt.Printf("%+v",token)
			fmt.Println("--------------")
			return []byte(SecretKey), nil
		})

	fmt.Println("Outside key func: inside Middleware:")
	fmt.Println("--------------")
	fmt.Printf("%+v",token)
	claims := token.Claims.(jwt.MapClaims)
	fmt.Printf("Token for user %v expires %v", claims["user"], claims["exp"])
	fmt.Println("--------------")


	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}

}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}



/*
Origin From：https://blog.csdn.net/wangshubo1989/article/details/74529333
*/