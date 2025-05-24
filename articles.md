# List of Articles

- Interservice Communication
- Application Architecture

## Interservice Communication

Communication is a very important aspect of any microservice architecture. Communication always involves 2 parts for the most part - async and sync commmuication. In both parts, latency is an important criteria, hence why grpc has been choosen as the method on commuincation to be adopted when communicating synchronously, however for asyn communication which is important as I would argue that its primary benefit is in decoupling your micro-service architecture, here will will go with an event driven architeure, where we will have a message broker that can publish, subscribe, and close connections to the event store, this will be implemented as an interface with the aforementioned methods to allow for flexibility of choosing/changing what event store you wish to use. In this project we'll create both a kafka store and a rabbitmq store that will both implement this interface, I'll advise to use rabbitmq in smaller systems and kafka in larger systems that will often require replayability(what does this mean? - and further explain). We'll implement both so we'll adjust accordignly as the applications requirements and needs evolve and change. We'll start with kafka, I know its a small system why go with something that big and complex(my arguement to you will be that if you can understand complex systems and the need for them, then it will be no prroblem to understand a smaller system like rabbitmq - and you're life with me much simpler with a simpler solution)

## Designing the architecture for E-commerce microservice project

Architecture is the the make or break all of any project and is often times overlooked or paid little attention to. This due to a variety of reasons like the inexperience of the engineers involved, lack of domain knowledge in the industry been worked on, unclear requirements been communicated, or the developers involved just building for the sake of building, as we are usually guilty of, we do love our problems and tinkering. While these are all valid reasons, I believe there should be no reason whatsovever for there to be no architectural document or diagram in place before starting a project, no matter how vague the document might be or prone to change it might be. Often times, the story of how companies end up having unmaintainable project goes something like this company A would like to deliver a project in 3 months time that they have promised their customers will be available, (often times without any technical validation that this is indeed possible), they then reach out to a "super" engineer and have he/she work in isolation (terrible idea) or alongside other developers, now super engineer knows he can solve the problem because he has worked on something similar or is familiar with the problem space (same is not guaranteed to be the case for other engineers). Super engineer proceeds to follow the provided ui prototype or FRD to build the project,assigns tasks to other developers, reviews their PR's before merging and all. Other engineers love and admire super engineer because he/she makes their lives easier they just have to work based "mindlessly" on the task assigned to them and provided their PR's match the product specifications and the super engineer is ok, they can close their laptops and thank God they don't have to deal with all that technical debt of figuring out why things should be refactored so the codebase is flexible enough to change down the line when the product team wants to make some customer happy, all they have to worry about is delivering a solution to a "specific" task assigned to them, they understand their problem and can offer a solution.

"The Problem" - every project has an archtural document of some sorts wether you're aware of it or not, the question or problem is wether its accessible or not. In the case of the project discussed earlier - such a document exists, just that most people are aren't aware of it as a matter of fact only only person is aware of it. The document in this case is the super engineer - who by no fault of his - is prone to forget, get a new job, or simply quit. Now on the upside, the company has a working product, but, super engineer leaves the company maybe his contract has expired or his off to complete his next super mission. And this is when the problem starts becoming evident, management now wants to make some changes - and they are glad the other developers were working with super engineer so they shouldn't be much of an issue, the firt couple of changes requested go fine (still without any architectural documentation), however the next couple of changes seem to take forever, management is complaining and feels as though developers are slacking of or the present engineers are simply not good enough, maybe the another super engineer, - now super engineer B is awesome he/she is able to speed of development by 25% primarily because his/her experience and/or talent enables him to be able quickly recognize patterns or just ignores what was there, doesn't even try to understand the existing codebase and build or improve on it, just quickly build a tool that he/integrates and solves whatever managment throws at him (still without any document). In the next 4 years, every developer hired either quits within 3 months or is unable to deliver.

All of this could have been avoided if there was a document for reference that was used modified, that all parties regardless of their level of technicality can refer to. If we were to trace the root source of the issue it will be super engineer A, how you start a project. and who leads the project determines the culture and practices that will adhered to in that project in the long run, and nothing spells that out more clearly than a architetural document and an accompanying software documentation. Yes, we could put all the blame on super engineer, but time constraints imposed by management might have led to him not taking out the time to write or maintain such a document. Regardless of who is to blame it should be evident to both parties that for any project to live a long and prosperous life such a document should not be an after thought, and a heavy emphasis as well as story points should be paid on writing and maintaing such documents.

We've said all this, to say regardless of the size of the project, or its importance, you as a developer should make a habit of writing and maintaining such documents before and during your project life cycle, your future self will thank you, or developers will appreciate the effort. Ok then, lets begin planning the architecture of our e-commerce project. We are using a micro-service arhitecture mainly for learning purposes, I'll advise any team to always start out with a monolith before implemting microservices, they are hard to maintain and quickly accumulate technical debt, just my humble opinion. Now, in the real world you'll usually be handed a document by the product team or some figma prototype to provide you with the specifications of the product stakeholders expect to be delivered. Now I've taken the liberty to provide such a document so we could go through it and then plan our architecture. The document is provided here. Now we'll consider the following before coming up with a draft of our architetural document

- Services & the requirements of each service
- Shared components, packages, ... code
- How are services discovered
- Communication between services
- Infrastructure

Then we'll move on to creating an architectural diagram and an accompanying document.

## Services & their requirements

what language is been used single, or multiple, use single, will consider using multiple later. what are your response signatures ...
Fortunately, or unfortunately this project will be written entirely in go meaning all services will be written in go, below are some of the services that we'll be creating.

- Auth Service: This will be responsible for authenticating users, creating users, and assigning/removing roles from users. These roles include Admin, Vendor, and Customer
- Notification Service: This service will be responsible for sending notifications which could be in-app(database and web sockets - a bit of push notification here), email, or sms
- Order Service: This will be the service that consumers will interact with to place their orders for products that have an inventory in the store of the vendor
- Product Service: Contains information about the different products that will are present on the platform, a central registry containing data on products - name, color, description, ....
- Payment Service: This will be the service responsible for processing payments for orders, vendor subscription plans, etc.
- Review & Rating Service: This will be the service that allows vendors, products, stores, to be reviewed and rated by different users
- Search & Recommend Service: This service is responsible for providing and storing search results
- Subscription & Traffic Service: This service is responsible for managing/creating subscription plans as well as defining and enforcing the traffic restriction logic for vendors
- Vendor Service: This is responsible for creating stores for vendors, and recording the inventories of the various products in the store. The service also allows vendors to manage/process the orders, their deliveries, and reversals for the products purchased from their stores

Now all these services have certain aspepts that are common to them and as such there should be folder that holds such code to be reused across services

## Shared Code

- Api Signature & Types: This is the expected response that each service is expected to deliver to the client. Make the life of frontend developers, mobile developers easy, ensure there is a consistent format for data across services. Also will contain structs, interfcaces that will be used across services
- Message Broker: Most or all services will be expecting some sort of asynchronous communication at some point its either their publising messages or subscribing to topics to receive messages. More on this later, but just know we need a message broker that will facilitate non-blocking communication between services
- Database: Although each service will have its own database(this is to ensure decoupling amongst services, and ensure scalability), its helpful to write database connection code that can be reused across services to connect to the same types of databases say postgres, or mongodb.
- env: Just a helper to retrieve environment variables
- events: This more or less acts as a dictionary that all services can refer to when the want to subscribe to a topic or know the events that are present on a topic. Just helps to prevent mistake and ensures a single source of truth across services.
- logger: This is just a logging mechanism that will be used across services, to ensure that all services produce logs of a standard format.
- observability: This will contain code that will facilate tracing via opentelemetry to provide insights into how data moves across services, and the issues or successes encountered
- proto: This is will contain the grpc generated code that will be used to create grpc clients and servers accross services
- Utils: Any utilities that will be used accross services like Pagination, Getting limit and offset from request queries, etc.

## Service Discovery

Services will have to be aware of each other in other to be able to communicate, will first start simple by implementing an api-gateway and a bit of hardcoding, then we'll discuss the pitfalls of this approach and how we can solve this issue especially in a multi-node environment like k8s with a service mesh implementation using istio probably

## Inter Service Communication

Service communication is either blocking or non-blocking, and we'll need to implement both based on what were trying to accomplish at the moment. Take for example a client makes a request to register a vendor. The following services will be involved, the auth service to create a user then the vendor service to create a vendor account with a store, and then the subscription service to create a subscription plan for the vendor, and then the payment service to generate a link for the user to complete payment for the subscription plan. In this scenario communication between these services will have to be blocking as one service will have to wait for a response from one service before proceeding to send the necessary information(userId, vendorId, subscriptionId) to the next service. In such a case there are a couple of ways blocking communication like this can be accomplished via http or grpc, will with grpc because of its low latency and strict type enforcemnt(use correct term here).
Now for non-blocking communication, is typically communication that really doesn't need an immediate reaction from the services involved. Take for example the user wishes to pay for the subscription with the link provided, this can happen at any time and the services simply have to be ready when that happens. Lets walk through how the process the user receives the link and then proceeds to make payment, the payment gateway will receive the payment and then notify the payment service through a webhook that the payment has been made. At this point the payment service doesn't need to make a blocking call to inform the subscription service that payment is made, rather it can make a non-blocking call publishing a message via the message broker and then the subsription service will listen for those messages and update the subscription accordingly. The main advantage of this is that both services can act independent of one another and there is really no coupling in this case.

## Infrastucture

Now, we'll not be focusing as much on infratructure initially so we'll start out simple with each service on the same server just a different port and using an api gateway to distribute incomming traffic to the service concerned. We'll do a bit of containarization here especially when implmenting databases, prometheus, grafana, loki, etc. Eventually, we'll move to a multinode environment like kubernetes, by this time you'll have understood why things work the way the do, and why we need kubernertes to scale this services. We'll then move to how to design and implement such a service in the cloud using AWS. And we'll round up with Gitops using tools like ArgoCD. (Part of this should be in conclusion, consider rethinking this)

## Monitoring & Observability

These is important for any system and essential for microservices as this will aid you in making critical decisions and in debugging. There are 3 essential parts involved:

- Loggging
- Tracing
- Metrics

## Conclusion

This is just an introduction to a couple of lectures we're will be discussing micro-services and system architecture, with the e-commerce project as a reference point. The relevant links will be shared and the repository shortly, still been worked(would like to make it as self-docmenting and explainatory as possible). Cheers and stay tuned for the next episode.

## Next Steps

- Inter service Communication with examples
- Mono repos in Microservice
- Service Discovery with Istio
- K8s and its adoption in micro services
- K8s and its core components
- ArgoCD & k8s
- Intro to k8s on the cloud
- Jenkins for automation\*
- Ansible for automation\*
- Linux is your friend\*

# Monorepos in Microservice architecture

(Also define mono repos ...).Generally, mono repos are used because they possess the unique advantage of ensuring that repositories remain seperate while been able to share code common to them in such a way that once a change is made they are all made immediately aware of. In most projects, particularly projects where different aspects are written in different languages or the primary language relies on (heavy bundling and dependepency graphs)[what does this mean exactly ?], yes I'm looking at you javascript, you'll need tools like NX or turbo repo or nx. Fortunately, this project will be written in a one programming language, which has built-in tooling that can be used to effectively set up a mono repo. Now to setup a monorepo in go we have 2 options:

- Option 1: Single Go module (go.mod in root)
  Simpler, best for small to medium teams.
  Shared packages and services import each other like:
  ```go
  import "ecommerce-mono/internal/broker"
  ```
- Option 2: Multiple modules (one go.mod per service)
  More scalable and decoupled.
  Each service (and proto/) has its own go.mod.
  Services depend on shared modules using relative paths or replace directives.
  ```bash
  // In /cmd/product-service/go.mod
  require (
    ecommerce-mono/internal/logger v0.0.0
  )
  replace ecommerce-mono/internal/logger => ../../internal/logger
  ```

We'll go with option 1, mainly because in this case we unfornuately are just one developer, ambitious but still one developer, and its easier for us to and ensures a better development workflow if we go with this option, besides refactoring to go with option 2 as the team grows or project requirement demands. As we established earlier(we didnt do this, but we ought to and confirm if correct), a monorepo is just a fancy way of saying you have one big repository with smaller repositories that are folders and this smaller repositories/folders share common packages, and code between them. If this is the case then lets simply define and discuss the folder structure of our monorepo before moving forward. Here is a how it will look below. We'll discuss each folder and its purpose as we proceed.

Folders

- Internal
- Proto - it seems this should contain generated code as well. also consider versioning and maintaining .proto files as a public contract. or the output should go to pkg
- cmd: entry points i.e executables
- scripts - CI/CD or dev scripts
- deploy - Helm charts / Dockerfiles / k8s manifests
- Readme.md

## Conclusion

## Next Steps

- Infrastucture & Infrastructure Planning
- Work on each service and their deployment

# Infrastucture & Infrastructure Planning

# Oauth Compliant Auth Service in Go

As with any concept, pictures do help alot. I have taken the liberty to create a diagram that illustrates what we would be building today. Now before we discuss the diagram its important that you the audience is familiar with some terms

- oauth
- jwt
- claims
- oidc
  Now in a typical stateless auth system which is what we would be building most applications will be staisfied with the jsw holding state like the user id name roles on the app, etc. and will likely make use of some oauth services like google, apple, etc. Now this is fine, now in the event that your system gets popular like lets say google or other systems need to integrate their entry point more often than not especially in the case of b2c applications will be the auth service. And its important to account for this when building your auth service. A key aspect in building scalable systems in accounting for extensibility. For any system to be extensible its has to enforce or adhere to some sought of standard, in the world of software engineering this standards are what we commonly refer to as protocols. That is why its easy for you to add google, apple as an auth provider in your application because they all follow a statndard that enables an easy and predictable api. And if we do follow this standard, oauth. Our application will eb extensible as well, and it will be easy for other applications to easily use our service

Diagram - > will show the auth service and its components register, login, verify, etc. and how they will interact with storage, then will show the middleware that will be exposed as a guard to be used by other services. Then the external oauth providers that will tie into the register, and login. And then our oauth that will be exposed for other applications to use as an oauth provider

Talk about and show code for the following

- login
- register
- middleware

Buil Flow

- normal routes
- oauth login for external provider as interface, and then map to have user pick, test out
- oauth exposed by our service

Conclusion

## Auth service

so we'll discuss

- app native authentication mechanism: you regular old login, register, and logoout routes
- oauth provider integration: example with github
- Exposing our very own oauth service: a focus on building extensible and compliant systems, a compliant system is an extensible one ...

The diagram below showcase the proposed architecture of our auth service

### App Native Auth Mechanism

### Oauth Provider Integration

### Creating and exposing our own oauth provider to be consumed by other

Conclusion
