import { api } from '/js/api.js';

window.addEventListener("DOMContentLoaded", function() {
    document.getElementById("login").addEventListener("click", login);
});


const login = async () => {
    const form = document.getElementById("login-form");
    const username = form.elements['username'].value;
    const password = form.elements['password'].value;

    const body = {
        username: username,
        password: password
    };

    try {
        await api.post('login', body);
        window.location.replace('/');
    } catch (e) {
        document.getElementById("error").innerHTML = (e.status === 401)
        ? "ユーザ名またはパスワードが異なります。" 
        : "ログインに失敗しました。";
    }
}