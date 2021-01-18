const userEmail = document.getElementById("user-email")
const userPassword = document.getElementById("user-password")
const userLoginButton = document.getElementById("user-login-button")
const content = document.getElementById("content")

if (sessionStorage.getItem("user") === null || sessionStorage.getItem("user") === "") {
    content.style.visibility = "visible"
    console.log("session is not logged in")
} else {
    console.log("session is logged in")
    verifyUser();
}

userLoginButton.addEventListener("click", () => {
    verifyUser();
})

function verifyUser() {
    let data = {
        UserEmail: userEmail.value,
        UserPassword: userPassword.value,
        UserSessionStorage: sessionStorage.getItem("user"),
    };
    fetch("/check_login", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result["Result"] === "ok") {
                console.log("good user login");
                sessionStorage.setItem("user", result["SessionLogin"])
                content.innerHTML = result["Content"];
                content.style.visibility = "visible"
            } else {
                console.log("bad user login")
            }
        });
    }).catch((error) => {
        console.log(error)
    });
}