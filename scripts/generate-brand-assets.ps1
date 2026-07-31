param(
	[Parameter(Mandatory = $true)]
	[string] $SourcePath
)

$ErrorActionPreference = 'Stop'

Add-Type -AssemblyName System.Drawing
Add-Type -ReferencedAssemblies 'System.Drawing' -TypeDefinition @'
using System;
using System.Collections.Generic;
using System.Drawing;
using System.Drawing.Imaging;
using System.Runtime.InteropServices;

public static class KoalaImageTools
{
	public static Bitmap RemoveConnectedLightBackground(Bitmap source)
	{
		var result = new Bitmap(source.Width, source.Height, PixelFormat.Format32bppArgb);
		using (var graphics = Graphics.FromImage(result))
		{
			graphics.DrawImageUnscaled(source, 0, 0);
		}

		var rectangle = new Rectangle(0, 0, result.Width, result.Height);
		var bitmapData = result.LockBits(rectangle, ImageLockMode.ReadWrite, PixelFormat.Format32bppArgb);
		try
		{
			var stride = bitmapData.Stride;
			var bytes = new byte[stride * result.Height];
			Marshal.Copy(bitmapData.Scan0, bytes, 0, bytes.Length);

			var width = result.Width;
			var height = result.Height;
			var visited = new bool[width * height];
			var queue = new Queue<int>(width * 4);

			Action<int, int> enqueue = (x, y) =>
			{
				var index = y * width + x;
				if (visited[index] || !IsBackgroundCandidate(bytes, stride, x, y))
				{
					return;
				}
				visited[index] = true;
				queue.Enqueue(index);
			};

			for (var x = 0; x < width; x++)
			{
				enqueue(x, 0);
				enqueue(x, height - 1);
			}
			for (var y = 0; y < height; y++)
			{
				enqueue(0, y);
				enqueue(width - 1, y);
			}

			while (queue.Count > 0)
			{
				var index = queue.Dequeue();
				var x = index % width;
				var y = index / width;
				if (x > 0) enqueue(x - 1, y);
				if (x + 1 < width) enqueue(x + 1, y);
				if (y > 0) enqueue(x, y - 1);
				if (y + 1 < height) enqueue(x, y + 1);
			}

			for (var y = 0; y < height; y++)
			{
				for (var x = 0; x < width; x++)
				{
					if (!visited[y * width + x])
					{
						continue;
					}
					bytes[y * stride + x * 4 + 3] = 0;
				}
			}

			Marshal.Copy(bytes, 0, bitmapData.Scan0, bytes.Length);
		}
		finally
		{
			result.UnlockBits(bitmapData);
		}
		return result;
	}

	private static bool IsBackgroundCandidate(byte[] bytes, int stride, int x, int y)
	{
		var offset = y * stride + x * 4;
		var blue = bytes[offset];
		var green = bytes[offset + 1];
		var red = bytes[offset + 2];
		var minimum = Math.Min(red, Math.Min(green, blue));
		var maximum = Math.Max(red, Math.Max(green, blue));
		return minimum >= 205 && maximum - minimum <= 18;
	}
}
'@

function New-ResizedPng {
	param(
		[System.Drawing.Image] $Source,
		[int] $Width,
		[int] $Height,
		[string] $Destination,
		[double] $Scale = 1.0,
		[System.Drawing.Color] $Background = [System.Drawing.Color]::Transparent
	)

	$bitmap = [System.Drawing.Bitmap]::new(
		$Width,
		$Height,
		[System.Drawing.Imaging.PixelFormat]::Format32bppArgb
	)
	try {
		$graphics = [System.Drawing.Graphics]::FromImage($bitmap)
		try {
			$graphics.Clear($Background)
			$graphics.CompositingMode = [System.Drawing.Drawing2D.CompositingMode]::SourceOver
			$graphics.CompositingQuality = [System.Drawing.Drawing2D.CompositingQuality]::HighQuality
			$graphics.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
			$graphics.PixelOffsetMode = [System.Drawing.Drawing2D.PixelOffsetMode]::HighQuality
			$graphics.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::HighQuality

			$targetWidth = [int][Math]::Round($Width * $Scale)
			$targetHeight = [int][Math]::Round($Height * $Scale)
			$left = [int][Math]::Round(($Width - $targetWidth) / 2)
			$top = [int][Math]::Round(($Height - $targetHeight) / 2)
			$destinationRectangle = [System.Drawing.Rectangle]::new(
				$left,
				$top,
				$targetWidth,
				$targetHeight
			)
			$graphics.DrawImage(
				$Source,
				$destinationRectangle,
				0,
				0,
				$Source.Width,
				$Source.Height,
				[System.Drawing.GraphicsUnit]::Pixel
			)
		}
		finally {
			$graphics.Dispose()
		}
		$bitmap.Save($Destination, [System.Drawing.Imaging.ImageFormat]::Png)
	}
	finally {
		$bitmap.Dispose()
	}
}

function New-PngIco {
	param(
		[string] $PngPath,
		[string] $Destination
	)

	$pngBytes = [System.IO.File]::ReadAllBytes($PngPath)
	$stream = [System.IO.MemoryStream]::new()
	$writer = [System.IO.BinaryWriter]::new($stream)
	try {
		$writer.Write([uint16] 0)
		$writer.Write([uint16] 1)
		$writer.Write([uint16] 1)
		$writer.Write([byte] 32)
		$writer.Write([byte] 32)
		$writer.Write([byte] 0)
		$writer.Write([byte] 0)
		$writer.Write([uint16] 1)
		$writer.Write([uint16] 32)
		$writer.Write([uint32] $pngBytes.Length)
		$writer.Write([uint32] 22)
		$writer.Write($pngBytes)
		[System.IO.File]::WriteAllBytes($Destination, $stream.ToArray())
	}
	finally {
		$writer.Dispose()
		$stream.Dispose()
	}
}

function New-OpenGraphImage {
	param(
		[System.Drawing.Image] $Logo,
		[string] $Destination
	)

	$bitmap = [System.Drawing.Bitmap]::new(
		1200,
		630,
		[System.Drawing.Imaging.PixelFormat]::Format32bppArgb
	)
	try {
		$graphics = [System.Drawing.Graphics]::FromImage($bitmap)
		try {
			$graphics.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::HighQuality
			$graphics.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
			$graphics.TextRenderingHint = [System.Drawing.Text.TextRenderingHint]::AntiAliasGridFit
			$graphics.Clear([System.Drawing.ColorTranslator]::FromHtml('#071525'))
			$rectangle = [System.Drawing.Rectangle]::new(0, 0, 1200, 630)
			$gradient = [System.Drawing.Drawing2D.LinearGradientBrush]::new(
				$rectangle,
				[System.Drawing.ColorTranslator]::FromHtml('#071525'),
				[System.Drawing.ColorTranslator]::FromHtml('#123D49'),
				25
			)
			try {
				$graphics.FillRectangle($gradient, $rectangle)
			}
			finally {
				$gradient.Dispose()
			}

			$graphics.DrawImage($Logo, [System.Drawing.Rectangle]::new(55, 55, 520, 520))

			$titleFont = [System.Drawing.Font]::new('Segoe UI', 56, [System.Drawing.FontStyle]::Bold)
			$subtitleFont = [System.Drawing.Font]::new('Segoe UI', 23, [System.Drawing.FontStyle]::Regular)
			$titleBrush = [System.Drawing.SolidBrush]::new([System.Drawing.Color]::White)
			$accentBrush = [System.Drawing.SolidBrush]::new(
				[System.Drawing.ColorTranslator]::FromHtml('#49D6C2')
			)
			$mutedBrush = [System.Drawing.SolidBrush]::new(
				[System.Drawing.ColorTranslator]::FromHtml('#CBD5E1')
			)
			try {
				$graphics.DrawString('Koala', $titleFont, $titleBrush, 585, 195)
				$koalaWidth = $graphics.MeasureString('Koala', $titleFont).Width
				$graphics.DrawString('GitHub', $titleFont, $accentBrush, 585 + $koalaWidth - 8, 195)
				$graphics.DrawString(
					'Discover and compare profile visualizers',
					$subtitleFont,
					$mutedBrush,
					590,
					290
				)
				$graphics.DrawString(
					'Activity graphs · Stats · README art',
					$subtitleFont,
					$mutedBrush,
					590,
					335
				)
			}
			finally {
				$titleFont.Dispose()
				$subtitleFont.Dispose()
				$titleBrush.Dispose()
				$accentBrush.Dispose()
				$mutedBrush.Dispose()
			}
		}
		finally {
			$graphics.Dispose()
		}
		$bitmap.Save($Destination, [System.Drawing.Imaging.ImageFormat]::Png)
	}
	finally {
		$bitmap.Dispose()
	}
}

$resolvedSource = (Resolve-Path -LiteralPath $SourcePath).Path
$repoRoot = (Resolve-Path -LiteralPath (Join-Path $PSScriptRoot '..')).Path
$outputDirectory = Join-Path $repoRoot 'static\assets\brand'
[void] (New-Item -ItemType Directory -Path $outputDirectory -Force)

$source = [System.Drawing.Bitmap]::FromFile($resolvedSource)
try {
	$master = [KoalaImageTools]::RemoveConnectedLightBackground($source)
	try {
		$masterPath = Join-Path $outputDirectory 'koalagithub-logo-1024.png'
		New-ResizedPng -Source $master -Width 1024 -Height 1024 -Destination $masterPath

		foreach ($size in @(512, 256, 128)) {
			New-ResizedPng `
				-Source $master `
				-Width $size `
				-Height $size `
				-Destination (Join-Path $outputDirectory "koalagithub-logo-$size.png")
		}

		foreach ($size in @(16, 32, 48)) {
			New-ResizedPng `
				-Source $master `
				-Width $size `
				-Height $size `
				-Destination (Join-Path $outputDirectory "favicon-$size.png")
		}

		New-ResizedPng `
			-Source $master `
			-Width 180 `
			-Height 180 `
			-Destination (Join-Path $repoRoot 'static\apple-touch-icon.png')
		New-ResizedPng `
			-Source $master `
			-Width 192 `
			-Height 192 `
			-Destination (Join-Path $outputDirectory 'icon-192.png')
		New-ResizedPng `
			-Source $master `
			-Width 512 `
			-Height 512 `
			-Destination (Join-Path $outputDirectory 'icon-512.png')

		$maskBackground = [System.Drawing.ColorTranslator]::FromHtml('#0F172A')
		New-ResizedPng `
			-Source $master `
			-Width 192 `
			-Height 192 `
			-Scale 0.8 `
			-Background $maskBackground `
			-Destination (Join-Path $outputDirectory 'icon-maskable-192.png')
		New-ResizedPng `
			-Source $master `
			-Width 512 `
			-Height 512 `
			-Scale 0.8 `
			-Background $maskBackground `
			-Destination (Join-Path $outputDirectory 'icon-maskable-512.png')

		New-PngIco `
			-PngPath (Join-Path $outputDirectory 'favicon-32.png') `
			-Destination (Join-Path $repoRoot 'static\favicon.ico')
		New-OpenGraphImage `
			-Logo $master `
			-Destination (Join-Path $outputDirectory 'og-image.png')
	}
	finally {
		$master.Dispose()
	}
}
finally {
	$source.Dispose()
}

Write-Output "Brand assets generated in $outputDirectory"
