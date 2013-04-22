package chatservice

type Entity interface {
	Id() idType
	SetId(id idType)

	Save() error
	Load(id idType) error
	Delete() error
}
