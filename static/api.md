
# REST API for AUMDUG

Following is a description of each of the urls to be hit for AUMDUG data along
with example data. No authentication is required, but lack of `for ;;` would
still be appreciated. 

The HTTP method for each action will be listed next to the url.

All requests and responses will be in JSON.

## Events

 __URL's:__

`http://aumdug.herokuapp.com/api/events`

`http://aumdug.herokuapp.com/api/events/{id}`

### Create New Event

This is for admins (re: leaders) to submit new events.

__Request__

`PUT http://aumdug.herokuapp.com/api/events`

HTML Body: 

```json
{ 
  "name":         string,
  "description":  string,
  "date":         "YYYY-MM-DDTHH:MM:SS-06:00"
}
```

"date" is of the format "YYYY-MM-DDTHH:MM:SS-06:00" for:

  * YYYY      = year, e.g. 2013
  * MM        = month, e.g. 11
  * DD        = day, e.g. 02
  * THH:MM:SS = T (literal), hours:minutes:seconds, e.g. for 3:30: 15:30:00
  * -06:00    = our offset from GMT, -06:00 is CST (us)

  So what does this mean? Give me a real date, with hours and minutes important
  for time, with -06:00 as the GMT time. E.g. "2013-23-11T20:30:00-06:00" would 
  be perfect for an event that will be occurring at 8:30 on November 11, 2013.

  Yes, this is obscure, but this is a standard UNIX time formatting that is most
  likely already implemented in a library for your platform of choice.

  __Caveat:__ This will probably be done mostly through a web interface and only be one
  person's problem, anyway. But here it is if you'd like to integrate adding
  events from the mobile apps.

__Reply__

```json
{
  "success": boolean,
  "event": {
    "id":           string,
    "name":         string,
    "description":  string,
    "date":         "YYYY-MM-DDTHH:MM:SS-06:00"
  }
}
```

"success" indicates whether the event was successfully created.
"id" should be stored in client 'models' in order to use for GET requests.
"date" see description in Create New Event > Request

### Get All Events

Get a list of all events.

__Request__

`GET http://aumdug.herokuapp.com/api/events`

__Reply__

```json
[
  "event": {
    "id":           string,
    "name":         string,
    "description":  string,
    "date":         "YYYY-MM-DDTHH:MM:SS-06:00"
  }
]
```

"id" should be stored in client 'models' in order to use for GET requests.
"date" see description in Create New Event > Request

### Get Event

Get a single event.

__Request__

`GET http://aumdug.herokuapp.com/api/events/{id}`

Where id is a string that is the unique id for an event. (hint: Get All Events
first, you'll have these)

__Reply__

```json
{
  "id":           string,
  "name":         string,
  "description":  string,
  "date":         "YYYY-MM-DDTHH:MM:SS-06:00"
}
```

"id" should be stored in client 'models' in order to use for GET requests.
"date" see description in Create New Event > Request

### Update Event

Update an event.

__Request__

`PUT http://aumdug.herokuapp.com/api/events/{id}`

Where id is a string that is the unique id for an event.

HTML Body: 

```json
{ 
  "name":         string,
  "description":  string,
  "date":         "YYYY-MM-DDTHH:MM:SS-06:00"
}
```

__Reply__

```json
{
  "id":           string,
  "name":         string,
  "description":  string,
  "date":         "YYYY-MM-DDTHH:MM:SS-06:00"
}
```

### Delete Event

Delete a single event.

__Request__

`DELETE http://aumdug.herokuapp.com/api/events/{id}`

Where is is a string that is the unique id for an event.

__Reply__

```json
{
  "success":  boolean
}
```

Where "success" indicates whether the event was successfully deleted
or not.







