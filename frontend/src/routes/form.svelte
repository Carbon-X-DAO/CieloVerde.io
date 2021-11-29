<script>
	let genders = ["Male", "Female", "Trans"]
</script>
<main>
	<form action="/submit" method="POST">
		<section>
		<h3>Personal Information</h3>
		<input type="text" id="fname" name="fname" placeholder="First name" class="wide-input"><br>
		<input type="text" id="lname" name="lname1" placeholder="Last name 1" class="narrow-input">
		<input type="text" id="lname" name="lname2" placeholder="Last name 2" class="narrow-input"><br>
		<input type="text" id="lname" name="addr" placeholder="Postal Address" class="narrow-input"s>
		<input type="text" id="lname" name="id_no" placeholder="Ident. Number" class="narrow-input"><br>
		<input type="text" id="lname" name="cell" placeholder="Telephone Number" class="wide-input">
		<input type="text" id="lname" name="email" placeholder="Email Address" class="wide-input"><br>
		<select name="gender" class="wide-select placeholder">
			<option value="label" disabled selected hidden>Gender</option>
			{#each genders as gender}
			<option value={gender}>{gender}</option>
			{/each}
		</select>
		<select name="age" class="wide-select placeholder">
			<option value="label" disabled selected hidden>Age</option>
			{#each Array(100) as _, i}
				{#if i>15}
				<option option={i}>{i}</option>
				{/if}
			{/each}
		</select><br>
		</section>
		<section>
		<h3>Marijuana Consumption</h3>
		<select name="daily_qty" width="10px" class="narrow-select placeholder">
			<option disabled selected hidden>Daily</option>
			<option value="first">Option 1</option>
			<option value="second">Option 2</option>
			<option value="second">Option 3</option>
			<option value="second">Option 4</option>
		</select>
		<select name="weekly_qty" class="narrow-select placeholder">
			<option disabled selected hidden>Weekly</option>
			<option value="first">Option 1</option>
			<option value="second">Option 2</option>
			<option value="second">Option 3</option>
			<option value="second">Option 4</option>
		</select>
		<select name="monthly_qty" class="narrow-select placeholder">
			<option disabled selected hidden>Monthly</option>
			<option value="first">Option 1</option>
			<option value="second">Option 2</option>
			<option value="second">Option 3</option>
			<option value="second">Option 4</option>
		</select>
		<p>
		<input type="submit" value="Submit">
		</p>
		</section>
	</form>
</main>

<style>
main {
	max-width: 500px;
	margin: 0 auto;
	font-size: 14px;
	/*width: 85%;*/
}

h3 {
	color: #ffffff;
	text-align: center;
}

form {
	/*border: 2px solid #eeeeee;*/
	border-radius: 2px;
	padding: 10px;
}

section {
	padding: 0px 10px;
	margin: auto;
}

select.narrow-select {
	width: 100px;
}

select.wide-select {
	width: 100px;
}

input.narrow-input {
	width: 140px;
}

input.wide-input {
	width: 200px;
}

select,input {
	/*border: 2px solid #BBBBBB;*/
	/*background-color: #E4E4E4;*/
	height: 20px;
	background-color: #ffffff;
	border: none;
	border-bottom: 2px solid #BBBBBB;
	margin: 0px 10px 10px 0px;
	color: #ffffff;
	margin-bottom: 30px;
}

input:focus-visible{
	outline: none;
}

select,::placeholder {
	/*color: #666666;*/
	color: rgb(102,102,102,.6);
	font-weight: 800;
	text-transform: uppercase;
}

/*
select,input,::placeholder{
	font-size: 16px;
}
*/

/*
.placeholder,::placeholder {
	font-size:12px;
}
*/

input[type=submit] {
	background-color: #726868;
	border: none;
	box-shadow: 3px 3px 1px #ccc;
	text-shadow: none;
	/*border-radius: 4px;*/
	color: #DDDDDD;
	height: 30px;
	width: 80px;
	margin: 0 auto;
}

p {
	display: flex;
}

input[type=submit]:focus, input[type=submit]:active, input[type=submit]:hover {
	background-color: #504534;
	cursor: pointer;
}

@media (max-width: 500px) {
	select.narrow-select {
		width: 100px;
	}

	select,input {
		font-size: 13px;
		margin-bottom: 20px;
	}
}

@media (max-width: 400px) {
	select.narrow-select {
		width: 90px;
	}

	select,input,::placeholder {
		font-size: 12px;
	}
}

@media (max-width: 375px) {
	select.narrow-select {
		width: 150px;
	}
}
</style>
