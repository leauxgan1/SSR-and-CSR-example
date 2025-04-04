const GET_URL = window.location + "getTodo";
const POST_URL = window.location + "submitTodo";
const REMOVE_URL = window.location + "removeTodo";

const todoTemplate = document.querySelector("#todoTemplate");

async function fetchTodos() {
	const response = await fetch(GET_URL, {method: "GET"});
	const json = response.json();
	return json;
}

function createAllTodos(parent,template,data) {
	if(!data) {
		return;
	}
	for (const todo of data) {
		parent.appendChild(createTodo(template,todo));
		console.log("Adding child: ", todo);
	}
}
function createTodo(template, data) {
	const createdTodo = template.content.firstElementChild.cloneNode(true);
	const titleElem = createdTodo.querySelector("h1");
	titleElem.textContent = data.todoTitle;

	const detailElem = createdTodo.querySelector("p");
	detailElem.textContent = data.todoDetails;
	const buttonElem = createdTodo.querySelector("button");
	buttonElem.dataset.todoID = data.todoID;
	buttonElem.addEventListener("click", async (e) => {
		const deleteId = e.target.dataset.todoID;
		console.log("ID OF ELEMENT TO BE DELETED: ", deleteId);
		try {
			const response = await fetch(REMOVE_URL, {
					method: "DELETE",
					headers: {
							"Content-Type": "application/json", 
							"Access-Control-Allow-Origin": "*",
							"ElementID": `${deleteId}`,
					},
			});
			if (!response.ok) {
				throw new Error(`http error: ${response.status}`);
			}
			console.log('Delete successful:', result);
		} catch (error) {
			console.log("ERROR FROM SERVER: " + error);
		}
		e.target.parentElement.remove();
	});

	return createdTodo;
}

async function setup() {
	const todoHolder = document.querySelector("#todoContainer");
	const todos = await fetchTodos();
	createAllTodos(todoHolder,todoTemplate,todos);
}

document.querySelector("form").addEventListener("submit", async (e) => {
	e.preventDefault();
	const formData = new FormData(e.target);
	try {
		const title = formData.get("todoTitle");
		const details = formData.get("todoDetails");
		const response = await fetch(POST_URL, {
				method: "POST",
				mode: "no-cors",
				body: JSON.stringify({ todoTitle: title, todoDetails: details }),
				headers: {
						"Content-Type": "application/json", 
						"Access-Control-Allow-Origin": "*",
				},
		});
		if(!response.ok) {
			throw new Error(`server error: ${response.error}`);
		}

		const data = await response.json();
		console.log("DATA FROM SERVER: ",data);
		const todos = document.getElementById("todoContainer");
		todos.appendChild(createTodo(todoTemplate,data));
	} catch (error) {
		//document.getElementById("result").innerText = `Error: ${error.message}`;
		console.log("ERROR FROM SERVER: " + error);
	}
});

document.addEventListener('DOMContentLoaded',() => {
	setup();
});
