## Overview 
This is a very lightweight golang server that exposes a /images endpoint where images can be uploaded and retrieved from AWS S3. 

The purpose of this server is not to solve any real-world problem, but rather to act as a proof-of-concept for what it means to containerize a golang application 
that will ultimately run on AWS ECS. 

# Misc. 
- curl -X POST -H "Content-Type: multipart/form-data" -F "file=@<path-to-file-for-upload>" http://localhost:3000/images 