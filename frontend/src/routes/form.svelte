<script>
	let genders = ["Male", "Female", "Trans"]
</script>
<main>
	<div class="content form">
		<form action="/submit" method="POST">
			<section>
			<h2>Personal Information</h2>
			<input type="text" id="fname" name="fname" placeholder="First name" class="wide-input">
			<input type="text" id="lname" name="lname1" placeholder="Last name 1" class="narrow-input">
			<input type="text" id="lname" name="lname2" placeholder="Last name 2" class="narrow-input">
			<input type="text" id="lname" name="addr" placeholder="Postal Address" class="narrow-input">
			<input type="text" id="lname" name="id_no" placeholder="Ident. Number" class="narrow-input">
			<input type="text" id="lname" name="cell" placeholder="Telephone Number" class="wide-input">
			<input type="text" id="lname" name="email" placeholder="Email Address" class="wide-input">
			<select name="gender" class="narrow-input select select-placeholder">
				<option value="label" disabled selected hidden>Gender</option>
				{#each genders as gender}
				<option value={gender}>{gender}</option>
				{/each}
			</select>
			<select name="age" class="narrow-input select select-placeholder">
				<option value="label" disabled selected hidden>Age</option>
				{#each Array(100) as _, i}
					{#if i>15}
					<option option={i}>{i}</option>
					{/if}
				{/each}
			</select>
			</section>
			<section>
			<h2>Marijuana Consumption</h2>
			<select name="daily_qty" width="10px" class="narrow-select select select-placeholder">
				<option disabled selected hidden>Daily</option>
				<option value="first">Option 1</option>
				<option value="second">Option 2</option>
				<option value="second">Option 3</option>
				<option value="second">Option 4</option>
			</select>
			<select name="weekly_qty" class="narrow-select select select-placeholder">
				<option disabled selected hidden>Weekly</option>
				<option value="first">Option 1</option>
				<option value="second">Option 2</option>
				<option value="second">Option 3</option>
				<option value="second">Option 4</option>
			</select>
			<select name="monthly_qty" class="narrow-select select select-placeholder">
				<option disabled selected hidden>Monthly</option>
				<option value="first">Option 1</option>
				<option value="second">Option 2</option>
				<option value="second">Option 3</option>
				<option value="second">Option 4</option>
			</select>
			
			<div class="btn-submit-wrap">
				<input type="submit" value="Submit" class="btn-filled">
			</div>

			</section>
		</form>
	</div>
</main>

<style>
.content {
	width: 100%;
    max-width: 660px;
    margin: 0 auto;
    padding: 0 15px;
    margin: 2rem auto 1rem;
}

h2 {
	font-size: 1.5rem;
	color: #ffffff;
	text-align: center;
	font-weight: 600;
}

/* hide default arrow */
select {
  -webkit-appearance: none;
  -moz-appearance: none;
  text-indent: 1px;
  text-overflow: '';
}
select::-ms-expand {
    display: none;
}

.select {
  background-image:
    linear-gradient(45deg, transparent 50%, #fff 50%),
    linear-gradient(135deg, #fff 50%, transparent 50%),
    linear-gradient(to right, #fff, #fff);
  background-position:
    calc(100% - 20px) calc(17px + 2px),
    calc(100% - 15px) calc(17px + 2px),
    calc(100% - 40px) 8px;
  background-size:
    5px 5px,
    5px 5px,
    1px 27px;
  background-repeat: no-repeat;
}

.select:focus {
  background-image:
    linear-gradient(45deg, var(--accent-color) 50%, transparent 50%),
    linear-gradient(135deg, transparent 50%, var(--accent-color) 50%),
    linear-gradient(to right, var(--accent-color), var(--accent-color));
  background-position:
    calc(100% - 15px) 17px,
    calc(100% - 20px) 17px,
    calc(100% - 40px) 8px;
  background-size:
    5px 5px,
    5px 5px,
    1px 27px;
  background-repeat: no-repeat;
  border-color: var(--accent-color);
  outline: 0;
}


.select:-moz-focusring {
  color: transparent;
  text-shadow: 0 0 0 #000;
}

input, select {
	background-color: transparent;
	border: none;
	border: 2px solid #fff;
	height: 48px;
	width: 100%;
	padding: 2px 10px;
	border-radius: 6px;
	box-sizing: border-box;
	font-size: 16px;
	color: #fff;
	margin-bottom: 15px;
	outline: transparent;
	font-weight: 600;
	cursor: pointer;
}

input::placeholder, 
select::placeholder {
	color: #fff;
	opacity: 0.8;
	font-size: 14px;
	font-weight: 300;
}
input:focus::placeholder,
select:focus::placeholder {
	opacity: 1;
}
input:focus,
select:focus {
	border-color: var(--accent-color);
}

.select-placeholder {
	font-weight: 300;
	font-size: 14px;
	opacity: 0.8;
}

.btn-submit-wrap {
	display: flex;
	justify-content: flex-end;
	width: 100%;
}

.btn-filled {
	background: #fff;
	color: #000;
	width: 160px;
}

@media (min-width: 620px) {
	section {
		display: flex;
		flex-wrap: wrap;
		justify-content: space-between;
	}
	.narrow-input {
		width: calc(50% - 10px);
	}
}

</style>
