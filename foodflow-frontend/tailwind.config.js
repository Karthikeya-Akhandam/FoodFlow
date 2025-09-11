/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
      "./pages/**/*.{js,ts,jsx,tsx,mdx}",
      "./components/**/*.{js,ts,jsx,tsx,mdx}",
      "./app/**/*.{js,ts,jsx,tsx,mdx}",
    ],
    theme: {
      extend: {
        colors: {
          primary: {
            DEFAULT: "var(--primary)",
            dark: "var(--primary-dark)",
            light: "var(--primary-light)",
          },
          accent: {
            DEFAULT: "var(--accent)",
            hover: "var(--accent-hover)",
          },
          secondary: "var(--secondary)",
          background: "var(--background)",
          foreground: "var(--foreground)",
          muted: "var(--muted)",
          border: "var(--border)",
          success: "var(--success)",
          warning: "var(--warning)",
          error: "var(--error)",
        },
        backgroundImage: {
          'gradient-primary': 'linear-gradient(135deg, var(--gradient-start), var(--gradient-end))',
          'hero-pattern': 'radial-gradient(circle at 25% 25%, rgba(5, 150, 105, 0.1) 0%, transparent 50%), radial-gradient(circle at 75% 75%, rgba(59, 130, 246, 0.1) 0%, transparent 50%)',
        },
        animation: {
          'float': 'float 6s ease-in-out infinite',
          'fade-in': 'fadeIn 0.8s ease-out forwards',
          'slide-in': 'slideIn 0.6s ease-out forwards',
          'slide-in-right': 'slideInRight 0.6s ease-out forwards',
          'pulse-glow': 'pulse-glow 2s ease-in-out infinite',
          'bounce': 'bounce 2s infinite',
          'shimmer': 'shimmer 2s infinite',
        },
        fontFamily: {
          sans: ["var(--font-geist-sans)", "Inter", "sans-serif"],
          mono: ["var(--font-geist-mono)", "Geist Mono", "monospace"],
          serif: ["Playfair Display", "serif"],
        },
      },
    },
    plugins: [],
  };
  