package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/rs/cors" 
    "employee/prisma/db"  
)

var client *db.PrismaClient

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Initialize Prisma Client
    client = db.NewClient()
    if err := client.Prisma.Connect(); err != nil {
        log.Fatalf("Failed to connect to Prisma: %v", err)
    }
}

func main() {
    defer func() {
        if err := client.Prisma.Disconnect(); err != nil {
            log.Fatalf("Failed to disconnect Prisma: %v", err)
        }
    }()
    router := mux.NewRouter()
    router.HandleFunc("/api/employees", createEmployeeHandler).Methods("POST")
    router.HandleFunc("/api/employees/{id}", updateEmployeeHandler).Methods("PUT")
    router.HandleFunc("/api/employees/{id}", deleteEmployeeHandler).Methods("DELETE")
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, 
        AllowCredentials: true,
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
    })
    log.Println("Server is running on port 8081...")
    if err := http.ListenAndServe(":8081", c.Handler(router)); err != nil {
        log.Fatalf("Could not start server: %v", err)
    }
}
func createEmployeeHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name     string `json:"name"`     
        Email    string `json:"email"`    
        Phone    string `json:"phone"`    
        Location string `json:"location"` 
    }
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    createdEmployee, err := client.Employee.CreateOne(
        db.Employee.Name.Set(input.Name),
        db.Employee.Email.Set(input.Email),
        db.Employee.Phone.Set(input.Phone),
        db.Employee.Location.Set(input.Location),
    ).Exec(r.Context())

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdEmployee)
}
// pagination code
func getEmployeesHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    pageStr := query.Get("page")
    limitStr := query.Get("limit")
    page := 1
    limit := 10
    if pageStr != "" {
        p, err := strconv.Atoi(pageStr)
        if err == nil && p > 0 {
            page = p
        }
    }
    if limitStr != "" {
        l, err := strconv.Atoi(limitStr)
        if err == nil && l > 0 {
            limit = l
        }
    }
    offset := (page - 1) * limit
    employees, err := client.Employee.FindMany().
        Take(limit).
        Skip(offset).
        Exec(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    var totalEmployees int
    countQuery := `SELECT COUNT(*) FROM "Employee"`
    if err := client.Prisma.QueryRaw(countQuery).Exec(r.Context(), &totalEmployees); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    response := struct {
        Data       []db.EmployeeModel `json:"data"`
        Total      int                `json:"total"`
        Page       int                `json:"page"`
        Limit      int                `json:"limit"`
        TotalPages int                `json:"totalPages"`
    }{
        Data:       employees,
        Total:      totalEmployees,
        Page:       page,
        Limit:      limit,
        TotalPages: (totalEmployees + limit - 1) / limit, 
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func updateEmployeeHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }
    var employee struct {
        Name     string `json:"name"`
        Email    string `json:"email"`
        Phone    string `json:"phone"`
        Location string `json:"location"`
    }
    if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    updatedEmployee, err := client.Employee.FindUnique(
        db.Employee.ID.Equals(id),
    ).Update(
        db.Employee.Name.Set(employee.Name),
        db.Employee.Email.Set(employee.Email),
        db.Employee.Phone.Set(employee.Phone),
        db.Employee.Location.Set(employee.Location),
    ).Exec(r.Context())
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(updatedEmployee)
}
func deleteEmployeeHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }
    _, err = client.Employee.FindUnique(
        db.Employee.ID.Equals(id),
    ).Delete().Exec(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
