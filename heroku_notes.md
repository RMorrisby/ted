
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

Sample Curl commands :

curl -XGET -i http://localhost:8080/is-alive

curl -XPOST -i -H "Content-Type: application/json" -d "{\"Name\": \"234\"}" http://localhost:8080/result

curl -XPOST -i -H "Content-Type: application/json" -d "{\"Name\": null}" http://localhost:8080/result
