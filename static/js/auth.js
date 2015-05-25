define(["utils"],
function(utils) {

    function loginCallback(data) {
        if (data.result === "ok") {
            $("#events #server-answer").empty();
            $("#server-answer").text("Авторизация прошла успешно.").css("color", "green");
            $("#logout-btn, #cabinet-btn").css("visibility", "visible");
            $("#password, #username").val("");

        } else if (data.result === "invalidCredentials") {
            $("#server-answer").text("Неверный логин.").css("color", "red");

        } else if (data.result === "badPassword") {
            $("#server-answer").text("Неверный пароль.").css("color", "red");
        } else if (data.result === "notEnabled") {
            $("#server-answer").text("Ваш аккаунт заблокирован.").css("color", "red");
        }
    };

    function logoutCallback(data) {
        if (data.result === "ok") {
            $("#server-answer").text("Вы вышли.").css("color", "green").css("visibility", "visible");
            location.href = "/"
            
        } else if (data.result === "badSid") {
            $("#server-answer").text("Invalid session ID.").css("color", "red");
        }
    };

    function jsonHandle(action, callback) {
        var js = {};
        js["action"] = action;

        if (action == "login") {
            js["login"] = $("#content #username").val();
            js["password"] = $("#content #password").val();

            // js["login"] = "admin";
            // js["password"] = "password";
        }

        utils.postRequest(js, callback, "/handler");
    };

    return {
        loginCallback: loginCallback,
        logoutCallback: logoutCallback,
        jsonHandle: jsonHandle
    };

});