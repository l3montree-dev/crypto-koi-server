package models

type Cryptogotchi struct {
	Base
	Name string `json:"name"`
	// the id of the token - might be changed in the future.
	// stored inside the blockchain
	TokenId string `json:"token_id"`
	// mapping to the record struct.
	// a record can be transformed to an event.
	Events []Record
}

func (c *Cryptogotchi) ToOpenseaNFT() OpenseaNFT {
	return OpenseaNFT{}
}
