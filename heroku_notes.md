
An account has been created in Heroku :

It created these URLs :
https://arcane-ravine-69473.herokuapp.com/ | https://git.heroku.com/arcane-ravine-69473.git



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


go.mod : 
This uses the Github URL as the identifier for this project, i.e.
module github.com/RMorrisby/ted


Procfile :
This is a Heroku file, instructing Heroku where to find the build web-application.
In this case, it / golang will put it into bin/