# Social Cloud Server

The server backend for the Social Cloud project.

## Running The Server

The server is setup to run on a Google App Engine instance for the social-cloud project.
First, setup the app.yaml file found in the main directory with the desired runtime of
the project. This has been specified as the go111 flexible runtime environment. This
environment is more flexible than previous environments and does not require the use of
any special functions when running the application. Of course, a Google Cloud Platform
project is required to be setup before running the application. To do so, go to this address:
https://cloud.google.com/ and follow the instructions to setup a Google Cloud Project and an
App Engine instance. This involves installing the GCP SDK and authenticating the installed SDK
with Google.

Second, other project dependencies, such as access to the Cloud SQL instance and the Cloud
Storage Bucket instance, are required to be setup to access the instances and their resources
through the Google App Engine instance. Further documentation is included within each dependency.
Note that you will also need to follow the instructions found on the GCP website to setup these
instances. These instances may also be accessed through the virtual shell provided by GCP.

Third, the project must be deployed to the Google App Engine instance through the use of the
gcloud command-line tool. First, the Google Cloud Platform SDK command-line tool must be installed.
Follow the instructions on the official website here: https://cloud.google.com/sdk/. Second, the
project must be deployed to the Google App Engine instance by executing the command
"gcloud app deploy" within the same directory as the app.yaml and main.go files. These can be
found within the main package within the src folder.

## Project Decomposition

The Social Cloud Server backend project has been decomposed into the following packages.

### src

Contains all the golang source code for the project. This includes access to each dependency
as well as definitions of internal models and functions.

### bucket

This package contains the golang source code needed to access the Cloud Storage Bucket instance
where uploaded images are stored. This is accomplished through the use of the storage package
found within the Google Cloud Platform SDK. This package allows the project to access the
Cloud Storage Bucket instance through the Google App Engine instance. This is accomplished by
creating a client through the storage package with the given context, which contains the credentials
needed to authenticate access to the Cloud Storage Bucket. With the client, the bucket is retrieved
by name and stored for later use. When uploading images to the Bucket, the user's email is appended
to the beginning of the file name so that each user may upload files of the same name without
overwriting other user's files.

### database

This package contains the golang source code needed to access the Cloud SQL instance where the
models and their data is stored. This is accomplished through the use of a third party PostreSQL
driver. The instance is connected to the Cloud SQL instance through the use of a username, password,
database name, and connection name associated with the Cloud SQL instance. These values are stored
in variables and given to the PostgreSQL driver in a formatted string. The formatted string specifies
the mentioned variables as well as the host of the instance, which is different than it normally
would be due to the project being hosted on a Google App Engine instance. The formatted string is
used to connect to the database and the database is then pinged to determine that it is alive. If
no errors arise, the database instance is stored for later use. The database package also provides
methods for formatting queries and executing queries on the database. These queries are performed
as complete transactions. Finally, the database package provides a function to rebuild each of the
models from the create and drop queries of each model.

### internal

This package contains the definitions for internal models and their API functions. This includes
the definition of the database models, the definitions of the request and response models, and
the definitions of the REST requests that are supported for each model. Each model within the
application is stored in its own package within this package. Within the package of each model
is contained the model's model package and api package. These packages store the following.

#### model

This package is responsible for defining the internal model. This includes a query that creates
the model within the Cloud SQL instance, a query that deletes the model within the Cloud SQL
instance, and the definition of the model that is returned on requests that return entities of
the given model. The structure that defines the internal representation of the model will include
JSON tags that are used when serializing the model into a JSON representation. The entire model
code is stored in one source file.

#### api

Each file within this package defines a separate handler that will handle a particular
request on a particular URL. Each handler has a request model, a response model, and
defined queries associated with it. Each handler also has a process function that handles
its given request type. Each handler is defined within its own golang source file.

This package is responsible for defining the portion of the REST API that is supported for
the given model. This usually includes create, get, and list endpoints that will create
and retrieve the model(s), respectively. This includes internal definitions of the models
that are expected to be received as a serialized JSON string within the corresponding request.
These are known as the request models. The request models include JSON tags that are used when
deserializing the JSON representation of the request. This also includes internal definitions
of the models that are returned by the endpoint after the request has been handled. These are
known as response models. The response models include JSON tags that are used when serializing
the JSON representation of the response. The API packages also include functions to return
these models and to process a given request. Each process function will verify the request is
the correct type and then process the request with the given arguments. This makes use of the
necessary SQL queries defined within the file and may require the use of the database, the
bucket, or the url-shortener.

### main

This package contains the golang source code that sets up and runs the server. The main.go
file defines the main function that establishes the connections to each of the required
dependencies, registers the list of routes and their corresponding handlers, and then spins
up the server to listen and handle requests. This all accomplished through the use of the
other packages. If any connection should fail, the setup is aborted.

### server

This package is responsible for defining the generic server code used throughout the application,
mostly within each api package of each internal package. This package includes a package called
endpoint that defines the generic handler structure. This is the structure that is implemented
by each of the handlers within the api package of each internal package. There is also a generic
listener structure defined within this package that uses the generic handler structure to listen
for and handle requests. The responsibility of the listener structure is to deserialize the JSON
representation of a request into the associated handler's request model, call the handler's process
function to process the request with the request instance, and then serialize the given response
instance into its JSON representation. The response is then returned the client. This package also
includes the routes.go file, which simply defines all the routes and their associated handlers.
Finally, the server.go file contains the definition of the server, which includes server functions
for registering the handler of each route with a listener and for serving and listening for requests
on the machine.

### url-shortener

This package contains the golang source code that interacts with Bitly for shortening long URLs.
This includes definitions of the request and response structures that are originally defined by
Bitly. These models include JSON tags to allow for seralization and deserialization of the request
and response, respsectively. This package also defines variables that store the URL and API key
necessary to complete the request. This API key is associated with a Bitly project account. Finally,
this package includes the definition of a function that performs the request and returns the response,
which is the shortened URL.

### vendor

This package includes the external libraries that are installed into the project. These libraries are
managed through the use of the golang dep package manager. The golang dep package manager relies on
the Gopkg.toml file for specifying the project's dependencies. These dependencies are then installed
through the use of the "dep ensure" command in the project's root directory, which is the directory
containing the Gopkg.toml.

### Other files of interest

Gopkg.toml
    -> This file defines the dependencies of the project, including the external PostgreSQL driver and
       the Google Cloud SDK that includes the Cloud Storage library for accessing the project's Bucket.
