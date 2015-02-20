define(["utils"],
function(utils) {

    function registerCallback(data) {

        if (data == null) {
            $("#server-answer").text("Сервер не отвечает.").css("color", "red");

        } else if (data.result === "ok") {
            $("#tab3 form").trigger("reset");
            $("#server-answer").text("Регистрация прошла успешно.").css("color", "green");
            //jsonHandle("login", loginCallback);

        } else if (data.result === "loginExists") {
            $("#server-answer").text("Такой логин уже существует.").css("color", "red");
            $("#password").val("");

        } else if (data.result === "badLogin") {
            $("#server-answer").text("Логин может содержать буквы и/или "
                + "цифры и иметь длину от 2 до 36 символов.").css("color", "red");
            $("#password").val("");

        } else if (data.result === "badPassword") {
            $("#server-answer").text("Пароль должен иметь длину от 6 "
                + "до 36 символов.").css("color", "red");
            $("#password").val("");
        }
    };

    function loginCallback(data) {
        if (data.result === "ok") {
            //$("#content, #navigation").css("display", "none");
            $("#server-answer").text("Авторизация прошла успешно.").css("color", "green");
            $("#logout-btn, #cabinet-btn").css("visibility", "visible");
            $("#password, #username").val("");

        } else if (data.result === "invalidCredentials") {
            $("#server-answer").text("Неверный логин.").css("color", "red");

        } else if (data.result === "badPassword") {
            $("#server-answer").text("Неверный пароль.").css("color", "red");
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
        if (action == "logout") {
            var js = {
                "action": "logout",
            }

        } else if (action == "login") {
            var js = {
                "action": action,
                "login": $("#tab2 #username").val(),
                "password": $("#password").val(),
            };

        } else if (action == "register") {
            var js = {
                "action": action,
                "login": $("#tab3 #username").val(),
                "password": $("#tab3 #password").val(),
                "fname": $("#tab3 #fname").val(),
                "lname": $("#tab3 #lname").val(),
                "pname": $("#tab3 #pname").val(),
            };
        }

        utils.postRequest(js, callback, "/handler");
    };

    return {
        registerCallback: registerCallback,
        loginCallback: loginCallback,
        logoutCallback: logoutCallback,
        jsonHandle: jsonHandle
    };

});