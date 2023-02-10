const image = document.createElement("img")
const params = new Proxy(new URLSearchParams(window.location.search), {
    get: (searchParams, prop) => searchParams.get(prop)
})
image.id = ""
image.className = ""
image.src = params.img_url

document.body.appendChild(image)