var searchBttn = document.getElementById("searchSubmit");
var searchAllBttn = document.getElementById("searchAll");
var searchElement = document.getElementById("search");
const imageContainer = document.getElementById("imageContainer")
const imagesContainer = document.getElementById("imagesContainer")


async function searchImage() {
    const url = 'http://localhost:8080/images/'
    const searchId = searchElement.value
    const options = {
        method: "GET"
    }

    if (searchId == "") {
        console.log("HTTP-Error: Bad Request")
        imagesContainer.innerHTML = ''
        imageContainer.innerHTML = 'Not Found'
        return
    }

    const response = await fetch(url + searchId, options)

    if (response.status === 200) {
        const imageBlob = await response.blob()
        const imageObjectURL = URL.createObjectURL(imageBlob);

        const image = document.createElement('img')
        image.src = imageObjectURL
        image.className = 'single-image'

        imageContainer.innerHTML = ''
        imagesContainer.innerHTML = ''
        imageContainer.append(image)
    } else {
        imageContainer.innerHTML = 'Not Found'
    }
}

async function searchImages() {
    const url = 'http://localhost:8080/images'
    const options = {
        method: "GET"
    }

    const response = await fetch(url, options)

    if (response.status !== 200) {
        imagesContainer.innerHTML = 'Not Found'
        imageContainer.innerHTML = ''
        return
    }

    imageContainer.innerHTML = ''
    imagesContainer.innerHTML = ''
    const images = await response.json()

    for (let i = 0; i < images.length; i++) {
        const image = document.createElement('img')
        image.src = "data:" + images[i].contentType + ";base64," + images[i].data;
        image.className = 'image'

        const imageBox = document.createElement('div')
        imageBox.className = 'img-box'
        const transparentBox = document.createElement('div')
        transparentBox.className = 'transparent-box'
        const caption = document.createElement('div')
        caption.className = 'caption'

        const imageName = document.createElement('p')
        imageName.innerText = images[i].name
        caption.append(imageName)

        transparentBox.append(caption)

        imageBox.append(image)
        imageBox.append(transparentBox)

        imagesContainer.append(imageBox)
    }
}

searchBttn.onclick = searchImage;
searchAllBttn.onclick = searchImages;
