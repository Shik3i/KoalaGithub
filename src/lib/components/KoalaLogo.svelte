<script lang="ts">
	let { size = 38, class: className = '' } = $props<{ size?: number; class?: string }>();
</script>

<svg
	width={size}
	height={size}
	viewBox="0 0 100 100"
	fill="none"
	xmlns="http://www.w3.org/2000/svg"
	class={className}
	aria-label="KoalaGitHub Mascot"
>
	<!-- Outer Ears -->
	<circle cx="22" cy="28" r="18" fill="#94A3B8" />
	<circle cx="78" cy="28" r="18" fill="#94A3B8" />
	<!-- Inner Ears -->
	<circle cx="22" cy="28" r="11" fill="#E2E8F0" />
	<circle cx="78" cy="28" r="11" fill="#E2E8F0" />
	<!-- Head -->
	<circle cx="50" cy="52" r="34" fill="#64748B" />
	<!-- Face patch -->
	<ellipse cx="50" cy="56" rx="26" ry="22" fill="#CBD5E1" />
	<!-- Eyes -->
	<ellipse cx="38" cy="48" rx="4" ry="5" fill="#0F172A" />
	<ellipse cx="62" cy="48" rx="4" ry="5" fill="#0F172A" />
	<circle cx="39" cy="46" r="1.5" fill="#FFFFFF" />
	<circle cx="63" cy="46" r="1.5" fill="#FFFFFF" />
	<!-- Nose -->
	<ellipse cx="50" cy="58" rx="7" ry="10" fill="#1E293B" />
	<!-- Cheeks -->
	<circle cx="30" cy="58" r="4.5" fill="#F472B6" opacity="0.6" />
	<circle cx="70" cy="58" r="4.5" fill="#F472B6" opacity="0.6" />
	<!-- Tech Accent: GitHub Octocat ear mark / glasses frame -->
	<path d="M 40 76 Q 50 82 60 76" stroke="#4F46E5" stroke-width="3.5" stroke-linecap="round" fill="none" />
</svg>
