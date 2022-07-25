<!-- ========================= SCRIPT -->
<script lang="ts">
	import BlackBox from '$components/containers/black-box.svelte';
	import Logo from '$components/layouts/logo.svelte';
	import ConfirmationButton from '$components/buttons/confirmation-button.svelte';
	import Warning from '$components/inputs/warning.svelte';
	import Eye from '$components/inputs/eye.svelte';
	import * as sanitizer from '$utils/sanitizer';
	import { uppercase_first_character } from '$utils/str';
	import { encrypt_password } from '$utils/password';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { browser } from '$app/env';
	import { _ } from 'svelte-i18n';
	import { apiAuth } from '$utils/api';

	let loading = false;

	let registration_attempts = 0;

	let new_password_blur = false;
	let confirm_new_password_blur = false;

	let new_password = '';
	let new_password_warning = '';
	let confirm_new_password = '';
	let confirm_new_password_warning = '';

	let response_warning = '';

	let password_token_url_parameter: string | null;

	if (browser) {
		password_token_url_parameter = $page.url.searchParams.get('password_token');
	}

	$: warnings = new_password_warning.length || confirm_new_password_warning.length;

	$: disabled =
		warnings > 0 ||
		!new_password.length ||
		!confirm_new_password.length ||
		!new_password_blur ||
		!confirm_new_password_blur;

	let show_password = false;
	$: new_password_input_type = show_password ? 'text' : 'password';

	function handle_register() {
		return new Promise((resolve) => {
			registration_attempts++;
			let inputs_corrupted = false;

			if (check_new_password()) inputs_corrupted = true;
			if (check_confirm_new_password()) inputs_corrupted = true;

			if (inputs_corrupted) return resolve(false);
			show_password = false;

			setTimeout(async () => {
				const res = await fetch(apiAuth('/v1/internal/apply-token-password'), {
					method: 'POST',
					headers: {
						'content-type': 'application/json',
						accept: 'application/json'
					},
					body: JSON.stringify({
						password_token: password_token_url_parameter,
						new_password: await encrypt_password(new_password)
					})
				});

				if (res.ok) {
					await res
						.json()
						.then(() => {
							resolve(true);
							goto('/login' + '?from=' + 'recover/apply');
						})
						.catch(() => {
							notifies_response_warning($_('auth.server_error'));
						});
				} else {
					notifies_response_warning($_('auth.server_error'));
				}
				resolve(false);
			}, 1000);
		});
	}

	function notifies_response_warning(warning: string) {
		response_warning = uppercase_first_character(warning);
	}

	// ----------------------------------------------------------------- sanitizing
	function check_new_password(
		event: (Event & { currentTarget: EventTarget & HTMLInputElement }) | null = null
	): boolean {
		response_warning = '';
		if (event) new_password = event.currentTarget.value;

		if (registration_attempts || new_password_blur)
			new_password_warning = sanitizer.password(new_password);

		return new_password_warning.length > 0;
	}
	function check_confirm_new_password(
		event: (Event & { currentTarget: EventTarget & HTMLInputElement }) | null = null
	): boolean {
		response_warning = '';
		if (event) confirm_new_password = event.currentTarget.value;

		let warning = sanitizer.confirm_password(new_password, confirm_new_password);
		if (confirm_new_password.length > 0 && !warning.length) confirm_new_password_blur = true;

		if (registration_attempts || confirm_new_password_blur) confirm_new_password_warning = warning;

		return confirm_new_password_warning.length > 0;
	}
</script>

<!-- ========================= HTML -->
<svelte:head>
	<title>{$_('title.recover_password')}</title>
</svelte:head>
<BlackBox title={$_('auth.forgot_password_header')}>
	<Logo alone />
	<form action="" class="pt-1 w-full">
		<label for="password" class="required">{$_('auth.password')}</label>
		<div class="relative">
			<input
				type={new_password_input_type}
				placeholder={$_('auth.password')}
				name="password"
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
		<label for="confirm password" class="required">{$_('auth.confirm_password')}</label>
		<div class="relative">
			<input
				type={new_password_input_type}
				placeholder={$_('auth.confirm_password')}
				name="confirm password"
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
		<ConfirmationButton
			name={$_('auth.change_password')}
			handler={handle_register}
			bind:loading
			bind:disabled
		/>
		<Warning centered content={response_warning} color="red" />
	</form>
	<p class="extra-link mt-4">
		<a href="/login">{$_('auth.back_to')} <span class="underline">{$_('auth.to_login')}</span></a>
	</p>
</BlackBox>

<!-- ========================= CSS -->
<style lang="postcss">
	label {
		@apply ml-2;
	}
</style>
