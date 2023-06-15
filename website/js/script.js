var bttn = document.getElementById("uploadSubmit");

function postImage()
{
    var uploadFile = document.getElementById("uploadFile");
    var file = uploadFile.files[0];

    let formData = new FormData();
    formData.append("file", file);

    console.log(file)

    const options = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            name: formData
        })
    };

    fetch('https://reqres.in/api/users/23', options)
        .then(res => res.json())
        .then(data => console.log(data))
        .catch(error => console.log('ERROR'))
}

// https://learn.javascript.ru/formdata

bttn.onclick = postImage;