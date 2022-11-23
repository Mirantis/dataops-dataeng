package awspipe

// import (
//     "context"
//     "fmt"
//     "log"

//     "github.com/aws/aws-sdk-go-v2/aws"
//     "github.com/aws/aws-sdk-go-v2/config"
//     "github.com/aws/aws-sdk-go-v2/service/dynamodb"
// )

// func main() {
//     // Using the SDK's default configuration, loading additional config
//     // and credentials values from the environment variables, shared
//     // credentials, and shared configuration files
//     cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("&config "))
//     if err != nil {
//         log.Fatalf("unable to load SDK config, %v", err)
//     }