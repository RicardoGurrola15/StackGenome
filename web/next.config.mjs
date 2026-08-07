/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  // Disable source maps in production
  productionBrowserSourceMaps: false,
  // Image optimization — only allow our own domain
  images: {
    remotePatterns: [],
    unoptimized: true, // required for static export with next/image
  },
}

export default nextConfig
