- Based Source : "Scribble - A tiny Golang JSON database"

- Upgrade :
  - Changed the primary key from "Title" to "ID"
    By default Scribble uses "Title" as the primary key. I have added a new field "ID" as the primary key. This will help in updating the records based on the ID.
  - Added a new function "UpdateRecord" to update the record based on the ID.
