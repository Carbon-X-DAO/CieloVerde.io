<script context="module">
	export const prerender = true;
</script>

<svelte:head>
	<title>Home</title>
</svelte:head>

<div class="container">
	<div class="banner">
		<div class="banner__img">
			<div class="background"></div>
		</div>
		<div class="banner__info">
			<h1>Plant<br>Medicine</h1>
			<p>Lorem, ipsum dolor sit amet consectetur adipisicing elit. Ad vero ratione asperiores impedit voluptas iusto libero expedita placeat error soluta.</p>
			<button class="ghost-btn">Read More</button>
		</div>
		
	</div>
</div>


<style>
	section {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		flex: 1;
	}

	.banner__info {
		max-width: 100%;
		margin-top: 90px;
	}
	.banner__info h1,
	.banner__info h2,
	.banner__info h3,
	.banner__info h4,
	.banner__info h5,
	.banner__info h6 {
		margin-bottom: 0;
	}

	.banner__info .ghost-btn {
		margin-top: 24px;
		padding: 0.8rem 2rem;
	}

	.banner__img {
		position: absolute;
		right: 0;
		top: -50px;
		display: block;
		z-index: -1;
		width: 100%;
		height: 100vh;
		box-shadow: 0 0 10px 10px #86B27522 inset;
		border-radius: 3px;
	}

	/* overlay above image */
	.banner__img::before {
		content: "";
		position: absolute;
		display: block;
		height: 100%;
		width: 100%;
		background-color: rgb(10 8 8 / 45%);
	}

	.background {
		background-image: url('./flower.jpg');
		background-repeat: no-repeat;
		background-size: cover;
		height: inherit;
		width: inherit;
	}

	@media (min-width: 320px) {
		h1 {
			text-align: left;
			color: #ffffff;
			font-size: 38px;
			font-weight: bold;
		}
	}
	@media (min-width: 600px) {
		.banner__info {
			max-width: 50%;
			margin-top: 120px;
		}
		h1 {
			font-size: 88px;
		}
		.banner__img {
			width: 79vw;
			height: 700px;
		}
	}

</style>
