var token;
// REQUEST API
async function restlogin(body) {
    var myInit = {
        method: "POST",
        headers: headers(),
        mode: "cors",
        cache: "default",
        body: body,
    };
    const response = await fetch(
        `${window.location.protocol}//${window.location.host}/login`,
        myInit
    );
    // console.log(response.status)
    if (response.status == "200") {
        token = response.headers.get("Token");
        username = response.headers.get("Username");
        document.cookie = ("Username=" + username + "; path=/");
        document.cookie = ("Token=" + token + "; path=/");
        location.assign(`${window.location.protocol}//${window.location.host}/html/HomeScreen.html`)
    } else {
        return alert("Credenciais Inválidas");
    }
    const data = JSON.parse(await (await response.blob()).text());
    console.log(data)
}

// BOTÃO DE ENVIO 
function enviarLogin() {
    var username = document.getElementById("userEmail").value
    username = username.toLowerCase()
    var userpassword = document.getElementById("userPassword").value
    const body = { "username": username, "userpassword": userpassword }
    data = JSON.stringify(body)
    restlogin(data);
}
// SET TOKEN ON LOCALSTORAGE VS COOKIE
function headers() {
    var h = new Headers()
    return h
}