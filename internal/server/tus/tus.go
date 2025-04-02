package tus

import (
	"fmt"
	"log"
	"net/http"
	"shopnexus-go-service/config"
	"shopnexus-go-service/internal/logger"

	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"
)

func Init(mux *http.ServeMux) error {
	// store := s3store.New(config.GetConfig().S3.Bucket, s.Services.S3.Client)
	store := filestore.New("./uploads")

	// A locking mechanism helps preventing data loss or corruption from
	// parallel requests to a upload resource. A good match for the disk-based
	// storage is the filelocker package which uses disk-based file lock for
	// coordinating access.
	// More information is available at https://tus.github.io/tusd/advanced-topics/locks/.
	locker := filelocker.New("./uploads")

	// A storage backend for tusd may consist of multiple different parts which
	// handle upload creation, locking, termination and so on. The composer is a
	// place where all those separated pieces are joined together. In this example
	// we only use the file store but you may plug in multiple.
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	// Create a new HTTP handler for the tusd server by providing a configuration.
	// The StoreComposer property must be set to allow the handler to function.
	logger.Log.Info(fmt.Sprintf("Tus url: %s", config.GetConfig().App.TusUrl))
	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              config.GetConfig().App.TusUrl,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
		// Cors: &tusd.CorsConfig{
		// 	Disable:          false,
		// 	AllowOrigin:      regexp.MustCompile(".*"),
		// 	AllowMethods:     "GET, POST, PUT, DELETE, HEAD, PATCH, OPTIONS",
		// 	AllowHeaders:     "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Upload-Length, Upload-Offset, Tus-Resumable, Upload-Metadata, Connect-Protocol-Version, Tus-Version, Tus-Max-Size, Tus-Extension, X-HTTP-Method-Override, X-Requested-With",
		// 	ExposeHeaders:    "Upload-Offset, Location, Upload-Length, Tus-Version, Tus-Resumable, Tus-Max-Size, Tus-Extension, Upload-Metadata",
		// 	MaxAge:           "3600",
		// 	AllowCredentials: true,
		// },
	})
	if err != nil {
		log.Fatalf("unable to create handler: %s", err)
	}

	// Start another goroutine for receiving events from the handler whenever
	// an upload is completed. The event will contains details about the upload
	// itself and the relevant HTTP request.
	go func() {
		for {
			event := <-handler.CompleteUploads
			log.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()

	// Right now, nothing has happened since we need to start the HTTP server on
	// our own. In the end, tusd will start listening on and accept request at
	// http://localhost:8080/files
	// http.Handle("/files/", http.StripPrefix("/files/", handler))
	// http.Handle("/files", http.StripPrefix("/files", handler))
	// err = http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	log.Fatalf("unable to listen: %s", err)
	// }
	mux.Handle("/files/", http.StripPrefix("/files/", handler))
	mux.Handle("/files", http.StripPrefix("/files", handler))
	logger.Log.Info(fmt.Sprintf("Tus server started on %s", config.GetConfig().App.TusUrl))

	return nil
}
