document.querySelector("#create-form").addEventListener("submit", createTodo)

Array.from(document.querySelectorAll(".complete-btn")).forEach((element) => {
    element.addEventListener("click", updateTodo)
})

Array.from(document.querySelectorAll(".delete-btn")).forEach((element) => {
    element.addEventListener("click", deleteTodo)
})

async function createTodo(e) {
    e.preventDefault()

    const todoText = document.querySelector("#todo-text").value
    const todoPriority = document.querySelector("#todo-priority").value

    try {
        const res = await fetch("/api/todo", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                text: todoText,
                priority: todoPriority,
            })
        })

        if (res.status != 201) {
            throw new Error(`invalid response code: ${res.status}`)
        }
    }
    catch (e) {
        console.error(e)
    }

    e.target.reset()
    location.reload()
}

async function updateTodo(e) {
    const todoID = e.target.parentElement.id

    try {
        const res = await fetch(`/api/todo/${todoID}`, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                completed: true,
            })
        })

        if (res.status != 200) {
            throw new Error(`invalid response code: ${res.status}`)
        }
    }
    catch (e) {
        console.error(e)
    }

    location.reload()
}

async function deleteTodo(e) {
    const todoID = e.target.parentElement.id

    try {
        const res = await fetch(`/api/todo/${todoID}`, {
            method: "DELETE",
        })

        if (res.status != 204) {
            throw new Error(`invalid response code: ${res.status}`)
        }
    }
    catch (e) {
        console.error(e)
    }

    location.reload()
}
