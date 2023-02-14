const main = document.getElementById("mainframe")
const qualitites = document.getElementsByClassName("qualities")
const params = new Proxy(new URLSearchParams(window.location.search), {
    get: (searchParams, prop) => searchParams.get(prop)
})

const img = new Image()
img.onload = function () {
    main.prepend(img)
}
img.src = params.img_src
img.classList.add("main_image")


for (const v of qualitites) {
    if (params[v.id]) {
        v.innerText = params[v.id]
    }
}