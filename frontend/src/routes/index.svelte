<script context="module">
	export const prerender = true;
</script>

<svelte:head>
	<title>Home</title>
</svelte:head>

<section>
	<div class="bg-img"></div>
	<h1>Plant<br>&nbsp;&nbsp;Medicine</h1>
</section>

<style>
	section {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		flex: 1;
	}

	.bg-img {
		background-image: url('./flower.jpg');
		position: absolute;
		right: 0;
		/*top: 50px;*/
		bottom: 0px;
		display: block;
		z-index: -1;
		width: 800px;
		height: 600px;
		box-shadow: 0 0 10px 10px #86B27522 inset;
		border-radius: 3px;
	}

	h1 {
		text-align: left;
		color: #F5F5DCAA;
		position: absolute;
		top: 0px;
		left: 20px;
		font-size: 72px;
		font-weight: bold;
	}

	.welcome {
		position: relative;
		width: 100%;
		height: 0;
		padding: 0 0 calc(100% * 495 / 2048) 0;
	}

	.welcome img {
		position: absolute;
		width: 100%;
		height: 100%;
		top: 0;
		display: block;
	}
</style>
