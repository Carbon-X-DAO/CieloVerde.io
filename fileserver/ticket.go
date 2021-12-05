package fileserver

var ticketHTMLTemplate = `
<!DOCTYPE html>
<html lang="en">
    <script>
		window.alert("Save this page as PDF to save your ticket!")
	</script>
	<body>
		<div id="flyer" style="position: relative;">
			<img src="/tickets/flyer.jpg" width="1100px" height="1100px">
			<div id="frame" style="position: absolute; display: flex; align-items: center; justify-content: center; width: 170px; height: 170px; border: dashed red; bottom: 81px; left: 497px;">
				<img src="%s" style="width: 155px; height: 155px;">
			</div>
		</div>
	</body>
</html>
`
