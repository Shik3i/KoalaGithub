<script lang="ts">
	let { size = 38, class: className = '' } = $props<{ size?: number; class?: string }>();
</script>

<img
	src="/assets/brand/koalagithub-logo-128.png"
	srcset="
		/assets/brand/koalagithub-logo-128.png 128w,
		/assets/brand/koalagithub-logo-256.png 256w,
		/assets/brand/koalagithub-logo-512.png 512w
	"
	sizes={`${size}px`}
	width={size}
	height={size}
	class={className}
	alt=""
	aria-hidden="true"
	decoding="async"
/>
