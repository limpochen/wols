var td_id
var td_desc

function getRecents() {
    $.ajax({
        url: '/recents.json',
        dataType: 'json',
        success: function (data) {
            recents(data)
        }
    });
}

function recents(jsons) {
    jsons.sort((a, b) => (a.last < b.last) ? 1 : -1)
    $("#recents").empty()
    var i = 1
    for (l in jsons) {
        var $tr = $("<tr></tr>")
        $tr.append(`<td style="text-align:center;">${i}</td>`)
        $tr.append(`<td style="text-align:center;font-family:monospace;">${jsons[l].mac}</td>`)
        $tr.append(`<td style="text-align:center;">${jsons[l].count}</td>`)
        $tr.append(`<td style="text-align:center;">${jsons[l].last}</td>`)
        $tr.append(`<td title="Double click to edit it" id = "id${i}" ondblclick="ModifyDesc('id${i}', '${jsons[l].mac}')">${jsons[l].desc}</td>`)
        $td = $(`<td style="text-align:center;"></td>`)
        $td.append(`<button type="button" class="button-send pure-button pure-button-primary" title="Broadcast" onclick="SendMac('${jsons[l].mac}','${jsons[l].desc}')"><i class="fa fa-power-off"></i></button>`)
        $td.append(`<button type="button" class="button-warn pure-button" style="background:rgb(255, 96, 0);" title="Remove" onclick="RemoveMac('${jsons[l].mac}')"><i class="fa fa-trash-o"></i></button>`)
        $tr.append($td)
        $("#recents").append($tr)
        i++
    }
    $("#recents tr:odd").addClass("pure-table-odd")
}

function checkMac(macAddress) {
    var regex = /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/;
    return regex.test(macAddress);
}

function SendMac(macs, desc) {
    if (macs == null) {
        macs = $("#input-mac").val()
    }
    if (!checkMac(macs)) {
        notify("Wrong MAC", `"${macs}" </br>not a IEEE 802 MAC-48 address`)
        return
    }
    var obj = { op: "broadcast", mac: macs }
    $.post("/broadcast.html", obj, function (data, status) {
        notify("发送完成", `MAC:${macs}</br>${desc}`)
        recents(JSON.parse(data))
    });
    $("#input-mac").val("")
}

function RemoveMac(macs) {
    var obj = { op: "remove", mac: macs }
    $.post("/remove.html", obj, function (data, status) {
        notify("MAC removed", "MAC: "+macs)
        recents(JSON.parse(data))
    });
}

function ModifyDesc(byId, mac) {
    var desc = $(`#${byId}`).html()
    if (desc.substring(0, 5) === '<form') {
        return
    }
    //用来关闭其他处于编辑状态的行
    if (td_id != null) {
        CancelDesc(td_id, td_desc)
    }
    td_id = byId
    td_desc = desc

    $(`#${byId}`).html(`<form class="pure-form" style="position:relative;"></form>`)
    $(`#${byId} form`).append(`<input class="pure-input-1" type ="text" id="${byId}_desc" value = "${desc}" onkeydown="KeyDesc(event,'${byId}','${desc}','${mac}')">`)
    $(`#${byId} form`).append(`<input style="display: none;">`)
    $(`#${byId} form`).append(`<button id="${byId}_ok" type="button" class="button-ok pure-button" title="Save this description" onclick="SaveDesc('${byId}','${desc}','${mac}')"></button>`)
    $(`#${byId}_ok`).append(`<i class="fa fa-check"></i>`)
    $(`#${byId}_desc`).focus()
}

function KeyMac(evt) {
    if ($("#input-mac").val() == null) {
        return
    }
    if (evt.key == "Enter") {
        SendMac()
    }
}

function KeyDesc(evt, byId, oldDesc, mac) {

    switch (evt.key) {
        case "Enter":
            SaveDesc(byId, oldDesc, mac)
            break;
        case "Escape":
            CancelDesc(byId, oldDesc)
            break;
        default:
            return;
    }
}

function SaveDesc(byId, oldDesc, macs) {
    var newDesc = $(`#${byId}_desc`).val()
    if (newDesc == oldDesc) {
        $(`#${byId}`).html(oldDesc)
        return
    }

    var obj = { op: "modify", mac: macs, desc: newDesc }
    $.post("/modify.html", obj, function (data, status) {
        $(`#${byId}`).html(data)
        notify("Describe was modified", "New describe: "+newDesc)
    })
    td_id = null
    td_desc = null
}

function CancelDesc(byId, desc) {
    $(`#${byId}`).html(desc)
}
