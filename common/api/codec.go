package api

import (
	"bytes"
	"encoding/gob"
	//"bytes"
	//"encoding/gob"
)

// Data presents the data transported between server and client.
type Data struct {
	Name string        // service name
	Args []interface{} // request's or response's body except error
	Err  string        // remote server error
}

func encode(data Data) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

	//d, err := com.Codec().Marshal(data)
	//if err != nil {
	//	return nil, err
	//}
	//return d, nil
}

func decode(b []byte) (Data, error) {
	buf := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buf)
	var data Data
	if err := decoder.Decode(&data); err != nil {
		return Data{}, err
	}
	return data, nil

	//var c2s = new(Data)
	//if err := com.Codec().Unmarshal(b, c2s); err != nil {
	//	return Data{}, err
	//}
	//return  Data{}, nil
}
