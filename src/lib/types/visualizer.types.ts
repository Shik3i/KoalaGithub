export type VisualizerCategory =
	| 'stats'
	| 'languages'
	| 'activity'
	| 'streak'
	| 'trophies'
	| 'headers'
	| 'configured'
	| 'badges'
	| 'quotes';

export interface VisualizerTheme {
	key: string;
	label: string;
}

export type PreviewType = 'image' | 'iframe' | 'setup-card' | 'opt-in-counter';

export interface Visualizer {
	id: string;
	name: string;
	description: string;
	category: VisualizerCategory;
	tags: string[];
	previewType: PreviewType;
	imageFormat: 'svg' | 'png' | 'gif' | 'jpg' | 'auto';
	imageUrlTemplate: string;
	markdownTemplate: string;
	websiteUrl: string;
	repositoryUrl: string;
	themes: VisualizerTheme[];
	defaultTheme: string;
	requiresUsername: boolean;
	requiresExternalSetup: boolean;
	setupExplanation?: string;
	actionWorkflowYaml?: string;
	privacyNotice?: string;
	addedAt: string; // YYYY-MM-DD
	enabled: boolean;
	estimatedHeight?: number; // In pixels, to prevent CLS
	voteCount: number;
	userVoted: boolean;
	githubStars?: number;
	lastRefreshedAt?: string;
	nextRefreshAt?: string;
}

export type AppTheme = 'light' | 'dark' | 'system';

export type SortOption = 'most-voted' | 'most-stars' | 'alphabetical' | 'recently-added';
export type SortDirection = 'ascending' | 'descending';
export type WorkflowFilter = 'all' | 'workflow' | 'instant';
