package dynamodbrepository

const (
	USER_TABLE_NAME    string = "UserServiceTable"
	PRODUCT_TABLE_NAME string = "ProductServiceTable"
)

type DynamoDbCompositeKey struct {
	Pk string `dynamodbav:"Pk"`
	Sk string `dynamodbav:"Sk"`
}
