const searchBttn = document.getElementById("searchSubmit");
const searchAllBttn = document.getElementById("searchAll");
const searchElement = document.getElementById("searchInput");
const uploadBttn = document.getElementById("uploadSubmit");
const uploadFile = document.getElementById("uploadFile");

const url = 'http://localhost:8080'

// Notifications
const NotifyStatuses = { SUCCESS: 'success', ERROR: 'error', WARNING: 'warning' };

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

// Getting single image
async function searchImage() {
    const searchId = searchElement.value

    if (searchId == "") {
        pushNotify(NotifyStatuses.ERROR, `Enter image id`)
        return
    }

    const uri = `${url}/images/${searchId}`
    const response = await fetch(uri, { method: "GET" })

    if (response.status !== 200) {
        pushNotify(NotifyStatuses.ERROR, `Image with such id does not exist `)
        return
    }

    const imageJson = await response.json()


    new Fancybox([
        {
            src: "data:" + imageJson.contentType + ";base64," + imageJson.data,
            type: "image",
        },
    ], {hideScrollbar: false});
}

// Getting all images
async function searchImages() {
    const uri = `${url}/images`
    const response = await fetch(uri, { method: "GET" })

    if (response.status !== 200) {
        pushNotify(NotifyStatuses.ERROR, 'Unable to get images')
        return
    }

    const images = await response.json()
    let galleryItems = [];

    if (images.length == 0) {
        pushNotify(NotifyStatuses.WARNING, 'Storage is empty')
        return
    }

    for (let i = 0; i < images.length; i++) {
        let galleryItem = {
            src: "data:" + images[i].contentType + ";base64," + images[i].data,
            type: "image",
        }
        galleryItems.push(galleryItem)
    }

    new Fancybox(galleryItems, {hideScrollbar: false})
}


// Saving image script
function saveImage() {
    const uri = `${url}/images/add`
    const file = uploadFile.files[0];

    if (file == null) {
        pushNotify(NotifyStatuses.ERROR, `Choose file`)
        return
    }

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
        const id = await response.json()

        if (response.status === 201) {
            pushNotify(NotifyStatuses.SUCCESS, `Image successful uploaded with id=${id}`)
        } else {
            pushNotify(NotifyStatuses.ERROR, `Unable to upload file`)
        }
    }
}

uploadBttn.onclick = saveImage;
searchBttn.onclick = searchImage;
searchAllBttn.onclick = searchImages;
