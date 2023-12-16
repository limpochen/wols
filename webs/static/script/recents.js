
var recents = {
    recents: null, //json 格式

    Load: function () {
        this.getRecents();
        $("div#recents").empty();
        var i = 0;
        for (i in this.recents) {
            i++;
            break;
        }
        if (i > 0) {
            $("div#recents").html(this.recentHead("#", "MAC", "Cnt", "Last", "Desc", "Op"));
            this.recentBody();
        }

        $(`div#recents .desc`).on("mouseover", (evt) => {
            var ops = evt.currentTarget.id.split("-");
            if ($(`#formdesc-${ops[1]}`).is(':hidden')) {
                $(`#modifydesc-${ops[1]}`).show();
            } else {
                $(`#modifydesc-${ops[1]}`).hide();
            }
        })
        $(`div#recents .desc`).on("mouseleave", (evt) => {
            var ops = evt.currentTarget.id.split("-");
            $(`#modifydesc-${ops[1]}`).hide();
        })
    },

    // todo: 返回需判断是否读取成功
    getRecents: function () {
        $.ajax({
            async: false, // todo
            type: 'POST',
            url: 'recents',
            dataType: 'json',
            success: function (data) {
                recents.recents = data;
            },
            error: function () {
                return;
            }
        });
    },

    recentHead: function (idx, mac, count, last, desc, op) {
        return `<table>
                    <caption><h4>Recents:</h4></caption>
                    <thead>
                        <tr>
                            <th>${idx}</th>
                            <th>${mac}</th>
                            <th>${count}</th>
                            <th>${last}</th>
                            <th>${desc}</th>
                            <th>${op}</th>
                        </tr>
                    </thead>
                    <tbody>
                    </tbody>
                </table>`;
    },

    recentBody: function () {
        this.recents.sort((a, b) => (a.last < b.last) ? 1 : -1);
        for (idx in this.recents) {
            $("div#recents table tbody").append(
                `<tr id="tr-${idx}">
                    <td class="idx">${parseInt(idx) + 1}</td>
                    <td class="mac">${this.recents[idx].mac}</td>
                    <td class="count">${this.recents[idx].count}</td>
                    <td class="last">${this.recents[idx].last}</td>
                    <td id="desc-${idx}" class="desc" title="Double click to modify this description">
                        <span id="spandesc-${idx}"></span>
                        <form id="formdesc-${idx}">
                            <input id="newdesc-${idx}" type="text" class="newdesc" name="newdesc" value="${this.recents[idx].desc}">
                            <input type="text" style="display:none">
                            <button id="savedesc-${idx}" type="button" class="savedesc" title="Save this description">
                                <i class="icon-submit"></i>
                            </button>
                        </form>
                        <button id="modifydesc-${idx}" class="modifydesc" title="Modify this description"><i class="icon-modify"></i></button>
                    </td>
                    <td class="op">
                        <button id="sendmac-${idx}" class="sendmac" type="button" title="Broadcast this hardware address"><i class="icon-send"></i></button>
                        <button id="removemac-${idx}" class="removemac" type="button" title="Remove this hardware address"><i class="icon-remove"></i></button>
                        <button type="button" id="noremove-${idx}" title="Cancle to remove it" class="noremove"><i class="icon-noremove"></i></button>
                    </td>
                </tr>`
            );
            $(`#spandesc-${idx}`).text(this.recents[idx].desc);
            $(`#modifydesc-${idx}`).hide();
            $(`#formdesc-${idx}`).hide();
            $(`#noremove-${idx}`).hide();
        }

    },

    RemoveMac: function (idx) {
        if ($(`#tr-${idx}`).attr("class") != "preremove") {
            this.cleanStatus();
            $(`#sendmac-${idx}`).hide()
            $(`#noremove-${idx}`).show()
            $(`#tr-${idx}`).addClass("preremove");  //blink
            return;
        }

        $(`#tr-${idx}`).removeClass("preremove");
        $(`#tr-${idx}`).addClass("remove");
        let mac = $(`#tr-${idx} .mac`).html()
        let desc = $(`#tr-${idx} .desc`).text()
        setTimeout(() => {
            var obj = { mac: mac };
            $.post("/remove", obj, function (data, status) {
                tips.Notify("Remove " + status, `MAC: ${mac}(${desc})`);
                var resp = JSON.parse(data);
                if (resp.status == "error") {
                    tips.Notify("Error from Server", resp.extra, err, 10)
                }
                recents.Load();
            }).fail(function (n) {
                tips.Notify("Remove failed", `MAC: ${macs}(${newDesc})`, err, 20);
            });
        }, 100);
    },

    noRemove: function (idx) {
        $(`#sendmac-${idx}`).show()
        $(`#noremove-${idx}`).hide()
        $(`#tr-${idx}`).removeClass("preremove");
    },

    sendMacNew: function () {
        var regExp = /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/;
        var mac = $("#input-mac").val();
        var macDesc = "Sented from WEB";
        if (!regExp.test(mac)) {
            tips.Notify("Bad Address", `"${mac}" Not a IEEE 802 MAC-48 address`, err);
            return;
        }
        recents.SendMac(mac, macDesc);
    },

    SendMac: function (mac, macDesc) {
        var obj = { mac: mac, desc: macDesc };

        $.ajaxSettings.timeout = '3000';
        $.post("/broadcast", obj, function (data, status) {
            tips.Notify("Send " + status, `MAC: ${obj.mac}(${obj.desc})`);
            var resp = JSON.parse(data);
            if (resp.status == "error") {
                tips.Notify("Error from Server", resp.extra, err, 10)
            }
            $("#input-mac").val("");
            recents.Load();
        }).fail(function () {
            tips.Notify("Send failed", `MAC: ${obj.mac}(${obj.desc})`, err, 10);
        })
    },

    modifyDesc: function (idx) {
        this.cleanStatus();
        $(`#modifydesc-${idx}`).hide();
        $(`#spandesc-${idx}`).hide();
        $(`#formdesc-${idx}`).show();
        $(`#newdesc-${idx}`).focus();
        $(`#newdesc-${idx}`).select();
    },

    SaveDesc: function (idx) {
        var mac = $(`#tr-${idx} .mac`).html()
        var desc = $.trim($(`#newdesc-${idx}`).val())
        if (desc == $(`#spandesc-${idx}`).text() || desc == ""){
            recents.CancelDesc(idx);
            return;    
        } 
        var obj = { mac: mac, desc: desc }
        $.post("/modify", obj, function (data, status) {
            tips.Notify(`Modified ${status}`, `MAC: ${mac}(${desc})`);
            var resp = JSON.parse(data);
            if (resp.status == "error") {
                tips.Notify("Error from Server", resp.extra, err, 10)
            }
            recents.Load();
        }).fail(function (n) {
            tips.Notify("Modify failed", `MAC: ${mac}(${desc})`, err, 20);
        })
    },

    CancelDesc: function (idx) {
        $(`#formdesc-${idx}`).hide();
        $(`#spandesc-${idx}`).show();
    },

    cleanStatus: function () {
        var idx;
        for (idx in this.recents) {
            if ($(`#spandesc-${idx}`).is(':hidden')) {
                this.CancelDesc(idx);
            }
            if ($(`#tr-${idx}`).attr("class") == "preremove") {
                this.noRemove(idx)
            }
        }
    },

}
