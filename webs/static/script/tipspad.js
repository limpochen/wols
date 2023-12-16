const tip = 0
const err = 1

var tips = {
	count: 0,
	index: 0,

	Create: function() {
		$("body").append('<div id="msgpad"></div>')
		$("#msgpad").addClass("tipspad")
	},

	Notify:	function (title, content, type = tip, time = 6) {
		if (this.count == 0) {
			this.index = 0
		}

		this.count++
		this.index++
		var this_index = this.index

		var msgbox = $(`<div id="msg-${this_index}" class="tipsbox"></div>`)
		msgbox.append(`<h5 class="boxtitle"><i class="icon-notify"></i>&nbsp;<span class="title"></span></h5>`)
		msgbox.append(`<div class="line"></div>`)
		msgbox.append(`<p class="content"></p>`)
		$("#msgpad").prepend(msgbox)
		$(`#msg-${this_index} .title`).text(''+title)
		$(`#msg-${this_index} .content`).text(content)
		if (type == err) {
			$(".icon-notify").addClass("icon-err")
		} else {
			$(".icon-notify").addClass("icon-tip")
		}
		$(`#msg-${this_index}`).slideDown("fast")

		setTimeout(() => {
			$(`#msg-${this_index}`).fadeOut("slow", () => {
				$(`#msg-${this_index}`).remove()
				this.count--
			})
		}, time * 1000);

	},

}