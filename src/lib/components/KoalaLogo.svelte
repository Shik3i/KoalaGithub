<script lang="ts">
	let { size = 38, class: className = '' } = $props<{ size?: number; class?: string }>();
</script>

<img
	src="/assets/brand/koalagithub-logo-128.webp"
	srcset="
		/assets/brand/koalagithub-logo-48.webp 48w,
		/assets/brand/koalagithub-logo-128.webp 128w,
		/assets/brand/koalagithub-logo-256.webp 256w,
		/assets/brand/koalagithub-logo-512.webp 512w
	"
	sizes={`${size}px`}
	width={size}
	height={size}
	class={className}
	alt=""
	aria-hidden="true"
	decoding="async"
/>
