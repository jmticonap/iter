package iter

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	u "github.com/jmticonap/iter/utils"
)

func NewIter(mux *http.ServeMux) {
	APP_PORT := u.GetEnvOrDefault("APP_PORT", "80")

	// (rt.HttpRouterHandler(routesList.Routes))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", APP_PORT),
		Handler: mux,
		// IMPORTANTE: Configura timeouts para que JMeter no deje conexiones abiertas
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Canal para escuchar señales de interrupción del SO (Ctrl+C, Docker stop, etc.)
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")

		// Creamos un contexto con un tiempo límite para el cierre (ej. 15 segundos)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown error: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Println("Starting HTTP server on", server.Addr)

	// ListenAndServe siempre devuelve un error. Al cerrar, devuelve http.ErrServerClosed.
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	// Esperamos a que el proceso de Shutdown termine antes de salir de la función
	<-idleConnsClosed
	log.Println("Server stopped gracefully")
}
