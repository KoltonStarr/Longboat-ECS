# Overview 
This is a very lightweight golang server that exposes a /images endpoint where images can be uploaded and retrieved from AWS S3. 

The purpose of this server is not to solve any real-world problem, but rather to act as a proof-of-concept for what it means to containerize a golang application 
that will ultimately run on AWS ECS. 

## Local Development
If you have AWS credentials configured in your ~.aws directory then you will be able to run the server locally outside of a docker container. Simply run ```go run .``` at the root level of the go program and it will start up. 

If you want to run the server in a docker container then first build the container locally (I only have the image published in a private AWS ECR registry at the moment) and the run it using the follwing command: 
```docker run -e AWS_ACCESS_KEY_ID=<your access key> -e AWS_SECRET_ACCESS_KEY=<your secret key> -e AWS_REGION=us-east-1 -p 3000:8080 <name of golang server container>```

## Testing The Routes
```curl -X POST -H "Content-Type: multipart/form-data" -F "file=@<relative path to an image file>" -v http://localhost:3000/images```

## Production
There is evidently no need to do any sort of credential management when running this container on AWS ECS. The AWS SDK golang code that instantiates a new s3 client 
knows in what environment it is running. So the code will communicate with  the ECS task agent that is hanging about and get creds from that. 