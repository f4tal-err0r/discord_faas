package functions

import (
	"archive/tar"
	"compress/gzip"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	v1 "github.com/f4tal-err0r/discord_faas/api/v1"
	"github.com/f4tal-err0r/discord_faas/pkgs/runtimes"
	"github.com/f4tal-err0r/discord_faas/pkgs/security"
	"google.golang.org/protobuf/proto"
)

func (h *Handler) DeployFuncHandler(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value("claims").(security.Claims)

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max size
	if err != nil {
		log.Printf("HTTP %d: failed to deploy function: %v", http.StatusBadRequest, err)
		http.Error(w, fmt.Sprintf("failed to deploy function: %v", err), http.StatusBadRequest)
		return
	}

	// Get metadata part
	metadataFile := r.FormValue("metadata")
	if metadataFile == "" {
		log.Printf("HTTP %d: missing metadata part: %v", http.StatusBadRequest, err)
		http.Error(w, fmt.Sprintf("missing metadata part: %v", err), http.StatusBadRequest)
		return
	}

	var buildReq v1.BuildFunc
	if err := proto.Unmarshal([]byte(metadataFile), &buildReq); err != nil {
		log.Printf("HTTP %d: invalid metadata: %v", http.StatusBadRequest, err)
		http.Error(w, fmt.Sprintf("invalid metadata: %v", err), http.StatusBadRequest)
		return
	}

	funcHash, err := funcHash()
	if err != nil {
		log.Printf("HTTP %d: failed to generate function hash: %v", http.StatusInternalServerError, err)
		http.Error(w, fmt.Sprintf("failed to generate function hash: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Function upload init: GuildID=", claims.GuildID, " FuncName=", buildReq.Name, " FuncHash=", funcHash)

	// Get func part (tarball)
	funcFile, _, err := r.FormFile("func")
	if err != nil {
		log.Printf("HTTP %d: missing func part: %v", http.StatusBadRequest, err)
		http.Error(w, fmt.Sprintf("missing func part: %v", err), http.StatusBadRequest)
		return
	}
	defer funcFile.Close()

	targetDir := filepath.Join(h.cfg.Storage.Path, "build", claims.GuildID, funcHash)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		log.Printf("HTTP %d: failed to create target dir: %v", http.StatusInternalServerError, err)
		http.Error(w, fmt.Sprintf("failed to create target dir: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tarHandler(funcFile, targetDir); err != nil {
		log.Printf("HTTP %d: failed to extract func tarball: %v", http.StatusInternalServerError, err)
		http.Error(w, fmt.Sprintf("failed to extract func tarball: %v", err), http.StatusInternalServerError)
		return
	}

	if err := buildFunc(targetDir, buildReq.Runtime, buildReq.Name); err != nil {
		log.Printf("HTTP %d: failed to build function: %v", http.StatusInternalServerError, err)
		http.Error(w, fmt.Sprintf("failed to build function: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Function build complete: GuildID=", claims.GuildID, " FuncName=", buildReq.Name, " FuncHash=", funcHash)

	//Serialize Upload Response
	uploadResp := v1.UploadResp{
		Hash: funcHash,
		Name: buildReq.Name,
	}

	respBytes, err := json.Marshal(uploadResp)
	if err != nil {
		log.Printf("HTTP %d: failed to serialize upload response: %v", http.StatusInternalServerError, err)
		http.Error(w, fmt.Sprintf("failed to serialize upload response: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func funcHash() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func tarHandler(funcFile multipart.File, targetDir string) error {
	var tarReader *tar.Reader
	buf := make([]byte, 512)
	n, _ := funcFile.Read(buf)
	funcFile.Seek(0, io.SeekStart)
	if n > 2 && buf[0] == 0x1f && buf[1] == 0x8b {
		gzr, err := gzip.NewReader(funcFile)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %v", err)
		}
		defer gzr.Close()
		tarReader = tar.NewReader(gzr)
	} else {
		tarReader = tar.NewReader(funcFile)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to untar: %v", err)
		}

		targetPath := filepath.Join(targetDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create dir: %v", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent dir: %v", err)
			}
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file: %v", err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file: %v", err)
			}
			outFile.Close()
		}
	}
	return nil
}

func buildFunc(targetDir string, runtime string, name string) error {
	if err := runtimes.ProtectedFilesFunc(targetDir, runtime); err != nil {
		return err
	}
	return nil
}
