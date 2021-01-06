### Heroku quirks

It doesn't like any other .go files in the base directory. Importing your own files in Go is ... weird.

# URLs

TED -> https://arcane-ravine-69473.herokuapp.com/

Logs -> https://dashboard.heroku.com/apps/arcane-ravine-69473/logs
Go DB help -> https://devcenter.heroku.com/articles/getting-started-with-go#use-a-database

# Info

An account has been created in Heroku :

It created these URLs :
https://arcane-ravine-69473.herokuapp.com/ | https://git.heroku.com/arcane-ravine-69473.git

Heroku chooses the internal port the web-application must listen on, but to access 
it, the standard port 80 is used. I.e. these two commands both work :

curl http://arcane-ravine-69473.herokuapp.com
curl http://arcane-ravine-69473.herokuapp.com:80

HTTPS :
curl https://arcane-ravine-69473.herokuapp.com    <-- works
curl https://arcane-ravine-69473.herokuapp.com:80  <-- doesn't work



How Heroku works :

Your codebase needs to be in Git. You need to add Heroku as a Remote
(the command 'heroku create' in the root directory of your codebase will do this)

Inspect with git remote -v  :

heroku  https://git.heroku.com/arcane-ravine-69473.git (fetch)
heroku  https://git.heroku.com/arcane-ravine-69473.git (push)
origin  https://github.com/RMorrisby/ted.git (fetch)
origin  https://github.com/RMorrisby/ted.git (push)

Commit & push all of your code to origin as per normal ("git push")
Pushing specifically to Heroku ("git push heroku" or "git push heroku main")
will push the 'main' branch to Heroku. It will then build the application & start it.
Therefore we do not need to build the application beforehand.


Running locally :
go build -o bin/ted.exe -v

Then
heroku local

Note : Procfile needs to be adjusted slightly to run on Windows (change bin/ted to bin\ted.exe)
This change should NOT be committed!



go.mod : 
This uses the Github URL as the identifier for this project, i.e.
module github.com/RMorrisby/ted


Procfile :
This is a Heroku file, instructing Heroku where to find the build web-application.
In this case, it / golang will put it into bin/



### Testing

## Local

Sample Curl commands :

curl -XGET -i http://localhost:8080/is-alive

curl -XPOST -i -H "Content-Type: application/json" -d "{\"Name\": \"234\"}" http://localhost:8080/result

curl -XPOST -i -H "Content-Type: application/json" -d "{\"Name\": null}" http://localhost:8080/result

## Heroku

curl -XGET -i http://arcane-ravine-69473.herokuapp.com/is-alive

curl -XPOST -i -H "Content-Type: application/json" -d "{  \"name\": \"test_TED_data_1_b\",  \"category\": \"Test client\",  \"status\": \"PASSED\",  \"timestamp\": \"2020-12-31 09:09:44 UTC\",  \"testRunIdentifier\": \"2020-12-31 09:09:44 UTC\"  }" http://arcane-ravine-69473.herokuapp.com/result

curl -XPOST -i -H "Content-Type: application/json" -d "{\"Name\": null}" http://arcane-ravine-69473.herokuapp.com/result

### DB

Created with 
(all in local cmd) `heroku addons:create heroku-postgresql:hobby-dev`

`heroku config` will show
DATABASE_URL: postgres://odpwvcfjzdxhom:c3cf3a8184dac6aed218ba4540996ff3b2c31ba163728232ba3921a9295cc4a0@ec2-3-90-124-60.compute-1.amazonaws.com:5432/d8hic8bk4o5ecj

`heroku pg` show more DB statistics

Get Go's PostGres module : `go get github.com/lib/pq@v1`


## CSS notes

Colours : https://www.rapidtables.com/web/css/css-color.html#white

