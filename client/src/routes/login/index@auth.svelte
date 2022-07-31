<!-- ========================= SCRIPT -->
<script lang="ts">
	import BlackBox from '$components/containers/black-box.svelte';
	import Logo from '$components/layouts/logo.svelte';
	import External from './external.svelte';
	import ConfirmationButton from '$components/buttons/confirmation-button.svelte';
	import Separator from '$components/generics/separator.svelte';
	import Eye from '$components/inputs/eye.svelte';
	import Warning from '$components/inputs/warning.svelte';
	import * as cookies from '$utils/cookies';
	import * as sanitizer from '$utils/sanitizer'; 
	import { encrypt_password } from '$utils/password';
	import { page, session } from '$app/stores';
	import { browser } from '$app/env';
	import { goto } from '$app/navigation';
	import { _ } from 'svelte-i18n';
	import { apiAuth } from '$utils/api';

	let from_url_parameter: string | null;

	let loading = false;

	let login_attempts = 0;

	let email_blur = false;
	let password_blur = false;

	let email = '';
	let email_warning = '';
	let password = '';
	let password_warning = '';

	let response_warning = '';
	let response_success = '';

	if (browser) {
		from_url_parameter = $page.url.searchParams.get('from');
		if (from_url_parameter == 'recover/apply') {
			response_success = $_('auth.password_updated');
			setTimeout(() => {
				response_success = '';
			}, 5000);
		}
	}

	$: warnings = email_warning.length || password_warning.length;

	$: disabled = warnings > 0 || !email.length || !password.length || !email_blur || !password_blur;

	let show_password = false;
	$: password_input_type = show_password ? 'text' : 'password';

	async function handle_login() {
		return new Promise((resolve) => {
			login_attempts++;
			let inputs_corrupted = false;

			if (check_email()) inputs_corrupted = true;
			if (check_password()) inputs_corrupted = true;

			if (inputs_corrupted) return resolve(false);
			show_password = false;

			setTimeout(async () => {
				const res = await fetch(apiAuth('/v1/internal/login'), {
					method: 'POST',
					headers: {
						'content-type': 'application/json',
						accept: 'application/json'
					},
					body: JSON.stringify({
						email,
						password: await encrypt_password(password)
					})
				});

				if (res.ok) {
					await res
						.json()
						.then((body) => {
							if (cookies.labels.token in body) {
								const user = atob(body[cookies.labels.user_info]);
								const me = JSON.parse(user);
								if (me) {
									cookies.add_a_cookie(cookies.labels.token, body.token);
									cookies.add_a_cookie(cookies.labels.user_info, body[cookies.labels.user_info]);
									session.set({
										token: body.token,
										user: {
											id: me.id,
											username: me.username,
											firstname: me.firstname,
											lastname: me.lastname,
											email: me.email,
											external: me.external
										}
									});
									loading = false
									resolve(true);
									goto('/');
								} else {
									notifies_response_warning($_('auth.server_error'));
								}
							} else {
								notifies_response_warning($_('auth.server_error'));
							}
						})
						.catch(() => {
							notifies_response_warning($_('auth.server_error'));
						});
				} else {
					if (res.status == 403) notifies_response_warning($_('auth.login_failed'));
					else notifies_response_warning($_('auth.server_error'));
				}
				loading = false
				resolve(false);
			}, 500);
		});
	}

	function notifies_response_warning(warning: string) {
		response_warning = warning;
	}

	// ----------------------------------------------------------------- sanitizing
	function check_email(): boolean {
		response_warning = '';

		let warning = sanitizer.email(email);

		if (warning.length == 0 && email.length > 0) email_blur = true;

		if (login_attempts || email_blur) email_warning = warning;

		return email_warning.length > 0;
	}

	function check_password(
		event: (Event & { currentTarget: EventTarget & HTMLInputElement }) | null = null
	): boolean {
		response_warning = '';
		if (event) password = event.currentTarget.value;

		let warning = sanitizer.password(password);

		if (warning.length == 0 && password.length > 0) password_blur = true;

		if (login_attempts || password_blur) password_warning = warning;

		return password_warning.length > 0;
	}

	if (browser) {
		document.onkeypress = function (event) {
			if (event.keyCode == 13) {
				event.preventDefault()
				
				email_blur = true; // email
				check_email()
				password_blur = true; // password
				check_password()

				if (!disabled) {
					loading = true
					handle_login()
				}
			}
		};
	}
</script>

<!-- ========================= HTML -->
<svelte:head>
	<title>{$_('title.login')}</title>
</svelte:head>
<BlackBox title={$_('auth.login_header')}>
	<Logo alone />
	<External disabled={loading} />
	<Separator content={$_('auth.omniauth_separator')} />
	<form action="" class="pt-1 w-full">
		<label for="email" class="required truncate">{$_('auth.email')}</label>
		<input
			type="email"
			placeholder={$_('auth.email')}
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
		<label for="password" class="required truncate">{$_('auth.password')}</label>
		<div class="relative">
			<input
				type={password_input_type}
				placeholder={$_('auth.password')}
				name="password"
				value={password}
				on:input={check_password}
				on:blur={() => {
					password_blur = true;
					check_password();
				}}
				disabled={loading}
			/>
			<Eye bind:open={show_password} />
		</div>
		<Warning content={password_warning} color="red" />
		<p class="extra-link absolute right-10 pl-6">
			<a href="/recover/ask">{$_('auth.forgot_password')}</a>
		</p>
		<div class="mt-10">
			<ConfirmationButton
				name={$_('auth.login_action')}
				handler={handle_login}
				bind:loading
				bind:disabled
			/>
		</div>
		<Warning centered content={response_warning} color="red" />
		<Warning centered content={response_success} color="green" />
	</form>
	<p class="extra-link mt-2">
		<a href="/register">
			{$_('auth.join_question')} <span class="underline">{$_('auth.signup')}</span>
		</a>
	</p>
</BlackBox>

<!-- ========================= CSS -->
<style lang="postcss">
	label {
		@apply ml-2;
	}
</style>
