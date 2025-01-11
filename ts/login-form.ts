document.addEventListener("DOMContentLoaded", () => {
    const loginForm = document.getElementById("login-form");
    const userButton = document.getElementById("user-button");
    const unhideButton = document.getElementById("unhide-login");

    if (loginForm) {
        if (window.location.pathname === "/") {
            loginForm.classList.add("hidden");
        } else if (window.location.pathname === "/my/account") {
            loginForm.classList.remove("hidden");
        }
        if (unhideButton) {
            unhideButton.addEventListener("click", () => {
                loginForm.classList.toggle("hidden");
                userButton.classList.toggle("hidden");
            });
        }
    }
});
