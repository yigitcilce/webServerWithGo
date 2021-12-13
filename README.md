# Creating a basic RestFul API with Golang
Entry level web server by using golang, tests are done in Postman<br />

V2_with_mux branch handles [1],[2],[3],[4],[5],[6]<br />
V1_without_mux branch handles [1],[2],[5]

"Entries[]" is the simulation of a Database

Usage of methods :<br />
[1]  GET http://localhost:10000 -> welcomes you to the home page <br />
[2]  GET http://localhost:10000/all -> returns all entries in Entries[]<br />
[3]  GET http://localhost:10000/entry/{id} -> returns the entry which has {id} in Entries[]<br />
[4]  DELETE http://localhost:10000/entry/{id}  -> erase an element of Entries[] with the given ID, if it exist in the first place<br />

[5] POST http://localhost:10000/entry  with a body : -> creates another entry if Id isnt already exist in Entries[], sorts Entries[]<br />
{<br />
"Id": "5", <br />
"Title": "Some title", <br />
"desc": "Some description if you feel like it", <br />
"content": "Hi" <br />
}<br />

[6] PUT http://localhost:10000/entry/{id} with a body -> updates an element of Entries[] with the given ID, if it exist in the first place<br />
{<br />
"Id": "5", <br />
"Title": "Some title", <br />
"desc": "Some description if you feel like it", <br />
"content": "Hi" <br />
}<br />
