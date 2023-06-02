# merlin-sessions-api


#API DOCS:

    1.  /get_events
        method: GET
        query:
            tenant_name: string
            api_key: string
            to: int
            from: int
            event_name: string
            event_type: string
            browser_type: string
            wallet_address: string
            user_prop : json
            event_prop: json
        example: /get_events?api_key=B6Fb2l73hE&tenant_name=merlin&event_type=click&event_name=Click&user_prop={"myCustomValue3":"true"}
    
    2.  /get_unic_events
        method: GET
        query:
            tenant_name: string
            api_key: string
    
    3.  /get_unic_events_name
        method: GET
        query:
            tenant_name: string
            api_key: string
    
    4.  /get_unic_browsers
        method: GET
        query:
            tenant_name: string
            api_key: string
            
    
    5. /add_event
        method: POST
        query:
            api_key: string
        data:
            {
            "enviornment_type":"",
            "event_name": "",
            "event_type": "",
            "auth_project": {
                "tenant_id": "",
                "tenant_name": "",
                "api_key": ""
                },
                "user_ids": {
                    "user_id": "",
                    "wallet_address": ""
                },
                "other_admin_variables": {
                    "insert_id": "",
                    "timestamp": "",
                    "app_version": ""
                },
                "event_properties": {},
                "automatically_tracked": {
                    "browser_type": "",
                    "os_name": "",
                    "os_version": "",
                    "device_type": "",
                    "timestapm": "",
                    "url": "",
                    "referrer": "",
                    "country": "",
                    "language": ""
                },
                "user_properties": {}
            }

---

#METHOD DOCS

* ./main.go
     - > main():
        function for run server
* ./initializers/db.go
    - > GetDb():
        function for connect to DB
* ./initializers/redis.go
    - > GetRedis():
        function for connect to Redis
* ./handlers/handler.go
    - > GetEvent():
        function for getting information on db and return as json response
    - > PostEvent():
        function for add information to redis and return as json response
    