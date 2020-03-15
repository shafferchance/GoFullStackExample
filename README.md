# Go Server
&nbsp;&nbsp;Simply connects to MongoDB and holds the connection context. Holds endpoints that can push and pull data from the mongodb database.
## Endpoints
- /
    - GET
        - Able to parse URL and determine which static file to send. There is a special clause handling set-up for the root case to keep from sending wrong html file, but in theory this allows breaking away from the index.html convention if one desired to do so
- /input
    - GET
        - Will retrieve all entries from the `new_db` of `text` collection  
    - POST
        - Will add entry to `new_db` of type `text` collection
## Notes
&nbsp;&nbsp;While making this it was discovered that go uses structs to interact with the database. Thus these functions are able to generized, but with much more difficulty due to having to initialize custom structs. However this is workable.<br>
&nbsp;&nbsp;In addition, sending JSON is much more difficult as there is no string templating without an import, creating some difficulty with readability that can be conquored with a some understanding of string concatenation.<br>
&nbsp;&nbsp;Finally, the driver for Go is quite in line with the mongo shell at the small cost of the need of an execution context. Making learning it quite easy at-least from a high level view.<br>
&nbsp;&nbsp;With some time this entire process could probably be generized and a little magic. In addtion, GO provides a much more performant option at the cost of some weird syntax. So I would declare this experiment a success.
# JS
&nbsp;&nbsp;Nothing special to report here. Fetch and regular DOM APIs work perfectly :)
