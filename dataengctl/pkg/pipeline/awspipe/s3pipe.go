package awspipe

// import (
//     "context"
//     "fmt"
//     "log"

//     "github.com/aws/aws-sdk-go-v2/aws"
//     "github.com/aws/aws-sdk-go-v2/config"
//     "github.com/aws/aws-sdk-go-v2/service/dynamodb"
// )

// // Documentation on Snowpipe
// // https://docs.snowflake.com/en/user-guide/data-load-snowpipe-intro.html
// // List of Supported Clouds
// // AWS - S3
// // Azure - Azure Blob Storage
// // GCP -  

// func pipelineClient () {
//     // Using the SDK's default configuration, loading additional config
//     // and credentials values from the environment variables, shared
//     // credentials, and shared configuration files
//     // Clients will need to be configurable by cloud. 
//     cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2"))
//     if err != nil {
//         log.Fatalf("unable to load SDK config, %v", err)
//     }