package fileserver

var ticketHTMLTemplate = `
<!DOCTYPE html>
<html lang="en">
    <script>
		window.alert("¡Tome una captura de pantalla del código QR para guardar boleto!")
	</script>
	<body>
	<div id="flyer" style="position: relative;">
		<img src="/tickets/flyer.jpg" width="1100px" height="1100px">
		<div id="frame" style="position: absolute;display: flex;align-items: center;justify-content: center;width: 170px;height: 170px;border: dashed red;bottom: 80px;left: 495px;border-width: thick;">
			<img src="%s" style="width: 160px; height: 160px;">
		</div>
	</div>
	</body>
</html>
`
