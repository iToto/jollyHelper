# JollyHelper

JollyHelper is a Christmas exchange application that takes in dictionary in the format:

`Name : email@address.com`

Then it will randomly assign each email address a name (that does not belong to the original email address)

## Setup

* Copy and remane the sample environment file

```sh
» cp env.example .env
```

* Configure the .env file to your needs

* Install [Heroku toolbelt](https://toolbelt.heroku.com)

## Run

You can run the application with:

```sh
» foreman start -e .env
```