# GoPingSites
![Build and Test Action](https://github.com/turnkey-commerce/go-ping-sites/actions/workflows/build-test.yaml/badge.svg)
 [![Go Report Card](https://goreportcard.com/badge/github.com/turnkey-commerce/go-ping-sites)](https://goreportcard.com/report/github.com/turnkey-commerce/go-ping-sites) 

GoPingSites is a tool to monitor multiple websites, written in Go.

Features include:
* Setup multiple sites to monitor with a configurable ping frequency for each site.
* Setup multiple contacts (per site) to notify about downtime and when service is restored.
* Notifications optionally sent via email and/or text messaging.
* Easy web user interface for dashboard, configurations, and uptime reports.
* History saved to a SQLite database.
* Easy installation and production deployment.

![Dashboard Page](https://github.com/turnkey-commerce/go-ping-sites/blob/master/screenshots/dashboard.png)

## Steps to setup locally
1. Use the make distribute command to build the project and will copy the necessary runtime files to $GOPATH/dist/go-ping-sites.

  ```
  > make distribute
  ```
2. Go to the directory $GOPATH/dist/go-ping-sites.

  ```
  > cd $GOPATH/dist/go-ping-sites
  ```
3. Copy the config_sample.toml to config.toml and optionally copy db-seed_sample.toml to db-seed.toml.
4. Edit the config.toml for your email/text settings and db-seed.toml to setup the initial sites for the application.
5. Run the application.

  ```
  > go-ping-sites
  ```
6. Browse to localhost:8000 and login as the **admin** user with the password **adminpassword**.
7. Go to the **Profile** tab and provide a new password for the admin account.

## Rebuilding the application
You can rebuild the application again without deleting the current configuration and database by just using the make without the distribute argument.

```
> make
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

For more information on how to use the site please refer to the [User Guide](https://github.com/turnkey-commerce/go-ping-sites/wiki/User-Guide).

For more information on how to deploy the application for production refer to the [Installation](https://github.com/turnkey-commerce/go-ping-sites/wiki/Installation)
and [Deployment](https://github.com/turnkey-commerce/go-ping-sites/wiki/Deployment) guides.
