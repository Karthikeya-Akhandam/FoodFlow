import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Providers } from "../lib/providers";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata = {
  title: "FoodFlow - Connecting Food Donors with Organizations",
  description: "A platform that connects organizations offering surplus food with organizations that need it. Reduce food waste and help communities through our credit-based matching system.",
  keywords: "food donation, food waste, community service, NGO, food distribution",
  authors: [{ name: "Karthikeya Akhandam" }],
  openGraph: {
    title: "FoodFlow - Connecting Food Donors with Organizations",
    description: "Reduce food waste and help communities through our credit-based matching system.",
    type: "website",
  },
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  );
}
