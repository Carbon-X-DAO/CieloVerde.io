<script>
	let genders = ["Male", "Female", "Trans"]
</script>
<main>
	<form action="/submit" method="POST">
		<input type="text" id="fname" name="fname" placeholder="First name">
		<input type="text" id="lname" name="lname1" placeholder="Last name 1">
		<input type="text" id="lname" name="lname2" placeholder="Last name 2"><br>
		<input type="text" id="lname" name="addr" placeholder="Postal Address"><br>
		<input type="text" id="lname" name="id_no" placeholder="Identification Number"><br>
		<input type="text" id="lname" name="cell" placeholder="Telephone Number">
		<input type="text" id="lname" name="email" placeholder="Email Address"><br>
		<select name="gender">
			<option value="label" disabled selected hidden>Gender</option>
			{#each genders as gender}
			<option value={gender}>{gender}</option>
			{/each}
		</select>
		<select name="age">
			<option value="label" disabled selected hidden>Age</option>
			{#each Array(100) as _, i}
				{#if i>15}
				<option option={i}>{i}</option>
				{/if}
			{/each}
		</select><br>
		<select name="daily_qty">
			<option disabled selected hidden>Daily Quantity</option>
			<option value="first">Option 1</option>
			<option value="second">Option 2</option>
			<option value="second">Option 3</option>
			<option value="second">Option 4</option>
		</select>
		<select name="weekly_qty">
			<option disabled selected hidden>Weekly Quantity</option>
			<option value="first">Option 1</option>
			<option value="second">Option 2</option>
			<option value="second">Option 3</option>
			<option value="second">Option 4</option>
		</select>
		<select name="monthly_qty">
			<option disabled selected hidden>Monthly Quantity</option>
			<option value="first">Option 1</option>
			<option value="second">Option 2</option>
			<option value="second">Option 3</option>
			<option value="second">Option 4</option>
		</select>
		<p>
		<input type="submit" value="Submit">
		</p>
	</form>
</main>

<style>
main {
	padding: 1em;
	max-width: 1000px;
	margin: 0 auto;
}
</style>
