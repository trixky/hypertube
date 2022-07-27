<!-- ========================= SCRIPT -->
<script lang="ts" context="module">
	import type { Load } from '@sveltejs/kit';

	async function getUserMovies(
		fetch: (info: RequestInfo, init?: RequestInit | undefined) => Promise<Response>,
		userId: number,
		page = 1
	) {
		// Get user movies
		const response = await fetch(apiUser(`/v1/user/${userId}/movies?page=${page}`), {
			method: 'GET',
			credentials: 'include'
		});
		if (!response.ok || response.status >= 400) {
			return {
				status: response.ok ? response.status : 500
			};
		}
		const results = (await response.json()) as {
			page: number;
			results: number;
			totalResults: number;
			medias: Result[];
		};

		return { medias: results.medias.map(addUserTitle), totalResults: results.totalResults };
	}

	// Preload search results and genres
	// -- and insert them on startup in the client
	export const load: Load = async ({ fetch, session, params }) => {
		const userId = parseInt(params.id);
		if (isNaN(userId) || userId < 1) {
			return {
				status: 400
			};
		}

		if (userId == session.user?.id) {
			return {
				status: 200,
				props: {
					user: session.user,
					...(await getUserMovies(fetch, session.user.id))
				}
			};
		}

		// Get user informations
		const url_search_params = new URLSearchParams();
		url_search_params.append('id', params.id);
		const url = apiUser('/v1/user?' + url_search_params.toString());
		const response = await fetch(url, {
			method: 'GET',
			credentials: 'include',
			headers: {
				'Content-type': 'application/json; charset=UTF-8'
			}
		});
		if (!response.ok || response.status >= 400) {
			return {
				status: response.ok ? response.status : 500
			};
		}
		const body = await response.json();
		if (!(cookies.labels.user_info in body)) {
			return {
				status: 500
			};
		}
		const user_64 = body[cookies.labels.user_info];
		const user_json = atob(user_64);
		const user = JSON.parse(user_json);

		return {
			status: 200,
			props: {
				user,
				...(await getUserMovies(fetch, user.id))
			}
		};
	};
</script>

<script lang="ts">
	import { _ } from 'svelte-i18n';
	import BlackBox from '$components/containers/black-box.svelte';
	import Warning from '$components/inputs/warning.svelte';
	import Eye from '$components/inputs/eye.svelte';
	import * as sanitizer from '$utils/sanitizer';
	import ConfirmationButton from '$components/buttons/confirmation-button.svelte';
	import * as cookies from '$utils/cookies';
	import { uppercase_first_character } from '$utils/str';
	import { encrypt_password } from '$utils/password';
	import { apiUser } from '$utils/api';
	import { session } from '$app/stores';
	import InfoLine from './info-line.svelte';
	import MediaList from '$components/generics/MediaList.svelte';
	import type { Result } from '$types/Media';
	import { addUserTitle } from '$utils/media';

	export let user: NonNullable<App.Session['user']>;

	$: {
		user;
		modification_mode = false;
	}

	$: its_me = $session.user?.id == user.id;
	let modification_mode = false;
	$: img_src = modification_mode ? '/return.png' : '/pen.png';
	$: img_alt = modification_mode ? $_('auth.cancel') : $_('auth.modify');

	let loading = false;

	let patch_attemps = 0;

	let user_does_not_exist = false;

	$: current_username = user.username;
	$: current_firstname = user.firstname;
	$: current_lastname = user.lastname;
	$: current_email = user.email;

	$: if (its_me) user_does_not_exist = false;

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

	let response_update_warning = '';
	let response_update_success = '';

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

	$: can_be_empty = Boolean(
		current_username.length ||
			current_firstname.length ||
			current_lastname.length ||
			current_email.length
	);

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

			const url_search_params = new URLSearchParams();

			setTimeout(async () => {
				url_search_params.append('username', username);
				url_search_params.append('firstname', firstname);
				url_search_params.append('lastname', lastname);
				url_search_params.append('email', email);
				url_search_params.append('current_password', await encrypt_password(current_password));
				url_search_params.append('new_password', await encrypt_password(new_password));

				const url = apiUser('/v1/me?' + url_search_params.toString());
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
							if (cookies.labels.user_info in body) {
								cookies.add_a_cookie(cookies.labels.user_info, body[cookies.labels.user_info]);
								$session.token = body[cookies.labels.user_info];
								if ($session.user) {
									if (username) {
										$session.user.username = username;
									}
									if (firstname) {
										$session.user.firstname = firstname;
									}
									if (lastname) {
										$session.user.lastname = lastname;
									}
									if (email) {
										$session.user.email = email;
									}
								}
								clear_all_inputs();
								notifies_response_success($_('auth.profile_updated'));
								resolve(true);
							} else {
								notifies_update_response_warning($_('auth.server_error'));
							}
						})
						.catch(() => {
							notifies_update_response_warning($_('auth.server_error'));
						});
				} else {
					if (res.status == 409) {
						emails_already_in_use.push(email);
						check_email();
					} else if (res.status == 403) {
						current_password_warning = $_('auth.invalid_password');
					} else {
						notifies_update_response_warning($_('auth.server_error'));
					}
				}
				resolve(false);
			}, 1000);
		});
	}

	function notifies_update_response_warning(warning: string) {
		response_update_warning = uppercase_first_character(warning);
	}

	function notifies_response_success(success: string) {
		response_update_success = uppercase_first_character(success);
		setTimeout(() => {
			response_update_success = '';
		}, 5000);
	}

	// ----------------------------------------------------------------- sanitizing
	function check_username(): boolean {
		response_update_warning = '';
		const warning = sanitizer.name(username, true);
		if (!warning.length) username_blur = true;
		if (patch_attemps || username_blur) username_warning = warning;
		return username_warning.length > 0;
	}
	function check_firstname(): boolean {
		response_update_warning = '';
		const warning = sanitizer.name(firstname, true);
		if (!warning.length) firstname_blur = true;
		if (patch_attemps || firstname_blur) firstname_warning = warning;
		return firstname_warning.length > 0;
	}
	function check_lastname(): boolean {
		response_update_warning = '';
		const warning = sanitizer.name(lastname, true);
		if (!warning.length) lastname_blur = true;
		if (patch_attemps || lastname_blur) lastname_warning = warning;
		return lastname_warning.length > 0;
	}
	function check_email(): boolean {
		response_update_warning = '';
		const warning = sanitizer.email(email, true);
		if (!warning.length) email_blur = true;
		if (patch_attemps || email_blur) {
			if (emails_already_in_use.includes(email)) email_warning = $_('auth.email_already_in_use');
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

	function check_current_password(
		event: (Event & { currentTarget: EventTarget & HTMLInputElement }) | null = null
	): boolean {
		response_update_warning = '';
		const warning = sanitizer.password(current_password, true);
		if (!warning.length) current_password_blur = true;
		if (event) current_password = event.currentTarget.value;
		if (check_if_all_password_are_empty()) return false;

		if (patch_attemps || current_password_blur)
			current_password_warning = sanitizer.password(current_password, false);

		return current_password_warning.length > 0;
	}
	function check_new_password(
		event: (Event & { currentTarget: EventTarget & HTMLInputElement }) | null = null
	): boolean {
		response_update_warning = '';
		const warning = sanitizer.password(new_password, true);
		if (!warning.length) new_password_blur = true;
		if (event) new_password = event.currentTarget.value;
		if (check_if_all_password_are_empty()) return false;

		if (patch_attemps || new_password_blur)
			new_password_warning = sanitizer.password(new_password, false);

		return new_password_warning.length > 0;
	}
	function check_confirm_new_password(
		event: (Event & { currentTarget: EventTarget & HTMLInputElement }) | null = null
	): boolean {
		response_update_warning = '';
		const warning = sanitizer.confirm_password(new_password, confirm_new_password, true);
		if (!warning.length) confirm_new_password_blur = true;
		if (event) confirm_new_password = event.currentTarget.value;
		if (check_if_all_password_are_empty()) return false;

		if (patch_attemps || confirm_new_password_blur)
			confirm_new_password_warning = sanitizer.confirm_password(new_password, confirm_new_password);

		return confirm_new_password_warning.length > 0;
	}

	function handle_keydown(event: KeyboardEvent) {
		if (event.keyCode == 27 || event.key == 'Escape') {
			modification_mode = false;
		}
	}

	// * User movies

	export let medias: Result[] = [];
	export let totalResults: number = 0;
	let loadingMovies = false;

	let page = 2;
	async function loadMore() {
		loadingMovies = true;
		const response = await fetch(apiUser(`/v1/user/${user.id}/movies?page=${page}`), {
			method: 'GET',
			credentials: 'include'
		});
		if (!response.ok || response.status >= 400) {
			return {
				status: response.ok ? response.status : 500
			};
		}
		const results = (await response.json()) as {
			page: number;
			results: number;
			totalResults: number;
			medias: Result[];
		};
		medias.push(...results.medias.map(addUserTitle));
		medias = medias;
		page += 1;
		loadingMovies = false;
	}
</script>

<svelte:window on:keydown={handle_keydown} />

<!-- ========================= HTML -->
<svelte:head>
	<title>hypertube :: User Profile</title>
</svelte:head>
<div class="flex flex-col flex-grow w-full h-full text-white">
	<div class="flex justify-center items-center ">
		<BlackBox title={its_me ? $_('auth.my_profile') : $_('auth.profile')}>
			{#if !user_does_not_exist}
				{#if its_me}
					<button class="absolute right-5 top-4" on:click={handle_pen}>
						<img class="invert" src={img_src} width="18px" height="18px" alt={img_alt} />
					</button>
				{/if}
				<form class="pt-1 w-full">
					<div>
						<InfoLine
							centered={!modification_mode}
							label={$_('auth.username')}
							bind:value={current_username}
							{can_be_empty}
						/>
						{#if modification_mode}
							<input
								type="text"
								placeholder={$_('auth.username')}
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
						<InfoLine
							centered={!modification_mode}
							label={$_('auth.first_name')}
							bind:value={current_firstname}
							{can_be_empty}
						/>
						{#if modification_mode}
							<input
								type="text"
								placeholder={$_('auth.first_name')}
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
						<InfoLine
							centered={!modification_mode}
							label={$_('auth.last_name')}
							bind:value={current_lastname}
							{can_be_empty}
						/>
						{#if modification_mode}
							<input
								type="text"
								placeholder={$_('auth.last_name')}
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
					{#if its_me}
						<div>
							<InfoLine
								centered={!modification_mode}
								label={$_('auth.email')}
								bind:value={current_email}
								{can_be_empty}
							/>
							{#if modification_mode && $session.user?.external === 'none'}
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
							{/if}
						</div>
					{/if}
					{#if modification_mode && $session.user?.external === 'none'}
						<div id="passwords">
							<InfoLine label={$_('auth.password')} no_value />
							<div class="relative">
								<input
									type={password_input_type}
									placeholder={$_('auth.current_password')}
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
									placeholder={$_('auth.new_password')}
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
									placeholder={$_('auth.new_password')}
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
						<ConfirmationButton
							name={$_('auth.update')}
							handler={handle_update}
							bind:loading
							bind:disabled
						/>
						<Warning centered content={response_update_warning} color="red" />
						<Warning centered content={response_update_success} color="green" />
					{/if}
				</form>
			{:else}
				<p class="text-white">{$_('auth.user_does_not_exists')}</p>
			{/if}
		</BlackBox>
	</div>
	<div class="flex-grow bg-black">
		{#if totalResults == 0}
			<p>{$_('media.user_no_results')}</p>
		{:else}
			<MediaList list={medias} {totalResults} {loadMore} loading={loadingMovies} />
		{/if}
	</div>
</div>
