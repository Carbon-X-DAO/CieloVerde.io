<script context="module">
	export const prerender = true;
	import { page } from '$app/stores';
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

	<ul class="cards">
		<li class="card" class:active={$page.path === '/'}><a sveltekit:prefetch href="/">
			<span class="fa-stack fa-2x">
				<i class="fa fa-circle fa-stack-2x"></i>
				<i class="fa fa-leaf fa-stack-1x fa-inverse"></i>
			</span>
			<p>Home <i class="fa fa-chevron-right"></i></p>
		</a></li>
		<li class="card" class:active={$page.path === '/about'}><a sveltekit:prefetch href="/about">
			<span class="fa-stack fa-2x">
				<i class="fa fa-circle fa-stack-2x"></i>
				<i class="fa fa-globe fa-stack-1x fa-inverse"></i>
			</span>
			<p>About <i class="fa fa-chevron-right"></i></p>
		</a></li>
		<li class="card" class:active={$page.path === '/form'}><a sveltekit:prefetch href="/form">
			<span class="fa-stack fa-2x">
				<i class="fa fa-circle fa-stack-2x"></i>
				<i class="fa fa-pagelines fa-stack-1x fa-inverse"></i>
			</span>
			<p>Form <i class="fa fa-chevron-right"></i></p>
		</a></li>
	</ul>

</div>


<style>
	.cards {
		padding-left: 0;
		margin-top: 70px;
	}
	.card {
		list-style: none;
		text-align: center;
		background-color: rgb(35 35 35);
		margin-bottom: 15px;
		border: 2px solid #ffffff;
		border-radius: 6px;
		cursor: pointer;
		border-color: var(--primary-color);
		/* box-shadow: 0px 11px 9px 0px var(--primary-color); */
		transition: all 0.2s ease-in-out;
	}
	.card .fa-chevron-right {
		font-size: 1rem;
	}

	.card:hover {
		background-color: var(--primary-color);
		border-color: #fff;
	}
	.card:hover a {
		transition: all 0.2s ease-in-out;
		color: #383832;
	}

	.card a {
		color: var(--primary-color);
		display: block;
		padding: 15px;
		text-decoration: none;
	}
	.card p {
		font-size: 20px;
		text-transform: uppercase;
		font-weight: 700;
	}

	@media (min-width: 620px) {
		.cards {
			display: flex;
			justify-content: center;
			margin: 0 -15px;
			margin-top: 70px;
		}
		.card {
			width: calc(100% / 3 - 30px);
			margin: 0 15px;
			background-color: rgb(35 35 35 / 79%);
		}
	}

	.banner__info {
		max-width: 100%;
		margin-top: 90px;
	}
	.banner__info h1 {
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
		background-image: url('/static/images/background.jpg');
		background-repeat: no-repeat;
		background-size: cover;
		height: inherit;
		width: inherit;
		background-position: bottom left;
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
