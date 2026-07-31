import type { Visualizer } from "$lib/types/visualizer.types";
import rawData from "./visualizers.json";

export const VISUALIZERS: Visualizer[] = rawData.map((item) => ({
  ...item,
  imageFormat: (item.imageFormat || "svg") as Visualizer["imageFormat"],
  category: item.category as Visualizer["category"],
  previewType: item.previewType as Visualizer["previewType"],
  voteCount: 0,
  userVoted: false,
}));
