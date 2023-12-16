var count = 0
var index = 0

function createtips() {
	$("body").append('<div id="msgpad"></div>')
	$("#msgpad").addClass("tipspad")
}

function notify(title, content) {
	if (count==0) {
		index = 0
	}

	count++
	index++
	var this_index = index

	$("#msgpad").prepend(`<div id="msg_${this_index}" class="tipsbox"></div>`)
	$(`#msg_${this_index}`).append(`<h4>${title}</h4>`)
	$(`#msg_${this_index}`).append(`<p>${content}</p>`)
	$(`#msg_${this_index}`).slideDown("slow")

	setTimeout(() => {
		$(`#msg_${this_index}`).fadeOut("slow", () => {
			$(`#msg_${this_index}`).remove()
			count--	
		})
	}, 3000);
}
