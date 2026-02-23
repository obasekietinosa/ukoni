import sharp from 'sharp';
import pngToIco from 'png-to-ico';
import fs from 'fs/promises';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const PUBLIC_DIR = path.resolve(__dirname, '../public');
const LOGO_PATH = path.join(PUBLIC_DIR, 'logo.svg');

async function main() {
  console.log('Generating favicons...');

  // 1. Generate PNGs
  const sizes = [16, 32, 180, 192, 512];

  for (const size of sizes) {
    const filename = size === 180 ? 'apple-touch-icon.png' :
                     size === 192 ? 'android-chrome-192x192.png' :
                     size === 512 ? 'android-chrome-512x512.png' :
                     `favicon-${size}x${size}.png`;

    const outputPath = path.join(PUBLIC_DIR, filename);

    // Resize with padding to keep aspect ratio and center
    await sharp(LOGO_PATH)
      .resize(size, size, {
        fit: 'contain',
        background: { r: 0, g: 0, b: 0, alpha: 0 }
      })
      .toFile(outputPath);

    console.log(`Generated ${filename}`);
  }

  // 2. Generate favicon.ico from 16x16 and 32x32
  try {
    const buf = await pngToIco([
        path.join(PUBLIC_DIR, 'favicon-16x16.png'),
        path.join(PUBLIC_DIR, 'favicon-32x32.png')
    ]);
    await fs.writeFile(path.join(PUBLIC_DIR, 'favicon.ico'), buf);
    console.log('Generated favicon.ico');
  } catch (err) {
    console.error('Error generating favicon.ico:', err);
    process.exit(1);
  }

  // 3. Generate site.webmanifest
  const manifest = {
    name: "Ukoni",
    short_name: "Ukoni",
    icons: [
      {
        src: "/android-chrome-192x192.png",
        sizes: "192x192",
        type: "image/png"
      },
      {
        src: "/android-chrome-512x512.png",
        sizes: "512x512",
        type: "image/png"
      }
    ],
    theme_color: "#ffffff",
    background_color: "#ffffff",
    display: "standalone"
  };

  await fs.writeFile(path.join(PUBLIC_DIR, 'site.webmanifest'), JSON.stringify(manifest, null, 2));
  console.log('Generated site.webmanifest');
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
