# go-ping-sites
GoPingSites is intended to be a tool to monitor multiple websites, written in Go.
## Steps to setup locally
1. Use the make distribute command to build the project and will copy the necessary runtime files to $GOPATH/dist/go-ping-sites.
```
> make distribute
```
2. Go to the directory $GOPATH/dist/go-ping-sites.
```
> cd $GOPATH/dist/go-ping-sites
```
3. Edit the config.toml for your email/text settings and db-seed.toml to setup the initial sites for the application.
4. Run the application.
```
> go-ping-sites
```
5. Browse to localhost:8000 and login as admin with the password adminpassword.
6. Go to Settings and click the "Users" link to edit the admin user with a new passsword.

## Rebuilding the application
Once the application has been setup locally using the steps given above, you probably want to rebuild the application but not overwrite your settings.  You can do that by running make without the "distribute" argument.
```
> make distribute
```
Or to run the application after building:
```
> make run
```
### Resetting the database
You can restart with a clean database by simply deleting the old one and restarting.
```
> rm go-ping-sites.db
> rm go-ping-sites-auth.db
> go-ping-sites
New Database, creating Schema...            
Seeding initial sites with db-seed.toml ...
```
It will recreate a new one once you rerun go-ping-sites, based on the settings in db-seed.toml.
