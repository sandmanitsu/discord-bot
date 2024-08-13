package disk

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func GetService() (*drive.Service, error) {
	path := path.Join("./", "google", "", "credentials.json")
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// setting config OAuth2
	config, err := google.ConfigFromJSON(b, drive.DriveReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	srv, err := drive.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func getClient(config *oauth2.Config) *http.Client {
	path := path.Join("./", "google", "", "token.json")
	tok, err := tokenFromFile(path)
	if err != nil {
		log.Fatalf("Unable to retrieve token from file: %v", err)
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// todo to other file
func ListFilesInFolder(srv *drive.Service, folderID string) {
	query := fmt.Sprintf("'%s' in parents", folderID)
	r, err := srv.Files.List().Q(query).Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
}
