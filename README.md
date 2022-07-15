# *Email Service* 

--- 

Microservice written in Golang and Deployed in Kubernetes. Is a Part of the `Store` Project, Allows to send Customized Email Notifications, Currently is responsible for notifying about Order States, and Other Sort of Events, related to the Project, such as discounts.

--- 

# *Requirements* 

`Docker` - `1.4.1 or lower`

`gRPC Library` - `Library, that allows to handle gRPC Code and generates source code from the files` 


## "Deployment Options Requirements"


### *Using Docker-Compose* 

`Docker-Compose` ~ `3.8 or higher` 

### *Using Kubernetes* 

Checkout the guide for Kubernetes Version of this Project at https://github.com/LovePelmeni/EmailService/docs/KUBER.md 

--- 

# *Usage* 

Clone this Repo

```
    $ git clone https://github.com/LovePelmeni/EmailService.git
```

(Highly Recommend to run the Project firstly in Docker-Compose and test it out using gRPC, Before going to Kubernetes)

### *Docker-Compose Project Version Usage* 

Go to the Root Directory and Run... 

```
    $ docker-compose up -d 
```

### *Kubernetes Project Version Usage* 

Follow this link in order to get more info https://github.com/LovePelmeni/EmailService/docs/KUBER.md 


--- 

## *External Links* 

`Email for contributions` ~ `kirklimushin@gmail.com` & `klimkiruk@gmail.com`