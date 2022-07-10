<!-- ========================= HTML -->
<!-- https://loading.io/spinner/progress-pie/-circle-round-pie-progress-percent-ratio-pizza-color-palette -->

<!-- removed svg element attributes: -->
<!-- xmlns="http://www.w3.org/2000/svg" -->
<!-- xmlns:xlink="http://www.w3.org/1999/xlink" -->
<script lang="ts">
	export let inversed = false
	$: color = inversed ? "#ffffff" : "#000000"
</script>

<svg
	class="inline-block"
	width="20px"
	height="20px"
	viewBox="0 0 100 100"
	preserveAspectRatio="xMidYMid"
>
	<path
		fill="none"
		d="M50 30A20 20 0 0 1 50 70A20 20 0 0 1 50 30"
		stroke={color}
		stroke-width="40"
	>
		<animate
			attributeName="stroke-dasharray"
			repeatCount="indefinite"
			dur="1.25"
			values="0 0 0 126;0 0 126 126;0 126 126 126"
			keyTimes="0;0.5;1"
			keySplines="0.5 0 0.5 1;0.5 0 0.5 1"
			calcMode="spline"
		/>
		<animate
			attributeName="stroke"
			values={Array(5).fill(color).join(";")}
			keyTimes="0;0.25;0.5;0.75;1"
			dur="5s"
			calcMode="discrete"
			repeatCount="indefinite"
		/>
	</path>
</svg>
