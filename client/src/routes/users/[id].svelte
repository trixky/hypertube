<script lang="ts">
	import BlackBox from '../../components/containers/black-box.svelte';
	import Warning from '../../components/inputs/warning.svelte';
	import Eye from '../../components/inputs/eye.svelte';
	import * as sanitzer from '../../utils/sanitizer';
	import ConfirmationButton from '../../components/buttons/confirmation-button.svelte';
	import * as cookies from '../../utils/cookies';
	import { uppercase_first_character } from '../../utils/str';
	import { encrypt_password } from '../../utils/password';
	import { page } from '$app/stores';
	import { browser } from '$app/env';
	import { me_store } from '../../stores/me';

	let me: any = undefined;

	if (browser) me = cookies.get_me();

	$: its_me = $me_store.id.toString() === $page.params.id;
	let modification_mode = false;

	let loading = false;

	let patch_attemps = 0;

	$: current_username = its_me ? $me_store.username : '?';
	$: current_firstname = its_me ? $me_store.firstname : '?';
	$: current_lastname = its_me ? $me_store.lastname : '?';
	$: current_email = its_me ? $me_store.email : '?';

	let username_blur = false;
	let firstname_blur = false;
	let lastname_blur = false;
	let email_blur = false;
	let current_password_blur = false;
	let new_password_blur = false;
	let confirm_new_password_blur = false;

	let username = '';
	let username_warning = '';
	let firstname = '';
	let firstname_warning = '';
	let lastname = '';
	let lastname_warning = '';
	let email = '';
	let email_warning = '';
	let current_password = '';
	let current_password_warning = '';
	let new_password = '';
	let new_password_warning = '';
	let confirm_new_password = '';
	let confirm_new_password_warning = '';

	let response_warning = '';
	let response_success = '';

	$: warnings =
		username_warning.length ||
		firstname_warning.length ||
		lastname_warning.length ||
		email_warning.length ||
		current_password_warning.length ||
		new_password_warning.length ||
		confirm_new_password_warning.length;

	$: disabled =
		warnings > 0 ||
		((!username.length || !username_blur) &&
			(!firstname.length || !firstname_blur) &&
			(!lastname.length || !lastname_blur) &&
			(!email.length || !email_blur) &&
			(!current_password.length || !current_password_blur) &&
			(!new_password.length || !new_password_blur) &&
			(!confirm_new_password.length || !confirm_new_password_blur));

	let show_password = false;
	$: password_input_type = show_password ? 'text' : 'password';

	let emails_already_in_use: Array<string> = [];

	function handle_pen() {
		modification_mode = !modification_mode;
	}

	function clear_all_inputs() {
		username = '';
		username_warning = '';
		firstname = '';
		firstname_warning = '';
		lastname = '';
		lastname_warning = '';
		email = '';
		email_warning = '';
		current_password = '';
		current_password_warning = '';
		new_password = '';
		new_password_warning = '';
		confirm_new_password = '';
		confirm_new_password_warning = '';
	}

	function handle_update() {
		return new Promise((resolve) => {
			patch_attemps++;

			let inputs_corrupted = false;

			if (check_username()) inputs_corrupted = true;
			if (check_firstname()) inputs_corrupted = true;
			if (check_lastname()) inputs_corrupted = true;
			if (check_email()) inputs_corrupted = true;
			if (check_current_password()) inputs_corrupted = true;
			if (check_new_password()) inputs_corrupted = true;
			if (check_confirm_new_password()) inputs_corrupted = true;

			if (inputs_corrupted) return resolve(false);
			show_password = false;

			var url_search_params = new URLSearchParams();

			setTimeout(async () => {
				url_search_params.append('username', username);
				url_search_params.append('firstname', firstname);
				url_search_params.append('lastname', lastname);
				url_search_params.append('email', email);
				url_search_params.append('current_password', await encrypt_password(current_password));
				url_search_params.append('new_password', await encrypt_password(new_password));

				let url = 'http://localhost:7170/v1/me' + '?' + url_search_params.toString();

				const res = await fetch(url, {
					method: 'PATCH',
					credentials: 'include',
					headers: {
						'Content-type': 'application/json; charset=UTF-8'
					}
				});

				if (res.ok) {
					await res
						.json()
						.then((body) => {
							if (body.hasOwnProperty(cookies.labels.me)) {
								cookies.add_a_cookie(cookies.labels.me, body.me);
								me_store.refresh_from_cookies();
								clear_all_inputs();
								notifies_response_success('profile updated');
								resolve(true);
							} else {
								notifies_response_warning(
									'An error occured on server side with your infos, please try again'
								);
							}
							body.token;
						})
						.catch(() => {
							notifies_response_warning(
								'An error occured in the response from the server side, please try again'
							);
						});
				} else {
					if (res.status == 409) {
						emails_already_in_use.push(email);
						check_email();
					} else if (res.status == 403) {
						current_password_warning = 'invalid password';
					} else notifies_response_warning('An error occured on server side, please try again');
				}
				resolve(false);
			}, 1000);
		});
	}

	function notifies_response_warning(warning: string) {
		response_warning = uppercase_first_character(warning);
	}

	function notifies_response_success(success: string) {
		response_success = uppercase_first_character(success);
		setTimeout(() => {
			response_success = '';
		}, 5000);
	}

	// ----------------------------------------------------------------- sanitizing
	function check_username(): boolean {
		response_warning = '';
		const warning = sanitzer.name(username, true);
		if (!warning.length) username_blur = true;
		if (patch_attemps || username_blur) username_warning = warning;
		return username_warning.length > 0;
	}
	function check_firstname(): boolean {
		response_warning = '';
		const warning = sanitzer.name(firstname, true);
		if (!warning.length) firstname_blur = true;
		if (patch_attemps || firstname_blur) firstname_warning = warning;
		return firstname_warning.length > 0;
	}
	function check_lastname(): boolean {
		response_warning = '';
		const warning = sanitzer.name(lastname, true);
		if (!warning.length) lastname_blur = true;
		if (patch_attemps || lastname_blur) lastname_warning = warning;
		return lastname_warning.length > 0;
	}
	function check_email(): boolean {
		response_warning = '';
		const warning = sanitzer.email(email, true);
		if (!warning.length) email_blur = true;
		if (patch_attemps || email_blur) {
			if (emails_already_in_use.includes(email))
				email_warning = 'Email is already in use, please choose another';
			else email_warning = warning;
		}

		return email_warning.length > 0;
	}
	function check_if_all_password_are_empty(): boolean {
		if (!current_password.length && !new_password && !confirm_new_password) {
			current_password_warning = '';
			new_password_warning = '';
			confirm_new_password_warning = '';
			return true;
		}
		return false;
	}

	function check_current_password(event: any = null): boolean {
		response_warning = '';
		const warning = sanitzer.password(current_password, true);
		if (!warning.length) current_password_blur = true;
		if (event) current_password = event.target.value;
		if (check_if_all_password_are_empty()) return false;

		if (patch_attemps || current_password_blur)
			current_password_warning = sanitzer.password(current_password, false);

		return current_password_warning.length > 0;
	}
	function check_new_password(event: any = null): boolean {
		response_warning = '';
		const warning = sanitzer.password(new_password, true);
		if (!warning.length) new_password_blur = true;
		if (event) new_password = event.target.value;
		if (check_if_all_password_are_empty()) return false;

		if (patch_attemps || new_password_blur)
			new_password_warning = sanitzer.password(new_password, false);

		return new_password_warning.length > 0;
	}
	function check_confirm_new_password(event: any = null): boolean {
		response_warning = '';
		const warning = sanitzer.confirm_password(new_password, confirm_new_password, true);
		if (!warning.length) confirm_new_password_blur = true;
		if (event) confirm_new_password = event.target.value;
		if (check_if_all_password_are_empty()) return false;

		if (patch_attemps || confirm_new_password_blur)
			confirm_new_password_warning = sanitzer.confirm_password(new_password, confirm_new_password);

		return confirm_new_password_warning.length > 0;
	}
</script>

<!-- ========================= HTML -->
<BlackBox>
	<h1 class="mt-2 mb-1 text-2xl text-white">Profile</h1>
	{#if its_me}
		<button class="absolute right-4 top-3" on:click={handle_pen}>
			<img class="invert" src="/pen.png" width="20px" height="20px" alt="pen" />
		</button>
	{/if}
	<form action="" class="pt-1">
		<div>
			<label class="inline-block" for="username">Username</label>
			<p class="value">{current_username}</p>
			{#if modification_mode}
				<input
					type="text"
					placeholder="Username"
					name="username"
					bind:value={username}
					on:input={check_username}
					on:blur={() => {
						username_blur = true;
						check_username();
					}}
					disabled={loading}
				/>
				<Warning content={username_warning} color="red" />
			{/if}
		</div>
		<div>
			<label class="inline-block" for="firstname">Firstname</label>
			<p class="value">{current_firstname}</p>
			{#if modification_mode}
				<input
					type="text"
					placeholder="Firstname"
					name="firstname"
					bind:value={firstname}
					on:input={check_firstname}
					on:blur={() => {
						firstname_blur = true;
						check_firstname();
					}}
					disabled={loading}
				/>
				<Warning content={firstname_warning} color="red" />
			{/if}
		</div>
		<div>
			<label class="inline-block" for="lastname">Lastname</label>
			<p class="value">{current_lastname}</p>
			{#if modification_mode}
				<input
					type="text"
					placeholder="Lastname"
					name="lastname"
					bind:value={lastname}
					on:input={check_lastname}
					on:blur={() => {
						lastname_blur = true;
						check_lastname();
					}}
					disabled={loading}
				/>
				<Warning content={lastname_warning} color="red" />
			{/if}
		</div>
		<div>
			<label class="inline-block" for="email">Email</label>
			<p class="value">{current_email}</p>
			{#if modification_mode && me?.external === 'none'}
				<input
					type="email"
					placeholder="Email"
					name="email"
					bind:value={email}
					on:input={check_email}
					on:blur={() => {
						email_blur = true;
						check_email();
					}}
					disabled={loading}
				/>
				<Warning content={email_warning} color="red" />
			{/if}
		</div>
		{#if modification_mode && me?.external === 'none'}
			<div id="passwords">
				<label class="inline-block" for="password">Password</label>
				<div class="relative">
					<input
						type={password_input_type}
						placeholder="Current password"
						name="current_password"
						value={current_password}
						on:input={check_current_password}
						on:blur={() => {
							current_password_blur = true;
							check_current_password();
						}}
						disabled={loading}
					/>
					<Eye bind:open={show_password} />
				</div>
				<Warning content={current_password_warning} color="red" />
				<div class="relative mt-3">
					<input
						type={password_input_type}
						placeholder="New password"
						name="new_password"
						value={new_password}
						on:input={check_new_password}
						on:blur={() => {
							new_password_blur = true;
							check_new_password();
						}}
						disabled={loading}
					/>
					<Eye bind:open={show_password} />
				</div>
				<Warning content={new_password_warning} color="red" />
				<div class="relative mt-3">
					<input
						type={password_input_type}
						placeholder="New password"
						name="confirm_new_password"
						value={confirm_new_password}
						on:input={check_confirm_new_password}
						on:blur={() => {
							confirm_new_password_blur = true;
							check_confirm_new_password();
						}}
						disabled={loading}
					/>
					<Eye bind:open={show_password} />
				</div>
				<Warning content={confirm_new_password_warning} color="red" />
			</div>
		{/if}
		{#if modification_mode}
			<ConfirmationButton name="update" handler={handle_update} bind:loading bind:disabled />
			<Warning centered content={response_warning} color="red" />
			<Warning centered content={response_success} color="green" />
		{/if}
	</form>
</BlackBox>

<!-- ========================= CSS -->
<style lang="postcss">
	.value {
		@apply inline text-gray-400;
	}
</style>
