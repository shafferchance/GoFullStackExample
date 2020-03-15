window.addEventListener("load", () => {
    document.querySelector(".refresh").addEventListener("click", handleRequet);
    document.querySelector("input[type='submit']").addEventListener("click", handleSubmit);
    document.querySelector("input[type='text']").addEventListener("keydown", e => {
        if (e.key === "Enter") {
            handleSubmit(e);
        }
    });
});

const handleRequet = e => {
    e.preventDefault();

    fetch(`/input`)
        .then(arr => arr.json())
        .then(data => {
            const lst = document.querySelector("ul");
            const newLst = document.createElement("ul");
            for (const ele of data["data"]) {
                let li = document.createElement('li');
                li.innerText = ele;
                newLst.append(li);
            }
            lst.replaceWith(newLst);
            document.querySelector("input[type='text']").focus()
        })
}

const handleSubmit = e => {
    e.preventDefault();
    
    const txt = document.querySelector("input[type='text']");
    console.log(e)

    fetch(`/input`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({text: txt.value})
    }).then(name => name.text())
    .then(imp => console.log(imp))
    .catch(err => console.log(err))

    // li.innerText = e.target.value;
    // lst.append(lst);
    txt.value = "";
    handleRequet(e);
}