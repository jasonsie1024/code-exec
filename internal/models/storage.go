package models

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"

	"cloud.google.com/go/storage"
)

func StorageRead(object *storage.ObjectHandle, v interface{}) error {
	reader, err := object.NewReader(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		err = reader.Close()
	}()

	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		log.Println(err)
	}

	return err
}

func StorageSave(object *storage.ObjectHandle, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	writer := object.NewWriter(context.Background())
	defer func() {
		err = writer.Close()
	}()

	_, err = io.Copy(writer, bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
	}

	return err
}
