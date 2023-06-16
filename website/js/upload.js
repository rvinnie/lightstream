var uploadBttn = document.getElementById("uploadSubmit");
var uploadFile = document.getElementById("uploadFile");

function saveImage() {
    const url = 'http://localhost:8080/images/add'
    const file = uploadFile.files[0];

    let reader = new FileReader();
    reader.readAsArrayBuffer(file)

    reader.onload = function () {
        const options = {
            method: 'POST',
            headers: {
                'Content-Type': file.type,
                'Filename': file.name,
            },
            body: reader.result
        };

        fetch(url, options)
            .then(res => res.json())
            .then(data => console.log(data))
    }
}

uploadBttn.onclick = saveImage;
