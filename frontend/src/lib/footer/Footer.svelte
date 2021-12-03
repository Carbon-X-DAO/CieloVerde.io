<script>
    let year = new Date().getFullYear()  // returns the current year
</script>

<footer class="footer">
	<div class="container">
		<div class="row">
			<div class="col">
				<a href="mailto:example@gmail.com"><i class="fa fa-envelope"></i> example@gmail.com</a>
				<a href="mailto:github.com"><i class="fa fa-github"></i> example@gmail.com</a>
			</div>
			<div class="col">
				<p class="copyright">With <i class="fa fa-heart"></i> by © CIELO VERDE {year}</p>
			</div>
		</div>
	</div>
</footer>

<style>
    footer {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		background: #232323;
    	border-top: 2px solid var(--primary-color);
		padding: 15px 0;
	}

	footer .col {
		margin-bottom: 10px;
	}
	footer .col:last-child {
		margin-bottom: 0;
	}

	footer .col a {
		display: flex;
		align-items: center;
		line-height: 1.4;
		padding: 6px;
    	padding-left: 0;
	}
	.col a .fa {
		margin-right: 5px;
	}

	footer a {
		color: #fff;
		margin-right: 10px;
		text-decoration: none;
	}
	footer a:hover {
		color: var(--accent-color);
	}

	.row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		flex-direction: column;
    	align-items: flex-start;
	}

	.copyright {
		font-weight: 700;
		font-size: 0.8rem;
		line-height: 1;
	}

	@media (min-width: 480px) {
		footer {
			padding: 10px 0;
		}
	}
	@media (min-width: 620px) {
		footer .col a {
			display: inline-flex;
		}
		.row {
			flex-direction: row;
			align-items: center;
		}
		.row .col {
			margin-bottom: 0;
		}
	}
</style>