import "./globals.css";
import type { Metadata } from "next";
import { Fredoka } from "next/font/google";
import ConditionalNavBar from "./components/ConditionalNavBar";

const fredoka = Fredoka({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Four-Colour Theorem Map Coloring",
  description:
    "A solver application for the famous Four-Colour Map Theorem. Initially built collaboratively during DevelopEd 2.0 Hackathon (2023) and later further developed.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={fredoka.className}>
      <body>
        <ConditionalNavBar />
        {children}
      </body>
    </html>
  );
}
