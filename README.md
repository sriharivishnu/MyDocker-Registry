# Shopify Challenge Submission
This repository was made for the challenge detailed [here](https://docs.google.com/document/d/1eg3sJTOwtyFhDopKedRD6142CFkDfWp1QvRKXNTPIOc/edit)

When I first saw that the challenge was to create an image repository, I immediately thought of a Docker image repository instead of a *photo* image repository :) So, I thought it would be a fun and interesting challenge if I could create a simpler version of DockerHub, a place where users can store their **Docker** images! My Challenge submission is a backend service which supports authentication, pushing / pulling docker images, and searching image repositories of all users.


## Prerequisites

- Docker installed on your system
- Python3 
- S3 bucket name + AWS access key + AWS secret access key (to run locally)
- Terraform (to deploy)
- Go (to run tests)


## Quickstart
1. To install the CLI, `python3 -m pip install ./cli`
2. Run `mydocker --help` to verify the installation and see a full list of available commands
3. For example, run `mydocker signup` to sign up to be able to push images!


```
Usage: mydocker [OPTIONS] COMMAND [ARGS]...

  CLI application that provides methods for testing the custom docker
  repository

Options:
  --help  Show this message and exit.

Commands:
  config        Set URL, reset settings, or get current user
  create        Creates a repository
  images        Retrieves the images for a given repository
  login         Log in to push images to a repository
  pull          Pulls an image from the repository.
  push          Pushes an image to the repository.
  repositories  Retrieves the repositories for a given user.
  search        Search for image repository
  signup        Sign up to push images to a repository
```

NOTE:
By default, the CLI points to the deployed API URL, but can be configured with `mydocker config --host <URL>` (when running locally, URL would be `http://localhost:5000/v1`)

## Sample Usage of mydocker
After installation, run `mydocker signup` to sign up. A config file will be created in your home directory `~/.mydocker.json`. 
```
WARNING: If connecting to production server, please do not use a sensitive password: https has not been configured yet for the deployed server.
```

Let's first grab a test image from the real DockerHub since the one I built does not have all the images just yet :)

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


## Directory Structure

There are three subdirectories in this repository
- CLI: A custom CLI which can be used to push / pull images from the repository along with other features as well!
- backend: The Go Gin server that is used for pulling / pushing images 
- infrastructure: The Infrastructure that is used for deploying the server and related resources


## Running Locally (docker-compose)

1. Set the AWS credentials for the server in the docker-compose.yml file
2. `docker-compose up` will run the server and the database (MariaDB) locally

3. `python3 -m pip install ./cli` will install the CLI onto your machine
4. Run `mydocker --help` to verify that the CLI was installed successfully

## Running Tests (Requires Go)

Run `go test -v ./... ` in the backend directory

## Cleaning Up
`docker-compose down --rmi all` to clean up all images and networks
`python3 -m pip uninstall mydocker` to uninstall the CLI

## Deploying
To deploy the infrastructure to AWS, I used Terraform which can be installed from [here](https://learn.hashicorp.com/tutorials/terraform/install-cli)

First, configure aws on your machine, or set the environment variable `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` to your access keys. Details on how to obtain the AWS access keys can be found [here](https://docs.aws.amazon.com/powershell/latest/userguide/pstools-appendix-sign-up.html)

Feel free to modify the variables in the terraform.tfvars to your application needs.

1. `cd infrastructure`
2. Run `terraform plan` to view the plan that Terraform will use to create the resources
3. Run `terraform apply` to deploy the stack
4. Run `cd ../backend && git archive -o eb.zip` to quickly create a zip of the backend directory
5. Log into the AWS CLI and upload the zip to Elastic Beanstalk

To destroy the resources, run `terraform destroy`

## My Journey

I learned a lot while building this project! My initial idea was to be able to use the official Docker CLI itself to be able to push and pull images; however when trying to implement this feature, I learned about how Docker images are stored as *diffs* instead of an entire folder for each layer of the image (which makes sense for storage optimization). So, keeping track of the diffs and image trees is very complex and would take much more time to implement (view https://github.com/distribution/distribution as an example)

So I decided to simplify my approach: Use a custom CLI that can take an image specified by a user and package up the image into a TAR that contains *all* the parent layers, and all tags + versions. I then compress the tar and upload it to my repository.


Another aspect that I learned a lot about was using Terraform to represent my infrastructure as code. It was such a great experience and I found it much more to my liking then doing everything through the AWS console interface.

## Future Features

With more time, here are some cool ideas I would love to integrate with this project!

- Maintaining an image tree similar to DockerHub
- Being able to search for images from what's in the images itself
- Private repositories

## Contact
Overall, I had a lot of fun building this project, and I hope you enjoy it as well! For any questions/inquiries, technical support, or if you just want to chat, feel free to contact me at srihari.vishnu@gmail.com

