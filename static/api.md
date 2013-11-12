
# REST API for AUMDUG

Following is a description of each of the urls to be hit for AUMDUG data along
with example data. No authentication is required, but lack of `for ;;` would
still be appreciated. 

The HTTP method for each action will be listed next to the url.

All requests and responses will be in JSON.

### Events

 __URL's:__

`http://aumdug.herokuapp.com/api/events`

`http://aumdug.herokuapp.com/api/events/{id}`

##### Create New Event

This is for admins (re: leaders) to submit new event.

__Request__

PUT `http://aumdug.herokuapp.com/api/events`

Body: 

<!-- TODO date format -->
{
  "name":         string,
  "description":  string,
  "date":         time.Time
}

__Reply format__

<!-- TODO date format -->
{
  "success": boolean,
  "event": {
    "name":         string,
    "description":  string,
    "date":         "0001-01-01T00:00:00Z",
    "created":      "2013-11-12T13:36:23.089274459-06:00"}
  }
}

"success" indicates whether the event was

##### Get All Events

Get a list of all events.

__Request__

GET `http://aumdug.herokuapp.com/api/events`

__Reply format__

[
  "event": {
    "name": string,
    "description": string,
    "date": "0001-01-01T00:00:00Z",
    "created":"2013-11-12T13:36:23.089274459-06:00"}
  }
]
