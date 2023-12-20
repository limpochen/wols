function randMac() {
    var rmac = ""
    $("#input-mac").val("")
    for (i = 0; i <= 5; i++) {
        rmac += ('0' + ("%02X", Math.floor((Math.random() * 254)) + 1).toString(16).toUpperCase()).slice(-2);

        if (i < 5) {
            rmac += "-";
        }
    }
    $("#input-mac").val(rmac)
}

// Start here:
$(document).ready(function () {
    $("#input_mac").focus();

    tips.Create();
    recents.Load();
    $("#new-send form").append('<button id="button-rand" type="button" onclick="randMac()">rand...</button>');

    $("#main-send").on("click", () => {
        if ($.trim($("#input-mac").val()) == "") {
            return;
        }
        recents.sendMacNew();
    })

    $("#input-mac").on("keydown", (event) => {
        if ($.trim($("#input-mac").val()) == "") {
            return;
        }

        if (event.key == "Enter") {
            recents.sendMacNew();
        }
    })

    $(`div#recents`).on('dblclick', 'td', (e) => {
        var ops = e.currentTarget.id.split("-");
        if (ops[0] != 'desc') {
            return;
        }
        recents.modifyDesc(ops[1]);
    })

    $(`div#recents`).on('click', 'button', (e) => {
        var ops = e.currentTarget.id.split("-");
        switch (ops[0]) {
            case "sendmac":
                recents.SendMac($(`#tr-${ops[1]} .mac`).text(), $(`#spandesc-${ops[1]}`).text());
                break;

            case "removemac":
                recents.RemoveMac(ops[1]);
                break;

            case "noremove":
                recents.noRemove(ops[1])
                break;

            case "savedesc":
                recents.SaveDesc(ops[1]);
                break;

            case "modifydesc":
                recents.modifyDesc(ops[1]);
            default:
        }
    })

    $(`div#recents`).on("keydown", 'input', (evt) => {
        var ops = evt.currentTarget.id.split("-");
        switch (evt.key) {
            case "Enter":
                recents.SaveDesc(ops[1]);
                break;
            case "Escape":
                recents.CancelDesc(ops[1]);
                break;
            default:
                return;
        }
    })

})
