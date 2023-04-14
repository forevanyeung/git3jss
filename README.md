# git3jss
Rewrite of the git2jss library by badstreff in Go. Serves as a Action you can use in Github Pipelines, and adds tools to create new script or extension attributes from templates and cleanup unused scripts. 

```
git3jss config -config 
               -server 
               -username 
               -password 
               -scripts-dir 
               -extensions-dir

git3jss download 
git3jss download -type scripts, extensions
git3jss download -config config.json

git3jss sync
git3jss sync scripts
git3jss sync extensions
git3jss sync -config config.json

git3jss review scripts
git3jss review 

git3jss new script
git3jss new extension
```
