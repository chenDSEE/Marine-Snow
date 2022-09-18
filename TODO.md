### adapt for gin version:
- [ ] all file in *framework/middleware*
- [ ] *grace_exit* example
- [ ] *link_call* example
- [x] *request_response* example
- [x] update *request_response* example to display the Handler be called
- [ ] *route_register* example


### feature:
- [ ] support distribute cron job across network
- [ ] support distribute cron job across local machine
- [ ] An interface for each, so that framework can get related information about the registered App.

### Directory management
- [ ] There should a directory for binary file, runtime file(log, runtime storage, and so on), default configuration.
- [ ] How to specify directory for each app.
- [ ] remove *storage* directory.
- [ ] remove *config* directory.


### remove file for gin version:
**context related:**
- [x] *framework/context.go*
- [x] *framework/core.go*
- [x] *framework/group_route.go*
- [x] *framework/request.go*
- [x] *framework/response.go*

**route related:**
- [x] *framework/route.go*
- [x] *framework/trieTree.go*