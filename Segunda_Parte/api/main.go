package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

// Estructura principal de la aplicación
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize configura la base de datos y las rutas
func (a *App) Initialize() {
	// Configuración de base de datos desde variables de entorno
	dbHost := getEnv("DB_HOST", "db")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "tasks")
	dbPort := getEnv("DB_PORT", "5432")

	// Cadena de conexión a PostgreSQL
	dbURI := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	var err error
	// Intentar conectar a la base de datos con reintentos
	for i := 0; i < 5; i++ {
		a.DB, err = gorm.Open(postgres.Open(dbURI), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Intento %d: Error al conectar a la base de datos: %v", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos después de múltiples intentos: %v", err)
	}

	// Automigración de modelos
	a.DB.AutoMigrate(&Task{})

	// Configuración del router
	a.Router = mux.NewRouter()
	a.setupRoutes()
}

// setupRoutes configura las rutas de la API
func (a *App) setupRoutes() {
	a.Router.HandleFunc("/api/tasks", a.getTasks).Methods("GET")
	a.Router.HandleFunc("/api/tasks", a.createTask).Methods("POST")
	a.Router.HandleFunc("/api/tasks/{id}", a.getTask).Methods("GET")
	a.Router.HandleFunc("/api/tasks/{id}", a.updateTask).Methods("PUT")
	a.Router.HandleFunc("/api/tasks/{id}", a.deleteTask).Methods("DELETE")

	// Middleware para CORS
	a.Router.Use(corsMiddleware)
}

// corsMiddleware maneja los encabezados CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Run inicia el servidor HTTP
func (a *App) Run(addr string) {
	log.Printf("Servidor iniciado en %s", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

// Handlers para las rutas

func (a *App) getTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []Task

	result := a.DB.Find(&tasks)
	if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (a *App) createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	task.CreatedAt = time.Now()

	if result := a.DB.Create(&task); result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func (a *App) getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var task Task
	if result := a.DB.First(&task, id); result.Error != nil {
		respondWithError(w, http.StatusNotFound, "Task not found")
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (a *App) updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var task Task
	if result := a.DB.First(&task, id); result.Error != nil {
		respondWithError(w, http.StatusNotFound, "Task not found")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var updatedTask Task
	if err := decoder.Decode(&updatedTask); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	task.Title = updatedTask.Title
	task.Description = updatedTask.Description
	task.Completed = updatedTask.Completed

	if result := a.DB.Save(&task); result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (a *App) deleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var task Task
	if result := a.DB.First(&task, id); result.Error != nil {
		respondWithError(w, http.StatusNotFound, "Task not found")
		return
	}

	if result := a.DB.Delete(&task); result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Funciones auxiliares

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	app := App{}
	app.Initialize()

	port := getEnv("PORT", "8080")
	app.Run(":" + port)
}
