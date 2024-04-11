package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func main() {
	fmt.Println("Server Started.")
	var wg sync.WaitGroup
	wg.Add(2)
	go createTCPServer("8080", &wg)
	go createHTTPServer("8081", &wg)

	wg.Wait()
	fmt.Println("Server Exited.")
}

func createTCPServer(port string, wg *sync.WaitGroup) {
	defer wg.Done()
	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		fmt.Println("Starting TCP Server...")
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("error on connection", err)
			continue
		}

		fmt.Println("\nhandling tcp request....  ")
		go handleConnection(conn)
	}

	defer fmt.Println("Killing TCP Server...")
}

func createHTTPServer(port string, wg *sync.WaitGroup) {
	balance := float64(200.00)
	http.HandleFunc("/authorization", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Processing HTTP Request...")
		decoder := json.NewDecoder(r.Body)

		var data map[string]interface{}
		err := decoder.Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Now you can use the data map
		fmt.Println("request data", data)
		val, found := data["amount"]
		if found {
			strVal := val.(string)

			floatVal, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				fmt.Println("fail to parse float = ", err)
				return
			}

			fmt.Println("float val", floatVal)
			if balance < floatVal {
				// Set the Content-Type header to application/json
				http.Error(w, "insufficient balance", http.StatusUnauthorized)
			} else {
				balance = balance - floatVal
			}
		}

		fmt.Println("Returning HTTP Request...")
		// Create a response map
		response := map[string]string{
			"message": "Data received successfully",
			"balance": fmt.Sprintf("%.2f", balance),
		}

		// Set the Content-Type header to application/json
		w.Header().Set("Content-Type", "application/json")

		// Write the response map as JSON to the ResponseWriter
		encoder := json.NewEncoder(w)
		err = encoder.Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("Starting HTTP Server...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Killing HTTP Server..., error = ", err)
		wg.Done()
	}

	fmt.Println("Killing HTTP Server...")
	wg.Done()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 128)
	count := 1
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Received TCP Data : %s", buf[:n])
		resp, err := sendToHttpServer(string(buf[:]))
		approved := true
		if err != nil {
			approved = false
			conn.Write([]byte(fmt.Sprintf("transaction error = %v \n", err)))
		} else {
			fmt.Println("current balance = ", resp.Balance)
			conn.Write([]byte(fmt.Sprintf("transaction approved = %v , balance = %v \n", approved, resp.Balance)))
		}

		count++
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}

type AuthPayload struct {
	AccountID string `json:"account_id"`
	Merchant  string `json:"merchant"`
	Amount    string `json:"amount"`
}

type AuthResponse struct {
	Balance string `json:"balance"`
	Error   string `json:"error"`
}

func payloadToJson(payload string) ([]byte, error) {
	reqPayload := AuthPayload{}
	payload = strings.ReplaceAll(payload, "\u0000", "")
	queryParams := strings.Split(payload, "&")

	for _, qp := range queryParams {
		query := strings.Split(qp, "=")
		fmt.Println("query value  =", query[1])
		trimmed := strings.TrimSpace(query[1])
		if query[0] == "account" {
			reqPayload.AccountID = trimmed
		} else if query[0] == "merchant" {
			reqPayload.Merchant = trimmed
		} else if query[0] == "amount" {
			reqPayload.Amount = trimmed
		}
	}

	reqBody, err := json.Marshal(reqPayload)
	if err != nil {
		return nil, err
	}

	return reqBody, nil
}

func sendToHttpServer(body string) (*AuthResponse, error) {

	fmt.Println("Converting TCP payload to JSON HTTP...")
	reqBody, err := payloadToJson(body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Sending HTTP Request..., body = ", string(reqBody))
	req, err := http.NewRequest("POST", "http://localhost:8081/authorization", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error from HTTP Request..., error = ", err)
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Processing HTTP Response..., body = ", string(respBody))
	var response AuthResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
