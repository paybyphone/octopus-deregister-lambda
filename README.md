# Octopus Deregister Lambda

An [AWS Lambda][lambda] function to deregister an [EC2 instance][ec2] from an
[Octopus Deploy][octopus] server when the instance is terminated.

It expects Octopus credentials to be stored in an [S3 bucket][s3]. The EC2 instance needs to
be tagged with the Octopus machine name.

## Building

Octopus Deregister is written in [Go][go], and uses [aws-lambda-go][lambda-go] to bridge
Lambda's Python runtime to Go.

Run tests using `go test ./...`.

Aws-lambda-go distributes a [Docker][docker] container with build and packaging tools. Run
```
docker run --rm -v $GOPATH:$GOPATH -e GOPATH=$GOPATH -w `pwd` eawsy/aws-lambda-go
```
The output `handler.zip` can then be uploaded to Lambda.

## Function Configuration

The following configuration parameters are passed as lambda environment variables

|Name                  |Purpose                                         |Required|
|----------------------|------------------------------------------------|--------|
|OCTOPUS_API_KEY_BUCKET|The S3 bucket name to access the Octopus API key|Yes     |
|OCTOPUS_API_KEY_PATH  |The S3 object key to access the Octopus API key |Yes     |
|OCTOPUS_URI           |URI for the Octopus server                      |Yes     |
|DEBUG_LOG             |Enables debug logging if non-blank              |No      |

## EC2 Configuration

The lambda function expects the EC2 instance to have the tag `octopus_name` set to the 
Octopus display name for the instance.


[lambda]:    https://aws.amazon.com/lambda/
[ec2]:       https://aws.amazon.com/ec2/
[s3]:        https://aws.amazon.com/s3/
[octopus]:   https://octopus.com
[go]:        https://golang.org
[lambda-go]: https://github.com/eawsy/aws-lambda-go
[docker]:    https://www.docker.com
