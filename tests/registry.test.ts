import { describe, it, expect } from 'vitest';
import { VISUALIZERS } from '../src/lib/data/visualizers';
import type { VisualizerCategory } from '../src/lib/types/visualizer.types';

const VALID_CATEGORIES: VisualizerCategory[] = [
	'stats',
	'languages',
	'activity',
	'streak',
	'trophies',
	'headers',
	'configured',
	'badges',
	'quotes'
];

describe('Visualizer Registry Data Integrity', () => {
	it('should contain entries and have unique IDs', () => {
		expect(VISUALIZERS.length).toBeGreaterThan(0);
		const ids = VISUALIZERS.map((v) => v.id);
		const uniqueIds = new Set(ids);
		expect(uniqueIds.size).toBe(ids.length);
	});

	it('should have valid HTTPS website and GitHub repository URLs', () => {
		for (const item of VISUALIZERS) {
			expect(item.websiteUrl).toMatch(/^https:\/\//);
			expect(item.repositoryUrl).toMatch(/^https:\/\/(www\.)?github\.com\//);
		}
	});

	it('should have valid category and non-empty tags', () => {
		for (const item of VISUALIZERS) {
			expect(VALID_CATEGORIES).toContain(item.category);
			expect(Array.isArray(item.tags)).toBe(true);
			expect(item.tags.length).toBeGreaterThan(0);
		}
	});

	it('should have unique theme keys and valid default theme when themes exist', () => {
		for (const item of VISUALIZERS) {
			if (item.themes.length > 0) {
				const themeKeys = item.themes.map((t) => t.key);
				const uniqueKeys = new Set(themeKeys);
				expect(uniqueKeys.size).toBe(themeKeys.length);
				expect(themeKeys).toContain(item.defaultTheme);
			}
		}
	});

	it('should use supported placeholders in templates', () => {
		for (const item of VISUALIZERS) {
			if (item.requiresUsername) {
				expect(item.imageUrlTemplate + item.markdownTemplate).toContain('{username}');
			}
		}
	});

	it('should provide setup explanations for configured entries', () => {
		for (const item of VISUALIZERS) {
			if (item.requiresExternalSetup) {
				expect(item.setupExplanation).toBeDefined();
				expect(item.setupExplanation!.length).toBeGreaterThan(10);
				expect(item.actionWorkflowYaml).toBeDefined();
				for (const line of item.actionWorkflowYaml!.split('\n')) {
					if (line.includes('uses:')) {
						expect(line, `Action must be pinned to a commit SHA: ${line}`).toMatch(
							/^\s*-\s+uses:\s+[^@]+@[0-9a-f]{40}\s*$/
						);
					}
				}
			}
		}
	});

	it('keeps the two snake generators separate', () => {
		const platane = VISUALIZERS.find((item) => item.id === 'platane-contribution-snake');
		const snakeAndCommits = VISUALIZERS.find((item) => item.id === 'dahan8473-snake-and-commits');
		expect(platane?.repositoryUrl).toBe('https://github.com/Platane/snk');
		expect(platane?.imageUrlTemplate).toContain('/output/github-snake.svg');
		expect(snakeAndCommits?.repositoryUrl).toBe('https://github.com/dahan8473/snake-and-commits');
		expect(snakeAndCommits?.imageUrlTemplate).toContain('/output/snake.svg');
	});
});
