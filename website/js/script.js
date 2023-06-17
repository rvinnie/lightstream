const searchBttn = document.getElementById("searchSubmit");
const searchAllBttn = document.getElementById("searchAll");
const searchElement = document.getElementById("search");
const imageContainer = document.getElementById("imageContainer")
const imagesContainer = document.getElementById("imagesContainer")
const uploadBttn = document.getElementById("uploadSubmit");
const uploadFile = document.getElementById("uploadFile");
const url = 'http://localhost:8080'

// Notifications
const NotifyStatuses = { SUCCESS: 'success', ERROR: 'error' };

function pushNotify(status, title) {
    let myNotify = new Notify({
        status: status,
        title: title,
        effect: 'slide',
        autoclose: true,
        autotimeout: 3000,
        type: 3
    })
}

// Searching images scripts
async function searchImage() {
    const searchId = searchElement.value

    if (searchId == "") {
        console.log("HTTP-Error: Bad Request")
        imagesContainer.innerHTML = ''
        imageContainer.innerHTML = ''
        pushNotify(NotifyStatuses.ERROR, `Enter resource id `)
        return
    }

    const uri = `${url}/images/${searchId}`
    const options = {
        method: "GET"
    }

    const response = await fetch(uri, options)

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
        pushNotify(NotifyStatuses.ERROR, `Image with such id does not exist `)
    }
}

async function searchImages() {
    const uri = `${url}/images`
    const options = {
        method: "GET"
    }

    const response = await fetch(uri, options)

    if (response.status !== 200) {
        imagesContainer.innerHTML = ''
        imageContainer.innerHTML = ''
        pushNotify(NotifyStatuses.ERROR, 'Unable to get images')
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


// Saving image script
function saveImage() {
    const uri = `${url}/images/add`
    const file = uploadFile.files[0];

    let reader = new FileReader();
    reader.readAsArrayBuffer(file)

    reader.onload = async function () {
        const options = {
            method: 'POST',
            headers: {
                'Content-Type': file.type,
                'Filename': file.name,
            },
            body: reader.result
        };

        const response = await fetch(uri, options)

        if (response.status === 201) {
            pushNotify(NotifyStatuses.SUCCESS, `Image successful uploaded`)
        } else {
            pushNotify(NotifyStatuses.ERROR, `Unable to upload file`)
        }
    }
}

uploadBttn.onclick = saveImage;
searchBttn.onclick = searchImage;
searchAllBttn.onclick = searchImages;
