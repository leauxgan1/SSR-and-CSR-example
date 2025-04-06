# Comparison of SSR and CSR
This project was created out of a curiosity on the difference in implementation between a Client-Side Rendering architecture and a Server-Side Rendering one. 
It is one thing to understand the high-level differences between the rendering techniques, but understanding the necessary implementation details is another thing that may otherwise go overlooked.
To keep things simple and focus on the design differences, a simple Todo application was created to focus on simple endpoints and a simple data structure to maintain, that being a list of todos for a user to update and check off over time.
# Dependencies / How to setup and run
The file server for the main index.html file is written in golang (version 1.24.1) and also serves as an API for an in-memory store of todos that can be queried, added, or removed based on user requests.
To test each server, you may either
 - Install go version 1.24.1  OR
 - Use the provided nix flake to enter a dev shell using "nix develop" that contains the correct go version

Then, simply enter the respective directory and run "go run main.go".
Each server will be opened on localhost:8080 and provide an index.html, but the differences will be explained below.
## Server Side Rendering
This method of rendering html dynamically involves first rendering the html on the server and then sending the html as a string to the client to be rendered.
This means that the client side doesnt require much logic and simply provides a way to receive html or html partials and render them into the current page.
For example, an HTTP POST request could look like this:
 - User presses a button to add a new todo to their todo list
 - An HTTP POST request is sent to the backend server, containing the title and details of the todo note
 - The server receives this request and parses each of these arguments before creating a new todo object and adding it to its internal store
 - The server may then either A) render the entire list of todos anew or B) render the single todo as an html partial
 - The server then sends back this information, allowing the html partial to be added to the page to reflect the changed state on the server
### SSR Conclusion
Before making conclusions on the properties of an SSR architecture, the implementation of this pattern should be acknowledged.
More specifically, the frontend library HTMX was used to handle the necessary AJAX calls to send requests to the server and insert their responses into the document as html.
The use of this library significantly reduced the overall amount of code necessary to to implement the main features of the app.
However, the total amount of code, including the handler code from HTMX, is still reduced in comparison to the CSR approach. 
This is due to the fact that 
## Client Side Rendering
This method of rendering html dynamically involves instead only rendering on the client side using the web's premier scripting language javascript.
This allows the server to only send the data required to properly render on the client rather than the entirety of the rendered data.
However, this requieres the data to be received and properly rendered using javascript on the client.
An example of how this might work is:
 - User presses a button to add a new todo to their todo list
 - An HTTP POST request is sent to the backend server, containing the title and details of the new todo note
 - The server receives this information and adds it to its internal store of todo objects
 - Either
  - A) The server responds with the information of the newly added todo object, and the client uses this to render a new html element in the todo list
  - B) The server does not respond with information, and the client simply uses the provided information from the user to render its own list
The latter option is only available when the user provides all of the necessary information to render the necessary html
When data is external, a request for that information is necessary as well as proper rendering of it into an html partial which is then added to the page
### CSR Conclusions
TODO
## Overall Conclusion
TODO
