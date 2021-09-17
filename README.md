# Shopify Challenge Submission

When I first saw that the challenge was to create an image repository, I immediately thought of a Docker image repository instead of a *photo* image repository :) So, I thought it would be a fun and interesting challenge if I could create a simpler version of DockerHub, a place where users can store their **Docker** images! So that's what I did: I created a service in which users can push / pull their images to be uploaded to a repository. They can also search for repositories with the search feature.


## Prerequisites

- Docker installed on your system
- Python3 
- S3 bucket name + AWS access key + AWS secret access key (to run locally)
- Terraform (to deploy)


## Directory Structure

There are three subdirectories in this repository
- CLI: A custom CLI which can be used to push / pull images from the repository along with other features as well!
- backend: The Go Gin server that is used for pulling / pushing images 
- infrastructure: The Infrastructure that is used for deploying the server and related resources


## Running Locally

1. Add a .env file in the root of the backend folder (a sample one is provided)
2. `docker-compose up` will run the server and the database (MariaDB) locally

3. `cd cli && python3 -m pip install .` will install the CLI onto your machine
4. Run `mydocker --help` to verify that the CLI was installed successfully


## Running Tests

Run `go test -v ./... ` in the backend directory

## Cleaning Up
`docker-compose down` to clean up all images and networks
`python3 -m pip uninstall mydocker` to uninstall the CLI

## Deploying
To deploy easily, I used Terraform and can be installed from [here](https://learn.hashicorp.com/tutorials/terraform/install-cli)

First, configure aws on your machine, or set the environment variable `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` to your access keys. Details on how to obtain the AWS access keys can be found [here](https://docs.aws.amazon.com/powershell/latest/userguide/pstools-appendix-sign-up.html)

Feel free to modify the variables in the terraform.tfvars to your application needs.

1. `cd infrastructure`
2. Run `terraform plan` to view the plan that Terraform will use to create the resources
3. Run `terraform apply` to deploy the stack
4. Run `cd ../backend && git archive -o eb.zip` to quickly create a zip of the backend directory
5. Log into the AWS CLI and upload the zip to Elastic Beanstalk

## Sample Usage
After installation, run `mydocker signup` to sign up. A config file will be created in your home directory `~/.mydocker.json`. 
```
WARNING: If connecting to production server, please do not use a sensitive password: https has not been configured yet for the deployed server.
```

Let's first grab a test image from the real DockerHub since the one we built does not have all the images just yet :)

`docker pull ubuntu:latest`

Then, when running `docker images`, you should see ubuntu image installed.

Let's tag this image with `docker tag ubuntu:latest <your-username>/test-ubuntu:v1`.

Now let's push this image to our custom image repository!

**Pushing the Image**
1. Create the repository under your account with `mydocker create test-ubuntu`
2. Then, run `mydocker push <your-username>/test-ubuntu:v1`
3. That's it! You have successfully pushed an image to the repository. Let's delete the local images we created with `docker rmi ubuntu:latest && docker rmi <your-username>/test-ubuntu:v1`. When running `docker images`, note that there is no ubuntu image

**Pulling an Image**
1. To test pulling the image, we can run `mydocker pull <your-username>/test-ubuntu:v1`
2. Run `docker images` to see that the image has been pulled locally!


## My Journey

I learned a lot while building this project! My initial idea was to be able to use the official Docker CLI itself to be able to push and pull images; however when trying to implement this feature, I learned about how Docker images are stored as *diffs* instead of just folder for each layer of the image (which makes sense for storage optimization). So, keeping track of the diffs and image trees is very complex and would take much more time to implement (view https://github.com/distribution/distribution as an example)

So I decided to simplify my approach: Use a custom CLI that can take an image specified by a user and package up the image into a TAR that contains *all* the parent layers, and all tags + versions. I then compress the tar and upload it to my repository.


Another aspect that I learned a lot about was using Terraform to represent my infrastructure as code. It was such a great experience and I found it much more to my liking then doing everything through the AWS console interface.


