# IamFeel PWA Icons

## Icon Requirements

PWA icons are required in multiple sizes for different devices and contexts:

- 72x72 - Minimum required size
- 96x96 - Android notification icon
- 128x128 - Chrome Web Store
- 144x144 - Windows pinned sites
- 152x152 - iPad home screen
- 192x192 - Android home screen (recommended)
- 384x384 - Splash screen
- 512x512 - High-res displays, splash screens

## Current Status

A placeholder SVG icon (`icon.svg`) has been created with the IamFeel branding (boxing fist + "IF" text).

## Generate PNG Icons

### Using ImageMagick (Recommended)

If you have ImageMagick installed:

```bash
cd web/static/icons

# Generate all required sizes
for size in 72 96 128 144 152 192 384 512; do
    convert icon.svg -resize ${size}x${size} icon-${size}x${size}.png
done
```

### Using online tools

1. Upload `icon.svg` to a PWA icon generator:
   - https://www.pwabuilder.com/imageGenerator
   - https://realfavicongenerator.net/

2. Download the generated icons and place them in this directory

### Manual creation

Create PNG files at each required size:
- Use Figma, Sketch, Adobe Illustrator, or similar
- Export as PNG at exact dimensions
- Name as `icon-{size}x{size}.png`

## Icon Design Guidelines

- **Safe zone**: Keep important content within 80% center area
- **Background**: Should work on any color (current: dark #0a0e12)
- **Simple**: Icon should be recognizable at 72x72
- **Maskable**: Works with any shape (circle, square, rounded square)
- **Brand colors**: Red (#ef4444) for accent, dark for background
