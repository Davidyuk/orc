define(["utils", "grid-utils", "blank"], function(utils, gridUtils, blank) {

    function AddPerson(dialogId, groupId) {
        $("#"+dialogId).empty();

        var block = $("<div/>");

        var lf = $("<label/>", { "text": "Фамилия" });
        var f = $("<input/>", { "id": 5, "for-saving": true, "required": true });
        block.append(lf).append(f);

        var li = $("<label/>", { "text": "Имя" });
        var i = $("<input/>", { "id": 6, "for-saving": true, "required": true });
        block.append(li).append(i);

        var lo = $("<label/>", { "text": "Отчество" });
        var o = $("<input/>", { "id": 7, "for-saving": true, "required": true });
        block.append(lo).append(o);

        var le = $("<label/>", { "text": "Email" });
        var e = $("<input/>", { "id": 4, "for-saving": true, "required": true });
        block.append(le).append(e);

        $("#"+dialogId).append(block);

        $("#"+dialogId).dialog({
            modal: true,
            toTop: "150",
            height: "auto",
            width: "auto",
            buttons: {
                "Добавить участника": function() {
                    var values = blank.getFormData(dialogId);
                    if (!values) {
                        return false;
                    }
                    var data = {"group_id": groupId, "data": values };
                    console.log("AddPerson: ", data);
                    utils.postRequest(
                        data,
                        function(response) {
                            gridUtils.showServerPromtInDialog($("#"+dialogId), response["result"]);
                        },
                        "/handler/addperson"
                    );
                },
                "Отмена": function() {
                    $(this).dialog("close");
                },
            }
        });
    }

    return {
        AddPerson: AddPerson,
    };

});